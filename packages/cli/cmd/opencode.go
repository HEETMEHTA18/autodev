package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/autodev-sh/autodev/core/osinfo"
	"github.com/autodev-sh/autodev/registry"
	"github.com/autodev-sh/autodev/scanner"
	"github.com/spf13/cobra"
)

// newOpenCodeCmd launches a new OpenCode session that is seeded with shared,
// project-wide knowledge so that every session understands the current state of
// the codebase — what has changed, what is in use, and what to do next.
func newOpenCodeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "opencode [task...]",
		Short: "Start a new OpenCode session with shared codebase knowledge",
		Long: `Start a new OpenCode session that is aware of the whole project.

Before launching, AutoDev builds a shared knowledge file that captures:
  - repository metadata (branch, commit, changed files, diff)
  - the detected stack (languages, frameworks, dependencies)
  - memory and past session logs from .autodevs/
  - the current AGENTS / skill rules

That knowledge file is handed to the new OpenCode session and appended to after
each session, so every session knows what was changed and what to do next.

  autodev opencode "fix the auth bug"
  autodev opencode            # start an interactive session`,
		Args:               cobra.ArbitraryArgs,
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				switch args[0] {
				case "--help", "-h", "help":
					return cmd.Help()
				}
			}
			return runOpenCode(strings.Join(args, " "))
		},
	}
}

// sharedKnowledgePath returns the canonical location of the shared knowledge
// file. It lives inside the project's .autodevs/ folder when present, otherwise
// in the user config directory so it follows the user across machines.
func sharedKnowledgePath() string {
	if _, err := os.Stat(".autodevs"); err == nil {
		return filepath.Join(".autodevs", "context", "opencode-knowledge.md")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "opencode-knowledge.md"
	}
	dir := filepath.Join(home, ".autodevs", "context")
	_ = os.MkdirAll(dir, 0755)
	return filepath.Join(dir, "opencode-knowledge.md")
}

// runOpenCode launches a new OpenCode session seeded with shared knowledge.
func runOpenCode(task string) error {
	agent, err := registry.Get("opencode")
	if err != nil {
		return err
	}
	status := registry.Detect(agent)
	if !status.Installed {
		fmt.Println()
		fmt.Println(renderAgentBox(
			agentTitleStyle.Render("  OpenCode Agent"),
			"",
			agentWarnStyle.Render("  ✗ OpenCode is not installed."),
			agentDimStyle.Render("  Install it with: autodev tools install opencode"),
			"",
		))
		fmt.Println()
		return nil
	}

	// 1. Build shared knowledge for this project.
	knowledgePath := sharedKnowledgePath()
	if err := os.MkdirAll(filepath.Dir(knowledgePath), 0755); err != nil {
		return err
	}
	knowledge, err := buildSharedKnowledge()
	if err != nil {
		return err
	}
	if err := os.WriteFile(knowledgePath, []byte(knowledge), 0644); err != nil {
		return err
	}

	// 2. Compose the session prompt so the agent grounds itself in the project.
	prompt := task
	if prompt == "" {
		prompt = "Inspect the current codebase, understand what has changed, and propose next steps."
	}
	seeded := fmt.Sprintf(
		"You are AutoDev, operating on this project.\n\n"+
			"Read the shared knowledge file at %s for the full state of this codebase "+
			"(git changes, detected stack, past session memory, and project rules). "+
			"Use it to ground every decision you make.\n\n"+
			"Task: %s",
		knowledgePath, prompt,
	)

	// 3. Print a small launch card.
	fmt.Println()
	fmt.Println(renderAgentBox(renderRouteCard(prompt, "coding", "OpenCode",
		"opencode run --title \"AutoDev session\"", agentDimStyle.Render("  Starting a new OpenCode session ..."))...))
	fmt.Println()

	// 4. Launch a brand-new OpenCode session (no --continue / --session).
	path, err := exec.LookPath(agent.Command)
	if err != nil {
		return fmt.Errorf("opencode is not installed. Run: autodev tools install opencode")
	}

	args := []string{"run", "--title", "AutoDev session", seeded}
	cmd := exec.Command(path, args...)
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		// The session may have been interrupted; still record the outcome.
		appendSessionOutcome(knowledgePath, task, err)
		return fmt.Errorf("opencode exited: %w", err)
	}

	appendSessionOutcome(knowledgePath, task, nil)
	return nil
}

// buildSharedKnowledge assembles a markdown knowledge bundle describing the
// current project so OpenCode sessions share a consistent mental model.
func buildSharedKnowledge() (string, error) {
	var b strings.Builder
	now := time.Now().Format(time.RFC3339)
	b.WriteString(fmt.Sprintf("# AutoDev Shared Codebase Knowledge\n\n_Updated %s_\n\n", now))

	// Project / OS context.
	info, _ := osinfo.Detect()
	if info != nil {
		b.WriteString(fmt.Sprintf("- **OS:** %s (%s)\n", info.Version, info.Arch))
	}
	if wd, err := os.Getwd(); err == nil {
		b.WriteString(fmt.Sprintf("- **Project root:** %s\n", wd))
	}

	// Git state.
	b.WriteString("\n## Repository state\n\n")
	gitLines := collectGitState()
	if len(gitLines) == 0 {
		b.WriteString("_Not a git repository._\n")
	} else {
		for _, l := range gitLines {
			b.WriteString(l + "\n")
		}
	}

	// Detected stack.
	b.WriteString("\n## Detected stack\n\n")
	if result, err := scanner.New(".").Scan(); err == nil && result != nil {
		b.WriteString("- **Languages:** " + joinOrNone(result.Languages) + "\n")
		b.WriteString("- **Frameworks:** " + joinOrNone(result.Frameworks) + "\n")
		b.WriteString("- **Package managers:** " + joinOrNone(result.PackageManagers) + "\n")
		b.WriteString("- **Databases:** " + joinOrNone(result.Databases) + "\n")
	} else {
		b.WriteString("_Could not scan project._\n")
	}

	// Memory and past session logs.
	b.WriteString("\n## AutoDev memory & past sessions\n\n")
	for _, f := range []string{".autodevs/memory/successes.md", ".autodevs/memory/patterns.md", ".autodevs/memory/mistakes.md"} {
		if data, err := os.ReadFile(f); err == nil {
			b.WriteString(fmt.Sprintf("### %s\n%s\n", filepath.Base(f), string(data)))
		}
	}

	// Session outcome history lives at the end of the knowledge file itself.
	if data, err := os.ReadFile(sharedKnowledgePath()); err == nil {
		content := string(data)
		// Carry forward the historical outcome log so new sessions keep context.
		if idx := strings.Index(content, "\n## Session outcomes\n"); idx != -1 {
			b.WriteString(content[idx:])
		}
	}

	return b.String(), nil
}

func collectGitState() []string {
	var lines []string
	run := func(name string, args ...string) string {
		out, err := exec.Command(name, args...).Output()
		if err != nil {
			return ""
		}
		return strings.TrimSpace(string(out))
	}
	if branch := run("git", "branch", "--show-current"); branch != "" {
		lines = append(lines, "- **Branch:** "+branch)
	}
	if commit := run("git", "log", "-1", "--oneline"); commit != "" {
		lines = append(lines, "- **HEAD:** "+commit)
	}
	status := run("git", "status", "--short")
	if status == "" {
		lines = append(lines, "- **Working tree:** clean")
	} else {
		lines = append(lines, "- **Working tree:** has uncommitted changes")
		lines = append(lines, "```")
		lines = append(lines, status)
		lines = append(lines, "```")
	}
	if stat := run("git", "diff", "--stat"); stat != "" {
		lines = append(lines, "\nChanged files (uncommitted diff):")
		lines = append(lines, "```")
		lines = append(lines, stat)
		lines = append(lines, "```")
	}
	return lines
}

func joinOrNone(items []string) string {
	if len(items) == 0 {
		return "_none detected_"
	}
	return strings.Join(items, ", ")
}

// appendSessionOutcome records the result of a finished session into the shared
// knowledge file so the next session can pick up where this one left off.
func appendSessionOutcome(knowledgePath, task string, runErr error) {
	outcome := "✓ completed"
	if runErr != nil {
		outcome = "✗ interrupted (" + runErr.Error() + ")"
	}
	entry := fmt.Sprintf("\n- **%s** — %s — %s\n", time.Now().Format("2006-01-02 15:04"), outcome, truncate(task, 120))

	data, err := os.ReadFile(knowledgePath)
	if err != nil {
		data = []byte("# AutoDev Shared Codebase Knowledge\n")
	}
	content := string(data)
	if idx := strings.Index(content, "\n## Session outcomes\n"); idx != -1 {
		content = content[:idx] + "\n## Session outcomes\n" + entry + content[idx+len("\n## Session outcomes\n"):]
	} else {
		content += "\n## Session outcomes\n" + entry
	}
	_ = os.WriteFile(knowledgePath, []byte(content), 0644)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
