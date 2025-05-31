package logc

const (
	LevelInfo = iota + 1
	LevelWarn
	LevelError
	LevelDebug
)

var ErrorLevels = map[int]string{
	LevelInfo:  "INFO",
	LevelWarn:  "WARN",
	LevelError: "ERROR",
	LevelDebug: "DEBUG",
}
