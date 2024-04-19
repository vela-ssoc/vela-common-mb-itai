package dong

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/vela-ssoc/vela-common-mb-itai/logback"
	"github.com/vela-ssoc/vela-common-mba/netutil"
)

type Client interface {
	Send(ctx context.Context, uids, gids []string, title, body string) error
}

func NewClient(cfg Configurer, cli netutil.HTTPClient, slog logback.Logger) Client {
	return &client{
		cfg:     cfg,
		client:  cli,
		slog:    slog,
		timeout: 5 * time.Second,
	}
}

type client struct {
	cfg     Configurer
	client  netutil.HTTPClient
	slog    logback.Logger
	timeout time.Duration
}

func (c *client) Send(parent context.Context, uids, gids []string, title, detail string) error {
	if len(uids) == 0 && len(gids) == 0 {
		c.slog.Warn("咚咚告警的用户 ID 和群组 ID 全部为空，跳过不发送")
		return nil
	}
	if parent == nil {
		parent = context.Background()
	}

	ctx, cancel := context.WithTimeout(parent, c.timeout)
	defer cancel()

	req, err := c.newRequest(ctx, uids, gids, title, detail)
	if err != nil {
		return err
	}
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer res.Body.Close()

	ret := new(responseBody)
	if err = json.NewDecoder(res.Body).Decode(ret); err != nil {
		return err
	}

	return ret.Error()
}

func (c *client) newRequest(ctx context.Context, uids, gids []string, title, detail string) (*http.Request, error) {
	addr, account, token, err := c.cfg.Load(ctx)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	body := &requestBody{
		UserIDs:  strings.Join(uids, ","),
		GroupIDs: strings.Join(gids, ","),
		Title:    title,
		Detail:   detail,
	}
	if err = json.NewEncoder(buf).Encode(body); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, addr, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Account", account)
	req.Header.Set("Token", token)

	return req, nil
}
