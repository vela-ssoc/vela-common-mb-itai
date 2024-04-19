package alarm

import (
	"context"
	"encoding/hex"
	"math/rand"
	"time"

	"github.com/vela-ssoc/vela-common-mb-itai/dal/model"
	"github.com/vela-ssoc/vela-common-mb-itai/gopool"
	"github.com/vela-ssoc/vela-common-mb-itai/integration/devops"
	"github.com/vela-ssoc/vela-common-mb-itai/integration/dong"
	"github.com/vela-ssoc/vela-common-mb-itai/integration/ntfmatch"
	"github.com/vela-ssoc/vela-common-mb-itai/logback"
	"github.com/vela-ssoc/vela-common-mb-itai/storage/v2"
)

type Alerter interface {
	EventSaveAndAlert(ctx context.Context, evt *model.Event) error
	RiskSaveAndAlert(ctx context.Context, rsk *model.Risk) error
}

func UnifyAlerter(store storage.Storer,
	match ntfmatch.Matcher,
	slog logback.Logger,
	dong dong.Client,
	dps devops.Client,
) Alerter {
	nano := time.Now().UnixNano()
	random := rand.New(rand.NewSource(nano))

	return &unifyAlert{
		store:  store,
		pool:   gopool.New(30, 1024, time.Minute),
		match:  match,
		slog:   slog,
		dong:   dong,
		dps:    dps,
		random: random,
	}
}

type unifyAlert struct {
	store  storage.Storer
	pool   gopool.Executor
	match  ntfmatch.Matcher
	slog   logback.Logger
	dong   dong.Client
	dps    devops.Client
	random *rand.Rand
}

func (ua *unifyAlert) EventSaveAndAlert(ctx context.Context, evt *model.Event) error {
	if evt == nil {
		ua.slog.Warn("event 为 nil 不作处理")
		return nil
	}

	// 入库前处理
	now := time.Now()
	if evt.OccurAt.IsZero() {
		evt.OccurAt = now
	}
	if evt.SendAlert {
		buf := make([]byte, 16)
		ua.random.Read(buf)
		evt.Secret = hex.EncodeToString(buf)
	}

	task := &eventTask{
		unify: ua,
		event: evt,
	}
	ua.pool.Submit(task)

	return nil
}

func (ua *unifyAlert) RiskSaveAndAlert(ctx context.Context, rsk *model.Risk) error {
	if rsk == nil {
		ua.slog.Warn("risk 为 nil 不作处理")
		return nil
	}

	// 入库前处理
	now := time.Now()
	rsk.Status = model.RSUnprocessed
	if rsk.OccurAt.IsZero() {
		rsk.OccurAt = now
	}
	if rsk.SendAlert {
		buf := make([]byte, 16)
		ua.random.Read(buf)
		rsk.Secret = hex.EncodeToString(buf)
	}

	task := &riskTask{
		unify: ua,
		risk:  rsk,
	}
	ua.pool.Submit(task)

	return nil
}
