package loglib

import (
	"context"
	"github.com/ZhuchkovAA/loglib/internal/service"
	"github.com/ZhuchkovAA/loglib/pkg/config"
	"github.com/ZhuchkovAA/loglib/pkg/models"
	"strconv"
	"time"
)

type Logger interface {
	Log(level int, msg string, fields ...*models.Field)
	Info(msg string, fields ...*models.Field)
	Warn(msg string, fields ...*models.Field)
	Error(msg string, fields ...*models.Field)
	Debug(msg string, fields ...*models.Field)
}

func String(key, value string) *models.Field  { return models.NewField(key, value) }
func Int(key string, value int) *models.Field { return models.NewField(key, strconv.Itoa(value)) }
func Error(err error) *models.Field           { return models.NewField("error", err.Error()) }
func Duration(key string, value time.Duration) *models.Field {
	return models.NewField("error", value.String())
}

func NewLogger(ctx context.Context, cfg config.Config, fn func(*models.Log) error) (Logger, error) {
	err := cfg.MustLoadConfig()
	if err != nil {
		return nil, err
	}

	storage := service.NewFallback(cfg.FallbackPath)

	logger := service.NewLogger(cfg.ServiceName, fn, storage, cfg.BufferSize)
	go func() {
		<-ctx.Done()
		logger.Close()
	}()

	go logger.Run()

	reSender := service.NewReSender(logger, storage)
	go reSender.Loop(ctx)

	return logger, nil
}
