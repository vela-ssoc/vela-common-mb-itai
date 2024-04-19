package ntfmatch

import (
	"context"
	"sync"

	"github.com/vela-ssoc/vela-common-mb-itai/dal/model"
	"github.com/vela-ssoc/vela-common-mb-itai/dal/query"
)

type Matcher interface {
	Reset()
	Event(ctx context.Context, evtCode string) *model.Subscriber
	Risk(ctx context.Context, rskCode string) *model.Subscriber
}

func NewMatch() Matcher {
	return &notifyMatch{}
}

type notifyMatch struct {
	mutex  sync.RWMutex
	loaded bool
	err    error
	subs   model.Subscribers
}

func (nm *notifyMatch) Reset() {
	nm.mutex.Lock()
	nm.loaded = false
	nm.mutex.Unlock()
}

func (nm *notifyMatch) Event(ctx context.Context, evtCode string) *model.Subscriber {
	return nm.load(ctx).Event(evtCode)
}

func (nm *notifyMatch) Risk(ctx context.Context, rskCode string) *model.Subscriber {
	return nm.load(ctx).Risk(rskCode)
}

func (nm *notifyMatch) load(ctx context.Context) model.Subscribers {
	nm.mutex.RLock()
	loaded, subs := nm.loaded, nm.subs
	nm.mutex.RUnlock()
	if loaded {
		return subs
	}

	return nm.slowLoad(ctx)
}

func (nm *notifyMatch) slowLoad(ctx context.Context) model.Subscribers {
	nm.mutex.Lock()
	defer nm.mutex.Unlock()
	if nm.loaded {
		return nm.subs
	}

	tbl := query.Notifier
	dats, err := tbl.WithContext(ctx).Find()
	if err != nil {
		nm.err = err
		nm.loaded = true
		return model.Subscribers{}
	}
	subs := model.Notifiers(dats).Subscribers()
	nm.subs = subs
	nm.loaded = true
	nm.err = nil

	return subs
}
