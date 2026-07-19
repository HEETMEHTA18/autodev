package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func newPonytailCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ponytail",
		Short: "Make your project overpowered with lazy senior dev rules",
		Long: `Instantly configures your current project with Ponytail rules.
This writes specific instructions to .cursorrules, .clinerules, and 
.github/copilot-instructions.md so any AI agent you use will think like the 
laziest senior dev in the room. The best code is the code you never wrote.`,
		Example: `  autodev ponytail`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPonytail()
		},
	}
	return cmd
}

func runPonytail() error {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700"))
	successStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FF87"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))

	fmt.Println()
	fmt.Println(titleStyle.Render("  👱‍♂️  AutoDev Ponytail: Lazy Senior Dev Mode"))
	fmt.Println(dimStyle.Render("  The best code is the code you never wrote."))
	fmt.Println()

	rulesContent := `# ⚡ PONYTAIL: LAZY SENIOR DEV RULES

The best code is the code you never wrote. Follow these steps for any task:
1. Does this need to exist? → no: skip it (YAGNI)
2. Stdlib does it? → use it
3. Native platform feature? → use it
4. Installed dependency? → use it
5. One line? → one line
6. Only then: the minimum that works

Do not write unnecessary boilerplate. Do not over-engineer.
`

	filesToSave := []string{
		".cursorrules",
		".clinerules",
		".windsurfrules",
		filepath.Join(".github", "copilot-instructions.md"),
	}

	for _, relPath := range filesToSave {
		dir := filepath.Dir(relPath)
		if dir != "." && dir != "" {
			_ = os.MkdirAll(dir, 0755)
		}
		
		// If the file exists, append instead of overwrite to preserve existing rules
		f, err := os.OpenFile(relPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to write %s: %w", relPath, err)
		}
		
		// Add a newline just in case existing content didn't end with one
		_, err = f.WriteString("\n\n" + rulesContent)
		if err != nil {
			_ = f.Close()
			return fmt.Errorf("failed to append to %s: %w", relPath, err)
		}

		if err := f.Close(); err != nil {
			return fmt.Errorf("failed to close %s after write: %w", relPath, err)
		}
		fmt.Printf("  %s Injected overpowered rules into %s\n", successStyle.Render("✓"), relPath)
	}

	fmt.Println()
	fmt.Println(successStyle.Render("  Your local workspace is now overpowered. AI assistants will write less code!"))
	fmt.Println()

	return nil
}
