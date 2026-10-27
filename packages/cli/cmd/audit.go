package cmd

import (
	"fmt"
	"time"

	"github.com/autodev-sh/autodev/scanner"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func newAuditCmd() *cobra.Command {
	var path string

	cmd := &cobra.Command{
		Use:   "audit [path]",
		Short: "Audit repository dependencies for security vulnerabilities",
		Long:  `Scan the repository's dependency files (package.json, requirements.txt, go.mod) and query the OSV (Open Source Vulnerabilities) database to find any known safety risks.`,
		Example: `  autodev audit
  autodev audit ./my-project`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				path = args[0]
			}
			if path == "" {
				path = "."
			}
			return runAudit(path)
		},
	}

	return cmd
}

func runAudit(path string) error {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#0070f3"))
	successStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00c853"))
	warnStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ee0000"))
	infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#0070f3")).Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))

	criticalStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ee0000"))
	highStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ee0000"))
	moderateStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#0070f3"))
	lowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))

	fmt.Println()
	fmt.Println(titleStyle.Render("  AutoDev Supply-Chain Safety Audit"))
	fmt.Println(dimStyle.Render("  Auditing dependencies against the OSV Vulnerability Database..."))
	fmt.Println()

	start := time.Now()
	results, err := scanner.AuditRepository(path)
	if err != nil {
		return fmt.Errorf("audit failed: %w", err)
	}
	elapsed := time.Since(start)

	if len(results) == 0 {
		fmt.Println(successStyle.Render("  ✓ No known security vulnerabilities found! All dependencies are safe."))
		fmt.Println()
		PrintGitHubCTA()
		return nil
	}

	totalVulns := 0
	for _, res := range results {
		totalVulns += len(res.Vulnerabilities)
	}

	fmt.Println(warnStyle.Render(fmt.Sprintf("  ✗ Found %d security vulnerabilities across %d packages:", totalVulns, len(results))))
	fmt.Println()

	for _, res := range results {
		pkgHeader := fmt.Sprintf("%s@%s (%s)", res.Package.Name, res.Package.Version, res.Package.Ecosystem)
		fmt.Println(infoStyle.Render(pkgHeader))

		for _, v := range res.Vulnerabilities {
			var sevBadge string
			switch v.Severity {
			case "CRITICAL":
				sevBadge = criticalStyle.Render("[CRITICAL]")
			case "HIGH":
				sevBadge = highStyle.Render("[HIGH]")
			case "MODERATE", "MEDIUM":
				sevBadge = moderateStyle.Render("[MODERATE]")
			default:
				sevBadge = lowStyle.Render(fmt.Sprintf("[%s]", v.Severity))
			}

			alias := ""
			if len(v.Aliases) > 0 {
				alias = fmt.Sprintf(" (%s)", v.Aliases[0])
			}

			fmt.Printf("  - %s %s%s: %s\n", sevBadge, v.ID, alias, v.Summary)
		}
		fmt.Println()
	}

	fmt.Println(dimStyle.Render(fmt.Sprintf("  Audited in %s", elapsed.Round(time.Millisecond))))
	fmt.Println()

	PrintGitHubCTA()
	return nil
}
