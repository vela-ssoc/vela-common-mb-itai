package devops

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/vela-ssoc/vela-common-mb-itai/dal/model"
	"github.com/vela-ssoc/vela-common-mba/netutil"
)

type Client interface {
	Send(ctx context.Context, title, body string, users []*model.Devops) error
}

func NewClient(cfg Configurer, cli netutil.HTTPClient) Client {
	return &devopsClient{
		cfg: cfg,
		cli: cli,
	}
}

type devopsClient struct {
	cfg Configurer
	cli netutil.HTTPClient
}

func (dc *devopsClient) Send(ctx context.Context, title, body string, users []*model.Devops) error {
	addr, err := dc.cfg.Load(ctx)
	if err != nil {
		return err
	}

	dat, _ := json.Marshal(users)
	req := &request{
		OriginName:     "security-alert",
		AlertType:      "security-alert",
		AlertObject:    title,
		AlertAttribute: "ssoc",
		Severity:       "disaster",
		Subject:        title,
		Body:           body,
		Notifier:       string(dat),
	}

	return dc.cli.JSON(ctx, http.MethodPost, addr.String(), req, &struct{}{}, nil)
}

type request struct {
	OriginName     string `json:"origin_name"`
	AlertType      string `json:"alert_type"`
	AlertObject    string `json:"alert_object"`
	AlertAttribute string `json:"alert_attribute"`
	Severity       string `json:"severity"`
	Subject        string `json:"subject"`
	Body           string `json:"body"`
	Notifier       string `json:"notifier"`
}
