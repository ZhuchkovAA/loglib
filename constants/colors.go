package logc

const (
	ColorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorGray   = "\033[90m"
)

func GetColorByLevel(level int) string {
	colors := map[int]string{
		LevelInfo:  colorBlue,
		LevelWarn:  colorYellow,
		LevelError: colorRed,
		LevelDebug: colorGreen,
	}

	if color, ok := colors[level]; ok {
		return color
	}

	return colorGray
}
