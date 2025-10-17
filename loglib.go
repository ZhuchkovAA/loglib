package loglib

import (
	"github.com/ZhuchkovAA/loglib/internal/app"
	"github.com/ZhuchkovAA/loglib/internal/config"
	"github.com/ZhuchkovAA/loglib/internal/domain/models"
	"time"
)

type Logger interface {
	Log(level int, msg string, fields ...models.Field)
	Info(msg string, fields ...models.Field)
	Warn(msg string, fields ...models.Field)
	Error(msg string, fields ...models.Field)
	Debug(msg string, fields ...models.Field)
}

type logLib struct {
	logger Logger
}

func String(key, value string) models.Field  { return models.Field{Key: key, Value: value} }
func Int(key string, value int) models.Field { return models.Field{Key: key, Value: value} }
func Error(err error) models.Field           { return models.Field{Key: "error", Value: err} }
func Duration(key string, value time.Duration) models.Field {
	return models.Field{Key: key, Value: value}
}

func New(logger Logger) *logLib {
	app.Run(config.Config{
		FallbackPath: "./",
	})

	return &logLib{
		logger: logger,
	}
}
