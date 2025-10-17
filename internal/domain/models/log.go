package models

import (
	"encoding/json"
	consts "github.com/ZhuchkovAA/loglib/constants"
	"google.golang.org/protobuf/types/known/structpb"
	"io"
	"log"
	"os"
	"time"
)

type LogEntry struct {
	Level     int
	Message   string
	Metadata  map[string]*structpb.Value
	Timestamp int64
}

func NewLogEntry(Level int, Message string, Metadata map[string]*structpb.Value, Timestamp int64) *LogEntry {
	return &LogEntry{
		Level:     Level,
		Message:   Message,
		Metadata:  Metadata,
		Timestamp: Timestamp,
	}
}

func (l *LogEntry) PrintColored() {
	metaJSON, err := json.MarshalIndent(l.Metadata, "", "  ")
	if err != nil {
		metaJSON = []byte("<invalid metadata>")
	}

	color := consts.GetColorByLevel(l.Level)

	var output io.Writer = os.Stdout
	if l.Level == consts.LevelWarn || l.Level == consts.LevelError {
		output = os.Stderr
	}

	_, err = consts.ColorFprintf(
		output,
		color,
		"[%s] %s\nMessage: %s\nMetadata: %s\n\n",
		consts.GetErrorLevelStr(l.Level),
		time.Unix(l.Timestamp, 0),
		l.Message,
		metaJSON,
	)

	if err != nil {
		log.Printf("[LOGGER_ERROR] Ошибка вывода лога: %v", err)
	}
}
