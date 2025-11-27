package models

type Field struct {
	Key   string
	Value any
}

func NewField(key string, value any) *Field {
	return &Field{Key: key, Value: value}
}
