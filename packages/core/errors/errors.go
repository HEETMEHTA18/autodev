package errors

import (
	"fmt"
)

type Code string

const (
	CodeUnknown       Code = "UNKNOWN"
	CodeValidation    Code = "VALIDATION"
	CodeConfig        Code = "CONFIG"
	CodeNotFound      Code = "NOT_FOUND"
	CodeGitHub        Code = "GITHUB"
	CodeDocker        Code = "DOCKER"
	CodeAI            Code = "AI"
	CodeScanner       Code = "SCANNER"
	CodeInstaller     Code = "INSTALLER"
	CodeNetwork       Code = "NETWORK"
	CodeFileSystem    Code = "FILESYSTEM"
	CodeSecurity      Code = "SECURITY"
	CodeCache         Code = "CACHE"
	CodePlugin        Code = "PLUGIN"
	CodeInternal      Code = "INTERNAL"
	CodeCancelled     Code = "CANCELLED"
)

type AppError struct {
	Code    Code
	Message string
	Err     error
	Details map[string]interface{}
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error { return e.Err }

func New(code Code, msg string) *AppError {
	return &AppError{Code: code, Message: msg}
}

func Wrap(code Code, msg string, err error) *AppError {
	return &AppError{Code: code, Message: msg, Err: err}
}

func (e *AppError) WithDetails(d map[string]interface{}) *AppError {
	e.Details = d
	return e
}

func ValidationError(msg string) *AppError  { return New(CodeValidation, msg) }
func ConfigError(msg string) *AppError       { return New(CodeConfig, msg) }
func NotFoundError(msg string) *AppError     { return New(CodeNotFound, msg) }
func GitHubError(msg string) *AppError       { return New(CodeGitHub, msg) }
func AIError(msg string) *AppError           { return New(CodeAI, msg) }
func ScannerError(msg string) *AppError      { return New(CodeScanner, msg) }
func InstallerError(msg string) *AppError    { return New(CodeInstaller, msg) }
func NetworkError(msg string) *AppError      { return New(CodeNetwork, msg) }
func SecurityError(msg string) *AppError     { return New(CodeSecurity, msg) }
func InternalError(msg string) *AppError     { return New(CodeInternal, msg) }

func ValidationErrorf(format string, args ...interface{}) *AppError  { return New(CodeValidation, fmt.Sprintf(format, args...)) }
func ConfigErrorf(format string, args ...interface{}) *AppError       { return New(CodeConfig, fmt.Sprintf(format, args...)) }
func NotFoundErrorf(format string, args ...interface{}) *AppError     { return New(CodeNotFound, fmt.Sprintf(format, args...)) }
func GitHubErrorf(format string, args ...interface{}) *AppError       { return New(CodeGitHub, fmt.Sprintf(format, args...)) }
func AIErrorf(format string, args ...interface{}) *AppError           { return New(CodeAI, fmt.Sprintf(format, args...)) }
func ScannerErrorf(format string, args ...interface{}) *AppError      { return New(CodeScanner, fmt.Sprintf(format, args...)) }
func InstallerErrorf(format string, args ...interface{}) *AppError    { return New(CodeInstaller, fmt.Sprintf(format, args...)) }

func WrapValidationError(err error) *AppError { return Wrap(CodeValidation, "validation failed", err) }
func WrapConfigError(err error) *AppError      { return Wrap(CodeConfig, "configuration error", err) }
func WrapGitHubError(err error) *AppError      { return Wrap(CodeGitHub, "GitHub API error", err) }

// CodeFromErr returns the error code, defaulting to CodeUnknown.
func CodeFromErr(err error) Code {
	if err == nil {
		return ""
	}
	if ae, ok := err.(*AppError); ok {
		return ae.Code
	}
	return CodeUnknown
}
