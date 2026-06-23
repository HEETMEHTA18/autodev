package cmd

import (
	"fmt"
	"os/exec"
	"regexp"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func newSecretsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secrets",
		Short: "Scan codebase for exposed secrets and credentials",
		Long:  `Scans the current repository for accidentally committed secrets like AWS Keys, GitHub Tokens, and OpenAI API keys using gitleaks (if installed) and built-in regex heuristics.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSecrets()
		},
	}
	return cmd
}

func runSecrets() error {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700"))
	successStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FF87"))
	warnStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF6B6B"))
	infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#4A90E2")).Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))

	fmt.Println()
	fmt.Println(titleStyle.Render("🕵️  AutoDev Secrets Scanner"))
	fmt.Println()

	foundSecrets := false

	if commandExists("gitleaks") {
		fmt.Println(infoStyle.Render("  Scanning with gitleaks..."))
		out, err := exec.Command("gitleaks", "detect", "-v", "--no-git").CombinedOutput()
		if err != nil {
			fmt.Println(warnStyle.Render("  [!] Secrets detected by gitleaks:"))
			fmt.Println(dimStyle.Render(string(out)))
			foundSecrets = true
		} else {
			fmt.Println(successStyle.Render("  ✓ No secrets detected by gitleaks"))
		}
	} else {
		fmt.Println(dimStyle.Render("  (gitleaks not installed. Falling back to built-in basic patterns)"))
		// Fallback simple scanner for Phase 1
		found := runBasicSecretScan()
		if found {
			foundSecrets = true
		} else {
			fmt.Println(successStyle.Render("  ✓ No common secrets detected in basic scan"))
		}
	}

	fmt.Println()
	if foundSecrets {
		fmt.Println(warnStyle.Render("  Status: VULNERABLE - Please revoke exposed credentials immediately!"))
	} else {
		fmt.Println(successStyle.Render("  Status: SECURE - No secrets exposed in repository."))
	}
	fmt.Println()
	return nil
}

func runBasicSecretScan() bool {
	// A basic regex fallback just to prove functionality.
	// In production, users should have gitleaks installed.
	warnStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF6B6B"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))

	patterns := map[string]*regexp.Regexp{
		"OpenAI API Key":    regexp.MustCompile(`sk-[a-zA-Z0-9]{48}`),
		"Anthropic API Key": regexp.MustCompile(`sk-ant-[a-zA-Z0-9_-]{70,}`),
		"GitHub Token":      regexp.MustCompile(`ghp_[a-zA-Z0-9]{36}`),
		"AWS Access Key":    regexp.MustCompile(`AKIA[0-9A-Z]{16}`),
		"Stripe Standard":   regexp.MustCompile(`sk_live_[0-9a-zA-Z]{24}`),
		"Slack Token":       regexp.MustCompile(`xox[baprs]-[0-9a-zA-Z]{10,48}`),
		"RSA Private Key":   regexp.MustCompile(`-----BEGIN RSA PRIVATE KEY-----`),
	}

	foundAny := false

	// Use grep equivalent via exec if possible, or just skip if too large
	// For a quick v1, we'll try running 'git grep' with these patterns if git is available
	if commandExists("git") {
		for name, pattern := range patterns {
			cmd := exec.Command("git", "grep", "--no-index", "-E", "-n", pattern.String())
			out, err := cmd.CombinedOutput()
			if err == nil && len(out) > 0 {
				fmt.Println(warnStyle.Render(fmt.Sprintf("  [!] Exposed %s found:", name)))
				fmt.Println(dimStyle.Render(string(out)))
				foundAny = true
			}
		}
	} else {
		fmt.Println(dimStyle.Render("  (git not found, skipping basic scan)"))
	}

	return foundAny
}
