package tests

import (
	"github.com/ZhuchkovAA/loglib"
	"github.com/ZhuchkovAA/loglib/config"
	consts "github.com/ZhuchkovAA/loglib/constants"
	"testing"
)

func TestLogTypes(t *testing.T) {
	cfg := config.Config{
		GRPCAddress:  "localhost:50051",
		FallbackPath: "test_fallback.log",
		ServiceName:  "test-service",
	}

	client, err := loglib.New(cfg)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	tests := []struct {
		name string
		fn   func() error
	}{
		{
			name: "Info",
			fn:   func() error { return client.Info("test message", loglib.String("env", "test")) },
		},
		{
			name: "Warn",
			fn:   func() error { return client.Warn("test message", loglib.String("env", "test")) },
		},
		{
			name: "Error",
			fn:   func() error { return client.Error("test message", loglib.String("env", "test")) },
		},
		{
			name: "Debug",
			fn:   func() error { return client.Debug("test message", loglib.String("env", "test")) },
		},
		{
			name: "Log with LevelDebug",
			fn:   func() error { return client.Log(consts.LevelDebug, "test message", loglib.String("env", "test")) },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fn(); err != nil {
				t.Errorf("%s failed: %v", tt.name, err)
			}
		})
	}
}
