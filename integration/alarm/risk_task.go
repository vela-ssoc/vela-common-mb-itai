package alarm

import (
	"context"
	"time"

	"github.com/vela-ssoc/vela-common-mb-itai/dal/model"
	"github.com/vela-ssoc/vela-common-mb-itai/dal/query"
)

type riskTask struct {
	unify *unifyAlert
	risk  *model.Risk
}

func (et *riskTask) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// 先入库保存
	rsk := et.risk
	tbl := query.Risk
	if err := tbl.WithContext(ctx).
		Create(rsk); err != nil || !rsk.SendAlert {
		return
	}

	// 发送告警
	key := rsk.FromCode
	sub := et.unify.match.Risk(ctx, key)
	if sub == nil || sub.Empty() {
		et.unify.slog.Infof("risk 风险 %s 没有订阅者", key)
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

func (et *riskTask) sendDong(ctx context.Context, dongs []string) {
	title, body := et.unify.store.RiskDong(ctx, et.risk, et.risk.Template)
	if err := et.unify.dong.Send(ctx, dongs, nil, title, body); err != nil {
		et.unify.slog.Warnf("发送风险 %s 失败：%s", dongs, err)
	} else {
		et.unify.slog.Infof("发送风险 %s 成功", dongs)
	}
}

func (et *riskTask) sendDevops(ctx context.Context, devs []*model.Devops) {
	// title, body := st.rend.RiskDong(ctx, st.dat, st.dat)
	err := et.unify.dps.Send(ctx, "告警", "内容", devs)
	et.unify.slog.Infof("发送 devops 风险 %s 结果：%v", devs, err)
}
