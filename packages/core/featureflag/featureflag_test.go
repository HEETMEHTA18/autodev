package featureflag

import (
	"os"
	"testing"
)

func TestDefaultEnabled(t *testing.T) {
	ff := New("")
	if !ff.IsEnabled(AIEnabled) {
		t.Error("expected AIEnabled to be enabled by default")
	}
}

func TestSetAndGet(t *testing.T) {
	ff := New("")
	ff.Set(AIEnabled, false)
	if ff.IsEnabled(AIEnabled) {
		t.Error("expected AIEnabled to be disabled after Set(false)")
	}

	ff.Set(AIEnabled, true)
	if !ff.IsEnabled(AIEnabled) {
		t.Error("expected AIEnabled to be enabled after Set(true)")
	}
}

func TestEnvIndividual(t *testing.T) {
	os.Setenv("AUTODEV_AI_ENABLED", "false")
	defer os.Unsetenv("AUTODEV_AI_ENABLED")

	ff := New("")
	if ff.IsEnabled(AIEnabled) {
		t.Error("expected AIEnabled to be disabled via env var")
	}
}

func TestEnvBulk(t *testing.T) {
	os.Setenv("AUTODEV_FEATURE", "AI_ENABLED=false,SECURITY_SCAN=true")
	defer os.Unsetenv("AUTODEV_FEATURE")

	ff := New("AUTODEV_FEATURE")
	if ff.IsEnabled(AIEnabled) {
		t.Error("expected AIEnabled to be disabled via bulk env")
	}
	if !ff.IsEnabled(SecurityScan) {
		t.Error("expected SecurityScan to be enabled via bulk env")
	}
}

func TestEnvBulkWithSpaces(t *testing.T) {
	os.Setenv("AUTODEV_FEATURE", " CACHE = false , MCP = true ")
	defer os.Unsetenv("AUTODEV_FEATURE")

	ff := New("AUTODEV_FEATURE")
	if ff.IsEnabled(Cache) {
		t.Error("expected Cache to be disabled")
	}
	if !ff.IsEnabled(MCP) {
		t.Error("expected MCP to be enabled")
	}
}

func TestAll(t *testing.T) {
	ff := New("")
	ff.Set(AIEnabled, false)
	all := ff.All()
	if val, ok := all[AIEnabled]; !ok || val {
		t.Error("expected AIEnabled=false in All()")
	}
	if val, ok := all[SecurityScan]; !ok || !val {
		t.Error("expected SecurityScan=true by default in All()")
	}
}

func TestIndividualEnvVarOverride(t *testing.T) {
	ff := New("")
	ff.Set(SecurityScan, false)

	os.Setenv("AUTODEV_SECURITY_SCAN", "true")
	defer os.Unsetenv("AUTODEV_SECURITY_SCAN")

	if !ff.IsEnabled(SecurityScan) {
		t.Error("individual env var should override Set value")
	}
}
