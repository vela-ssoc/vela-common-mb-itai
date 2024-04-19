package cmdb

import (
	"context"
	"net/url"

	"github.com/vela-ssoc/vela-common-mb-itai/storage/v2"
)

type Configurer interface {
	Load(ctx context.Context) (addr *url.URL, err error)
}

func NewConfigure(store storage.Storer) Configurer {
	return &fromDB{
		store: store,
	}
}

type fromDB struct {
	store storage.Storer
}

func (f *fromDB) Load(ctx context.Context) (*url.URL, error) {
	return f.store.CmdbURL(ctx)
}
