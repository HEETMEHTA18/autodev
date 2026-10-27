package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func newMigrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate legacy JSON configuration files to the new YAML schema",
		Long: `Scan the workspace and ~/.config/autodev for legacy .autodev.json files 
and convert them into the new unified .autodev.yaml configuration format.`,
		Example: `  autodev migrate`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMigrate()
		},
	}
	return cmd
}

func runMigrate() error {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#0070f3"))
	okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00c853")).Bold(true)
	warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ee0000"))

	fmt.Println()
	fmt.Println(titleStyle.Render("AutoDev Profile Config Migrator"))

	home, err := os.UserHomeDir()
	var pathsToMigrate []string
	if err == nil {
		pathsToMigrate = append(pathsToMigrate, filepath.Join(home, ".config", "autodev", ".autodev.json"))
		pathsToMigrate = append(pathsToMigrate, filepath.Join(home, ".config", "autodev", "config.json"))
	}
	pathsToMigrate = append(pathsToMigrate, ".autodev.json")

	migratedCount := 0

	for _, p := range pathsToMigrate {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			continue
		}

		fmt.Printf("  Found legacy config: %s\n", p)
		data, err := os.ReadFile(p)
		if err != nil {
			fmt.Printf("    %s Failed to read legacy file: %v\n", warnStyle.Render("✗"), err)
			continue
		}

		var parsed map[string]interface{}
		if err := json.Unmarshal(data, &parsed); err != nil {
			fmt.Printf("    %s Legacy file is not valid JSON: %v\n", warnStyle.Render("✗"), err)
			continue
		}

		// Convert to YAML
		yamlData, err := yaml.Marshal(parsed)
		if err != nil {
			fmt.Printf("    %s Failed to convert config to YAML: %v\n", warnStyle.Render("✗"), err)
			continue
		}

		ext := filepath.Ext(p)
		dir := filepath.Dir(p)
		baseWithoutExt := filepath.Base(p)
		baseWithoutExt = baseWithoutExt[:len(baseWithoutExt)-len(ext)]

		newPath := filepath.Join(dir, baseWithoutExt+".yaml")

		err = os.WriteFile(newPath, yamlData, 0644)
		if err != nil {
			fmt.Printf("    %s Failed to write new YAML config: %v\n", warnStyle.Render("✗"), err)
			continue
		}
		fmt.Printf("    %s Created YAML config: %s\n", okStyle.Render("✓"), newPath)

		// Backup the old JSON file
		backupPath := p + ".bak"
		err = os.Rename(p, backupPath)
		if err != nil {
			fmt.Printf("    %s Failed to backup legacy file: %v\n", warnStyle.Render("✗"), err)
		} else {
			fmt.Printf("    %s Legay file backed up to: %s\n", okStyle.Render("✓"), backupPath)
		}
		migratedCount++
	}

	fmt.Println()
	if migratedCount == 0 {
		fmt.Println("  No legacy JSON configuration files found to migrate.")
	} else {
		fmt.Printf(okStyle.Render("  ✓ Successfully migrated %d legacy configuration files!\n"), migratedCount)
	}
	fmt.Println()
	PrintGitHubCTA()
	return nil
}
