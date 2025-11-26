package service

import (
	"context"
	"github.com/ZhuchkovAA/loglib/pkg/models"
	"time"
)

type Storage interface {
	Load() ([]*models.Log, error)
}

type Sender interface {
	AddToQueue(log *models.Log)
}

type ReSender struct {
	sender  Sender
	storage Storage
	tick    *time.Ticker
}

func NewReSender(sender Sender, storage Storage) *ReSender {
	return &ReSender{
		sender:  sender,
		storage: storage,
		tick:    time.NewTicker(10 * time.Second),
	}
}

func (r *ReSender) Loop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-r.tick.C:
			r.resend()
		}

	}
}

func (r *ReSender) resend() {
	logs, err := r.storage.Load()
	if err != nil {
		return
	}

	for _, log := range logs {
		r.sender.AddToQueue(log)
	}
}
