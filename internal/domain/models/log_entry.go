package models

type LogEntry struct {
	Service   string
	Level     string
	Message   string
	Metadata  map[string]string
	Timestamp int64
}
