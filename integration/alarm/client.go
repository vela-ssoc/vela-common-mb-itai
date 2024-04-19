package alarm

import (
	"context"

	"github.com/vela-ssoc/vela-common-mb-itai/dal/model"
	"github.com/vela-ssoc/vela-common-mb-itai/gopool"
	"github.com/vela-ssoc/vela-common-mb-itai/integration/ntfmatch"
	"github.com/vela-ssoc/vela-common-mb-itai/logback"
)

type Client interface {
	RiskSubscriber(ctx context.Context, fc string) *model.Subscriber
	EventSubscriber(ctx context.Context, fc string) *model.Subscriber

	Dong(ctx context.Context, userIDs []string, title, body string) error
	Email(ctx context.Context, userIDs []string, title, body string) error
}

type alarmClient struct {
	pool  gopool.Executor
	match ntfmatch.Matcher
	slog  logback.Logger
}

func (ac *alarmClient) RiskSubscriber(ctx context.Context, fc string) *model.Subscriber {
	return ac.match.Risk(ctx, fc)
}

func (ac *alarmClient) EventSubscriber(ctx context.Context, fc string) *model.Subscriber {
	return ac.match.Event(ctx, fc)
}

func (ac *alarmClient) Send() {
}
