package catalog

import (
	"testing"
)

func TestCatalogResolution(t *testing.T) {
	c, err := Load()
	if err != nil {
		t.Fatalf("failed to load catalog: %v", err)
	}

	// Test profile resolution
	resolved, err := c.ResolveProfile("web-dev")
	if err != nil {
		t.Fatalf("failed to resolve web-dev profile: %v", err)
	}

	if len(resolved) == 0 {
		t.Fatalf("resolved list is empty")
	}

	// Verify that dependencies come before dependants
	// e.g. nodejs must come before react
	nodeIdx := -1
	reactIdx := -1
	for i, pkg := range resolved {
		if pkg.ID == "nodejs" {
			nodeIdx = i
		} else if pkg.ID == "react" {
			reactIdx = i
		}
	}

	if nodeIdx == -1 {
		t.Errorf("nodejs not found in resolved profile")
	}
	if reactIdx == -1 {
		t.Errorf("react not found in resolved profile")
	}
	if nodeIdx > reactIdx {
		t.Errorf("nodejs resolved after react: nodejs at %d, react at %d", nodeIdx, reactIdx)
	}
}

func TestAgentCatalogEntries(t *testing.T) {
	c, err := Load()
	if err != nil {
		t.Fatalf("failed to load catalog: %v", err)
	}

	agents := []string{"opencode", "claude", "codex", "gemini", "aider"}
	for _, id := range agents {
		pkg, ok := c.GetPackage(id)
		if !ok {
			t.Fatalf("agent package %q missing from catalog", id)
		}
		if pkg.Verify == "" {
			t.Errorf("agent %q has no verify command", id)
		}
		if pkg.Install.Linux.Method == "" || pkg.Install.Darwin.Method == "" || pkg.Install.Windows.Method == "" {
			t.Errorf("agent %q missing platform install steps", id)
		}
	}

	resolved, err := c.ResolveProfile("ai-agents")
	if err != nil {
		t.Fatalf("resolve ai-agents profile: %v", err)
	}
	if len(resolved) == 0 {
		t.Fatal("ai-agents profile resolved to nothing")
	}
}
