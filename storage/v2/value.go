package storage

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/vela-ssoc/vela-common-mb-itai/dal/model"
	"github.com/vela-ssoc/vela-common-mb-itai/dal/query"
	"gorm.io/gorm"
)

type valuer interface {
	id() string
	load(ctx context.Context) (*model.Store, error)
	forget() (shared bool)
	validate([]byte) error
}

type valueDB struct {
	// uid 数据库唯一 id
	uid string

	// share 是否与 broker 节点共享
	share bool

	// valid 校验方法
	valid func(id string, val []byte) error

	// callback 当 load 加载的数据为空或加载错误时会调用该方法
	callback []byte

	// filter 加载到数据后是否二次过滤处理
	filter func([]byte) []byte

	// mutex 数据锁
	mutex sync.RWMutex

	// loaded 是否已经加载
	loaded bool

	// value 数据
	value *model.Store

	// err 是否从数据库加载出错了
	err error
}

func (v *valueDB) id() string {
	return v.uid
}

func (v *valueDB) load(ctx context.Context) (*model.Store, error) {
	v.mutex.RLock()
	loaded, value, err := v.loaded, v.value, v.err
	v.mutex.RUnlock()
	if loaded {
		return value, err
	}

	return v.loadDB(ctx)
}

func (v *valueDB) forget() (shared bool) {
	v.mutex.Lock()
	v.loaded = false
	v.mutex.Unlock()
	return v.share
}

func (v *valueDB) validate(dat []byte) error {
	if vfn := v.valid; vfn != nil {
		return vfn(v.uid, dat)
	}
	return nil
}

func (v *valueDB) loadDB(ctx context.Context) (*model.Store, error) {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	if v.loaded {
		return v.value, v.err
	}

	value, err := v.loadValue(ctx)
	v.value = value
	v.err = err
	v.loaded = true

	return value, err
}

func (v *valueDB) loadValue(ctx context.Context) (*model.Store, error) {
	uid := v.uid
	var value []byte
	tbl := query.Store
	dat, err := tbl.WithContext(ctx).Where(tbl.ID.Eq(uid)).First()
	if err == nil && dat != nil {
		value = dat.Value
	}
	if len(value) == 0 {
		value = v.callback
	}
	if ffn := v.filter; len(value) != 0 && ffn != nil {
		value = ffn(value)
	}
	if len(value) != 0 {
		dat.Value = value
		return dat, nil
	}

	if err == nil || errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("尚未配置 store %s 数据", uid)
	}

	return nil, fmt.Errorf("查询 store %s 数据错误：%s", uid, err.Error())
}
