package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/autodev-sh/autodev/core/promptcapture"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func newReplayCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "replay",
		Short:   "Replay a prompt from the latest session",
		Long:    `Display a list of captured prompts from the latest session and select one to rerun against the codebase.`,
		Example: `  autodev replay`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runReplay()
		},
	}
	return cmd
}

func runReplay() error {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#0070f3"))
	accentStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00c853"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	warnStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#f5a623"))
	cmdStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#0070f3"))

	engine, err := promptcapture.NewEngine("")
	if err != nil {
		return err
	}

	err = engine.LoadLatestSession()
	if err != nil {
		return fmt.Errorf("failed to load latest session: %w", err)
	}

	if len(engine.Session.Events) == 0 {
		fmt.Println()
		fmt.Println(dimStyle.Render("  No prompts found in the latest session to replay."))
		fmt.Println()
		return nil
	}

	fmt.Println()
	fmt.Println(titleStyle.Render("  Replay Prompts from Session: " + engine.Session.SessionID))
	fmt.Println(dimStyle.Render("  Select a prompt below to rerun it against the codebase:"))
	fmt.Println()

	for i, ev := range engine.Session.Events {
		fmt.Printf("    %d. %s\n", i+1, ev.Prompt)
	}
	fmt.Println()

	fmt.Print(accentStyle.Render("  Select prompt to replay [1-" + fmt.Sprintf("%d", len(engine.Session.Events)) + "]: "))
	var idx int
	_, err = fmt.Scanln(&idx)
	if err != nil || idx < 1 || idx > len(engine.Session.Events) {
		return fmt.Errorf("invalid selection")
	}

	selectedPrompt := engine.Session.Events[idx-1].Prompt
	fmt.Println()
	fmt.Printf("  Replaying: %s\n", accentStyle.Render(selectedPrompt))
	fmt.Println()

	apiKey := os.Getenv("GEMINI_API_KEY")
	var responseText string
	if apiKey != "" {
		// Same two-phase discipline as `autodev chat`: refine, then answer.
		responseText, err = askGemini(selectedPrompt, engine.Session.Events, chatOptions{thinkingBudget: 1024})
		if err != nil {
			fmt.Println(warnStyle.Render(fmt.Sprintf("  Error querying Gemini: %v. Falling back to simulation.", err)))
			responseText = getSimulatedResponse(selectedPrompt)
		}
	} else {
		responseText = getSimulatedResponse(selectedPrompt)
	}

	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#ededed")).Render(responseText))
	fmt.Println()

	suggestedCommands, suggestedFiles := parseActions(responseText)

	var runCmds []promptcapture.ExecutedCommand
	var genFiles []promptcapture.GeneratedFile

	for _, f := range suggestedFiles {
		fmt.Printf("  Create/Write file %s? [y/N] ", cmdStyle.Render(f.Path))
		if askConfirm() {
			dir := filepath.Dir(f.Path)
			if dir != "." {
				_ = os.MkdirAll(filepath.Join(engine.Root, dir), 0755)
			}
			err := os.WriteFile(filepath.Join(engine.Root, f.Path), []byte(f.Content), 0644)
			if err != nil {
				fmt.Printf("    %s Failed to write file: %v\n", warnStyle.Render("✗"), err)
			} else {
				fmt.Printf("    %s Saved %s (%d bytes)\n", accentStyle.Render("✓"), f.Path, len(f.Content))
				genFiles = append(genFiles, promptcapture.GeneratedFile{
					FilePath:  f.Path,
					SizeBytes: int64(len(f.Content)),
					Action:    "created",
					Timestamp: time.Now(),
				})
			}
		}
	}

	for _, c := range suggestedCommands {
		fmt.Printf("  Run command: %s? [y/N] ", cmdStyle.Render(c))
		if askConfirm() {
			fmt.Printf("    Running `%s`...\n", c)
			execCmd, err := runCommandAndCapture(engine.Root, c)
			if err != nil {
				fmt.Printf("    %s Run failed: %v\n", warnStyle.Render("✗"), err)
			} else {
				fmt.Printf("    %s Finished (Exit Code: %d, Duration: %.2fs)\n",
					accentStyle.Render("✓"), execCmd.ExitCode, float64(execCmd.DurationMs)/1000.0)
				runCmds = append(runCmds, execCmd)
			}
		}
	}

	// Add event
	err = engine.AddEvent(selectedPrompt, responseText, runCmds, genFiles)
	if err != nil {
		fmt.Println(warnStyle.Render(fmt.Sprintf("  Warning: failed to log replayed prompt event: %v", err)))
	} else {
		fmt.Println(dimStyle.Render("  [captured replayed prompt, files, and commands]"))
	}
	fmt.Println()

	return nil
}
