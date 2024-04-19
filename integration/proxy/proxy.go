package proxy

import (
	"net/http/httputil"
	"net/url"
)

func New(rawURL, token string) (*httputil.ReverseProxy, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	rewriteFunc := func(pr *httputil.ProxyRequest) {
		pr.SetURL(u)
		pr.Out.Header.Set("Authorization", token)
	}
	px := &httputil.ReverseProxy{
		Rewrite: rewriteFunc,
	}

	return px, nil
}
