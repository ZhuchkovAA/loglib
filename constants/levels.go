package consts

const (
	LevelInfo = iota + 1
	LevelWarn
	LevelError
	LevelDebug
)

var errorLevels = map[int]string{
	LevelInfo:  "INFO",
	LevelWarn:  "WARN",
	LevelError: "ERROR",
	LevelDebug: "DEBUG",
}

func GetErrorLevelStr(level int) string {
	if levelStr, ok := errorLevels[level]; ok {
		return levelStr
	}

	return "UNKNOWN"
}
