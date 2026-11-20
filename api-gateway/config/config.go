package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
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
	return &cfg, nil
}

func (c *Config) Validate() error {
	if c.Server.ReadHeaderTimeout <= 0 || c.Server.ReadTimeout <= 0 || c.Server.WriteTimeout <= 0 || c.Server.IdleTimeout <= 0 {
		return errors.New("timeouts are invalid")
	}
	if c.Server.ListenAddr == "" {
		return errors.New("listen address is empty")
	}
	for _, value := range c.Routes {
		if value.Prefix == "" {
			return fmt.Errorf("route prefix cannot be empty")
		}

		if !strings.HasPrefix(value.Prefix, "/") {
			return fmt.Errorf("route prefix must start with /")
		}

		if value.Upstream == "" {
			return fmt.Errorf("route upstream cannot be empty")
		}

		switch value.Protocol {
		case "rest", "graphql", "soap", "grpc":
		default:
			return fmt.Errorf("unsupported protocol: %s", value.Protocol)
		}
	}
	return nil
}
