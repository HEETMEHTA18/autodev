package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/autodev-sh/autodev/installer"
	"github.com/autodev-sh/autodev/scanner"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func newCloneCmd() *cobra.Command {
	var skipConfirm bool

	cmd := &cobra.Command{
		Use:   "clone [repository-url] [target-directory]",
		Short: "Clone a Git repository, scan it, and install all missing dependencies",
		Long:  `Automatically clone any Git/GitHub repository, analyze its languages/frameworks, estimate setup/download time, and configure your local machine with all required tools.`,
		Example: `  autodev clone https://github.com/HEETMEHTA18/autodev
  autodev clone git@github.com:HEETMEHTA18/autodev.git my-project`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repoURL := args[0]
			var targetDir string
			if len(args) > 1 {
				targetDir = args[1]
			} else {
				targetDir = getRepoName(repoURL)
			}
			return runClone(repoURL, targetDir, skipConfirm)
		},
	}

	cmd.Flags().BoolVarP(&skipConfirm, "yes", "y", false, "skip confirmation prompts")
	return cmd
}

func getRepoName(url string) string {
	parts := strings.Split(strings.TrimSuffix(url, "/"), "/")
	last := parts[len(parts)-1]
	return strings.TrimSuffix(last, ".git")
}

func runClone(repoURL, targetDir string, skipConfirm bool) error {
	repoURL = strings.TrimSpace(repoURL)
	targetDir = strings.TrimSpace(targetDir)

	if strings.HasPrefix(repoURL, "-") {
		return fmt.Errorf("invalid repository URL: cannot start with a hyphen")
	}
	if targetDir != "" && strings.HasPrefix(targetDir, "-") {
		return fmt.Errorf("invalid target directory: cannot start with a hyphen")
	}

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#0070f3"))
	okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00c853")).Bold(true)
	warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ee0000"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	highlightStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#0070f3")).Bold(true)

	fmt.Println()
	fmt.Println(titleStyle.Render("AUTODEV CLONE & DEPLOY"))
	fmt.Printf("  Target:  %s\n", highlightStyle.Render(repoURL))
	fmt.Printf("  Folder:  %s\n\n", highlightStyle.Render(targetDir))

	// Ensure git is installed
	_, err := exec.LookPath("git")
	if err != nil {
		fmt.Println(warnStyle.Render("  [ERROR] git is not installed on your system."))
		fmt.Println(dimStyle.Render("  Installing git first..."))
		inst := installer.New(false)
		if err := inst.Install("git"); err != nil {
			return fmt.Errorf("failed to install git: %w", err)
		}
		fmt.Println(okStyle.Render("  ✓ git installed successfully"))
	}

	// 1. Clone
	fmt.Println(dimStyle.Render("  → Cloning repository..."))
	var gitArgs []string
	if targetDir != "" {
		gitArgs = []string{"clone", "--", repoURL, targetDir}
	} else {
		gitArgs = []string{"clone", "--", repoURL}
	}
	cloneCmd := exec.Command("git", gitArgs...)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr
	if err := cloneCmd.Run(); err != nil {
		return fmt.Errorf("git clone failed: %w", err)
	}
	fmt.Println(okStyle.Render("  ✓ Repository cloned successfully"))
	fmt.Println()

	// 2. Scan
	absPath, err := filepath.Abs(targetDir)
	if err != nil {
		absPath = targetDir
	}
	fmt.Println(dimStyle.Render("  → Analyzing project structure and stack..."))
	s := scanner.New(absPath)
	result, err := s.Scan()
	if err != nil {
		return fmt.Errorf("project scan failed: %w", err)
	}

	if len(result.Technologies) == 0 {
		fmt.Println(warnStyle.Render("  No technologies detected in this repository."))
		return nil
	}

	var techNames []string
	for _, t := range result.Technologies {
		techNames = append(techNames, t.Name)
	}
	fmt.Printf("  Stack:   %s\n\n", highlightStyle.Render(strings.Join(techNames, ", ")))

	// Map detected packages to standard installers
	inst := installer.New(false)
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

	// Average download & installation times in seconds for calculation
	runtimeTimes := map[string]int{
		"nodejs":  30,
		"python":  30,
		"go":      40,
		"rust":    90,
		"docker":  60,
		"bun":     15,
		"pnpm":    10,
		"flutter": 120,
		"java":    45,
		"php":     40,
		"ruby":    45,
	}

	var toInstall []string
	totalEstSeconds := 0

	for _, lang := range append(result.Languages, result.Infra...) {
		if rtName, ok := langToRuntime[lang]; ok {
			status := inst.CheckStatus(rtName)
			if !status.Installed {
				toInstall = append(toInstall, rtName)
				if est, exists := runtimeTimes[rtName]; exists {
					totalEstSeconds += est
				} else {
					totalEstSeconds += 20 // default fallback
				}
			}
		}
	}

	if len(toInstall) == 0 {
		fmt.Println(okStyle.Render("  ✓ All required tools and runtimes are already installed! Ready to build."))
		return nil
	}

	// Format estimated time
	var timeStr string
	if totalEstSeconds >= 60 {
		timeStr = fmt.Sprintf("%dm %ds", totalEstSeconds/60, totalEstSeconds%60)
	} else {
		timeStr = fmt.Sprintf("%ds", totalEstSeconds)
	}

	fmt.Println(titleStyle.Render("  PENDING INSTALLATIONS"))
	for _, rt := range toInstall {
		fmt.Printf("  - %-15s [est. %ds]\n", rt, runtimeTimes[rt])
	}
	fmt.Printf("\n  Total Estimated Setup Time: %s\n\n", highlightStyle.Render(timeStr))

	if !skipConfirm {
		fmt.Print("  Proceed with installation? [y/N] ")
		var answer string
		_, _ = fmt.Scanln(&answer)
		if answer != "y" && answer != "Y" && answer != "yes" {
			fmt.Println(dimStyle.Render("  Installation cancelled."))
			return nil
		}
	}

	// Install with loader progress bar counting down the total estimated time
	stopLoader := make(chan bool)
	go func() {
		ticks := 0
		for {
			select {
			case <-stopLoader:
				return
			default:
				remSeconds := totalEstSeconds - ticks
				if remSeconds < 0 {
					remSeconds = 0
				}
				var remStr string
				if remSeconds >= 60 {
					remStr = fmt.Sprintf("%dm %ds", remSeconds/60, remSeconds%60)
				} else {
					remStr = fmt.Sprintf("%ds", remSeconds)
				}

				// Hard Neo-Brutalist inline loader bar
				width := 30
				filled := 0
				if totalEstSeconds > 0 {
					filled = (ticks * width) / totalEstSeconds
				}
				if filled > width {
					filled = width
				}
				bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)

				fmt.Printf("\r  [%s] %d%% | Est. remaining: %-8s", bar, (ticks*100)/totalEstSeconds, remStr)
				time.Sleep(1 * time.Second)
				ticks++
			}
		}
	}()

	// Execute installs
	exitCode := 0
	for _, rtName := range toInstall {
		rt, _ := installer.GetRuntime(rtName)
		if err := inst.Install(rtName); err != nil {
			fmt.Printf("\n%s\n", warnStyle.Render(fmt.Sprintf("  ✗ Failed to install %s: %v", rt.Name, err)))
			exitCode = 1
		}
	}

	close(stopLoader)
	fmt.Println("\r" + strings.Repeat(" ", 60) + "\r") // Clear loader line

	if exitCode == 0 {
		fmt.Println(okStyle.Render("  ✓ Environment configured successfully! All runtimes are ready."))
	} else {
		return fmt.Errorf("some runtimes failed to install; please review errors above")
	}

	return nil
}
