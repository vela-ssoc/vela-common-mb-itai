package dong

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/vela-ssoc/vela-common-mb-itai/dal/model"
	"github.com/vela-ssoc/vela-common-mb-itai/dal/query"
	"gorm.io/gorm"
)

func NewConfig() Configurer {
	return &emcConfig{}
}

type Configurer interface {
	Forget()
	Load(ctx context.Context) (addr, account, token string, err error)
}

type emcConfig struct {
	mutex  sync.RWMutex
	loaded bool
	err    error
	data   *model.Emc
}

func (ec *emcConfig) Forget() {
	ec.mutex.Lock()
	ec.loaded = false
	ec.mutex.Unlock()
}

func (ec *emcConfig) Load(ctx context.Context) (string, string, string, error) {
	ec.mutex.RLock()
	loaded, data, err := ec.loaded, ec.data, ec.err
	ec.mutex.RUnlock()
	if loaded {
		if err != nil {
			return "", "", "", err
		}
		return data.Host, data.Account, data.Token, nil
	}

	return ec.slowLoad(ctx)
}

func (ec *emcConfig) slowLoad(ctx context.Context) (string, string, string, error) {
	ec.mutex.Lock()
	defer ec.mutex.Unlock()

	if ec.loaded {
		if err := ec.err; err != nil {
			return "", "", "", err
		}
		data := ec.data
		return data.Host, data.Account, data.Token, nil
	}

	tbl := query.Emc
	dat, err := tbl.WithContext(ctx).
		Where(tbl.Enable.Is(true)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = fmt.Errorf("没有找到咚咚服务号配置")
		}
		ec.err = err
		ec.data = nil
		ec.loaded = true
		return "", "", "", err
	}

	ec.err = nil
	ec.data = dat
	ec.loaded = true

	return dat.Host, dat.Account, dat.Token, nil
}
