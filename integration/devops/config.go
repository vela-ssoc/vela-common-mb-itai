package devops

import (
	"context"
	"net/url"

	"github.com/vela-ssoc/vela-common-mb-itai/storage/v2"
)

type Configurer interface {
	Load(ctx context.Context) (*url.URL, error)
}

func NewConfig(store storage.Storer) Configurer {
	return &config{
		store: store,
	}
}

type config struct {
	store storage.Storer
}

func (cfg *config) Load(ctx context.Context) (*url.URL, error) {
	return cfg.store.AlarmURL(ctx)
}
