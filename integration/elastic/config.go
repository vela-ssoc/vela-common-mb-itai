package elastic

import (
	"context"
	"encoding/base64"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	"github.com/vela-ssoc/vela-common-mb-itai/dal/query"
	"github.com/vela-ssoc/vela-common-mb-itai/problem"
)

func NewConfigure(name string) Configurer {
	return &elasticConfig{
		trip: http.DefaultTransport,
		name: name,
	}
}

type Configurer interface {
	Reset()
	Load(ctx context.Context) (http.Handler, error)
}

type elasticConfig struct {
	name   string
	trip   http.RoundTripper
	mutex  sync.RWMutex
	load   bool
	err    error
	handle http.Handler
}

func (ec *elasticConfig) Reset() {
	ec.mutex.Lock()
	ec.load = false
	ec.handle = nil
	ec.err = nil
	ec.mutex.Unlock()
}

func (ec *elasticConfig) Load(ctx context.Context) (http.Handler, error) {
	ec.mutex.RLock()
	load, h, err := ec.load, ec.handle, ec.err
	ec.mutex.RUnlock()
	if load {
		return h, err
	}

	return ec.slowLoad(ctx)
}

func (ec *elasticConfig) slowLoad(ctx context.Context) (http.Handler, error) {
	ec.mutex.Lock()
	defer ec.mutex.Unlock()
	if ec.load {
		return ec.handle, ec.err
	}

	tbl := query.Elastic
	data, err := tbl.WithContext(ctx).Where(tbl.Enable.Is(true)).First()
	if err != nil {
		ec.err = err
		ec.load = true
		return nil, err
	}

	hosts := data.Hosts
	if len(hosts) == 0 { // 兼容线上老数据
		hosts = []string{data.Host}
	}

	handle, err := ec.newNodes(hosts, data.Username, data.Password)
	ec.handle = handle
	ec.err = err
	ec.load = true

	return handle, err
}

func (ec *elasticConfig) newNodes(addrs []string, uname, passwd string) (http.Handler, error) {
	auth := "Basic " + base64.StdEncoding.EncodeToString([]byte(uname+":"+passwd))
	urls := make([]*url.URL, 0, len(addrs))
	for _, addr := range addrs {
		pu, err := url.Parse(addr)
		if err != nil {
			return nil, err
		}
		urls = append(urls, pu)
	}
	if len(urls) == 1 {
		node := ec.newNode(urls[0], auth)
		return node, nil
	}

	nodes := make([]http.Handler, 0, len(addrs))
	for _, u := range urls {
		node := ec.newNode(u, auth)
		nodes = append(nodes, node)
	}

	nano := time.Now().UnixNano()
	rd := &randomHandle{
		random:  rand.New(rand.NewSource(nano)),
		size:    len(nodes),
		handles: nodes,
	}

	return rd, nil
}

// newNode 初始化创建代理，支持 BasicAuth
func (ec *elasticConfig) newNode(u *url.URL, auth string) http.Handler {
	node := &nodeHandle{name: ec.name}
	rewriteFunc := func(r *httputil.ProxyRequest) {
		r.SetXForwarded()
		r.SetURL(u)
		r.Out.Header.Set("Authorization", auth)
	}

	px := &httputil.ReverseProxy{
		Transport:      ec.trip,
		Rewrite:        rewriteFunc,
		ModifyResponse: node.modifyResponse,
		ErrorHandler:   node.errorHandler,
	}
	node.px = px

	return node
}

type nodeHandle struct {
	name string
	px   *httputil.ReverseProxy
}

func (n *nodeHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	n.px.ServeHTTP(w, r)
}

func (*nodeHandle) modifyResponse(w *http.Response) error {
	if w.StatusCode == http.StatusUnauthorized {
		w.StatusCode = http.StatusBadRequest
	}
	return nil
}

func (n *nodeHandle) errorHandler(w http.ResponseWriter, r *http.Request, err error) {
	if e, ok := err.(*net.OpError); ok {
		// 隐藏后端服务 IP
		e.Addr = nil
		e.Net += " elasticsearch"
		err = e
	}

	pd := &problem.Detail{
		Type:     n.name,
		Title:    "代理错误",
		Status:   http.StatusBadRequest,
		Detail:   err.Error(),
		Instance: r.RequestURI,
	}
	_ = pd.JSON(w)
}

type randomHandle struct {
	random  *rand.Rand
	size    int
	handles []http.Handler
}

func (rd *randomHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	n := rd.random.Intn(rd.size)
	h := rd.handles[n]
	h.ServeHTTP(w, r)
}
