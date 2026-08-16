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

const (
	groupDiagnose  = "diagnose"
	groupSetup     = "setup"
	groupBuild     = "build"
	groupAI        = "ai"
	groupIntegrate = "integrations"
	groupSecurity  = "security"
	groupUtilities = "utilities"
)

var rootCmd = &cobra.Command{
	Use:     "autodev",
	Short:   "Understand, set up and operate your development environment.",
	Version: "0.6.0",
	Long: `AutoDev is a developer environment control center.

Use the interactive command center when you are new to AutoDev, or use the
focused commands below when you already know what you want to do.`,
	Example: `  # Recommended first steps
  autodev
  autodev doctor
  autodev scan
  autodev setup
  autodev agent

  # Preview changes before executing them
  autodev setup --dry-run

  # Machine-readable automation
  autodev doctor --json`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Name() == "start" && cmd.Parent() != nil && cmd.Parent().Name() == "mcp" {
			return
		}
		if cmd.Name() == "help" || (cmd.Name() == "autodev" && len(args) == 0) {
			return
		}
		AutoGenerateRulesSilent(".")
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := catalog.Load()
		if err != nil {
			return fmt.Errorf("failed to load catalog: %w", err)
		}
		return tui.RunProfessional(c)
	},
}

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

	rootCmd.AddGroup(
		&cobra.Group{ID: groupDiagnose, Title: "Diagnose:"},
		&cobra.Group{ID: groupSetup, Title: "Setup & Environment:"},
		&cobra.Group{ID: groupBuild, Title: "Build & Delivery:"},
		&cobra.Group{ID: groupAI, Title: "AI & Agents:"},
		&cobra.Group{ID: groupIntegrate, Title: "Integrations:"},
		&cobra.Group{ID: groupSecurity, Title: "Security:"},
		&cobra.Group{ID: groupUtilities, Title: "Utilities:"},
	)

	commands := []*cobra.Command{
		newScanCmd(), newSetupCmd(), newGitHubCmd(), newDoctorCmd(), newReportCmd(),
		newInstallCmd(), newUpdateCmd(), newCleanCmd(), newSkillsCmd(), newExportCmd(),
		newProfileCmd(), newUICmd(), newCloneCmd(), newAuditCmd(), newMCPCmd(), newCreateCmd(),
		newBenchmarkCmd(), newContainerizeCmd(), newMigrateCmd(), newCICmd(), newPluginCmd(),
		newCanvasCmd(), newPonytailCmd(), newSecurityCmd(), newSecretsCmd(), newHardenCmd(),
		newUpgradeCmd(), newToolsCmd(), newRunCmd(), newSessionCmd(), newAgentCmd(), newOpenCodeCmd(),
	}
	for _, command := range commands {
		command.GroupID = commandGroup(command.Name())
		rootCmd.AddCommand(command)
	}

	promptsCmd := newPromptsCmd()
	promptsCmd.GroupID = groupUtilities
	rootCmd.AddCommand(promptsCmd)

	for _, command := range []*cobra.Command{newChatCmd(), newCaptureCmd(), newDaemonCmd(), newReplayCmd(), newExportPromptsCmd(), newSyncCmd()} {
		command.Hidden = true
		rootCmd.AddCommand(command)
	}
}

func commandGroup(name string) string {
	switch name {
	case "doctor", "scan", "report", "benchmark":
		return groupDiagnose
	case "setup", "install", "update", "clean", "profile", "upgrade", "tools", "clone", "migrate":
		return groupSetup
	case "create", "run", "containerize", "ci", "build":
		return groupBuild
	case "agent", "chat", "session", "skills", "opencode":
		return groupAI
	case "github", "mcp", "plugin", "canvas", "ui":
		return groupIntegrate
	case "audit", "security", "secrets", "harden":
		return groupSecurity
	default:
		return groupUtilities
	}
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

func PrintGitHubCTA() {
	if jsonOut {
		return
	}
	starStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Bold(true)
	linkStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF87")).Underline(true)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#555555"))
	fmt.Println()
	fmt.Println(dimStyle.Render("  ──────────────────────────────────────────────────────────"))
	fmt.Printf("  %s Star the repo to support AutoDev: %s\n", starStyle.Render("⭐ Love this tool?"), linkStyle.Render("https://github.com/HEETMEHTA18/autodev"))
	fmt.Println(dimStyle.Render("  ──────────────────────────────────────────────────────────"))
	fmt.Println()
}
