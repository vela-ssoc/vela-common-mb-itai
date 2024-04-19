package alarm

import (
	"context"
	"time"

	"github.com/vela-ssoc/vela-common-mb-itai/dal/model"
	"github.com/vela-ssoc/vela-common-mb-itai/integration/devops"
	"github.com/vela-ssoc/vela-common-mb-itai/integration/dong"
	"github.com/vela-ssoc/vela-common-mb-itai/logback"
	"github.com/vela-ssoc/vela-common-mb-itai/storage"
)

type sendTask struct {
	dat   any
	sub   *model.Subscriber
	store storage.Storer
	slog  logback.Logger
	dong  dong.Client
	dps   devops.Client
}

func (st *sendTask) Run() {
	// 发送告警
	if dongs := st.sub.Dong; len(dongs) != 0 {
		st.sendDong(dongs)
	}
	if devs := st.sub.Devops; len(devs) != 0 {
		st.sendDevops(devs)
	}
}

func (st *sendTask) sendDong(dongs []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	title, body := st.store.EventDong(ctx, st.dat)
	if err := st.dong.Send(ctx, dongs, nil, title, body); err != nil {
		st.slog.Warnf("发送风险 %s 失败：%s", dongs, err)
	} else {
		st.slog.Infof("发送风险 %s 成功")
	}
}

func (st *sendTask) sendDevops(devs []*model.Devops) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// title, body := st.rend.RiskDong(ctx, st.dat, st.dat)
	err := st.dps.Send(ctx, "告警", "内容", devs)
	st.slog.Infof("发送风险 %s 结果：%v", devs, err)
}
