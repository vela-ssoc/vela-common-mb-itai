package storage

import (
	"context"
	"sync"

	"github.com/vela-ssoc/vela-common-mb-itai/dal/query"
)

type Valuer interface {
	ID() string
	Value(ctx context.Context) ([]byte, error)
	Shared() bool
	Invalid([]byte) bool
	Reset()
}

type storeValue struct {
	id     string
	shared bool
	valid  func([]byte) bool
	mutex  sync.RWMutex
	done   bool
	value  []byte
	err    error
}

func (sv *storeValue) Invalid(dat []byte) bool {
	if fn := sv.valid; fn != nil {
		return fn(dat)
	}
	return false
}

func (sv *storeValue) ID() string {
	return sv.id
}

func (sv *storeValue) Value(ctx context.Context) ([]byte, error) {
	sv.mutex.RLock()
	done, value, err := sv.done, sv.value, sv.err
	sv.mutex.RUnlock()
	if done {
		return value, err
	}
	return sv.slowLoad(ctx)
}

func (sv *storeValue) Shared() bool {
	return sv.shared
}

func (sv *storeValue) Reset() {
	sv.mutex.Lock()
	sv.done = false
	sv.mutex.Unlock()
}

func (sv *storeValue) slowLoad(ctx context.Context) ([]byte, error) {
	sv.mutex.Lock()
	defer sv.mutex.Unlock()

	if sv.done {
		return sv.value, sv.err
	}

	tbl := query.Store
	dat, err := tbl.WithContext(ctx).Where(tbl.ID.Eq(sv.id)).First()
	if err != nil {
		sv.err = err
		sv.value = nil
		sv.done = true
		return nil, err
	}

	val := dat.Value
	sv.err = nil
	sv.value = val
	sv.done = true

	return val, nil
}
