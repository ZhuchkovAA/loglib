package loglib

import (
	"encoding/json"
	"fmt"
	"github.com/ZhuchkovAA/loglib/config"
	"github.com/ZhuchkovAA/loglib/constants"
	"github.com/ZhuchkovAA/loglib/internal/actions"
	"github.com/ZhuchkovAA/loglib/internal/domain/models"
	"io"
	"log"
	"os"
	"time"
)

type Logger interface {
	Log(level int, message string, metadata map[string]string)
	Info(message string, metadata map[string]string)
	Warn(message string, metadata map[string]string)
	Error(message string, metadata map[string]string)
	Debug(message string, metadata map[string]string)
}

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

func (c *Client) Log(level int, message string, metadata map[string]string) {
	timeUnix := time.Now().Unix()

	levelStr, ok := consts.ErrorLevels[level]
	if !ok {
		levelStr = "UNKNOWN"
	}

	c.queue <- models.LogEntry{
		Service:   c.sender.Service,
		Level:     levelStr,
		Message:   message,
		Metadata:  metadata,
		Timestamp: timeUnix,
	}

	PrintColored(level, message, metadata, timeUnix)
}

func (c *Client) Info(message string, metadata map[string]string) {
	c.Log(consts.LevelInfo, message, metadata)
}

func (c *Client) Warn(message string, metadata map[string]string) {
	c.Log(consts.LevelWarn, message, metadata)
}

func (c *Client) Error(message string, metadata map[string]string) {
	c.Log(consts.LevelError, message, metadata)
}

func (c *Client) Debug(message string, metadata map[string]string) {
	c.Log(consts.LevelDebug, message, metadata)
}

func (c *Client) run() {
	for entry := range c.queue {
		if err := c.sender.Send(entry); err != nil {
			_ = c.fallback.Save(entry)
		}
	}
}

func PrintColored(level int, message string, metadata map[string]string, timestamp int64) {
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
