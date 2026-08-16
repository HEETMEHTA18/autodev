package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildSharedKnowledge(t *testing.T) {
	// Move into a temp project dir with a fake .autodevs tree.
	tmp := t.TempDir()
	old, _ := os.Getwd()
	defer os.Chdir(old)
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	os.MkdirAll(".autodevs/memory", 0755)
	os.MkdirAll(".autodevs/context", 0755)
	os.WriteFile(".autodevs/memory/successes.md", []byte("# Successes\n- fixed auth\n"), 0644)
	os.WriteFile("go.mod", []byte("module x\n"), 0644)

	kb, err := buildSharedKnowledge()
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if !strings.Contains(kb, "AutoDev Shared Codebase Knowledge") {
		t.Fatalf("missing header: %s", kb)
	}
	if !strings.Contains(kb, "Successes") {
		t.Fatalf("missing memory: %s", kb)
	}
	if !strings.Contains(kb, "Detected stack") {
		t.Fatalf("missing stack: %s", kb)
	}
	if !strings.Contains(kb, "Repository state") {
		t.Fatalf("missing git section: %s", kb)
	}
	t.Logf("knowledge bundle OK (%d bytes)", len(kb))
}

func TestAppendSessionOutcome(t *testing.T) {
	path := filepath.Join(t.TempDir(), "knowledge.md")
	os.WriteFile(path, []byte("# KB\n"), 0644)
	appendSessionOutcome(path, "fix bug", nil)
	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), "Session outcomes") {
		t.Fatalf("missing outcomes: %s", data)
	}
	if !strings.Contains(string(data), "completed") {
		t.Fatalf("missing outcome: %s", data)
	}
}
