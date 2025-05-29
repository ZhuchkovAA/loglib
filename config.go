package loglib

import (
	"fmt"
	"strings"
)

type Config struct {
	GRPCAddress  string
	FallbackPath string
	ServiceName  string
}

func (c *Config) MustLoadConfig() error {
	var missing []string

	if c.GRPCAddress == "" {
		missing = append(missing, "GRPC address")
	}
	if c.FallbackPath == "" {
		missing = append(missing, "fallback path")
	}
	if c.ServiceName == "" {
		missing = append(missing, "service name")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing config values: %s", strings.Join(missing, ", "))
	}

	return nil
}
