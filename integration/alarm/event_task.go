package alarm

import (
	"context"
	"time"

	"github.com/vela-ssoc/vela-common-mb-itai/dal/model"
	"github.com/vela-ssoc/vela-common-mb-itai/dal/query"
)

type eventTask struct {
	unify *unifyAlert
	event *model.Event
}

func (et *eventTask) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// 先入库保存
	evt := et.event
	tbl := query.Event
	if err := tbl.WithContext(ctx).
		Create(et.event); err != nil || !evt.SendAlert {
		return
	}

	// 发送告警
	key := evt.FromCode
	sub := et.unify.match.Event(ctx, key)
	if sub == nil || sub.Empty() {
		et.unify.slog.Infof("event 事件 %s 没有订阅者", key)
		return
	}

	// 发送告警
	if dongs := sub.Dong; len(dongs) != 0 {
		et.sendDong(ctx, dongs)
	}
	if devs := sub.Devops; len(devs) != 0 {
		et.sendDevops(ctx, devs)
	}
}

func (et *eventTask) sendDong(ctx context.Context, dongs []string) {
	title, body := et.unify.store.EventDong(ctx, et.event, "")
	if err := et.unify.dong.Send(ctx, dongs, nil, title, body); err != nil {
		et.unify.slog.Warnf("发送事件 %s 失败：%s", dongs, err)
	} else {
		et.unify.slog.Infof("发送事件成功")
	}
}

func (et *eventTask) sendDevops(ctx context.Context, devs []*model.Devops) {
	// title, body := st.rend.RiskDong(ctx, st.dat, st.dat)
	err := et.unify.dps.Send(ctx, "告警", "内容", devs)
	et.unify.slog.Infof("发送 devops 事件 %s 结果：%v", devs, err)
}
