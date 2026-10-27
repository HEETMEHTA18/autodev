package cmd

import (
	"strings"
	"testing"
)

// ── grade bands ──────────────────────────────────────────────────────────────

func TestScoreGrade(t *testing.T) {
	cases := []struct {
		score int
		want  string
	}{
		{100, "A"}, {90, "A"}, {89, "B"},
		{75, "B"}, {74, "C"}, {60, "C"},
		{59, "D"}, {40, "D"}, {39, "F"}, {0, "F"},
	}
	for _, c := range cases {
		if got := scoreGrade(c.score); got != c.want {
			t.Errorf("scoreGrade(%d) = %q, want %q", c.score, got, c.want)
		}
	}
}

// ── weighted score ───────────────────────────────────────────────────────────

func TestComputeScore(t *testing.T) {
	allOK := []checkResult{
		{Name: "Exposed Secrets Scanner", Status: "ok"},
		{Name: "Linter & Code Format Status", Status: "ok"},
	}
	if got := computeScore(allOK); got != 100 {
		t.Errorf("all ok: score = %d, want 100", got)
	}

	allErr := []checkResult{
		{Name: "Exposed Secrets Scanner", Status: "error"},
		{Name: "Git Configuration (.gitignore)", Status: "error"},
	}
	if got := computeScore(allErr); got != 0 {
		t.Errorf("all error: score = %d, want 0", got)
	}

	allWarn := []checkResult{{Name: "Dependency Lockfiles", Status: "warning"}}
	if got := computeScore(allWarn); got != 50 {
		t.Errorf("all warning: score = %d, want 50", got)
	}

	// Weighted mix: secrets ok (w3) + git warning (w2) + linter ok (w1)
	// earned = 3*1 + 2*0.5 + 1*1 = 5 of 6 → 83.33 → 83
	mixed := []checkResult{
		{Name: "Exposed Secrets Scanner", Status: "ok"},
		{Name: "Git Configuration (.gitignore)", Status: "warning"},
		{Name: "Linter & Code Format Status", Status: "ok"},
	}
	if got := computeScore(mixed); got != 83 {
		t.Errorf("weighted mix: score = %d, want 83", got)
	}

	// Unknown check names fall back to weight 1 so new doctor checks score.
	unknown := []checkResult{{Name: "Some Future Check", Status: "ok"}}
	if got := computeScore(unknown); got != 100 {
		t.Errorf("unknown check: score = %d, want 100", got)
	}

	if got := computeScore(nil); got != 0 {
		t.Errorf("empty results: score = %d, want 0", got)
	}
}

// ── badge ────────────────────────────────────────────────────────────────────

func TestScoreBadgeURL(t *testing.T) {
	url := scoreBadgeURL(87, "B")
	want := "https://img.shields.io/badge/AutoDev_Score-87%2F100_B-0070f3"
	if url != want {
		t.Errorf("scoreBadgeURL = %q, want %q", url, want)
	}

	colors := map[string]string{
		"A": "00c853", "B": "0070f3", "C": "f5a623", "D": "f5a623", "F": "ee0000",
	}
	for grade, want := range colors {
		if got := scoreBadgeColor(grade); got != want {
			t.Errorf("scoreBadgeColor(%q) = %q, want %q", grade, got, want)
		}
	}
}

func TestScoreBadgeMarkdown(t *testing.T) {
	md := scoreBadgeMarkdown(92, "A")
	for _, want := range []string{
		"[![AutoDev Score](",
		"img.shields.io/badge/AutoDev_Score-92%2F100_A-00c853",
		"](https://github.com/HEETMEHTA18/autodev)",
	} {
		if !strings.Contains(md, want) {
			t.Errorf("badge markdown missing %q\ngot: %s", want, md)
		}
	}
}

// ── summary + CI gate ────────────────────────────────────────────────────────

func TestSummarizeChecks(t *testing.T) {
	s := summarizeChecks([]checkResult{
		{Status: "ok"}, {Status: "ok"}, {Status: "warning"}, {Status: "error"},
	})
	if s.Passed != 2 || s.Warnings != 1 || s.Errors != 1 {
		t.Errorf("summarize = %+v, want {2 1 1}", s)
	}
}

func TestCheckScoreMinimum(t *testing.T) {
	if err := checkScoreMinimum(10, 0); err != nil {
		t.Errorf("min=0 should never fail, got %v", err)
	}
	if err := checkScoreMinimum(80, 80); err != nil {
		t.Errorf("score == min should pass, got %v", err)
	}
	if err := checkScoreMinimum(79, 80); err == nil {
		t.Error("score < min should fail")
	}
}

// ── command wiring ───────────────────────────────────────────────────────────

func TestScoreCmdFlags(t *testing.T) {
	cmd := newScoreCmd()

	if cmd.Short == "" || cmd.Long == "" {
		t.Error("score command needs Short and Long descriptions")
	}

	f := cmd.Flags().Lookup("badge")
	if f == nil || f.DefValue != "false" {
		t.Error("expected --badge flag defaulting to false")
	}
	m := cmd.Flags().Lookup("min")
	if m == nil || m.DefValue != "0" {
		t.Error("expected --min flag defaulting to 0")
	}

	if g := commandGroup("score"); g != groupDiagnose {
		t.Errorf("score group = %q, want %q", g, groupDiagnose)
	}
}
