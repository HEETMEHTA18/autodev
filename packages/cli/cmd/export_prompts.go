package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/autodev-sh/autodev/core/promptcapture"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func newExportPromptsCmd() *cobra.Command {
	var outFile string
	var format string

	cmd := &cobra.Command{
		Use:     "export-prompts",
		Aliases: []string{"export"},
		Short:   "Export captured prompts to a file",
		Long:  `Export all captured prompts across all sessions into a single Markdown or JSON file.`,
		Example: `  autodev export-prompts -o prompts_backup.md -f markdown
  autodev export-prompts -o prompts.json -f json
  autodev prompts export -o prompts.json -f json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runExportPrompts(outFile, format)
		},
	}

	cmd.Flags().StringVarP(&outFile, "output", "o", "exported_prompts.md", "output file path")
	cmd.Flags().StringVarP(&format, "format", "f", "markdown", "output format (markdown or json)")
	return cmd
}

func runExportPrompts(outFile string, format string) error {
	accentStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00c853"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))

	engine, err := promptcapture.NewEngine("")
	if err != nil {
		return err
	}

	promptsLogPath := filepath.Join(engine.Root, ".autodevs", "prompts", "prompts.json")
	if _, err := os.Stat(promptsLogPath); os.IsNotExist(err) {
		fmt.Println()
		fmt.Println(dimStyle.Render("  No captured prompts found to export."))
		fmt.Println()
		return nil
	}

	data, err := os.ReadFile(promptsLogPath)
	if err != nil {
		return fmt.Errorf("failed to read prompts log: %w", err)
	}

	format = strings.ToLower(strings.TrimSpace(format))

	if format == "json" {
		err = os.WriteFile(outFile, data, 0644)
		if err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
	} else if format == "markdown" || format == "md" {
		var prompts []promptcapture.GlobalPrompt
		if err := json.Unmarshal(data, &prompts); err != nil {
			return fmt.Errorf("failed to parse prompts: %w", err)
		}

		var sb strings.Builder
		sb.WriteString("# AutoDev Exported Prompts\n\n")
		sb.WriteString(fmt.Sprintf("*Exported on: %s*\n\n", time.Now().Format("2006-01-02 15:04:05")))
		sb.WriteString("| Timestamp | Session ID | Prompt |\n")
		sb.WriteString("| --- | --- | --- |\n")

		for _, p := range prompts {
			cleanPrompt := strings.ReplaceAll(p.Prompt, "\n", " ")
			sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n",
				p.Timestamp.Format("2006-01-02 15:04:05"),
				p.SessionID,
				cleanPrompt,
			))
		}

		err = os.WriteFile(outFile, []byte(sb.String()), 0644)
		if err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
	} else {
		return fmt.Errorf("invalid format: must be 'markdown' or 'json'")
	}

	fmt.Println()
	fmt.Printf("  %s Prompts successfully exported to %s (%s format)\n",
		accentStyle.Render("✓"),
		outFile,
		format,
	)
	fmt.Println()
	return nil
}
