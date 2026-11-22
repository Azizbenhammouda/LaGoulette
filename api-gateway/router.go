package main

import (
	"api-gateway/config"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// waqt fama handler for /api/users/, w yji request for /api/users,
//
//	it automatically redirects the request to /api/users/.
//
// its servermux behaviour
func NewRouter(cfg *config.Config) (http.Handler, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthHandler)
	for _, route := range cfg.Routes {
		prefix := route.Prefix

		if !strings.HasSuffix(prefix, "/") {
			prefix += "/"
		}
		proxy, err := NewProxy(route.Upstream)
		if err != nil {
			return nil, err
		}
		mux.Handle(prefix, proxy)
	}

	return mux, nil
}
func NewProxy(upstream string) (*httputil.ReverseProxy, error) {
	targeturl, err := url.Parse(upstream)
	if err != nil {
		return nil, err
	}
	Reverse := &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			pr.SetURL(targeturl) // bch nbadel el url eli jey ml request li targeturl
			// zada yzid el path
			pr.Out.Host = targeturl.Host // bch badel el host eli jey ml request yabda
			//  el host howa el apigateway bch nabdlou li el service
			// normalement setUrl deja ybdel el host header lil host mtaa target url
			// zdtha just be sure
			pr.SetXForwarded() //make sure the services recieve the actual clients Ip
			// the host the client originally asked for and weather its http or https
		},
	}
	return Reverse, nil
}
