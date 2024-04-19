package storage

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"net"
	"net/url"
	"sync"

	"github.com/vela-ssoc/vela-common-mb-itai/dal/model"
)

type Storer interface {
	Forget(id string) (shared bool)
	Validate(id string, data []byte) error

	LocalAddr(ctx context.Context) (string, error)
	CmdbURL(ctx context.Context) (*url.URL, error)
	SsoURL(ctx context.Context) (*url.URL, error)
	AlarmURL(ctx context.Context) (*url.URL, error)
	Startup(ctx context.Context) (*model.Startup, error)
	DeployScript(ctx context.Context, goos string, v any) *bytes.Buffer

	LoginDong(ctx context.Context, v any) (title, body string)
	RiskDong(ctx context.Context, v any, customID string) (title, body string)
	RiskHTML(ctx context.Context, v any) *bytes.Buffer
	EventDong(ctx context.Context, v any, customID string) (title, body string)
	EventHTML(ctx context.Context, v any) *bytes.Buffer
	// RiskYW(ctx context.Context, v any) (title, body string)
	// EventYW(ctx context.Context, v any) (title, body string)
}

const (
	uidLocalAddr       = "global.local.addr"
	uidCmdbURL         = "global.cmdb.url"
	uidSsoURL          = "global.sso.url"
	uidAlarmURL        = "global.alarm.url"
	uidStartupParam    = "global.startup.param"
	uidDeployLinux     = "global.deploy.linux.tmpl"
	uidDeployWindows   = "global.deploy.windows.tmpl"
	uidDeployDarwin    = "global.deploy.darwin.tmpl"
	uidDingTitle       = "global.ding.title"
	uidDingTmpl        = "global.ding.tmpl"
	uidRiskDongTitle   = "global.risk.dong.title"
	uidRiskDongTmpl    = "global.risk.dong.tmpl"
	uidRiskEmailTitle  = "global.risk.email.title"
	uidRiskEmailTmpl   = "global.risk.email.tmpl"
	uidRiskHTMLTmpl    = "global.risk.html.tmpl"
	uidEventDongTitle  = "global.risk.dong.title"
	uidEventDongTmpl   = "global.risk.dong.tmpl"
	uidEventEmailTitle = "global.risk.email.title"
	uidEventEmailTmpl  = "global.risk.email.tmpl"
	uidEventHTMLTmpl   = "global.risk.html.tmpl"
)

func NewStore() Storer {
	stores := make(map[string]valuer, 32)
	sdb := &storeDB{stores: stores}

	// 项目内置的参数
	{
		sdb.stores[uidLocalAddr] = &valueDB{uid: uidLocalAddr, valid: sdb.validIP}
	}
	{
		id := uidCmdbURL
		callback := "http://yunwei.eastmoney.com:81/api/v0.1/ci/s"
		val := sdb.createHTTP(id, true, callback)
		stores[id] = val
	}
	{
		id := uidSsoURL
		callback := "https://eastmoney-office.eastmoney.com/bd-cas/applogin?devTyp=pc"
		val := sdb.createHTTP(id, false, callback)
		stores[id] = val
	}
	{
		id := uidAlarmURL
		callback := "http://yunwei.eastmoney.com:81/api/v0.1/alerts"
		val := sdb.createHTTP(id, true, callback)
		stores[id] = val
	}
	{
		callback := `
{
    "node":{
        "dns":"114.114.114.114",
        "prefix":"share"
    },
    "logger":{
        "skip":1,
        "level":"error",
        "caller":true,
        "format":"text",
        "console":false,
        "filename":"logs/vela.error.log"
    },
    "console":{
        "script":"resource/script",
        "address":"127.0.0.1:306",
        "network":"tcp"
    }
}`
		id := uidStartupParam
		value := &valueDB{
			uid:      id,
			share:    true,
			callback: []byte(callback),
		}
		sdb.stores[id] = &startupValue{under: value}
	}

	{
		tid, bid := uidDingTitle, uidDingTmpl
		title := sdb.createTemplate(tid, false, sdb.filterCRLF)
		body := sdb.createTemplate(bid, false, sdb.filterCRLF)
		sdb.stores[tid] = title
		sdb.stores[bid] = body
	}
	{
		tid, bid := uidRiskDongTitle, uidRiskDongTmpl
		title := sdb.createTemplate(tid, true, sdb.filterCRLF)
		body := sdb.createTemplate(bid, true, sdb.filterCRLF)
		sdb.stores[tid] = title
		sdb.stores[bid] = body
	}
	{
		tid, bid := uidRiskEmailTitle, uidRiskEmailTmpl
		title := sdb.createTemplate(tid, true, nil)
		body := sdb.createTemplate(bid, true, nil)
		sdb.stores[tid] = title
		sdb.stores[bid] = body
	}
	{
		id := uidRiskHTMLTmpl
		tmpl := sdb.createTemplate(id, false, nil)
		sdb.stores[id] = tmpl
	}
	{
		tid, bid := uidEventDongTitle, uidEventDongTmpl
		title := sdb.createTemplate(tid, true, sdb.filterCRLF)
		body := sdb.createTemplate(bid, true, sdb.filterCRLF)
		sdb.stores[tid] = title
		sdb.stores[bid] = body
	}
	{
		tid, bid := uidEventEmailTitle, uidEventEmailTmpl
		title := sdb.createTemplate(tid, true, nil)
		body := sdb.createTemplate(bid, true, nil)
		sdb.stores[tid] = title
		sdb.stores[bid] = body
	}
	{
		id := uidEventHTMLTmpl
		tmpl := sdb.createTemplate(id, false, nil)
		sdb.stores[id] = tmpl
	}

	return sdb
}

type storeDB struct {
	mutex  sync.RWMutex
	stores map[string]valuer
}

func (sdb *storeDB) Forget(id string) (shared bool) {
	val := sdb.getValue(id)
	return val.forget()
}

func (sdb *storeDB) Validate(id string, data []byte) error {
	val := sdb.getValue(id)
	return val.validate(data)
}

func (sdb *storeDB) LocalAddr(ctx context.Context) (string, error) {
	val := sdb.getValue(uidLocalAddr)
	data, err := val.load(ctx)
	if err == nil && data != nil && len(data.Value) != 0 {
		return string(data.Value), nil
	}
	return "", err
}

func (sdb *storeDB) CmdbURL(ctx context.Context) (*url.URL, error) {
	return sdb.httpURL(ctx, uidCmdbURL)
}

func (sdb *storeDB) SsoURL(ctx context.Context) (*url.URL, error) {
	return sdb.httpURL(ctx, uidSsoURL)
}

func (sdb *storeDB) AlarmURL(ctx context.Context) (*url.URL, error) {
	return sdb.httpURL(ctx, uidAlarmURL)
}

func (sdb *storeDB) Startup(ctx context.Context) (*model.Startup, error) {
	id := uidStartupParam
	val := sdb.getValue(id)
	if sv, ok := val.(*startupValue); ok {
		return sv.loadStartup(ctx)
	}
	return nil, fmt.Errorf("store %s 不是有效的 startup 配置", id)
}

func (sdb *storeDB) DeployScript(ctx context.Context, goos string, v any) *bytes.Buffer {
	var id string
	switch goos {
	case "linux":
		id = uidDeployLinux
	case "windows":
		id = uidDeployWindows
	case "darwin":
		id = uidDeployDarwin
	default:
		id = "global.deploy." + goos + ".tmpl"
	}

	buf := new(bytes.Buffer)
	if err := sdb.templateRender(ctx, id, buf, v); err != nil {
		buf.WriteString(err.Error())
	}

	return buf
}

func (sdb *storeDB) LoginDong(ctx context.Context, v any) (string, string) {
	tb, bb := new(bytes.Buffer), new(bytes.Buffer)
	if err := sdb.templateRender(ctx, uidDingTitle, tb, v); err != nil {
		tb.WriteString(err.Error())
	}
	if err := sdb.templateRender(ctx, uidDingTmpl, bb, v); err != nil {
		bb.WriteString(err.Error())
	}

	return tb.String(), bb.String()
}

func (sdb *storeDB) RiskDong(ctx context.Context, v any, customID string) (string, string) {
	tid, bid := uidRiskDongTitle, uidRiskDongTmpl
	if customID != "" {
		bid = customID
	}

	tb, bb := new(bytes.Buffer), new(bytes.Buffer)
	if err := sdb.templateRender(ctx, tid, tb, v); err != nil {
		tb.WriteString(err.Error())
	}
	if err := sdb.templateRender(ctx, bid, bb, v); err != nil {
		bb.WriteString(err.Error())
	}

	return tb.String(), bb.String()
}

func (sdb *storeDB) RiskHTML(ctx context.Context, v any) *bytes.Buffer {
	bb := new(bytes.Buffer)
	if err := sdb.templateRender(ctx, uidRiskHTMLTmpl, bb, v); err != nil {
		bb.WriteString(err.Error())
	}

	return bb
}

func (sdb *storeDB) EventDong(ctx context.Context, v any, customID string) (string, string) {
	tid, bid := uidEventDongTitle, uidEventDongTmpl
	if customID != "" {
		bid = customID
	}

	tb, bb := new(bytes.Buffer), new(bytes.Buffer)
	if err := sdb.templateRender(ctx, tid, tb, v); err != nil {
		tb.WriteString(err.Error())
	}
	if err := sdb.templateRender(ctx, bid, bb, v); err != nil {
		bb.WriteString(err.Error())
	}

	return tb.String(), bb.String()
}

func (sdb *storeDB) EventHTML(ctx context.Context, v any) *bytes.Buffer {
	bb := new(bytes.Buffer)
	if err := sdb.templateRender(ctx, uidEventHTMLTmpl, bb, v); err != nil {
		bb.WriteString(err.Error())
	}

	return bb
}

func (sdb *storeDB) templateRender(ctx context.Context, id string, w io.Writer, v any) error {
	val := sdb.getValue(id)
	if tmpl, ok := val.(tmplValuer); ok {
		return tmpl.rend(ctx, w, v)
	}
	return fmt.Errorf("store %s 不是一个渲染模板", id)
}

func (sdb *storeDB) httpURL(ctx context.Context, id string) (*url.URL, error) {
	val := sdb.getValue(id)
	hv, ok := val.(httpValuer)
	if !ok {
		return nil, fmt.Errorf("store %s 不是一个 http 协议参数", id)
	}

	u, err := hv.loadURL(ctx)
	if err != nil {
		return nil, err
	}

	clone := &url.URL{
		Scheme:      u.Scheme,
		Opaque:      u.Opaque,
		User:        u.User,
		Host:        u.Host,
		Path:        u.Path,
		RawPath:     u.RawPath,
		OmitHost:    u.OmitHost,
		ForceQuery:  u.ForceQuery,
		RawQuery:    u.RawQuery,
		Fragment:    u.Fragment,
		RawFragment: u.RawFragment,
	}

	return clone, nil
}

func (sdb *storeDB) getValue(id string) valuer {
	sdb.mutex.RLock()
	v, ok := sdb.stores[id]
	sdb.mutex.RUnlock()
	if ok {
		return v
	}
	return sdb.createAndGet(id)
}

func (sdb *storeDB) createAndGet(id string) valuer {
	other := sdb.createTemplate(id, true, nil)
	sdb.mutex.Lock()
	v, ok := sdb.stores[id]
	if !ok {
		sdb.stores[id] = other
	}
	sdb.mutex.Unlock()
	if ok {
		return v
	}

	return other
}

func (sdb *storeDB) createTemplate(id string, shared bool, filter func([]byte) []byte) tmplValuer {
	under := &valueDB{
		uid:      id,
		share:    shared, // 默认与 broker 共享
		valid:    sdb.validateTmpl,
		callback: []byte("{{ . }}"),
		filter:   filter,
	}

	return &valueTmpl{
		under: under,
	}
}

func (sdb *storeDB) createHTTP(id string, shared bool, cb string) httpValuer {
	under := &valueDB{
		uid:      id,
		share:    shared,
		valid:    sdb.validHTTP,
		callback: []byte(cb),
	}
	return &valueHTTP{
		under: under,
	}
}

func (sdb *storeDB) validateTmpl(id string, data []byte) error {
	_, err := template.New(id).Parse(string(data))
	if err != nil {
		err = fmt.Errorf("store %s 不是一个正确的模板：%s", id, err.Error())
	}
	return err
}

func (sdb *storeDB) validIP(id string, dat []byte) error {
	ip := net.ParseIP(string(dat))
	if len(ip) == 0 || ip.IsUnspecified() || ip.IsLoopback() {
		return fmt.Errorf("store %s 不是一个合法的 IP", id)
	}
	return nil
}

func (sdb *storeDB) validHTTP(id string, data []byte) error {
	u, err := url.Parse(string(data))
	if err != nil {
		return fmt.Errorf("store %s 不是一个合法的 URL", id)
	}
	scheme := u.Scheme
	if scheme != "http" && scheme != "https" {
		return fmt.Errorf("store %s 不是一个 http 协议的 URL", id)
	}

	return nil
}

func (sdb *storeDB) filterCRLF(dat []byte) []byte {
	return bytes.ReplaceAll(dat, []byte("\n"), nil)
}
