package config

import (
	"fmt"
	"strings"
)

type Config struct {
	ServiceName  string
	FallbackPath string
	BufferSize   int // Буфер очереди (default: 128)
}

func (c *Config) MustLoadConfig() error {
	var missing []string

	if c.ServiceName == "" {
		missing = append(missing, "service name")
	}

	if c.FallbackPath == "" {
		missing = append(missing, "fallback path")
	}

	if c.BufferSize == 0 {
		c.BufferSize = 128
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing config values: %s", strings.Join(missing, ", "))
	}

	return nil
}
