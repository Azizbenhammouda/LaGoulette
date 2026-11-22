package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Protocol string

const (
	ProtocolREST    Protocol = "rest"
	ProtocolGraphQL Protocol = "graphql"
	ProtocolGRPC    Protocol = "grpc"
	ProtocolSOAP    Protocol = "soap"
)

type ServerConfig struct {
	ListenAddr        string        `yaml:"listen_addr"`
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout"`
	ReadTimeout       time.Duration `yaml:"read_timeout"`
	WriteTimeout      time.Duration `yaml:"write_timeout"`
	IdleTimeout       time.Duration `yaml:"idle_timeout"`
}
type Route struct {
	Prefix                string        `yaml:"prefix"`
	Upstream              string        `yaml:"upstream"`
	Protocol              Protocol      `yaml:"protocol"`
	ResponseHeaderTimeout time.Duration `yaml:"response_header_timeout"`
}
type Config struct {
	Server ServerConfig `yaml:"server"`
	Routes []Route      `yaml:"routes"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("parsing file: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	return &cfg, nil
}

func (c *Config) Validate() error {
	if c.Server.ReadHeaderTimeout <= 0 {
		return fmt.Errorf("read header timeout must be > 0")
	}
	if c.Server.ReadTimeout <= 0 {
		return fmt.Errorf("read timeout must be > 0")
	}
	if c.Server.WriteTimeout <= 0 {
		return fmt.Errorf("write timeout must be > 0")
	}
	if c.Server.IdleTimeout <= 0 {
		return fmt.Errorf("idle timeout must be > 0")
	}
	if c.Server.ListenAddr == "" {
		return fmt.Errorf("listen address is empty")
	}

	for i, route := range c.Routes {
		if route.Prefix == "" {
			return fmt.Errorf("route[%d]: prefix cannot be empty ", i)
		}

		if !strings.HasPrefix(route.Prefix, "/") {
			return fmt.Errorf("route[%d] :prefix must start with /", i)
		}

		if route.Upstream == "" {
			return fmt.Errorf("route[%d] :upstream cannot be empty", i)
		}

		switch route.Protocol {
		case ProtocolREST, ProtocolGraphQL, ProtocolSOAP, ProtocolGRPC:
		default:
			return fmt.Errorf("route[%d] :unsupported protocol: %q", i, route.Protocol)
		}
		if route.ResponseHeaderTimeout <= 0 {
			return fmt.Errorf("route[%d]: response_header_timeout must be > 0", i)
		}
	}
	return nil
}
