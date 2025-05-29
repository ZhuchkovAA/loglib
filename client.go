package loglib

import (
	"encoding/json"
	"github.com/ZhuchkovAA/loglib/internal/actions"
	"github.com/ZhuchkovAA/loglib/internal/domain/models"
	"log"
	"time"
)

type Client struct {
	queue    chan models.LogEntry
	sender   *actions.Sender
	fallback *actions.Fallback
	resender *actions.Resender
}

func New(cfg Config) (*Client, error) {
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

func (c *Client) Log(level, message string, metadata map[string]string) {
	timeUnix := time.Now().Unix()

	c.queue <- models.LogEntry{
		Service:   c.sender.Service,
		Level:     level,
		Message:   message,
		Metadata:  metadata,
		Timestamp: timeUnix,
	}

	metaJSON, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		metaJSON = []byte("<invalid metadata>")
	}

	log.Printf("[%s] %s\nMessage: %s\nMetadata: %s\n\n", level, time.Unix(timeUnix, 0), message, metaJSON)
}

func (c *Client) run() {
	for entry := range c.queue {
		if err := c.sender.Send(entry); err != nil {
			_ = c.fallback.Save(entry)
		}
	}
}
