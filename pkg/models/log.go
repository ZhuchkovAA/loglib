package models

import (
	"encoding/json"
	"github.com/ZhuchkovAA/loglib/pkg/constants"
	"io"
	"log"
	"os"
	"time"
)

type Log struct {
	Service   string
	Level     int
	Message   string
	Metadata  map[string]string
	Timestamp int64
}

func NewLog(service string, level int, message string, metadata map[string]string) *Log {
	return &Log{
		Service:   service,
		Level:     level,
		Message:   message,
		Metadata:  metadata,
		Timestamp: time.Now().Unix(),
	}
}

func (l *Log) PrintColored() {
	metaJSON, err := json.MarshalIndent(l.Metadata, "", "  ")
	if err != nil {
		metaJSON = []byte("<invalid metadata>")
	}

	color := constants.GetColorByLevel(l.Level)

	var output io.Writer = os.Stdout
	if l.Level == constants.LevelWarn || l.Level == constants.LevelError {
		output = os.Stderr
	}

	_, err = constants.ColorFprintf(
		output,
		color,
		"[%s] %s\nMessage: %s\nMetadata: %s\n\n",
		constants.GetErrorLevelStr(l.Level),
		time.Unix(l.Timestamp, 0),
		l.Message,
		metaJSON,
	)

	if err != nil {
		log.Printf("[LOGGER_ERROR] Ошибка вывода лога: %v", err)
	}
}
