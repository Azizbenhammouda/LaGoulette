package main

import (
	"api-gateway/config"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	mux, err := NewRouter(cfg)
	if err != nil {
		log.Fatalf("starting proxy: %v", err)
	}
	s := &http.Server{
		Addr:              cfg.Server.ListenAddr,
		ReadHeaderTimeout: cfg.Server.ReadHeaderTimeout,
		WriteTimeout:      cfg.Server.WriteTimeout,
		ReadTimeout:       cfg.Server.ReadTimeout,
		IdleTimeout:       cfg.Server.IdleTimeout,
		Handler:           mux,
	}
	fmt.Println("gateway running..")
	err = s.ListenAndServe()
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			log.Println("server closed")
		} else {
			log.Fatal(err)
		}
	}
}
