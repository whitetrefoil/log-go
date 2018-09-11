package logger

import (
	"fmt"
	"io"
	"os"
	"time"
)

const LogTime = "01-02 15:04:05"

func timestamp() string {
	return time.Now().Format(LogTime)
}

type Logger struct {
	Tag string
	Out io.Writer
}

func (l *Logger) Logf(format string, a ...interface{}) {
	fmt.Fprintln(l.Out, fmt.Sprintf("%s [%s]", timestamp(), l.Tag), fmt.Sprintf(format, a...))
}

func (l *Logger) Logln(a ...interface{}) {
	fmt.Fprintln(l.Out, fmt.Sprintf("%s [%s]", timestamp(), l.Tag), fmt.Sprint(a...))
}

func NewLogger(tag string) *Logger {
	return &Logger{tag, os.Stdout}
}
