package main

import (
	"api-gateway/config"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// When a subtree handler exists for /api/users/, and a request arrives for /api/users,
//
//	it automatically redirects the request to /api/users/.
//
// its servermux behaviour
func NewRouter(cfg *config.Config) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthHandler)
	for _, route := range cfg.Routes {
		prefix := route.Prefix

		if !strings.HasSuffix(prefix, "/") {
			prefix += "/"
		}

		mux.HandleFunc("GET "+prefix, func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(
				w,
				"matched %s -> %s (%s)\n",
				route.Prefix,
				route.Upstream,
				route.Protocol,
			)
		})
	}

	return mux
}
func NewProxy(cfg *config.Config) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(cfg.Routes[0].Upstream)
	if err != nil {
		return nil, err
	}
	Reverse := httputil.NewSingleHostReverseProxy(url)
	return Reverse, nil
}
