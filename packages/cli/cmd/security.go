package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func newSecurityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "security",
		Short: "Run local security checks using installed tools",
		Long:  `Run a comprehensive security audit of your codebase using industry-standard tools like gosec, govulncheck, npm audit, and more.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSecurity()
		},
	}
	return cmd
}

func runSecurity() error {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700"))
	successStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FF87"))
	warnStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF6B6B"))
	infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#4A90E2")).Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))

	fmt.Println()
	fmt.Println(titleStyle.Render("🛡️  AutoDev Security Doctor"))
	fmt.Println(dimStyle.Render("  Scanning local environment with available security tools..."))
	fmt.Println()

	var recommendations []string
	criticalCount := 0
	highCount := 0
	mediumCount := 0

	// 1. Go Tools
	if _, err := os.Stat("go.mod"); err == nil {
		fmt.Println(infoStyle.Render("  📦 Checking Go Modules..."))
		if err := runExternalCommand("go", "mod", "verify"); err != nil {
			recommendations = append(recommendations, "Fix module verification failures (run 'go mod verify')")
			highCount++
		} else {
			fmt.Println(successStyle.Render("    ✓ go mod verify passed"))
		}

		// go mod tidy -diff is only available in go1.21+ but we can just run go mod tidy and check if git status changes
		// For safety and compatibility, we'll just skip the -diff flag if it fails, but let's try it.
		// Actually, go mod tidy doesn't have -diff in all versions. We'll just suggest running go mod tidy if it's messy.

		if commandExists("govulncheck") {
			fmt.Println(infoStyle.Render("  🔍 Running govulncheck..."))
			out, err := exec.Command("govulncheck", "./...").CombinedOutput()
			if err != nil {
				outStr := string(out)
				if strings.Contains(outStr, "Vulnerability") {
					recommendations = append(recommendations, "Address govulncheck findings")
					criticalCount++ // Assume govulncheck findings are critical/high
				} else {
					fmt.Println(warnStyle.Render(fmt.Sprintf("    ✗ govulncheck failed to execute: %v", err)))
				}
			} else {
				fmt.Println(successStyle.Render("    ✓ govulncheck passed"))
			}
		} else {
			fmt.Println(dimStyle.Render("    - govulncheck not found (skip)"))
			recommendations = append(recommendations, "Install govulncheck for Go vulnerability scanning")
		}

		if commandExists("gosec") {
			fmt.Println(infoStyle.Render("  🔒 Running gosec..."))
			out, err := exec.Command("gosec", "-quiet", "./...").CombinedOutput()
			if err != nil {
				outStr := string(out)
				lines := strings.Split(outStr, "\n")
				for _, line := range lines {
					if strings.Contains(line, "Severity: HIGH") {
						highCount++
					} else if strings.Contains(line, "Severity: MEDIUM") {
						mediumCount++
					}
				}
				if strings.Contains(outStr, "Severity: HIGH") || strings.Contains(outStr, "Severity: MEDIUM") {
					recommendations = append(recommendations, "Address gosec static analysis findings")
				} else {
					fmt.Println(warnStyle.Render(fmt.Sprintf("    ✗ gosec failed to execute: %v", err)))
				}
			} else {
				fmt.Println(successStyle.Render("    ✓ gosec passed"))
			}
		} else {
			fmt.Println(dimStyle.Render("    - gosec not found (skip)"))
			recommendations = append(recommendations, "Install gosec for Go static analysis")
		}
	}

	// 2. Node.js Tools
	if _, err := os.Stat("package.json"); err == nil {
		if commandExists("npm") {
			fmt.Println(infoStyle.Render("  📦 Running npm audit..."))
			out, err := exec.Command("npm", "audit").CombinedOutput()
			if err != nil {
				outStr := string(out)
				if strings.Contains(outStr, "critical") {
					criticalCount++
				}
				if strings.Contains(outStr, "high") {
					highCount++
				}
				if strings.Contains(outStr, "moderate") {
					mediumCount++
				}
				if strings.Contains(outStr, "critical") || strings.Contains(outStr, "high") || strings.Contains(outStr, "moderate") {
					recommendations = append(recommendations, "Run 'npm audit fix' to resolve Node.js vulnerabilities")
				} else {
					fmt.Println(warnStyle.Render(fmt.Sprintf("    ✗ npm audit failed to execute: %v", err)))
				}
			} else {
				fmt.Println(successStyle.Render("    ✓ npm audit passed"))
			}
		}
	}

	// 3. Python Tools
	if _, err := os.Stat("requirements.txt"); err == nil {
		if commandExists("pip-audit") {
			fmt.Println(infoStyle.Render("  🐍 Running pip-audit..."))
			out, err := exec.Command("pip-audit").CombinedOutput()
			if err != nil {
				if strings.Contains(string(out), "Found") || strings.Contains(string(out), "Vulnerabilities") {
					highCount++
					recommendations = append(recommendations, "Address pip-audit vulnerability findings")
				} else {
					fmt.Println(warnStyle.Render(fmt.Sprintf("    ✗ pip-audit failed to execute: %v", err)))
				}
			} else {
				fmt.Println(successStyle.Render("    ✓ pip-audit passed"))
			}
		} else {
			fmt.Println(dimStyle.Render("    - pip-audit not found (skip)"))
		}
	}

	// 4. Rust Tools
	if _, err := os.Stat("Cargo.toml"); err == nil {
		if commandExists("cargo-audit") || commandExists("cargo") { // cargo audit is a subcommand
			fmt.Println(infoStyle.Render("  🦀 Running cargo audit..."))
			out, err := exec.Command("cargo", "audit").CombinedOutput()
			if err != nil {
				if strings.Contains(string(out), "error") || strings.Contains(string(out), "Vulnerabilities") || err != nil {
					// cargo audit exits with 1 on vulns
					highCount++
					recommendations = append(recommendations, "Address cargo audit findings")
				} else {
					fmt.Println(warnStyle.Render(fmt.Sprintf("    ✗ cargo audit failed to execute: %v", err)))
				}
			} else {
				fmt.Println(successStyle.Render("    ✓ cargo audit passed"))
			}
		}
	}

	// 5. General / Trivy / Semgrep
	if commandExists("trivy") {
		fmt.Println(infoStyle.Render("  🚢 Running trivy fs scanner..."))
		out, err := exec.Command("trivy", "fs", "--quiet", "--severity", "CRITICAL,HIGH", ".").CombinedOutput()
		if err != nil {
			if len(out) > 0 {
				criticalCount++
				recommendations = append(recommendations, "Address Trivy filesystem vulnerabilities")
			} else {
				fmt.Println(warnStyle.Render(fmt.Sprintf("    ✗ trivy failed to execute: %v", err)))
			}
		} else {
			fmt.Println(successStyle.Render("    ✓ trivy passed"))
		}
	}

	if commandExists("semgrep") {
		fmt.Println(infoStyle.Render("  🔎 Running semgrep..."))
		out, err := exec.Command("semgrep", "scan", "--quiet", "--error").CombinedOutput()
		if err != nil {
			if len(out) > 0 {
				highCount++
				recommendations = append(recommendations, "Address Semgrep static analysis findings")
			} else {
				fmt.Println(warnStyle.Render(fmt.Sprintf("    ✗ semgrep failed to execute: %v", err)))
			}
		} else {
			fmt.Println(successStyle.Render("    ✓ semgrep passed"))
		}
	}

	// Calculate Score
	score := 100
	score -= (criticalCount * 15)
	score -= (highCount * 5)
	score -= (mediumCount * 2)
	if score < 0 {
		score = 0
	}

	fmt.Println()
	scoreColor := successStyle
	if score < 70 {
		scoreColor = warnStyle
	} else if score < 90 {
		scoreColor = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700"))
	}

	fmt.Printf("  %s %s\n", titleStyle.Render("Security Score:"), scoreColor.Render(fmt.Sprintf("%d/100", score)))
	fmt.Println()

	if criticalCount > 0 || highCount > 0 || mediumCount > 0 {
		fmt.Printf("  %s %d\n", warnStyle.Render("Critical:"), criticalCount)
		fmt.Printf("  %s %d\n", lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B")).Render("High:"), highCount)
		fmt.Printf("  %s %d\n", lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Render("Medium:"), mediumCount)
		fmt.Println()
	}

	if len(recommendations) > 0 {
		fmt.Println(titleStyle.Render("  Recommendations:"))
		for _, rec := range recommendations {
			fmt.Printf("  %s %s\n", warnStyle.Render("✗"), rec)
		}
	} else {
		fmt.Println(successStyle.Render("  ✓ No security recommendations at this time. Great job!"))
	}

	fmt.Println()
	return nil
}

func runExternalCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%v: %s", err, stderr.String())
	}
	return nil
}
