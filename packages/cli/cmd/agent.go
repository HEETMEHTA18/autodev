package cmd

import (
	"fmt"
	"strings"

	"github.com/autodev-sh/autodev/registry"
	"github.com/spf13/cobra"
)

// newAgentCmd routes a task description to the best-fit agent and launches it.
func newAgentCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "agent <task...>",
		Short: "Route a task to the best-fit AI agent",
		Long: `Analyze a task description and launch the most suitable AI agent.

  autodev agent "fix the authentication bug"          → OpenCode (coding)
  autodev agent "research the latest postgres release"→ Gemini (research)
  autodev agent "review this pull request"            → Codex (review)
  autodev agent "harden this codebase"                → AutoDevs security

Install a missing agent automatically with: autodev tools install <agent>`,
		Args:               cobra.MinimumNArgs(1),
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRouter(strings.Join(args, " "))
		},
	}
}

// classifyTask picks the agent id for a task description.
func classifyTask(task string) string {
	t := strings.ToLower(task)

	if hasAny(t, "security", "vuln", "cve", "harden", "malware", "audit", "exploit", "inject") {
		return "security"
	}
	if hasAny(t, "review", "pull request", "pr", "code quality", "lint", "refactor check") {
		return "codex"
	}
	if hasAny(t, "research", "explain", "search", "summarize", "find", "why", "compare", "docs", "document") {
		return "gemini"
	}
	if hasAny(t, "complex", "architecture", "design", "plan", "migrate", "full system") {
		return "claude"
	}
	return "opencode"
}

func hasAny(s string, needles ...string) bool {
	for _, n := range needles {
		if strings.Contains(s, n) {
			return true
		}
	}
	return false
}

// renderRouteCard builds the routing decision card.
func renderRouteCard(task, taskType, name, command, extra string) []string {
	display := task
	if r := []rune(task); len(r) > 34 {
		display = string(r[:34]) + "…"
	}
	lines := []string{
		agentTitleStyle.Render("  Task Routing"),
		"",
		fmt.Sprintf("  %-14s %s", agentDimStyle.Render("Task"), agentDimStyle.Render(display)),
		fmt.Sprintf("  %-14s %s", agentDimStyle.Render("Type"), agentRoleStyle.Render(taskType)),
		fmt.Sprintf("  %-14s %s", agentDimStyle.Render("Agent"), name),
		fmt.Sprintf("  %-14s %s", agentDimStyle.Render("Command"), command),
	}
	if extra != "" {
		lines = append(lines, "", extra)
	}
	return lines
}

func runRouter(task string) error {
	id := classifyTask(task)
	if id == "security" {
		return runRouterSecurity(task)
	}

	agent, err := registry.Get(id)
	if err != nil {
		return err
	}

	s := registry.Detect(agent)
	extra := ""
	if !s.Installed {
		extra = fmt.Sprintf("%s\n  %s",
			agentWarnStyle.Render("  ✗ "+agent.Name+" is not installed."),
			agentDimStyle.Render("  autodev tools install "+agent.ID))
	} else {
		extra = agentDimStyle.Render("  Starting " + agent.Name + " ...")
	}

	fmt.Println()
	fmt.Println(renderAgentBox(renderRouteCard(task, string(agent.Role), agent.Name, agent.Command, extra)...))
	fmt.Println()

	if !s.Installed {
		return nil
	}
	return runAgent(agent.ID, []string{task})
}

func runRouterSecurity(task string) error {
	lines := renderRouteCard(task, "security", "AutoDevs Security", "autodev security",
		agentDimStyle.Render("  Starting AutoDevs security scan ..."))

	fmt.Println()
	fmt.Println(renderAgentBox(lines...))
	fmt.Println()
	return runSecurity()
}
