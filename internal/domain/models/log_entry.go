package models

import "google.golang.org/protobuf/types/known/structpb"

type LogEntry struct {
	Service   string
	Level     string
	Message   string
	Metadata  map[string]*structpb.Value
	Timestamp int64
}
