package cmd

import (
	"strings"
	"testing"

	"github.com/autodev-sh/autodev/core/promptcapture"
)

// ── chat flags ───────────────────────────────────────────────────────────────

func TestChatFlags(t *testing.T) {
	cmd := newChatCmd()

	cases := []struct {
		name     string
		defValue string
	}{
		{"thinking-budget", "1024"},
		{"show-thoughts", "false"},
		{"no-refine", "false"},
		{"no-grounding", "false"},
	}
	for _, c := range cases {
		f := cmd.Flags().Lookup(c.name)
		if f == nil {
			t.Errorf("expected --%s flag", c.name)
			continue
		}
		if f.DefValue != c.defValue {
			t.Errorf("--%s default = %q, want %q", c.name, f.DefValue, c.defValue)
		}
	}
}

// ── thinking budget clamp ────────────────────────────────────────────────────

func TestClampThinkingBudget(t *testing.T) {
	cases := []struct{ in, want int }{
		{-100, 0},
		{0, 0},
		{1024, 1024},
		{maxThinkingBudget, maxThinkingBudget},
		{999999, maxThinkingBudget},
	}
	for _, c := range cases {
		if got := clampThinkingBudget(c.in); got != c.want {
			t.Errorf("clampThinkingBudget(%d) = %d, want %d", c.in, got, c.want)
		}
	}
}

// ── system prompts ───────────────────────────────────────────────────────────

func TestAssistantSystemPromptDiscipline(t *testing.T) {
	p := assistantSystemPrompt()
	for _, want := range []string{
		"Scope discipline",
		"finance question must not turn into medical",
		"Anti-hallucination",
		"Never invent facts",
		"say so plainly instead of guessing",
		"Bounded thinking",
		"Think only as long as needed",
		"[RUN_COMMAND:",
		"[WRITE_FILE:",
	} {
		if !strings.Contains(p, want) {
			t.Errorf("assistant system prompt missing %q", want)
		}
	}
}

func TestRefineSystemPromptMarkers(t *testing.T) {
	p := refineSystemPrompt()
	if !strings.Contains(p, refineStartMarker) || !strings.Contains(p, refineEndMarker) {
		t.Error("refine prompt must require both markers")
	}
	for _, want := range []string{"strictly inside the original topic domain", "admit when it does not know"} {
		if !strings.Contains(p, want) {
			t.Errorf("refine prompt missing %q", want)
		}
	}
}

// ── refined prompt extraction ────────────────────────────────────────────────

func TestExtractRefinedPrompt(t *testing.T) {
	good := "some chatter\n<<REFINED>>\nAnswer this in 3 bullets.\n<<END_REFINED>>\nmore chatter"
	got, err := extractRefinedPrompt(good)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "Answer this in 3 bullets." {
		t.Errorf("got %q", got)
	}

	if _, err := extractRefinedPrompt("no markers here"); err == nil {
		t.Error("expected error for missing markers")
	}
	if _, err := extractRefinedPrompt("<<REFINED>> text without end"); err == nil {
		t.Error("expected error for missing end marker")
	}
	if _, err := extractRefinedPrompt("<<REFINED>><<END_REFINED>>"); err == nil {
		t.Error("expected error for empty refined prompt")
	}
}

// ── contents builder ─────────────────────────────────────────────────────────

func TestBuildContents(t *testing.T) {
	history := []promptcapture.PromptEvent{
		{Prompt: "hello", Response: "hi"},
		{Prompt: "", Response: ""}, // skipped
	}
	contents := buildContents("do the work", history)

	if len(contents) != 3 {
		t.Fatalf("len(contents) = %d, want 3 (1 history pair + current)", len(contents))
	}
	if contents[0].Role != "user" || contents[0].Parts[0].Text != "hello" {
		t.Errorf("unexpected first content: %+v", contents[0])
	}
	if contents[1].Role != "model" || contents[1].Parts[0].Text != "hi" {
		t.Errorf("unexpected second content: %+v", contents[1])
	}
	last := contents[len(contents)-1]
	if last.Role != "user" || last.Parts[0].Text != "do the work" {
		t.Errorf("last content must be the current prompt, got %+v", last)
	}
}

// ── response parsing ─────────────────────────────────────────────────────────

func TestParseGeminiResponse(t *testing.T) {
	body := []byte(`{
		"candidates": [{
			"content": {"parts": [
				{"text": "thought noise", "thought": true},
				{"text": "part one "},
				{"text": "part two"}
			]},
			"groundingMetadata": {"groundingChunks": [
				{"web": {"uri": "https://a.example", "title": "Doc A"}},
				{"web": {"uri": "https://a.example", "title": "Doc A dup"}},
				{"web": {"uri": "https://b.example", "title": ""}}
			]}
		}]
	}`)

	reply, err := parseGeminiResponse(body, true)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if reply.Text != "part one \npart two" {
		t.Errorf("text = %q", reply.Text)
	}
	if len(reply.Thoughts) != 1 || reply.Thoughts[0] != "thought noise" {
		t.Errorf("thoughts = %v", reply.Thoughts)
	}
	if len(reply.Sources) != 2 {
		t.Errorf("sources = %v (want deduped, empty title → uri)", reply.Sources)
	}
	if reply.Sources[1] != "https://b.example — https://b.example" {
		t.Errorf("empty title should fall back to uri, got %q", reply.Sources[1])
	}

	// thoughts dropped when not requested
	reply, err = parseGeminiResponse(body, false)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(reply.Thoughts) != 0 {
		t.Errorf("expected no thoughts, got %v", reply.Thoughts)
	}

	// blocked request
	blocked := []byte(`{"promptFeedback": {"blockReason": "SAFETY"}}`)
	if _, err := parseGeminiResponse(blocked, false); err == nil || !strings.Contains(err.Error(), "SAFETY") {
		t.Errorf("expected block-reason error, got %v", err)
	}

	// no candidates
	if _, err := parseGeminiResponse([]byte(`{"candidates": []}`), false); err == nil {
		t.Error("expected error for empty candidates")
	}
}
