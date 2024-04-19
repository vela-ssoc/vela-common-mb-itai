package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strings"
	"sync"
	"text/template"

	"github.com/vela-ssoc/vela-common-mb-itai/dal/model"
)

type Storer interface {
	// Reset 重置指定 ID 的 value 缓存
	Reset(id string)

	// Shared 是否与 broker 节点共享，manager 节点用到该功能。
	Shared(id string) bool

	// Invalid 校验某个 ID 的参数是否无效
	Invalid(id string, val []byte) bool
	LocalAddr(ctx context.Context) (string, error)
	CmdbURL(ctx context.Context) (string, error)
	SsoURL(ctx context.Context) (string, error)
	AlarmURL(ctx context.Context) (string, error)
	Startup(ctx context.Context) (*model.Startup, error)
	Deploy(ctx context.Context, goos string, v any) *bytes.Buffer
	LoginDong(ctx context.Context, v any) (title, body string)
	RiskDong(ctx context.Context, v any) (title, body string)
	RiskHTML(ctx context.Context, v any) *bytes.Buffer
	EventDong(ctx context.Context, v any) (title, body string)
	EventHTML(ctx context.Context, v any) *bytes.Buffer
	RiskYW(ctx context.Context, v any) (title, body string)
	EventYW(ctx context.Context, v any) (title, body string)
}

func NewStore() Storer {
	ret := &storeDB{values: make(map[string]Valuer, 16)}

	{
		val := ret.newValue("global.local.addr", false, validIP)
		ret.values[val.ID()] = val
		ret.localAddr = val
	}
	{
		callback := "http://yunwei.eastmoney.com:81/api/v0.1/ci/s"
		val := ret.newHTTP("global.cmdb.url", true, callback)
		ret.values[val.ID()] = val
		ret.cmdbURL = val
	}
	{
		callback := "https://eastmoney-office.eastmoney.com/bd-cas/applogin?devTyp=pc"
		val := ret.newHTTP("global.sso.url", false, callback)
		ret.values[val.ID()] = val
		ret.ssoURL = val
	}
	{
		callback := "http://yunwei.eastmoney.com:81/api/v0.1/alerts"
		val := ret.newHTTP("global.alarm.url", true, callback)
		ret.values[val.ID()] = val
		ret.alarmURL = val
	}
	{
		val := ret.newValue("global.startup.param", true, validJSON)
		ret.values[val.ID()] = val
		ret.startupParam = val
	}

	// ============ 部署脚本
	{
		val := ret.newTmpl("global.deploy.linux.tmpl", true)
		ret.values[val.ID()] = val
		ret.deployLinuxTmpl = val
	}
	{
		val := ret.newTmpl("global.deploy.windows.tmpl", true)
		ret.values[val.ID()] = val
		ret.deployWindowsTmpl = val
	}
	{
		val := ret.newTmpl("global.deploy.darwin.tmpl", true)
		ret.values[val.ID()] = val
		ret.deployDarwinTmpl = val
	}
	{
		val := ret.newTmpl("global.deploy.android.tmpl", true)
		ret.values[val.ID()] = val
		ret.deployAndroidTmpl = val
	}

	// 咚咚登录验证码
	{
		val := ret.newTmpl("global.ding.title", false)
		ret.values[val.ID()] = val
		ret.loginDongTitle = val
	}
	{
		val := ret.newTmpl("global.ding.tmpl", false)
		ret.values[val.ID()] = val
		ret.loginDongBody = val
	}

	//  ============ event/risk
	{
		val := ret.newTmpl("global.risk.dong.title", true)
		ret.values[val.ID()] = val
		ret.riskDongTitle = val
	}
	{
		val := ret.newTmpl("global.risk.dong.tmpl", true, ret.trimCRLF)
		ret.values[val.ID()] = val
		ret.riskDongBody = val
	}
	{
		val := ret.newTmpl("global.risk.email.title", true)
		ret.values[val.ID()] = val
		ret.riskEmailTitle = val
	}
	{
		val := ret.newTmpl("global.risk.email.tmpl", true)
		ret.values[val.ID()] = val
		ret.riskEmailBody = val
	}
	{
		val := ret.newTmpl("global.risk.html.tmpl", false)
		ret.values[val.ID()] = val
		ret.riskHTML = val
	}
	{
		val := ret.newTmpl("global.event.dong.title", true)
		ret.values[val.ID()] = val
		ret.eventDongTitle = val
	}
	{
		val := ret.newTmpl("global.event.dong.tmpl", true, ret.trimCRLF)
		ret.values[val.ID()] = val
		ret.eventDongBody = val
	}
	{
		val := ret.newTmpl("global.event.email.tmpl", true)
		ret.values[val.ID()] = val
		ret.eventEmailTitle = val
	}
	{
		val := ret.newTmpl("global.event.email.tmpl", true)
		ret.values[val.ID()] = val
		ret.eventEmailBody = val
	}
	{
		val := ret.newTmpl("global.event.html.tmpl", false)
		ret.values[val.ID()] = val
		ret.eventHTML = val
	}
	{
		val := ret.newTmpl("global.event.yunwei.title", true)
		ret.values[val.ID()] = val
		ret.riskYWTitle = val
	}
	{
		val := ret.newTmpl("global.event.yunwei.tmpl", true)
		ret.values[val.ID()] = val
		ret.riskYWBody = val
	}
	{
		val := ret.newTmpl("global.risk.yunwei.title", true)
		ret.values[val.ID()] = val
		ret.eventYWTitle = val
	}
	{
		val := ret.newTmpl("global.risk.yunwei.tmpl", true)
		ret.values[val.ID()] = val
		ret.eventYWBody = val
	}

	return ret
}

type storeDB struct {
	mutex  sync.RWMutex
	values map[string]Valuer

	localAddr    Valuer
	cmdbURL      httpValuer
	ssoURL       httpValuer
	alarmURL     httpValuer
	startupParam Valuer

	deployLinuxTmpl   storeRender
	deployWindowsTmpl storeRender
	deployDarwinTmpl  storeRender
	deployAndroidTmpl storeRender

	loginDongTitle storeRender
	loginDongBody  storeRender

	riskDongTitle   storeRender
	riskDongBody    storeRender
	riskEmailTitle  storeRender
	riskEmailBody   storeRender
	riskYWTitle     storeRender
	riskYWBody      storeRender
	riskHTML        storeRender
	eventDongTitle  storeRender
	eventDongBody   storeRender
	eventEmailTitle storeRender
	eventEmailBody  storeRender
	eventYWTitle    storeRender
	eventYWBody     storeRender
	eventHTML       storeRender
}

func (sdb *storeDB) Reset(id string) {
	if val := sdb.values[id]; val != nil {
		val.Reset()
	}
}

func (sdb *storeDB) Shared(id string) bool {
	if val := sdb.values[id]; val != nil {
		return val.Shared()
	}
	return false
}

func (sdb *storeDB) Invalid(id string, dat []byte) bool {
	if val := sdb.values[id]; val != nil {
		return val.Invalid(dat)
	}
	return false
}

func (sdb *storeDB) LocalAddr(ctx context.Context) (string, error) {
	val, err := sdb.localAddr.Value(ctx)
	return string(val), err
}

func (sdb *storeDB) CmdbURL(ctx context.Context) (string, error) {
	return sdb.cmdbURL.Addr(ctx)
}

func (sdb *storeDB) SsoURL(ctx context.Context) (string, error) {
	return sdb.ssoURL.Addr(ctx)
}

func (sdb *storeDB) AlarmURL(ctx context.Context) (string, error) {
	return sdb.alarmURL.Addr(ctx)
}

func (sdb *storeDB) Startup(ctx context.Context) (*model.Startup, error) {
	val, err := sdb.startupParam.Value(ctx)
	if err != nil {
		return nil, err
	}
	ret := new(model.Startup)
	err = json.Unmarshal(val, ret)
	return ret, err
}

func (sdb *storeDB) LoginDong(ctx context.Context, v any) (string, string) {
	title := sdb.loginDongTitle.Rend(ctx, v)
	body := sdb.loginDongBody.Rend(ctx, v)
	return title.String(), body.String()
}

func (sdb *storeDB) RiskDong(ctx context.Context, v any) (string, string) {
	title := sdb.riskDongTitle.Rend(ctx, v)
	body := sdb.riskDongBody.Rend(ctx, v)
	return title.String(), body.String()
}

func (sdb *storeDB) RiskHTML(ctx context.Context, v any) *bytes.Buffer {
	return sdb.riskHTML.Rend(ctx, v)
}

func (sdb *storeDB) EventDong(ctx context.Context, v any) (string, string) {
	title := sdb.eventDongTitle.Rend(ctx, v)
	body := sdb.eventDongBody.Rend(ctx, v)
	return title.String(), body.String()
}

func (sdb *storeDB) EventHTML(ctx context.Context, v any) *bytes.Buffer {
	return sdb.eventHTML.Rend(ctx, v)
}

func (sdb *storeDB) RiskYW(ctx context.Context, v any) (string, string) {
	title := sdb.riskYWTitle.Rend(ctx, v)
	body := sdb.riskYWBody.Rend(ctx, v)
	return title.String(), body.String()
}

func (sdb *storeDB) EventYW(ctx context.Context, v any) (string, string) {
	title := sdb.eventYWTitle.Rend(ctx, v)
	body := sdb.eventYWBody.Rend(ctx, v)
	return title.String(), body.String()
}

func (sdb *storeDB) Deploy(ctx context.Context, goos string, v any) *bytes.Buffer {
	var rend storeRender
	switch strings.ToLower(goos) {
	case "linux":
		rend = sdb.deployLinuxTmpl
	case "windows":
		rend = sdb.deployWindowsTmpl
	case "darwin":
		rend = sdb.deployDarwinTmpl
	case "android":
		rend = sdb.deployAndroidTmpl
	default:
		msg := fmt.Sprintf("不支持的 %s 操作系统部署脚本", goos)
		return bytes.NewBufferString(msg)
	}

	return rend.Rend(ctx, v)
}

func (sdb *storeDB) newValue(id string, shared bool, valid func([]byte) bool) Valuer {
	return &storeValue{
		id:     id,
		shared: shared,
		valid:  valid,
	}
}

func (sdb *storeDB) newTmpl(id string, shared bool, filter ...func([]byte) []byte) storeRender {
	val := &storeValue{
		id:     id,
		shared: shared,
		valid:  validTmpl,
	}
	return &storeTemplate{
		value:  val,
		filter: filter,
	}
}

func (sdb *storeDB) newHTTP(id string, shared bool, cb string) httpValuer {
	val := &storeValue{
		id:     id,
		shared: shared,
		valid:  validHTTP,
	}
	return &httpValue{
		value:    val,
		callback: cb,
	}
}

func (sdb *storeDB) trimCRLF(dat []byte) []byte {
	return bytes.ReplaceAll(dat, []byte("\n"), nil)
}

func validIP(dat []byte) bool {
	ip := net.ParseIP(string(dat))
	return ip.IsUnspecified() || ip.IsLoopback()
}

func validJSON(dat []byte) bool {
	return !json.Valid(dat)
}

func validTmpl(dat []byte) bool {
	_, err := template.New("valid").Parse(string(dat))
	return err != nil
}

func validHTTP(dat []byte) bool {
	u, err := url.Parse(string(dat))
	if err != nil || u.Host == "" {
		return true
	}
	scheme := u.Scheme
	return scheme != "http" && scheme != "https"
}
