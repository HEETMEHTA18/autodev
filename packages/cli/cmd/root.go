// Package cmd defines all AutoDev CLI commands using Cobra.
package cmd

import (
	"fmt"
	"os"

	"github.com/autodev-sh/autodev/catalog"
	"github.com/autodev-sh/autodev/cli/tui"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	verbose bool
	noColor bool
	dryRun  bool
	jsonOut bool
)

// rootCmd — running `autodev` with no args opens the interactive TUI.
var rootCmd = &cobra.Command{
	Use:     "autodev",
	Short:   "Set up any development environment in one command.",
	Version: "0.4.3",
	Long: `
  █████╗ ██╗   ██╗████████╗ ██████╗ ██████╗ ███████╗██╗   ██╗
  ██╔══██╗██║   ██║╚══██╔══╝██╔═══██╗██╔══██╗██╔════╝██║   ██║
  ███████║██║   ██║   ██║   ██║   ██║██║  ██║█████╗  ██║   ██║
  ██╔══██║██║   ██║   ██║   ██║   ██║██║  ██║██╔══╝  ╚██╗ ██╔╝
  ██║  ██║╚██████╔╝   ██║   ╚██████╔╝██████╔╝███████╗ ╚████╔╝ 
  ╚═╝  ╚═╝ ╚═════╝    ╚═╝    ╚══════╝ ╚═════╝ ╚══════╝  ╚═══╝ 

  The App Store for Developers.
  Run with no arguments to open the interactive installer.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Avoid running for MCP communication or version/help checks to guarantee performance
		if cmd.Name() == "start" && cmd.Parent() != nil && cmd.Parent().Name() == "mcp" {
			return
		}
		if cmd.Name() == "help" || (cmd.Name() == "autodev" && len(args) == 0) {
			return
		}
		// Run silently
		AutoGenerateRulesSilent(".")
	},
	// When called with no subcommand → open the TUI
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := catalog.Load()
		if err != nil {
			return fmt.Errorf("failed to load catalog: %w", err)
		}
		return tui.Run(c)
	},
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: .autodev.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable color output")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "preview actions without executing")
	rootCmd.PersistentFlags().BoolVar(&jsonOut, "json", false, "output results as JSON")

	_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	_ = viper.BindPFlag("no_color", rootCmd.PersistentFlags().Lookup("no-color"))

	rootCmd.AddCommand(
		newScanCmd(),
		newSetupCmd(),
		newGitHubCmd(),
		newDoctorCmd(),
		newReportCmd(),
		newInstallCmd(),
		newUpdateCmd(),
		newCleanCmd(),
		newSkillsCmd(),
		newExportCmd(),
		newProfileCmd(),
		newUICmd(),
		newCloneCmd(),
		newAuditCmd(),
		newMCPCmd(),
		newCreateCmd(),
		newBenchmarkCmd(),
		newContainerizeCmd(),
		newMigrateCmd(),
		newCICmd(),
		newPluginCmd(),
		newCanvasCmd(),
		newPonytailCmd(),
		newSecurityCmd(),
		newSecretsCmd(),
		newHardenCmd(),
		newUpgradeCmd(),
	)

	// Prompts manager command (displays history & groups prompt capture subcommands)
	promptsCmd := newPromptsCmd()
	rootCmd.AddCommand(promptsCmd)

	// Keep top-level commands for compatibility but hide them from main help list
	chatCmd := newChatCmd()
	chatCmd.Hidden = true
	rootCmd.AddCommand(chatCmd)

	captureCmd := newCaptureCmd()
	captureCmd.Hidden = true
	rootCmd.AddCommand(captureCmd)

	daemonCmd := newDaemonCmd()
	daemonCmd.Hidden = true
	rootCmd.AddCommand(daemonCmd)

	replayCmd := newReplayCmd()
	replayCmd.Hidden = true
	rootCmd.AddCommand(replayCmd)

	exportPromptsCmd := newExportPromptsCmd()
	exportPromptsCmd.Hidden = true
	rootCmd.AddCommand(exportPromptsCmd)

	syncCmd := newSyncCmd()
	syncCmd.Hidden = true
	rootCmd.AddCommand(syncCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName(".autodev")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		home, _ := os.UserHomeDir()
		viper.AddConfigPath(home + "/.config/autodev")
	}
	viper.SetEnvPrefix("AUTODEV")
	viper.AutomaticEnv()
	_ = viper.ReadInConfig()
}

// PrintGitHubCTA prints a friendly CTA requesting users to star the GitHub repo.
func PrintGitHubCTA() {
	if jsonOut {
		return
	}
	starStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Bold(true)
	linkStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF87")).Underline(true)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#555555"))

	fmt.Println()
	fmt.Println(dimStyle.Render("  ──────────────────────────────────────────────────────────"))
	fmt.Printf("  %s Star the repo to support AutoDev: %s\n",
		starStyle.Render("⭐ Love this tool?"),
		linkStyle.Render("https://github.com/HEETMEHTA18/autodev"))
	fmt.Println(dimStyle.Render("  ──────────────────────────────────────────────────────────"))
	fmt.Println()
}
