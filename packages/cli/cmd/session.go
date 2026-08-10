package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/autodev-sh/autodev/registry"
	"github.com/spf13/cobra"
)

// newSessionCmd manages agent sessions (launch, list, stop).
func newSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "session",
		Short: "Manage AI agent sessions",
		Long: `Manage AI agent sessions. A session is a running agent process tracked by
AutoDevs.

  autodev session            — show agent status and running sessions
  autodev session new <agent>— launch an agent (same as autodev run)
  autodev session list       — list running agent sessions
  autodev session stop <agent> — stop a running agent session`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSessionStatus()
		},
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:                "new <agent> [args...]",
			Short:              "Start a new agent session",
			Long:               "Launch an agent CLI in a managed terminal (alias of autodev run).",
			Args:               cobra.MinimumNArgs(1),
			DisableFlagParsing: true,
			RunE: func(cmd *cobra.Command, args []string) error {
				return runAgent(args[0], args[1:])
			},
		},
		&cobra.Command{
			Use:   "list",
			Short: "List running agent sessions",
			Long:  "List which registered agents currently have running processes.",
			RunE: func(cmd *cobra.Command, args []string) error {
				return runSessionList()
			},
		},
		&cobra.Command{
			Use:   "stop <agent>",
			Short: "Stop a running agent session",
			Long:  "Stop all processes matching the agent's binary.",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return runSessionStop(args[0])
			},
		},
	)
	return cmd
}

func runSessionStatus() error {
	lines := []string{
		agentTitleStyle.Render("  AutoDevs Sessions"),
		"",
		fmt.Sprintf("  %-15s %s", agentRoleStyle.Render("AGENT"), agentRoleStyle.Render("STATUS")),
		agentDivider,
	}

	for _, a := range registry.All() {
		s := registry.Detect(a)
		running := registry.Running(a)

		var status string
		switch {
		case !s.Installed:
			status = agentWarnStyle.Render("✗ not installed")
		case running:
			status = agentOKStyle.Render("● running")
		default:
			status = agentDimStyle.Render("○ idle  " + s.Version)
		}
		lines = append(lines, fmt.Sprintf("  %-15s %s", a.Name, status))
	}

	lines = append(lines,
		agentDivider,
		"",
		agentDimStyle.Render("  autodev session new <agent>    start a session"),
		agentDimStyle.Render("  autodev session list           show running sessions"),
		agentDimStyle.Render("  autodev session stop <agent>   stop a session"),
	)

	fmt.Println()
	fmt.Println(renderAgentBox(lines...))
	fmt.Println()
	return nil
}

func runSessionList() error {
	lines := []string{
		agentTitleStyle.Render("  Running Agent Sessions"),
		"",
	}

	found := false
	for _, a := range registry.All() {
		if !registry.Running(a) {
			continue
		}
		found = true
		lines = append(lines, fmt.Sprintf("  %s %-15s %s", agentOKStyle.Render("●"), a.Name, agentDimStyle.Render(a.Command)))
	}
	if !found {
		lines = append(lines, agentDimStyle.Render("  No agent sessions running."))
	}

	fmt.Println()
	fmt.Println(renderAgentBox(lines...))
	fmt.Println()
	return nil
}

func runSessionStop(id string) error {
	agent, err := registry.Get(id)
	if err != nil {
		return err
	}

	if _, err := exec.LookPath("pkill"); err != nil {
		return fmt.Errorf("pkill not available on this system")
	}

	lines := []string{agentTitleStyle.Render("  Stop Session — " + agent.Name)}
	cmd := exec.Command("pkill", "-f", "^"+agent.Command)
	if err := cmd.Run(); err != nil {
		lines = append(lines, "", agentDimStyle.Render("  No running "+agent.Name+" session found."))
	} else {
		lines = append(lines, "", fmt.Sprintf("  %s %s session stopped", agentOKStyle.Render("✓"), agent.Name))
	}

	fmt.Println()
	fmt.Println(renderAgentBox(lines...))
	fmt.Println()
	return nil
}

// normalizeAgentID maps user-friendly names to registry ids.
func normalizeAgentID(input string) string {
	input = strings.ToLower(strings.TrimSpace(input))
	switch input {
	case "claude-code", "claudecode", "claude_code", "claude code":
		return "claude"
	case "open-code", "opencode-ai", "open code":
		return "opencode"
	case "google", "gemini-cli", "gemini cli":
		return "gemini"
	}
	return input
}
