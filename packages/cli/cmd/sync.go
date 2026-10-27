package cmd

import (
	"fmt"

	"github.com/autodev-sh/autodev/core/promptcapture"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func newSyncCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync queued offline prompt events to DevMentor API",
		Long:  `Scan the offline queue under .autodevs/analytics/queue.json and sync events to DevMentor API.`,
		Example: `  autodev sync`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSync()
		},
	}
	return cmd
}

func runSync() error {
	accentStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00c853"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	warnStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#f5a623"))

	engine, err := promptcapture.NewEngine("")
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println(dimStyle.Render("  Connecting to DevMentor API..."))

	count, err := engine.SyncOfflineEvents()
	if err != nil {
		fmt.Println(warnStyle.Render(fmt.Sprintf("  Offline sync failed: %v", err)))
		fmt.Println()
		return err
	}

	if count == 0 {
		fmt.Println(accentStyle.Render("  ✓ All events are already synchronized. No offline events queued."))
	} else {
		fmt.Printf("  %s Successfully synchronized %d events to DevMentor telemetry API.\n",
			accentStyle.Render("✓"),
			count,
		)
	}
	fmt.Println()
	return nil
}
