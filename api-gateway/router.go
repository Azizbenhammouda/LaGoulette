package main

import (
	"api-gateway/config"
	"fmt"
	"net/http"
	"strings"
)

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
