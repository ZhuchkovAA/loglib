package service

import (
	"encoding/json"
	"github.com/ZhuchkovAA/loglib/internal/domain/models"
	"os"
	"sync"
)

type Fallback struct {
	mu   sync.Mutex
	path string
}

func NewFallback(path string) *Fallback {
	return &Fallback{path: path}
}

func (f *Fallback) Save(entry *models.LogEntry) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	file, err := os.OpenFile(f.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	b, _ := json.Marshal(entry)
	line := string(b) + "\n"
	_, err = file.WriteString(line)
	return err
}
