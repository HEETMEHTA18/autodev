package registry

import "testing"

func TestAllAgents(t *testing.T) {
	agents := All()
	if len(agents) != 5 {
		t.Errorf("expected 5 agents, got %d", len(agents))
	}
	seen := map[string]bool{}
	for _, a := range agents {
		if a.ID == "" || a.Command == "" || a.Manager == "" || a.Pkg == "" {
			t.Errorf("agent %+v missing required fields", a)
		}
		if seen[a.ID] {
			t.Errorf("duplicate agent id %q", a.ID)
		}
		seen[a.ID] = true
	}
}

func TestGet(t *testing.T) {
	a, err := Get("opencode")
	if err != nil {
		t.Fatalf("Get(opencode) error: %v", err)
	}
	if a.Command != "opencode" {
		t.Errorf("command = %q, want %q", a.Command, "opencode")
	}
	if _, err := Get("nope"); err == nil {
		t.Error("expected error for unknown agent")
	}
}

func TestDetectMissing(t *testing.T) {
	s := Detect(Agent{ID: "definitely-not-installed", Command: "definitely-not-installed-bin"})
	if s.Installed {
		t.Error("expected Installed=false for missing binary")
	}
}
