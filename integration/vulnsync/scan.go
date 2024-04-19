package vulnsync

import (
	"context"
	"time"

	"github.com/vela-ssoc/vela-common-mb-itai/dal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Scan 全表扫描并预处理漏洞数据
// sync 扫描前是否先同步漏洞库
func (syn *Synchro) Scan(ctx context.Context, sync bool) error {
	if !syn.fulling.CompareAndSwap(false, true) {
		return ErrSynchronizing
	}
	defer syn.fulling.Store(false)

	nonce := time.Now().UnixNano()
	if sync {
		if err := syn.vulnerability(ctx, true, nonce); err != nil {
			return err
		}
	}

	// 1. 对 sbom_project 表全表扫描预处理
	var pts []*model.SBOMProject
	syn.db.Select("id").FindInBatches(&pts, 100, func(tx *gorm.DB, _ int) error {
		projectIDs := make([]int64, 0, len(pts))
		for _, pt := range pts {
			projectIDs = append(projectIDs, pt.ID)
		}
		syn.projects(projectIDs, nonce)
		return nil
	})

	// 2. 对 sbom_minion 表全表扫描预处理
	const limit = 100
	var offset int
	for {
		var mps []*model.SBOMProject
		syn.db.Distinct("minion_id").Order("minion_id").Offset(offset).Limit(limit).Find(&mps)
		offset += limit
		if len(mps) == 0 {
			break
		}

		minionIDs := make([]int64, 0, len(mps))
		for _, mp := range mps {
			minionIDs = append(minionIDs, mp.MinionID)
		}

		syn.minions(minionIDs, nonce)

		if len(minionIDs) < limit {
			break
		}
	}

	return nil
}

func (syn *Synchro) projects(pids []int64, nonce int64) {
	rawSQL := "SELECT SUM(critical_num)   AS critical_num,  " +
		" SUM(critical_score)             AS critical_score," +
		" SUM(high_num)                   AS high_num,      " +
		" SUM(high_score)                 AS high_score,    " +
		" SUM(medium_num)                 AS medium_num,    " +
		" SUM(medium_score)               AS medium_score,  " +
		" SUM(low_num)                    AS low_num,       " +
		" SUM(low_score)                  AS low_score,     " +
		" SUM(total_num)                  AS total_num,     " +
		" SUM(total_score)                AS total_score    " +
		" FROM sbom_component " +
		" WHERE project_id = ?"
	for _, pid := range pids {
		sct := &scoreCounter{Nonce: nonce}
		syn.db.Raw(rawSQL, pid).Scan(sct)
		cols := sct.Columns(nonce)
		syn.db.Model(&model.SBOMProject{}).Where("id = ?", pid).Updates(cols)
	}
}

func (syn *Synchro) minions(minionIDs []int64, nonce int64) {
	rawSQL := "SELECT SUM(critical_num)   AS critical_num,  " +
		" SUM(critical_score)             AS critical_score," +
		" SUM(high_num)                   AS high_num,      " +
		" SUM(high_score)                 AS high_score,    " +
		" SUM(medium_num)                 AS medium_num,    " +
		" SUM(medium_score)               AS medium_score,  " +
		" SUM(low_num)                    AS low_num,       " +
		" SUM(low_score)                  AS low_score,     " +
		" SUM(total_num)                  AS total_num,     " +
		" SUM(total_score)                AS total_score    " +
		" FROM sbom_project " +
		" WHERE minion_id = ?"

	for _, mid := range minionIDs {
		sct := new(scoreCounter)
		syn.db.Raw(rawSQL, mid).Scan(sct)
		// 查询节点的 inet
		var min model.Minion
		syn.db.Select("inet").Where("id = ?", mid).First(&min)

		mon := &model.SBOMMinion{
			ID:            mid,
			Inet:          min.Inet,
			CriticalNum:   sct.CriticalNum,
			CriticalScore: sct.CriticalScore,
			HighNum:       sct.HighNum,
			HighScore:     sct.HighScore,
			MediumNum:     sct.MediumNum,
			MediumScore:   sct.MediumScore,
			LowNum:        sct.LowNum,
			LowScore:      sct.LowScore,
			TotalNum:      sct.TotalNum,
			TotalScore:    sct.TotalScore,
			Nonce:         nonce,
		}

		syn.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(mon)
	}
}
