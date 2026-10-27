package cmd

import (
	"fmt"
	"os"

	"github.com/autodev-sh/autodev/installer"
	"github.com/autodev-sh/autodev/scanner"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func newSetupCmd() *cobra.Command {
	var path string
	var skipConfirm bool

	cmd := &cobra.Command{
		Use:   "setup [path]",
		Short: "Detect and install all missing runtimes and dependencies",
		Long: `Scan the repository, detect all required technologies, check which ones 
are missing, and install them automatically. This is the main AutoDev command.`,
		Example: `  autodev setup
  autodev setup ./my-project
  autodev setup --dry-run
  autodev setup --yes`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				path = args[0]
			}
			if path == "" {
				path = "."
			}
			return runSetup(path, dryRun, skipConfirm)
		},
	}

	cmd.Flags().StringVarP(&path, "path", "p", ".", "path to scan and setup")
	cmd.Flags().BoolVarP(&skipConfirm, "yes", "y", false, "skip confirmation prompts")
	return cmd
}

func runSetup(path string, isDryRun, skipConfirm bool) error {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#0070f3"))
	okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00c853")).Bold(true)
	warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ee0000"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))

	fmt.Println()
	fmt.Println(titleStyle.Render("  AutoDev Setup"))
	if isDryRun {
		fmt.Println(warnStyle.Render("  ! dry-run: no changes will be made"))
	}
	fmt.Println()

	// Step 1: Scan
	fmt.Println(dimStyle.Render("  → Scanning repository..."))
	s := scanner.New(path)
	result, err := s.Scan()
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	if len(result.Technologies) == 0 {
		fmt.Println(warnStyle.Render("  No technologies detected in this directory."))
		return nil
	}

	fmt.Printf(okStyle.Render("  ✓ Detected %d technologies\n"), len(result.Technologies))
	fmt.Println()

	// Step 2: Check what's installed
	inst := installer.New(isDryRun)
	langToRuntime := map[string]string{
		"Node.js": "nodejs",
		"Python":  "python",
		"Go":      "go",
		"Rust":    "rust",
		"Docker":  "docker",
		"Bun":     "bun",
		"pnpm":    "pnpm",
		"Flutter": "flutter",
		"Java":    "java",
		"Kotlin":  "kotlin",
		"PHP":     "php",
		"Ruby":    "ruby",
	}

	var toInstall []string
	fmt.Println(titleStyle.Render("  Checking installation status:"))
	for _, lang := range append(result.Languages, result.Infra...) {
		if rtName, ok := langToRuntime[lang]; ok {
			status := inst.CheckStatus(rtName)
			if status.Installed {
				fmt.Printf("  %s %-20s %s\n", okStyle.Render("✓"), lang, dimStyle.Render(status.Version))
			} else {
				fmt.Printf("  %s %-20s %s\n", warnStyle.Render("✗"), lang, warnStyle.Render("not installed"))
				toInstall = append(toInstall, rtName)
			}
		}
	}
	fmt.Println()

	if len(toInstall) == 0 {
		fmt.Println(okStyle.Render("  ✓ All required tools are already installed!"))
		return nil
	}

	fmt.Printf(titleStyle.Render("  %d tools to install: %v\n"), len(toInstall), toInstall)
	fmt.Println()

	if !skipConfirm && !isDryRun {
		fmt.Print("  Proceed with installation? [y/N] ")
		var answer string
		_, _ = fmt.Scanln(&answer)
		if answer != "y" && answer != "Y" && answer != "yes" {
			fmt.Println(dimStyle.Render("  Installation cancelled."))
			return nil
		}
	}

	// Step 3: Install
	exitCode := 0
	for _, rtName := range toInstall {
		rt, _ := installer.GetRuntime(rtName)
		fmt.Printf(titleStyle.Render("\n  Installing %s...\n"), rt.Name)
		if err := inst.Install(rtName); err != nil {
			fmt.Fprintln(os.Stderr, warnStyle.Render(fmt.Sprintf("  ✗ Failed to install %s: %v", rt.Name, err)))
			exitCode = 1
		} else {
			fmt.Println(okStyle.Render(fmt.Sprintf("  ✓ %s installed successfully", rt.Name)))
		}
	}

	fmt.Println()
	if exitCode == 0 {
		fmt.Println(okStyle.Render("  ✓ Setup complete! Run 'autodev doctor' to verify your environment."))
		trackCLIMetric("setup_complete")
	} else {
		fmt.Println(warnStyle.Render("  Setup completed with errors. See above for details."))
		os.Exit(exitCode)
	}
	PrintGitHubCTA()
	return nil
}
