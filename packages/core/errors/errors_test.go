package errors

import (
	"errors"
	"testing"
)

func TestAppError(t *testing.T) {
	err := New(CodeValidation, "invalid input")
	if err.Code != CodeValidation {
		t.Errorf("expected CodeValidation, got %s", err.Code)
	}
	if err.Message != "invalid input" {
		t.Errorf("expected 'invalid input', got '%s'", err.Message)
	}
}

func TestAppErrorWrapping(t *testing.T) {
	original := errors.New("underlying error")
	err := Wrap(CodeGitHub, "API call failed", original)

	if !errors.Is(err, original) {
		t.Error("expected errors.Is to find original error")
	}
	if err.Code != CodeGitHub {
		t.Errorf("expected CodeGitHub, got %s", err.Code)
	}
}

func TestAppErrorString(t *testing.T) {
	err := New(CodeConfig, "missing key")
	expected := "[CONFIG] missing key"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}

	wrapped := Wrap(CodeNetwork, "timeout", errors.New("connection refused"))
	if wrapped.Error() == "" {
		t.Error("expected non-empty error string")
	}
}

func TestConvenienceConstructors(t *testing.T) {
	tests := []struct {
		name string
		fn   func(string) *AppError
		code Code
	}{
		{"ValidationError", ValidationError, CodeValidation},
		{"ConfigError", ConfigError, CodeConfig},
		{"NotFoundError", NotFoundError, CodeNotFound},
		{"GitHubError", GitHubError, CodeGitHub},
		{"AIError", AIError, CodeAI},
		{"ScannerError", ScannerError, CodeScanner},
		{"InstallerError", InstallerError, CodeInstaller},
		{"NetworkError", NetworkError, CodeNetwork},
		{"SecurityError", SecurityError, CodeSecurity},
		{"InternalError", InternalError, CodeInternal},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn("test message")
			if err.Code != tt.code {
				t.Errorf("expected code %s, got %s", tt.code, err.Code)
			}
		})
	}
}

func TestWrapConvenience(t *testing.T) {
	original := errors.New("original error")

	err := WrapValidationError(original)
	if err.Code != CodeValidation {
		t.Errorf("expected CodeValidation, got %s", err.Code)
	}

	err = WrapConfigError(original)
	if err.Code != CodeConfig {
		t.Errorf("expected CodeConfig, got %s", err.Code)
	}

	err = WrapGitHubError(original)
	if err.Code != CodeGitHub {
		t.Errorf("expected CodeGitHub, got %s", err.Code)
	}
}

func TestCodeFromErr(t *testing.T) {
	if CodeFromErr(nil) != "" {
		t.Error("expected empty code for nil error")
	}

	appErr := New(CodeValidation, "test")
	if CodeFromErr(appErr) != CodeValidation {
		t.Errorf("expected CodeValidation, got %s", CodeFromErr(appErr))
	}

	stdErr := errors.New("standard error")
	if CodeFromErr(stdErr) != CodeUnknown {
		t.Errorf("expected CodeUnknown, got %s", CodeFromErr(stdErr))
	}
}

func TestWithDetails(t *testing.T) {
	err := New(CodeValidation, "invalid input").
		WithDetails(map[string]interface{}{"field": "email"})
	if err.Details["field"] != "email" {
		t.Errorf("expected details to contain field=email")
	}
}

func TestFormattedConstructors(t *testing.T) {
	err := ValidationErrorf("field %s is required", "name")
	if err.Message != "field name is required" {
		t.Errorf("unexpected message: %s", err.Message)
	}

	err = ConfigErrorf("key %s not found in %s", "port", "config.yaml")
	if err.Message != "key port not found in config.yaml" {
		t.Errorf("unexpected message: %s", err.Message)
	}
}
