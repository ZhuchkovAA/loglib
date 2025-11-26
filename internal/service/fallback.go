package service

import (
	"bufio"
	"encoding/json"
	"github.com/ZhuchkovAA/loglib/pkg/models"
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

func (f *Fallback) Save(entry *models.Log) error {
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

func (f *Fallback) Load() ([]*models.Log, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	file, err := os.OpenFile(f.path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	logs := make([]*models.Log, 0)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		var entry models.Log
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}

		logs = append(logs, &entry)
	}

	if err := scanner.Err(); err != nil {
		return logs, err
	}

	file.Truncate(0)
	file.Seek(0, 0)

	return logs, nil
}
