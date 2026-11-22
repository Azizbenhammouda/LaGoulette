package main

import (
	"api-gateway/config"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
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
		proxy, err := NewProxy(route)
		if err != nil {
			return nil, err
		}
		mux.Handle(prefix, proxy)
	}

	return mux, nil
}
func NewProxy(cfr config.Route) (*httputil.ReverseProxy, error) {
	targeturl, err := url.Parse(cfr.Upstream)
	if err != nil {
		return nil, err
	}
	transport := &http.Transport{
		// time to wait to establish a tcp connection
		DialContext: (&net.Dialer{
			Timeout: time.Second * 5,
		}).DialContext,
		// time to wait for the upstream response header to arrive
		ResponseHeaderTimeout: cfr.ResponseHeaderTimeout,
	}
	Reverse := &httputil.ReverseProxy{
		// rewrite concept el global ybdel el request eli jaya ml client
		//  w el proxy yebaath request jdida
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
		Transport: transport,
	}
	return Reverse, nil
}
