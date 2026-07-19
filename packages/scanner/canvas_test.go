package scanner

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestCanvasGenerate(t *testing.T) {
	gen := NewCanvasGenerator(".")
	canvas, err := gen.Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if canvas.SchemaVersion != CanvasSchemaVersion {
		t.Errorf("schema_version = %q, want %q", canvas.SchemaVersion, CanvasSchemaVersion)
	}

	if canvas.GeneratedAt.IsZero() {
		t.Error("generated_at is zero")
	}

	if canvas.Project.Name == "" {
		t.Error("project.name is empty")
	}

	if len(canvas.Files) == 0 {
		t.Error("files is empty — should have found at least one .go file")
	}

	if canvas.AI.OneLiner == "" {
		t.Error("ai_context.one_liner is empty")
	}

	if len(canvas.AI.RecommendedReading) == 0 {
		t.Error("ai_context.recommended_reading is empty")
	}

	t.Logf("Canvas: %d files, %d LOC, %d layers, %d findings",
		len(canvas.Files), canvas.Project.TotalLOC,
		len(canvas.Architecture.Layers), len(canvas.Security.Findings))
	t.Logf("AI: %s", canvas.AI.OneLiner)
}

func TestCanvasSave(t *testing.T) {
	dir := t.TempDir()
	gen := NewCanvasGenerator(".")
	canvas, err := gen.Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Override root path for save
	gen.RootPath = dir
	err = gen.Save(canvas)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	savedPath := filepath.Join(dir, ".autodevs", "canvas", "canvas.json")
	data, err := os.ReadFile(savedPath)
	if err != nil {
		t.Fatalf("Read saved canvas failed: %v", err)
	}

	var loaded Canvas
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("Unmarshal saved canvas failed: %v", err)
	}

	if loaded.SchemaVersion != CanvasSchemaVersion {
		t.Errorf("loaded schema_version = %q, want %q", loaded.SchemaVersion, CanvasSchemaVersion)
	}

	if len(loaded.Files) != len(canvas.Files) {
		t.Errorf("loaded files count = %d, want %d", len(loaded.Files), len(canvas.Files))
	}
}

func TestCanvasFileClassification(t *testing.T) {
	tests := []struct {
		path     string
		lang     string
		expected string
	}{
		{"src/components/Navbar.tsx", "TypeScript React", "component"},
		{"src/utils/helpers.ts", "TypeScript", "util"},
		{"src/api/routes.go", "Go", "route"},
		{"src/models/user.go", "Go", "model"},
		{"cmd/root.go", "Go", "source"},
		{"main.go", "Go", "entry"},
		{"config.yaml", "YAML", "config"},
		{"test/user_test.go", "Go", "test"},
	}

	for _, tt := range tests {
		result := classifyFile(tt.path, tt.lang, nil)
		if result != tt.expected {
			t.Errorf("classifyFile(%q, %q) = %q, want %q", tt.path, tt.lang, result, tt.expected)
		}
	}
}

func TestCanvasSecurityScan(t *testing.T) {
	content := []byte(`
package main

import "fmt"

func main() {
	password := "hardcoded123"
	eval("dangerous code")
	fmt.Println(password)
}
`)

	findings := scanFileSecurity("test.go", content)
	if len(findings) == 0 {
		t.Error("expected at least one security finding")
	}

	hasPassword := false
	hasEval := false
	for _, f := range findings {
		if f.Type == "hardcoded_secret" {
			hasPassword = true
		}
		if f.Type == "code_injection" {
			hasEval = true
		}
	}

	if !hasPassword {
		t.Error("expected hardcoded_secret finding")
	}
	if !hasEval {
		t.Error("expected code_injection finding")
	}
}

func TestCanvasCircularDeps(t *testing.T) {
	files := []FileCard{
		{Path: "a.go", Dependencies: []string{"b.go"}},
		{Path: "b.go", Dependencies: []string{"c.go"}},
		{Path: "c.go", Dependencies: []string{"a.go"}},
	}

	cycles := detectCircularDeps(files)
	if len(cycles) == 0 {
		t.Error("expected circular dependency detection")
	}
}

func TestCanvasTodos(t *testing.T) {
	content := []byte(`
// TODO: fix this later
func main() {
	// FIXME: broken
	// HACK: temporary
}
`)

	if !hasTodos(content) {
		t.Error("expected hasTodos to be true")
	}

	todos := extractTodos(content)
	if len(todos) < 3 {
		t.Errorf("expected at least 3 todos, got %d", len(todos))
	}
}

func TestCanvasImports(t *testing.T) {
	goContent := []byte(`
package main

import (
	"fmt"
	"net/http"
	"github.com/foo/bar"
)
`)
	imports := extractImports(goContent, "Go")
	if len(imports) < 3 {
		t.Errorf("expected at least 3 Go imports, got %d: %v", len(imports), imports)
	}

	jsContent := []byte(`
import React from 'react';
import { useState } from 'react';
const x = require('lodash');
`)
	imports = extractImports(jsContent, "JavaScript")
	if len(imports) < 3 {
		t.Errorf("expected at least 3 JS imports, got %d: %v", len(imports), imports)
	}
}

func TestCanvasFunctions(t *testing.T) {
	goContent := []byte(`
package main

func Hello() string {
	return "hello"
}

func testSomething() bool {
	return true
}
`)
	funcs := extractFunctions(goContent, "Go")
	if len(funcs) < 2 {
		t.Errorf("expected at least 2 Go functions, got %d", len(funcs))
	}
}
