package main

import (
	"errors"
	"log"
	"net/http"
	"time"
)

type ServerConfig struct {
	ListenAddr        string        `yaml:"listen_addr"`
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout"`
	ReadTimeout       time.Duration `yaml:"read_timeout"`
	WriteTimeout      time.Duration `yaml:"write_timeout"`
	IdleTimeout       time.Duration `yaml:"idle_timeout"`
}
type Route struct {
	Prefix   string `yaml:"prefix"`
	Upstream string `yaml:"upstream"`
	Protocol string `yaml:"protocol"`
}
type Config struct {
	Server ServerConfig `yaml:"server"`
	Routes []Route      `yaml:"routes"`
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthHandler)
	serverCfg := ServerConfig{
		ListenAddr:        ":8000",
		ReadHeaderTimeout: time.Second * 5,
		WriteTimeout:      time.Second * 15,
		ReadTimeout:       time.Second * 15,
		IdleTimeout:       time.Second * 60,
	}
	s := &http.Server{
		Addr:              serverCfg.ListenAddr,
		ReadHeaderTimeout: serverCfg.ReadHeaderTimeout,
		WriteTimeout:      serverCfg.WriteTimeout,
		ReadTimeout:       serverCfg.ReadTimeout,
		IdleTimeout:       serverCfg.IdleTimeout,
		Handler:           mux,
	}
	err := s.ListenAndServe()
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			log.Println("server closed")
		} else {
			log.Fatal(err)
		}
	}
}
