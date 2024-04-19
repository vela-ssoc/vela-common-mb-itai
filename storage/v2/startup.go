package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/vela-ssoc/vela-common-mb-itai/dal/model"
)

type startupValue struct {
	under   valuer
	mutex   sync.RWMutex
	loaded  bool
	startup *model.Startup
	err     error
}

func (s *startupValue) id() string {
	return s.under.id()
}

func (s *startupValue) load(ctx context.Context) (*model.Store, error) {
	return s.under.load(ctx)
}

func (s *startupValue) forget() bool {
	s.mutex.Lock()
	s.loaded = false
	shared := s.under.forget()
	s.mutex.Unlock()
	return shared
}

func (s *startupValue) validate(data []byte) error {
	uid := s.id()
	mod := new(model.Startup)
	if err := json.Unmarshal(data, mod); err != nil {
		return fmt.Errorf("store %s 不是有效的 startup 配置：%s", uid, err.Error())
	}
	return nil
}

func (s *startupValue) loadStartup(ctx context.Context) (*model.Startup, error) {
	s.mutex.RLock()
	loaded, startup, err := s.loaded, s.startup, s.err
	s.mutex.RUnlock()
	if loaded {
		return startup, err
	}

	return s.slowLoadStartup(ctx)
}

func (s *startupValue) slowLoadStartup(ctx context.Context) (*model.Startup, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.loaded {
		return s.startup, s.err
	}

	data, err := s.under.load(ctx)
	if err != nil {
		s.err = err
		s.loaded = true
		return nil, err
	}

	startup := new(model.Startup)
	err = json.Unmarshal(data.Value, startup)
	s.err = err
	s.startup = startup
	s.loaded = true

	return startup, err
}
