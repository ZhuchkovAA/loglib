package tests

import (
	"github.com/ZhuchkovAA/loglib"
	"github.com/ZhuchkovAA/loglib/config"
	"os"
	"testing"
	"time"
)

func TestFallback_Create(t *testing.T) {
	fallbackFile := "test_fallback.log"
	defer os.Remove(fallbackFile)

	cfg := config.Config{
		GRPCAddress:  "localhost:1", // сервис логов должен быть недоступен для теста
		FallbackPath: fallbackFile,
		ServiceName:  "test-service",
	}

	client, err := loglib.New(cfg)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Логируем сообщение
	client.Info("test message", map[string]string{"env": "test"})
	client.Warn("test message", map[string]string{"env": "test"})
	client.Error("test message", map[string]string{"env": "test"})
	client.Debug("test message", map[string]string{"env": "test"})

	// Подождём чтобы run() успел обработать очередь
	time.Sleep(500 * time.Millisecond)

	// Проверим, что файл fallback появился
	data, err := os.ReadFile(fallbackFile)
	if err != nil {
		t.Fatalf("failed to read fallback file: %v", err)
	}

	if len(data) == 0 {
		t.Errorf("expected data in fallback file, got empty file")
	}
}
