package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/autodev-sh/autodev/core/promptcapture"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func newPromptsCmd() *cobra.Command {
	var today bool

	cmd := &cobra.Command{
		Use:     "prompts",
		Aliases: []string{"history"},
		Short:   "Manage and view captured prompt history",
		Long: `View the history of captured prompts. Running this command opens the master prompts.md file
in a terminal pager, or prints a summary of today's active session if the --today flag is set.

Available Subcommands:
  chat           Start an interactive prompt session to capture and track prompts
  capture        Capture prompts from a terminal-based AI assistant
  daemon         Start background AI session prompt capture daemon
  replay         Replay a prompt from the latest session
  export         Export captured prompts to a file
  sync           Sync queued offline prompt events to DevMentor API`,
		Example: `  autodev prompts
  autodev prompts --today
  autodev prompts chat
  autodev prompts replay`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runHistory(today)
		},
	}

	cmd.Flags().BoolVar(&today, "today", false, "display only today's session prompts")

	// Add subcommands
	cmd.AddCommand(
		newChatCmd(),
		newCaptureCmd(),
		newDaemonCmd(),
		newReplayCmd(),
		newExportPromptsCmd(),
		newSyncCmd(),
	)

	return cmd
}

func runHistory(today bool) error {
	engine, err := promptcapture.NewEngine("")
	if err != nil {
		return err
	}

	// Import prompts from active Antigravity/agy editor session
	_ = promptcapture.ImportAntigravityPrompts(engine.Root)

	if today {
		return showTodayPrompts(engine)
	}

	rootMDPath := filepath.Join(engine.Root, ".autodevs", "prompts.md")
	if _, err := os.Stat(rootMDPath); os.IsNotExist(err) {
		rootMDPath = filepath.Join(engine.Root, "prompts.md")
	}

	if _, err := os.Stat(rootMDPath); os.IsNotExist(err) {
		fmt.Println()
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Render("  No prompt capture history found. Start a session with 'autodev chat' or use 'autodev capture'."))
		fmt.Println()
		return nil
	}

	// Try using less pager, fallback to standard output
	cmd := exec.Command("less", "-R", rootMDPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		data, err := os.ReadFile(rootMDPath)
		if err != nil {
			return err
		}
		fmt.Print(string(data))
	}
	return nil
}

func showTodayPrompts(engine *promptcapture.Engine) error {
	sessionsDir := filepath.Join(engine.Root, ".autodevs", "sessions")
	files, err := os.ReadDir(sessionsDir)
	if err != nil {
		return fmt.Errorf("failed to read sessions: %w", err)
	}

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#0070f3"))
	accentStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00c853"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))

	todayStr := time.Now().Format("2006-01-02")
	fmt.Println()
	fmt.Println(titleStyle.Render(fmt.Sprintf("  Captured Prompts for Today (%s)", todayStr)))
	fmt.Println()

	found := false
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".json") && strings.HasPrefix(f.Name(), todayStr) {
			data, err := os.ReadFile(filepath.Join(sessionsDir, f.Name()))
			if err != nil {
				continue
			}

			var session promptcapture.SessionLog
			if err := json.Unmarshal(data, &session); err != nil {
				continue
			}

			if len(session.Events) == 0 {
				continue
			}

			found = true
			fmt.Printf("  %s %s (%s)\n",
				accentStyle.Render("●"),
				session.SessionID,
				dimStyle.Render(strings.Join(session.Metadata.Languages, ", ")),
			)

			for _, ev := range session.Events {
				fmt.Printf("    %s %s\n",
					dimStyle.Render(ev.Timestamp.Format("15:04:05")),
					ev.Prompt,
				)
			}
			fmt.Println()
		}
	}

	if !found {
		fmt.Println(dimStyle.Render("  No prompts captured today. Run 'autodev chat' to capture some!"))
		fmt.Println()
	}

	return nil
}
