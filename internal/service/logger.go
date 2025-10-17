package service

import (
	consts "github.com/ZhuchkovAA/loglib/constants"
	"github.com/ZhuchkovAA/loglib/internal/domain/models"
	"google.golang.org/protobuf/types/known/structpb"
	"log"
	"time"
)

type clientSender interface {
	Send(*models.LogEntry) error
}

type clientFallback interface {
	Save(*models.LogEntry) error
}

type Logger struct {
	Queue    chan *models.LogEntry
	Sender   clientSender
	Fallback clientFallback
	ReSender *ReSender
}

func NewLogger(sender clientSender, fallback clientFallback) *Logger {
	return &Logger{
		Queue:    make(chan *models.LogEntry),
		Sender:   sender,
		Fallback: fallback,
	}
}

func (l *Logger) Log(level int, message string, fields ...models.Field) {
	timeUnix := time.Now().Unix()

	metadata := map[string]*structpb.Value{}
	for _, field := range fields {
		val, err := structpb.NewValue(field.Value)
		if err != nil {
			log.Println(err)
			return
		}
		metadata[field.Key] = val
	}

	logMessage := models.NewLogEntry(level, message, metadata, timeUnix)

	l.Queue <- logMessage
	logMessage.PrintColored()
}

func (l *Logger) Info(message string, fields ...models.Field) {
	l.Log(consts.LevelInfo, message, fields...)
}

func (l *Logger) Warn(message string, fields ...models.Field) {
	l.Log(consts.LevelWarn, message, fields...)
}

func (l *Logger) Error(message string, fields ...models.Field) {
	l.Log(consts.LevelError, message, fields...)
}

func (l *Logger) Debug(message string, fields ...models.Field) {
	l.Log(consts.LevelDebug, message, fields...)
}

func (l *Logger) Run() {
	for entry := range l.Queue {
		if err := l.Sender.Send(entry); err != nil {
			_ = l.Fallback.Save(entry)
		}
	}
}
