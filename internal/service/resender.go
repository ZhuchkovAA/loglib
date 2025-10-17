package service

import (
	"bufio"
	"encoding/json"
	"github.com/ZhuchkovAA/loglib/internal/domain/models"
	"os"
	"strings"
	"time"
)

type ReSender struct {
	path   string
	sender *Sender
	tick   *time.Ticker
}

func NewReSender(path string, sender *Sender) *ReSender {
	return &ReSender{
		path:   path,
		sender: sender,
		tick:   time.NewTicker(10 * time.Second),
	}
}

func (r *ReSender) Loop() {
	for range r.tick.C {
		r.resend()
	}
}

func (r *ReSender) resend() {
	file, err := os.OpenFile(r.path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	var remaining []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		var entry *models.LogEntry
		if err := json.Unmarshal([]byte(line), entry); err != nil {
			continue
		}

		if err := r.sender.Send(entry); err != nil {
			remaining = append(remaining, line)
		}
	}

	file.Truncate(0)
	file.Seek(0, 0)
	file.WriteString(strings.Join(remaining, "\n"))
}
