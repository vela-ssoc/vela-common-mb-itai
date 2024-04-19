package storage

import (
	"context"
	"fmt"
	"net/url"
	"sync"

	"github.com/vela-ssoc/vela-common-mb-itai/dal/model"
)

type httpValuer interface {
	valuer
	loadURL(ctx context.Context) (*url.URL, error)
}

type valueHTTP struct {
	under  valuer
	mutex  sync.RWMutex
	loaded bool
	addr   *url.URL
	err    error
}

func (v *valueHTTP) id() string {
	return v.under.id()
}

func (v *valueHTTP) load(ctx context.Context) (*model.Store, error) {
	return v.under.load(ctx)
}

func (v *valueHTTP) forget() bool {
	v.mutex.Lock()
	v.loaded, v.addr = false, nil
	shared := v.under.forget()
	v.mutex.Unlock()

	return shared
}

func (v *valueHTTP) validate(dat []byte) error {
	return v.under.validate(dat)
}

func (v *valueHTTP) loadURL(ctx context.Context) (*url.URL, error) {
	v.mutex.RLock()
	loaded, addr, err := v.loaded, v.addr, v.err
	v.mutex.RUnlock()
	if loaded {
		return addr, err
	}

	return v.loadSlow(ctx)
}

func (v *valueHTTP) loadSlow(ctx context.Context) (*url.URL, error) {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	if v.loaded {
		return v.addr, v.err
	}

	uid := v.under.id()
	data, err := v.under.load(ctx)
	if err != nil {
		v.err = err
		v.loaded = true
		return nil, err
	}

	u, exx := url.Parse(string(data.Value))
	if exx != nil {
		exx = fmt.Errorf("store %s 不是一个有效的 URL", uid)
		v.err = exx
		v.loaded = true
		return nil, err
	}
	scheme := u.Scheme
	if scheme != "http" && scheme != "https" {
		exx = fmt.Errorf("store %s 不是一个 http 协议的 URL", uid)
		v.err = exx
		v.loaded = true
		return nil, err
	}

	v.err = nil
	v.addr = u
	v.loaded = true

	return u, nil
}
