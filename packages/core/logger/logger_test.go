package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestLevels(t *testing.T) {
	tests := []struct {
		level Level
		want  string
	}{
		{LevelDebug, "DEBUG"},
		{LevelInfo, "INFO"},
		{LevelWarn, "WARN"},
		{LevelError, "ERROR"},
	}
	for _, tt := range tests {
		if got := tt.level.String(); got != tt.want {
			t.Errorf("Level(%d).String() = %q, want %q", tt.level, got, tt.want)
		}
	}
}

func TestLevelFromString(t *testing.T) {
	tests := []struct {
		input string
		want  Level
	}{
		{"debug", LevelDebug},
		{"DEBUG", LevelDebug},
		{"info", LevelInfo},
		{"warn", LevelWarn},
		{"warning", LevelWarn},
		{"error", LevelError},
		{"unknown", LevelInfo},
	}
	for _, tt := range tests {
		if got := LevelFromString(tt.input); got != tt.want {
			t.Errorf("LevelFromString(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestLoggerNoColor(t *testing.T) {
	var buf bytes.Buffer
	l := New(LevelDebug, true, &buf)
	l.Info("test %s", "message")

	output := buf.String()
	if !strings.Contains(output, "[INFO]") {
		t.Errorf("expected [INFO] in output, got: %s", output)
	}
	if !strings.Contains(output, "test message") {
		t.Errorf("expected message in output, got: %s", output)
	}
	if strings.Contains(output, "\033[") {
		t.Errorf("expected no color codes, got: %s", output)
	}
}

func TestLoggerWithColor(t *testing.T) {
	var buf bytes.Buffer
	l := New(LevelDebug, false, &buf)
	l.Info("hello")

	output := buf.String()
	if !strings.Contains(output, "INFO") {
		t.Errorf("expected INFO in output, got: %s", output)
	}
}

func TestLoggerLevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	l := New(LevelWarn, true, &buf)
	l.Debug("debug msg")
	l.Info("info msg")
	l.Warn("warn msg")
	l.Error("error msg")

	output := buf.String()
	if strings.Contains(output, "debug msg") {
		t.Errorf("expected debug to be filtered out")
	}
	if strings.Contains(output, "info msg") {
		t.Errorf("expected info to be filtered out")
	}
	if !strings.Contains(output, "warn msg") {
		t.Errorf("expected warn msg in output")
	}
	if !strings.Contains(output, "error msg") {
		t.Errorf("expected error msg in output")
	}
}

func TestLoggerWithPrefix(t *testing.T) {
	var buf bytes.Buffer
	l := New(LevelInfo, true, &buf)
	l2 := l.WithPrefix("[scan]")
	l2.Info("scanning...")

	output := buf.String()
	if !strings.Contains(output, "[scan]") {
		t.Errorf("expected prefix in output, got: %s", output)
	}
}

func TestLevelOrdering(t *testing.T) {
	if !(LevelDebug < LevelInfo) {
		t.Error("LevelDebug should be < LevelInfo")
	}
	if !(LevelInfo < LevelWarn) {
		t.Error("LevelInfo should be < LevelWarn")
	}
	if !(LevelWarn < LevelError) {
		t.Error("LevelWarn should be < LevelError")
	}
}

func TestSetLevel(t *testing.T) {
	l := New(LevelDebug, true, nil)
	l.SetLevel(LevelError)
	if l.level != LevelError {
		t.Errorf("SetLevel didn't work: got %v, want %v", l.level, LevelError)
	}
}

func TestPackageLevelFunctions(t *testing.T) {
	var buf bytes.Buffer
	old := Default
	defer func() { Default = old }()
	Default = New(LevelInfo, true, &buf)
	Info("package %s", "test")
	if !strings.Contains(buf.String(), "package test") {
		t.Errorf("package-level Info failed, got: %s", buf.String())
	}
}
