package tests

import (
	"context"
	"github.com/ZhuchkovAA/loglib"
	"github.com/ZhuchkovAA/loglib/pkg/config"
	"github.com/ZhuchkovAA/loglib/pkg/constants"
	"github.com/ZhuchkovAA/loglib/pkg/models"
	"testing"
)

func TestLogTypes(t *testing.T) {
	cfg := config.Config{
		FallbackPath: "test_fallback.log",
		ServiceName:  "test-service",
	}

	client, err := loglib.NewLogger(context.Background(), cfg, func(*models.Log) error { return nil })
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	tests := []struct {
		name string
		fn   func()
	}{
		{
			name: "Info",
			fn:   func() { client.Info("test message", loglib.String("env", "test")) },
		},
		{
			name: "Warn",
			fn:   func() { client.Warn("test message", loglib.String("env", "test")) },
		},
		{
			name: "Error",
			fn:   func() { client.Error("test message", loglib.String("env", "test")) },
		},
		{
			name: "Debug",
			fn:   func() { client.Debug("test message", loglib.String("env", "test")) },
		},
		{
			name: "Log with LevelDebug",
			fn:   func() { client.Log(constants.LevelDebug, "test message", loglib.String("env", "test")) },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fn()
		})
	}
}
