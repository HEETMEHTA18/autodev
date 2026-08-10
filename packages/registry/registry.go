// Package registry is the AutoDevs agent registry: a lightweight index of the
// external AI/developer CLIs AutoDevs discovers, installs and launches. The
// heavy install logic lives in the catalog package; this package only knows
// each agent's identity, role, package manager and how to probe it.
package registry

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Role is the classification used by the agent router.
type Role string

const (
	RoleCoding   Role = "coding"
	RoleResearch Role = "research"
	RoleReview   Role = "review"
	RoleSecurity Role = "security"
	RoleComplex  Role = "complex"
)

// Agent describes one external AI CLI managed by AutoDevs.
type Agent struct {
	ID      string // cli id, also the catalog package id
	Name    string
	Command string // binary on PATH
	Role    Role
	Manager string // npm | pip
	Pkg     string // package name for the manager
	Desc    string
}

// All returns the registered agents in a stable display order.
func All() []Agent {
	return []Agent{
		{ID: "opencode", Name: "OpenCode", Command: "opencode", Role: RoleCoding, Manager: "npm", Pkg: "opencode-ai", Desc: "Terminal coding agent"},
		{ID: "claude", Name: "Claude Code", Command: "claude", Role: RoleComplex, Manager: "npm", Pkg: "@anthropic-ai/claude-code", Desc: "Anthropic coding agent"},
		{ID: "codex", Name: "Codex", Command: "codex", Role: RoleReview, Manager: "npm", Pkg: "@openai/codex", Desc: "OpenAI coding agent"},
		{ID: "gemini", Name: "Gemini CLI", Command: "gemini", Role: RoleResearch, Manager: "npm", Pkg: "@google/gemini-cli", Desc: "Google research agent"},
		{ID: "aider", Name: "Aider", Command: "aider", Role: RoleCoding, Manager: "pip", Pkg: "aider-chat", Desc: "AI pair programmer"},
	}
}

// Get returns the agent with the given id, or an error.
func Get(id string) (Agent, error) {
	for _, a := range All() {
		if a.ID == id {
			return a, nil
		}
	}
	return Agent{}, fmt.Errorf("unknown agent %q (run 'autodev tools list')", id)
}

// Status is the result of probing an agent.
type Status struct {
	Agent     Agent
	Installed bool
	Version   string
}

// Detect reports whether an agent binary is on PATH and, when present, its
// reported version.
func Detect(a Agent) Status {
	path, err := exec.LookPath(a.Command)
	if err != nil {
		return Status{Agent: a}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var out bytes.Buffer
	cmd := exec.CommandContext(ctx, path, "--version")
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return Status{Agent: a, Installed: true}
	}
	return Status{Agent: a, Installed: true, Version: firstLine(out.String())}
}

// DetectAll probes every registered agent.
func DetectAll() []Status {
	statuses := make([]Status, 0, len(All()))
	for _, a := range All() {
		statuses = append(statuses, Detect(a))
	}
	return statuses
}

func firstLine(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.IndexByte(s, '\n'); i != -1 {
		return s[:i]
	}
	return s
}

// Running reports whether a process matching the agent binary is active.
func Running(a Agent) bool {
	if _, err := exec.LookPath("pgrep"); err != nil {
		return false
	}
	cmd := exec.Command("pgrep", "-f", "^"+a.Command)
	return cmd.Run() == nil
}
