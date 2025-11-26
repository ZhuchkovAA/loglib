package tests

import (
	"context"
	"errors"
	"fmt"
	"github.com/ZhuchkovAA/loglib"
	"github.com/ZhuchkovAA/loglib/pkg/config"
	"github.com/ZhuchkovAA/loglib/pkg/models"
	"os"
	"testing"
	"time"
)

func TestFallbackCreate(t *testing.T) {
	fallbackFile := "./test_fallback.log"
	defer os.Remove(fallbackFile)

	cfg := config.Config{
		FallbackPath: fallbackFile,
		ServiceName:  "test-service",
	}

	client, err := loglib.NewLogger(context.Background(), cfg, func(log *models.Log) error {
		fmt.Println("Начало обработки: ", log.Message, log.Level)
		return errors.New("Test error")
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Логируем сообщение
	client.Info("test message", loglib.String("env", "test"))
	client.Warn("test message", loglib.String("env", "test"))
	client.Error("test message", loglib.String("env", "test"))
	client.Debug("test message", loglib.String("env", "test"))
	client.Log(11, "test message", loglib.String("env", "test"))

	// Подождём чтобы run() успел обработать очередь
	time.Sleep(20500 * time.Millisecond)

	// Проверим, что файл fallback появился
	data, err := os.ReadFile(fallbackFile)
	if err != nil {
		t.Fatalf("failed to read fallback file: %v", err)
	}

	if len(data) == 0 {
		t.Errorf("expected data in fallback file, got empty file")
	}
}
