package config

import "time"

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

func Load(path string) (*Config, error)
