package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/autodev-sh/autodev/catalog"
	"github.com/autodev-sh/autodev/registry"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	agentOKStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF87")).Bold(true)
	agentWarnStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B"))
	agentDimStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	agentTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700"))
	agentRoleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#4A90E2"))
	agentBoxStyle   = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#4A90E2")).
			PaddingLeft(1).
			PaddingRight(1).
			Width(56)
	agentDivider = agentDimStyle.Render("  ────────────────────────────────────────────")
)

// renderAgentBox wraps content in AutoDevs' signature rounded box.
func renderAgentBox(lines ...string) string {
	return agentBoxStyle.Render(strings.Join(lines, "\n"))
}

// newToolsCmd manages the installed AI agent CLIs.
func newToolsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tools",
		Short: "List, install, and remove AI agent CLIs",
		Long: `Manage the AI/developer CLIs AutoDevs orchestrates (OpenCode, Claude Code,
Codex, Gemini CLI, Aider).

  autodev tools list              — show every agent and its install status
  autodev tools install <agent>   — install an agent through its package manager
  autodev tools remove <agent>    — uninstall an agent`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runToolsList()
		},
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "list",
			Short: "List AI agents with install status",
			Long:  "Show every registered AI agent, its role, and whether it is installed.",
			RunE: func(cmd *cobra.Command, args []string) error {
				return runToolsList()
			},
		},
		&cobra.Command{
			Use:   "install <agent>",
			Short: "Install an AI agent CLI",
			Long:  "Install an AI agent CLI through its package manager (npm or pip).",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return runToolsInstall(args[0])
			},
		},
		&cobra.Command{
			Use:   "remove <agent>",
			Short: "Remove an AI agent CLI",
			Long:  "Uninstall an AI agent CLI through its package manager (npm or pip).",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return runToolsRemove(args[0])
			},
		},
	)
	return cmd
}

func runToolsList() error {
	statuses := registry.DetectAll()

	lines := []string{
		agentTitleStyle.Render("  AutoDevs Agent Registry"),
		"",
		fmt.Sprintf("  %-15s %-9s %s", agentRoleStyle.Render("AGENT"), agentRoleStyle.Render("ROLE"), agentRoleStyle.Render("STATUS")),
		agentDivider,
	}

	installedCount := 0
	for _, s := range statuses {
		status := agentWarnStyle.Render("✗ not installed")
		if s.Installed {
			installedCount++
			status = agentOKStyle.Render("✓ " + s.Version)
		}
		lines = append(lines, fmt.Sprintf("  %-15s %-9s %s", s.Agent.Name, agentRoleStyle.Render(string(s.Agent.Role)), status))
	}

	lines = append(lines,
		agentDivider,
		"",
		agentDimStyle.Render(fmt.Sprintf("  %d/%d agents installed", installedCount, len(statuses))),
		"",
		agentDimStyle.Render("  autodev tools install <agent>   install an agent"),
		agentDimStyle.Render("  autodev run <agent>             launch an agent"),
		agentDimStyle.Render("  autodev session                 manage sessions"),
	)

	fmt.Println()
	fmt.Println(renderAgentBox(lines...))
	fmt.Println()
	return nil
}

func runToolsInstall(id string) error {
	agent, err := registry.Get(id)
	if err != nil {
		return err
	}

	c, err := catalog.Load()
	if err != nil {
		return err
	}
	pkg, ok := c.GetPackage(agent.ID)
	if !ok {
		return fmt.Errorf("no catalog entry for agent %q", agent.ID)
	}

	if pkg.IsInstalled() {
		s := registry.Detect(agent)
		lines := []string{
			agentTitleStyle.Render("  " + agent.Name + " — already installed"),
			"",
			fmt.Sprintf("  %s %s %s", agentOKStyle.Render("✓"), agent.Name, agentDimStyle.Render(s.Version)),
		}
		fmt.Println()
		fmt.Println(renderAgentBox(lines...))
		fmt.Println()
		return nil
	}

	fmt.Println()
	fmt.Printf("  %s\n", agentTitleStyle.Render("Installing "+agent.Name))
	if err := execInstall(pkg); err != nil {
		return fmt.Errorf("install %s: %w", agent.ID, err)
	}

	s := registry.Detect(agent)
	lines := []string{agentTitleStyle.Render("  " + agent.Name + " installed")}
	if !s.Installed {
		lines = append(lines,
			"",
			agentWarnStyle.Render("  ✗ Installed but not found on PATH"),
			agentDimStyle.Render("  Restart your shell or add the binary to PATH."),
		)
	} else {
		lines = append(lines,
			"",
			fmt.Sprintf("  %s %s %s", agentOKStyle.Render("✓"), agent.Name, agentDimStyle.Render(s.Version)),
			"",
			agentDimStyle.Render("  Launch it: autodev run "+agent.ID),
		)
	}
	fmt.Println()
	fmt.Println(renderAgentBox(lines...))
	fmt.Println()
	return nil
}

func runToolsRemove(id string) error {
	agent, err := registry.Get(id)
	if err != nil {
		return err
	}

	path, err := exec.LookPath(agent.Command)
	if err != nil {
		lines := []string{
			agentTitleStyle.Render("  " + agent.Name),
			"",
			agentWarnStyle.Render("  ✗ Not installed — nothing to remove."),
		}
		fmt.Println()
		fmt.Println(renderAgentBox(lines...))
		fmt.Println()
		return nil
	}

	var cmd *exec.Cmd
	switch agent.Manager {
	case "npm":
		cmd = exec.Command("npm", "uninstall", "-g", agent.Pkg)
	case "pip":
		cmd = exec.Command("pip3", "uninstall", "-y", agent.Pkg)
	default:
		return fmt.Errorf("no uninstall path for manager %q", agent.Manager)
	}

	fmt.Println()
	fmt.Printf("  %s\n", agentTitleStyle.Render("Removing "+agent.Name))
	fmt.Println(agentDimStyle.Render(fmt.Sprintf("  → %s (%s)", path, agent.Manager)))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("remove %s: %w", agent.ID, err)
	}

	lines := []string{
		agentTitleStyle.Render("  " + agent.Name + " removed"),
		"",
		fmt.Sprintf("  %s %s uninstalled", agentOKStyle.Render("✓"), agent.Name),
	}
	fmt.Println()
	fmt.Println(renderAgentBox(lines...))
	fmt.Println()
	return nil
}
