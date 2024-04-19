package vulnsync

import (
	"context"
	"errors"
	"sync/atomic"
	"time"

	"github.com/vela-ssoc/vela-common-mb-itai/dal/model"
	"github.com/vela-ssoc/vela-common-mb-itai/integration/sonatype"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrSynchronizing = errors.New("漏洞正在同步")

func New(db *gorm.DB, sona sonatype.Searcher) *Synchro {
	return &Synchro{
		db:   db,
		sona: sona,
	}
}

// Synchro 漏洞同步处理器
type Synchro struct {
	db      *gorm.DB          // 数据库连接
	sona    sonatype.Searcher // 漏洞库查询
	syncing atomic.Bool       // 是否正在同步漏洞数据库
	fulling atomic.Bool       // 是否正在全表扫描预处理数据
}

// Light 只同步更新漏洞数据
// 根据系统中的现存的 purl 查询漏洞，获取最新的漏洞数据
func (syn *Synchro) Light(ctx context.Context) error {
	nonce := time.Now().UnixNano()

	return syn.vulnerability(ctx, false, nonce)
}

// vulnerability 查询更新漏洞数据
// pre 是否在查询漏洞的同时更新组件中的漏洞统计信息（会增加程序的处理耗时）
func (syn *Synchro) vulnerability(ctx context.Context, pre bool, nonce int64) error {
	// 一般而言：漏洞同步全表扫描组件中的 purl，然后联网查询 purl 的漏洞数据。
	// 整个操作即耗性能又耗时间，全局同时有一个同步操作就够了，防止并发调用。
	if !syn.syncing.CompareAndSwap(false, true) {
		return ErrSynchronizing
	}
	defer syn.syncing.Store(false)

	const limit = 100
	var offset int
	for {
		var cps []*model.SBOMComponent
		if err := syn.db.Distinct("purl").Order("purl").
			Offset(offset).Limit(limit).
			Find(&cps).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			} else {
				return err
			}
		}
		offset += limit

		size := len(cps)
		if size == 0 {
			break
		}
		purls := make([]string, 0, size)
		for _, cp := range cps {
			purls = append(purls, cp.PURL)
		}

		// 查询漏洞库（若运行出错也不影响下一批漏洞查询）
		if vns, err := syn.sona.Search(ctx, purls); err == nil && len(vns) != 0 {
			_ = syn.save(pre, vns, nonce) // 保存入库
		}

		if size < limit { // 查询的条数小于 limit 说明没有数据了
			break
		}
	}

	return nil
}

// save 将查询出的漏洞数据保存到数据库
func (syn *Synchro) save(pre bool, vns []*model.SBOMVuln, nonce int64) error {
	size := len(vns)
	if size == 0 {
		return nil
	}

	hm := make(map[string]*scoreCounter, size)
	for _, vn := range vns {
		vn.Nonce = nonce
		if pre {
			purl := vn.PURL
			ct := hm[purl]
			if ct == nil {
				ct = &scoreCounter{Nonce: nonce}
				hm[purl] = ct
			}
			ct.Put(vn.Score)
		}
	}

	// 保存漏洞数据
	cla := clause.OnConflict{DoUpdates: clause.AssignmentColumns([]string{"nonce"})}
	err := syn.db.Clauses(cla).Create(vns).Error
	if !pre {
		return err
	}

	// FIXME: 在循环内执行更新 SQL 需要慎用，但暂时没有其他好的解决方案
	for purl, sct := range hm {
		cols := sct.Columns(nonce)
		syn.db.Model(&model.SBOMComponent{}).Where("purl = ?", purl).Updates(cols)
	}

	return nil
}

type scoreCounter struct {
	CriticalNum   int             `gorm:"column:critical_num"`
	CriticalScore model.CVSSScore `gorm:"column:critical_score"`
	HighNum       int             `gorm:"column:high_num"`
	HighScore     model.CVSSScore `gorm:"column:high_score"`
	MediumNum     int             `gorm:"column:medium_num"`
	MediumScore   model.CVSSScore `gorm:"column:medium_score"`
	LowNum        int             `gorm:"column:low_num"`
	LowScore      model.CVSSScore `gorm:"column:low_score"`
	TotalNum      int             `gorm:"column:total_num"`
	TotalScore    model.CVSSScore `gorm:"column:total_score"`
	Nonce         int64           `gorm:"column:nonce"`
}

func (sct *scoreCounter) Put(score model.CVSSScore) {
	lvl := score.Level()
	switch lvl {
	case model.CVSSCritical:
		sct.CriticalNum++
		sct.CriticalScore += score
	case model.CVSSHigh:
		sct.HighNum++
		sct.HighScore += score
	case model.CVSSMedium:
		sct.MediumNum++
		sct.MediumScore += score
	case model.CVSSLow:
		sct.LowNum++
		sct.LowScore += score
	default:
		return
	}
	sct.TotalNum++
	sct.TotalScore += score
}

func (sct *scoreCounter) Columns(nonce int64) map[string]any {
	return map[string]any{
		"critical_num":   sct.CriticalNum,
		"critical_score": sct.CriticalScore,
		"high_num":       sct.HighNum,
		"high_score":     sct.HighScore,
		"medium_num":     sct.MediumNum,
		"medium_score":   sct.MediumScore,
		"low_num":        sct.LowNum,
		"low_score":      sct.LowScore,
		"total_num":      sct.TotalNum,
		"total_score":    sct.TotalScore,
		"nonce":          nonce,
	}
}
