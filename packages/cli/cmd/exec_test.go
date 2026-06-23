package cmd

import (
	"testing"
)

func TestIsValidPackageName(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"curl", true},
		{"docker.io", true},
		{"python3-pip", true},
		{"g++", true},
		{"-unsafe-flag", false},
		{"pkg; rm -rf /", false},
		{"", false},
		{"pkg&", false},
		{"  spaces  ", false},
	}

	for _, test := range tests {
		result := isValidPackageName(test.name)
		if result != test.expected {
			t.Errorf("isValidPackageName(%q) = %v; want %v", test.name, result, test.expected)
		}
	}
}

func TestRunLinuxInstallValidation(t *testing.T) {
	err := runLinuxInstall([]string{"curl", "-o", "/etc/passwd"})
	if err == nil {
		t.Error("runLinuxInstall should fail on invalid package names")
	}
}
