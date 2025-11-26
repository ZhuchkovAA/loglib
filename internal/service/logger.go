package service

import (
	"github.com/ZhuchkovAA/loglib/pkg/constants"
	"github.com/ZhuchkovAA/loglib/pkg/models"
	"sync"
)

type LogStorage interface {
	Save(*models.Log) error
}

type Logger struct {
	ServiceName string
	queue       chan *models.Log
	SenderFunc  func(*models.Log) error
	Storage     LogStorage
	wg          sync.WaitGroup
	mu          sync.RWMutex
	closed      bool
}

func NewLogger(serviceName string, senderFunc func(*models.Log) error, storage LogStorage, bufferSize int) *Logger {
	return &Logger{
		ServiceName: serviceName,
		queue:       make(chan *models.Log, bufferSize),
		SenderFunc:  senderFunc,
		Storage:     storage,
	}
}

func (l *Logger) Log(level int, message string, fields ...*models.Field) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	metadata := make(map[string]string)

	for _, field := range fields {
		metadata[field.Key] = field.Value
	}

	logMessage := models.NewLog(l.ServiceName, level, message, metadata)

	l.AddToQueue(logMessage)

	logMessage.PrintColored()
}

func (l *Logger) Close() {
	l.mu.Lock()
	l.closed = true
	l.mu.Unlock()

	close(l.queue)
	l.wg.Wait()
}

func (l *Logger) AddToQueue(log *models.Log) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if l.closed {
		return
	}

	select {
	case l.queue <- log:
	default:
		_ = l.Storage.Save(log)
	}
}

func (l *Logger) Info(message string, fields ...*models.Field) {
	l.Log(constants.LevelInfo, message, fields...)
}

func (l *Logger) Warn(message string, fields ...*models.Field) {
	l.Log(constants.LevelWarn, message, fields...)
}

func (l *Logger) Error(message string, fields ...*models.Field) {
	l.Log(constants.LevelError, message, fields...)
}

func (l *Logger) Debug(message string, fields ...*models.Field) {
	l.Log(constants.LevelDebug, message, fields...)
}

func (l *Logger) Run() {
	l.wg.Add(1)
	defer l.wg.Done()

	for log := range l.queue {
		if err := l.SenderFunc(log); err != nil {
			_ = l.Storage.Save(log)
		}
	}
}
