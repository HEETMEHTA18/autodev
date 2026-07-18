package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	}
	return "UNKNOWN"
}

type Logger struct {
	mu       sync.Mutex
	level    Level
	noColor  bool
	writer   io.Writer
	prefixes []string
}

var Default = New(LevelInfo, false, os.Stdout)

func New(level Level, noColor bool, writer io.Writer) *Logger {
	return &Logger{level: level, noColor: noColor, writer: writer}
}

func levelColor(l Level) string {
	switch l {
	case LevelDebug:
		return "\033[36m"
	case LevelInfo:
		return "\033[32m"
	case LevelWarn:
		return "\033[33m"
	case LevelError:
		return "\033[31m"
	}
	return "\033[0m"
}

const reset = "\033[0m"
const dim = "\033[2m"

func (l *Logger) log(level Level, format string, args ...interface{}) {
	if level < l.level {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	ts := time.Now().Format("15:04:05.000")
	msg := fmt.Sprintf(format, args...)

	var prefix string
	for _, p := range l.prefixes {
		prefix += p + " "
	}

	if l.noColor {
		fmt.Fprintf(l.writer, "%s [%s] %s%s\n", ts, level.String(), prefix, msg)
		return
	}

	color := levelColor(level)
	fmt.Fprintf(l.writer, "%s%s %s%s%s %s%s%s\n",
		dim, ts, reset,
		color+level.String()+reset,
		dim, reset,
		prefix,
		msg)
}

func (l *Logger) Debug(format string, args ...interface{}) { l.log(LevelDebug, format, args...) }
func (l *Logger) Info(format string, args ...interface{})  { l.log(LevelInfo, format, args...) }
func (l *Logger) Warn(format string, args ...interface{})  { l.log(LevelWarn, format, args...) }
func (l *Logger) Error(format string, args ...interface{}) { l.log(LevelError, format, args...) }

func (l *Logger) WithPrefix(p string) *Logger {
	l2 := &Logger{level: l.level, noColor: l.noColor, writer: l.writer}
	l2.prefixes = append(l2.prefixes, l.prefixes...)
	l2.prefixes = append(l2.prefixes, p)
	return l2
}

func (l *Logger) SetLevel(level Level) { l.mu.Lock(); defer l.mu.Unlock(); l.level = level }
func (l *Logger) SetNoColor(v bool)    { l.mu.Lock(); defer l.mu.Unlock(); l.noColor = v }
func (l *Logger) SetWriter(w io.Writer) { l.mu.Lock(); defer l.mu.Unlock(); l.writer = w }

// Package-level convenience functions
func Debug(format string, args ...interface{}) { Default.log(LevelDebug, format, args...) }
func Info(format string, args ...interface{})  { Default.log(LevelInfo, format, args...) }
func Warn(format string, args ...interface{})  { Default.log(LevelWarn, format, args...) }
func Error(format string, args ...interface{}) { Default.log(LevelError, format, args...) }

// LevelFromString parses a log level string (debug, info, warn, error).
func LevelFromString(s string) Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn", "warning":
		return LevelWarn
	case "error":
		return LevelError
	}
	return LevelInfo
}
