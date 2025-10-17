package consts

import (
	"fmt"
	"io"
)

const (
	сolorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorGray   = "\033[90m"
)

var (
	colors = map[int]string{
		LevelInfo:  colorBlue,
		LevelWarn:  colorYellow,
		LevelError: colorRed,
		LevelDebug: colorGreen,
	}
)

func GetColorByLevel(level int) string {
	if color, ok := colors[level]; ok {
		return color
	}

	return colorGray
}

func ColorFprintf(w io.Writer, color string, format string, a ...any) (int, error) {
	format = fmt.Sprintf("%s%s%s", color, format, сolorReset)
	return fmt.Fprintf(w, format, a...)
}
