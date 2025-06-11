package loglib

import (
	"encoding/json"
	"fmt"
	"github.com/ZhuchkovAA/loglib/config"
	"github.com/ZhuchkovAA/loglib/constants"
	"github.com/ZhuchkovAA/loglib/internal/actions"
	"github.com/ZhuchkovAA/loglib/internal/domain/models"
	"google.golang.org/protobuf/types/known/structpb"
	"io"
	"log"
	"os"
	"time"
)

type Logger interface {
	Log(level int, msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Debug(msg string, fields ...Field)
}

type Field struct {
	Key   string
	Value any
}

func String(key, value string) Field                 { return Field{Key: key, Value: value} }
func Int(key string, value int) Field                { return Field{Key: key, Value: value} }
func Error(err error) Field                          { return Field{Key: "error", Value: err} }
func Duration(key string, value time.Duration) Field { return Field{Key: key, Value: value} }

type Client struct {
	queue    chan models.LogEntry
	sender   *actions.Sender
	fallback *actions.Fallback
	resender *actions.Resender
}

func New(cfg config.Config) (*Client, error) {
	err := cfg.MustLoadConfig()
	if err != nil {
		return nil, err
	}

	sender := actions.NewSender(cfg)

	client := &Client{
		queue:    make(chan models.LogEntry),
		sender:   sender,
		fallback: actions.NewFallback(cfg.FallbackPath),
		resender: actions.NewResender(cfg.FallbackPath, sender),
	}

	go client.run()

	return client, nil
}

func (c *Client) Log(level int, message string, fields ...Field) {
	timeUnix := time.Now().Unix()

	levelStr, ok := consts.ErrorLevels[level]
	if !ok {
		levelStr = "UNKNOWN"
	}

	metadata := map[string]*structpb.Value{}
	for _, field := range fields {
		val, err := structpb.NewValue(field.Value)
		if err != nil {
			log.Println(err)
			return
		}
		metadata[field.Key] = val
	}

	c.queue <- models.LogEntry{
		Service:   c.sender.Service,
		Level:     levelStr,
		Message:   message,
		Metadata:  metadata,
		Timestamp: timeUnix,
	}

	PrintColored(level, message, timeUnix, metadata)
}

func (c *Client) Info(message string, fields ...Field) {
	c.Log(consts.LevelInfo, message, fields...)
}

func (c *Client) Warn(message string, fields ...Field) {
	c.Log(consts.LevelWarn, message, fields...)
}

func (c *Client) Error(message string, fields ...Field) {
	c.Log(consts.LevelError, message, fields...)
}

func (c *Client) Debug(message string, fields ...Field) {
	c.Log(consts.LevelDebug, message, fields...)
}

func (c *Client) run() {
	for entry := range c.queue {
		if err := c.sender.Send(entry); err != nil {
			_ = c.fallback.Save(entry)
		}
	}
}

func PrintColored(level int, message string, timestamp int64, metadata map[string]*structpb.Value) {
	metaJSON, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		metaJSON = []byte("<invalid metadata>")
	}

	color := consts.GetColorByLevel(level)

	var output io.Writer = os.Stdout
	if level == consts.LevelWarn || level == consts.LevelError {
		output = os.Stderr
	}

	_, err = fmt.Fprintf(output, "%s[%s] %s\nMessage: %s\nMetadata: %s%s\n\n",
		color, consts.ErrorLevels[level], time.Unix(timestamp, 0), message, metaJSON, consts.ColorReset)

	if err != nil {
		log.Printf("[LOGGER_ERROR] Ошибка вывода лога: %v", err)
	}
}

var _ Logger = (*Client)(nil)
