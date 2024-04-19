package ssoauth

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/vela-ssoc/vela-common-mb-itai/logback"
	"github.com/vela-ssoc/vela-common-mba/netutil"
)

type Client interface {
	Auth(ctx context.Context, uname, passwd string) error
}

func NewClient(cfg Configurer, client netutil.HTTPClient, slog logback.Logger) Client {
	return &ssoClient{
		client: client,
		cfg:    cfg,
		slog:   slog,
	}
}

type ssoClient struct {
	client netutil.HTTPClient
	cfg    Configurer
	slog   logback.Logger
}

func (sc *ssoClient) Auth(ctx context.Context, uname, passwd string) error {
	addr, err := sc.cfg.Load(ctx)
	if err != nil {
		return err
	}

	sum := md5.Sum([]byte(passwd))
	pwd := hex.EncodeToString(sum[:])
	query := url.Values{"usrNme": []string{uname}, "passwd": []string{pwd}}
	// addr += query.Encode()
	if addr.RawQuery != "" {
		addr.RawQuery += "&"
	}
	addr.RawQuery += query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, addr.String(), nil)
	if err != nil {
		return err
	}

	res, err := sc.client.Do(req)
	if err != nil {
		return err
	}
	ret := new(ssoReply)
	if err = json.NewDecoder(res.Body).Decode(&ret); err != nil {
		return err
	}

	return ret.Error()
}
