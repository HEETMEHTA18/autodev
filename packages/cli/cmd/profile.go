package cmd

import (
	"fmt"
	"strings"

	"github.com/autodev-sh/autodev/catalog"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func newProfileCmd() *cobra.Command {
	var list bool

	cmd := &cobra.Command{
		Use:   "profile [profile-id]",
		Short: "Install a pre-defined developer profile (role-based tool set)",
		Long: `Install all tools for a specific developer role in one command.
Available profiles: web-dev, ml-engineer, flutter-dev, devops-engineer, 
backend-dev, android-dev, data-scientist, rust-dev, fullstack-ai`,
		Example: `  autodev profile web-dev
  autodev profile ml-engineer
  autodev profile flutter-dev --dry-run
  autodev profile --list`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := catalog.Load()
			if err != nil {
				return err
			}
			if list || len(args) == 0 {
				return printProfiles(c)
			}
			return runProfile(c, args[0])
		},
	}

	cmd.Flags().BoolVarP(&list, "list", "l", false, "list all available profiles")
	return cmd
}

func printProfiles(c *catalog.Catalog) error {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#0070f3"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	boldStyle := lipgloss.NewStyle().Bold(true)

	fmt.Println()
	fmt.Println(titleStyle.Render("  Developer Profiles"))
	fmt.Println()

	for _, prof := range c.Profiles {
		fmt.Printf("  %-20s  %s\n",
			boldStyle.Render(prof.ID),
			prof.Name,
		)
		fmt.Printf("     %s\n", dimStyle.Render(prof.Description))
		fmt.Printf("     %s\n\n", dimStyle.Render(strings.Join(prof.Packages, " · ")))
	}

	fmt.Println(dimStyle.Render("  Usage: autodev profile <id>"))
	fmt.Println()
	return nil
}

func runProfile(c *catalog.Catalog, profileID string) error {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#0070f3"))
	okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00c853")).Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ee0000"))

	prof, ok := c.GetProfile(profileID)
	if !ok {
		return fmt.Errorf("unknown profile %q — run 'autodev profile --list' to see available profiles", profileID)
	}

	fmt.Println()
	fmt.Printf("  %s %s\n", titleStyle.Render(prof.Name), dimStyle.Render("profile"))
	fmt.Println(dimStyle.Render("  " + prof.Description))
	fmt.Println()

	// Resolve with deps
	resolved, err := c.ResolveProfile(profileID)
	if err != nil {
		return err
	}

	fmt.Println(titleStyle.Render(fmt.Sprintf("  %d packages to install:", len(resolved))))
	for _, pkg := range resolved {
		fmt.Printf("  %-20s  %s\n", pkg.Name, dimStyle.Render(pkg.Description))
	}
	fmt.Println()

	if dryRun {
		fmt.Println(warnStyle.Render("  ! dry-run: no changes made"))
		return nil
	}

	fmt.Print("  Proceed? [y/N] ")
	var ans string
	_, _ = fmt.Scanln(&ans)
	if strings.ToLower(ans) != "y" && strings.ToLower(ans) != "yes" {
		fmt.Println(dimStyle.Render("  Cancelled."))
		return nil
	}

	for _, pkg := range resolved {
		fmt.Printf(titleStyle.Render("\n  Installing %s...\n"), pkg.Name)
		if pkg.IsInstalled() {
			fmt.Println(okStyle.Render("    ✓ Already installed"))
			continue
		}
		if err := installCatalogPackage(pkg); err != nil {
			fmt.Println(warnStyle.Render(fmt.Sprintf("  ✗ %s: %v", pkg.Name, err)))
		} else {
			fmt.Println(okStyle.Render(fmt.Sprintf("  ✓ %s installed", pkg.Name)))
		}
	}

	fmt.Println()
	fmt.Println(okStyle.Render("  ✓ Profile installation complete!"))
	fmt.Println(dimStyle.Render("  Run 'autodev doctor' to verify everything is working."))
	fmt.Println()
	return nil
}

// installCatalogPackage executes the platform install for a catalog package.
func installCatalogPackage(pkg *catalog.Package) error {
	return execInstall(pkg)
}

func runProfileNoConfirm(c *catalog.Catalog, profileID string) error {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#0070f3"))
	okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00c853")).Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ee0000"))

	prof, ok := c.GetProfile(profileID)
	if !ok {
		return fmt.Errorf("unknown profile %q", profileID)
	}

	fmt.Println()
	fmt.Printf("  %s %s\n", titleStyle.Render(prof.Name), dimStyle.Render("profile"))
	fmt.Println(dimStyle.Render("  " + prof.Description))
	fmt.Println()

	resolved, err := c.ResolveProfile(profileID)
	if err != nil {
		return err
	}

	fmt.Println(titleStyle.Render(fmt.Sprintf("  %d packages to install:", len(resolved))))
	for _, pkg := range resolved {
		fmt.Printf("  %-20s  %s\n", pkg.Name, dimStyle.Render(pkg.Description))
	}
	fmt.Println()

	for _, pkg := range resolved {
		fmt.Printf(titleStyle.Render("\n  Installing %s...\n"), pkg.Name)
		if pkg.IsInstalled() {
			fmt.Println(okStyle.Render("    ✓ Already installed"))
			continue
		}
		if err := installCatalogPackage(pkg); err != nil {
			fmt.Println(warnStyle.Render(fmt.Sprintf("  ✗ %s: %v", pkg.Name, err)))
		} else {
			fmt.Println(okStyle.Render(fmt.Sprintf("  ✓ %s installed", pkg.Name)))
		}
	}

	fmt.Println()
	fmt.Println(okStyle.Render("  ✓ Profile installation complete!"))
	fmt.Println()
	return nil
}
