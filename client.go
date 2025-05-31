package loglib

import (
	"encoding/json"
	"github.com/ZhuchkovAA/loglib/config"
	"github.com/ZhuchkovAA/loglib/constants"
	"github.com/ZhuchkovAA/loglib/internal/actions"
	"github.com/ZhuchkovAA/loglib/internal/domain/models"
	"log"
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
	log.SetFlags(0)

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

	levelStr, ok := logc.ErrorLevels[level]
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

	metaJSON, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		metaJSON = []byte("<invalid metadata>")
	}

	color := logc.GetColorByLevel(level)
	log.Printf("%s[%s] %s\nMessage: %s\nMetadata: %s%s\n\n",
		color, levelStr, time.Unix(timeUnix, 0), message, metaJSON, logc.ColorReset)
}

func (c *Client) Info(message string, metadata map[string]string) {
	c.Log(logc.LevelInfo, message, metadata)
}

func (c *Client) Warn(message string, metadata map[string]string) {
	c.Log(logc.LevelWarn, message, metadata)
}

func (c *Client) Error(message string, metadata map[string]string) {
	c.Log(logc.LevelError, message, metadata)
}

func (c *Client) Debug(message string, metadata map[string]string) {
	c.Log(logc.LevelDebug, message, metadata)
}

func (c *Client) run() {
	for entry := range c.queue {
		if err := c.sender.Send(entry); err != nil {
			_ = c.fallback.Save(entry)
		}
	}
}

var _ Logger = (*Client)(nil)
