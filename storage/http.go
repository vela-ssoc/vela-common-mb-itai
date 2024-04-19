package storage

import (
	"context"
	"strings"
	"sync"

	"gorm.io/gorm"
)

type httpValuer interface {
	Valuer
	Addr(ctx context.Context) (string, error)
}

type httpValue struct {
	value    Valuer
	callback string // 当数据库中不存在时，会使用该值作为配置
	mutex    sync.RWMutex
	done     bool
	addr     string
	err      error
}

func (hv *httpValue) ID() string {
	return hv.value.ID()
}

func (hv *httpValue) Value(ctx context.Context) ([]byte, error) {
	return hv.value.Value(ctx)
}

func (hv *httpValue) Shared() bool {
	return hv.value.Shared()
}

func (hv *httpValue) Invalid(dat []byte) bool {
	return hv.value.Invalid(dat)
}

func (hv *httpValue) Reset() {
	hv.mutex.Lock()
	hv.done = false
	hv.value.Reset()
	hv.mutex.Unlock()
}

func (hv *httpValue) Addr(ctx context.Context) (string, error) {
	hv.mutex.RLock()
	done, addr, err := hv.done, hv.addr, hv.err
	hv.mutex.RUnlock()
	if done {
		return addr, err
	}

	return hv.loadSlow(ctx)
}

func (hv *httpValue) loadSlow(ctx context.Context) (string, error) {
	hv.mutex.Lock()
	defer hv.mutex.Unlock()
	if hv.done {
		return hv.addr, hv.err
	}

	var addr string
	if val, err := hv.value.Value(ctx); err == nil {
		addr = string(val)
	} else {
		if err == gorm.ErrRecordNotFound && hv.callback != "" {
			addr = hv.callback
		} else {
			hv.err = err
			hv.done = true
			return "", err
		}
	}
	if strings.Contains(addr, "?") {
		addr += "&"
	} else {
		addr += "?"
	}
	hv.addr = addr
	hv.err = nil
	hv.done = true

	return addr, nil
}
