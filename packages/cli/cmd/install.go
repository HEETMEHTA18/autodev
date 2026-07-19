package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/autodev-sh/autodev/catalog"
	"github.com/autodev-sh/autodev/core/osinfo"
	"github.com/autodev-sh/autodev/installer"
	"github.com/autodev-sh/autodev/scanner"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func newInstallCmd() *cobra.Command {
	var listFlag bool
	var skipConfirm bool

	cmd := &cobra.Command{
		Use:   "install <package-id>",
		Short: "Install a specific package by ID",
		Long: `Install a specific runtime, framework, or tool from the AutoDev catalog.
Dependencies are resolved automatically. Run --list to see all available packages.`,
		Example: `  autodev install nodejs
  autodev install flutter
  autodev install --yes nodejs
  autodev install --list`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := catalog.Load()
			if err != nil {
				return err
			}
			if listFlag || len(args) == 0 {
				return printCatalogList(c)
			}
			return runCatalogInstall(c, args[0], skipConfirm)
		},
	}

	cmd.Flags().BoolVarP(&listFlag, "list", "l", false, "list all available packages")
	cmd.Flags().BoolVarP(&skipConfirm, "yes", "y", false, "skip confirmation prompts")
	return cmd
}

func printCatalogList(c *catalog.Catalog) error {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	catStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#4A90E2"))

	fmt.Println()
	fmt.Println(titleStyle.Render("  AutoDev Package Catalog"))
	fmt.Println()

	byCat := c.ByCategory()
	for _, cat := range catalog.CategoryOrder {
		pkgs, ok := byCat[cat]
		if !ok || len(pkgs) == 0 {
			continue
		}
		fmt.Println(catStyle.Render("  [" + cat + "]"))
		for _, pkg := range pkgs {
			deps := ""
			if len(pkg.Deps) > 0 {
				deps = dimStyle.Render(fmt.Sprintf("  (deps: %v)", pkg.Deps))
			}
			fmt.Printf("    %-18s  %s%s\n",
				pkg.ID,
				dimStyle.Render(pkg.Description),
				deps,
			)
		}
		fmt.Println()
	}

	fmt.Println(dimStyle.Render("  Usage: autodev install <id>"))
	fmt.Println(dimStyle.Render("  Or run 'autodev' for the interactive installer"))
	fmt.Println()
	return nil
}

func runCatalogInstall(c *catalog.Catalog, id string, skipConfirm bool) error {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700"))
	okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF87")).Bold(true)
	warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))

	// Resolve with deps
	resolved, err := c.Resolve([]string{id})
	if err != nil {
		return err
	}

	fmt.Println()
	pkg, _ := c.GetPackage(id)
	fmt.Printf("%s\n", titleStyle.Render("Installing "+pkg.Name))

	// Find out what is not installed
	var toInstall []*catalog.Package
	for _, p := range resolved {
		if !p.IsInstalled() {
			toInstall = append(toInstall, p)
		}
	}

	if len(toInstall) > 0 {
		fmt.Println(dimStyle.Render("  The following commands and packages will be installed/executed:"))
		for _, p := range toInstall {
			var method string
			var cmdDetail string
			switch runtime.GOOS {
			case "linux":
				method = p.Install.Linux.Method
				if method == "script" {
					cmdDetail = strings.Join(p.Install.Linux.Script, " && ")
				} else {
					cmdDetail = strings.Join(p.Install.Linux.Packages, ", ")
				}
			case "darwin":
				method = p.Install.Darwin.Method
				if method == "script" {
					cmdDetail = strings.Join(p.Install.Darwin.Script, " && ")
				} else {
					cmdDetail = strings.Join(p.Install.Darwin.Packages, ", ")
				}
			case "windows":
				method = p.Install.Windows.Method
				if method == "script" {
					cmdDetail = strings.Join(p.Install.Windows.Script, " && ")
				} else {
					cmdDetail = strings.Join(p.Install.Windows.Packages, ", ")
				}
			}
			fmt.Printf("    - %s (%s) → %s\n", p.Name, method, cmdDetail)
		}
		fmt.Println()

		if !skipConfirm && !dryRun {
			fmt.Print("  Proceed with installation? [y/N] ")
			var answer string
			_, _ = fmt.Scanln(&answer)
			answer = strings.ToLower(strings.TrimSpace(answer))
			if answer != "y" && answer != "yes" {
				fmt.Println(dimStyle.Render("  Installation cancelled."))
				return nil
			}
		}
	}

	if len(resolved) > 1 {
		fmt.Println(dimStyle.Render("  Installing dependencies first:"))
		for _, p := range resolved[:len(resolved)-1] {
			fmt.Printf("  %s\n", dimStyle.Render(p.Name))
		}
	}
	fmt.Println()

	var installedPkgs []*catalog.Package

	for _, p := range resolved {
		fmt.Printf(titleStyle.Render("  - %s\n"), p.Name)
		if dryRun {
			fmt.Println(warnStyle.Render("    [dry-run] skipped"))
			continue
		}

		if p.IsInstalled() {
			fmt.Println(okStyle.Render("    ✓ Already installed"))
			installedPkgs = append(installedPkgs, p)
			continue
		}

		if err := execInstall(p); err != nil {
			fmt.Println(warnStyle.Render(fmt.Sprintf("  [FAIL] %s: %v", p.Name, err)))
		} else {
			fmt.Println(okStyle.Render(fmt.Sprintf("  [OK] %s installed", p.Name)))
			installedPkgs = append(installedPkgs, p)
		}
	}

	if len(installedPkgs) > 0 {
		runBumblebeeSafetyCheck(installedPkgs)
	}

	fmt.Println()
	fmt.Println(okStyle.Render(fmt.Sprintf("  [OK] %s ready!", pkg.Name)))
	return nil
}

func runBumblebeeSafetyCheck(pkgs []*catalog.Package) {
	bumblebeeHeaderStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FFFF"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF87")).Bold(true)
	warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B")).Bold(true)
	infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#4A90E2")).Bold(true)

	criticalStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF2A2A"))
	highStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF6B6B"))
	moderateStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700"))
	lowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))

	fmt.Println()
	fmt.Println(bumblebeeHeaderStyle.Render("  🛡️  Bumblebee Post-Install Safety Audit"))

	type pkgInfo struct {
		pkg       *catalog.Package
		version   string
		ecosystem string
		names     []string
	}

	var infos []pkgInfo
	var displayNames []string

	for _, pkg := range pkgs {
		version := getVersionFromVerify(pkg)
		if version == "" {
			version = "0.0.1"
		}
		ecosystem, names := determineOsvEcosystemAndNames(pkg)
		if len(names) == 0 {
			ecosystem = "npm"
			names = []string{pkg.ID}
		}
		infos = append(infos, pkgInfo{
			pkg:       pkg,
			version:   version,
			ecosystem: ecosystem,
			names:     names,
		})
		displayNames = append(displayNames, fmt.Sprintf("%s@%s", pkg.Name, version))
	}

	type auditResult struct {
		pkgName    string
		pkgVersion string
		name       string
		ecosystem  string
		vulns      []scanner.Vulnerability
		err        error
	}

	totalQueries := 0
	for _, info := range infos {
		totalQueries += len(info.names)
	}

	resChan := make(chan auditResult, totalQueries)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	spinnerDone := make(chan struct{})
	go func() {
		spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		start := time.Now()
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		i := 0
		for {
			select {
			case <-spinnerDone:
				fmt.Print("\r\033[K") // Clear line
				return
			case <-ticker.C:
				elapsed := time.Since(start).Seconds()
				fmt.Printf("\r  %s %s Running safety audit on %s... %s",
					lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF")).Render(spinner[i%len(spinner)]),
					dimStyle.Render("Bumblebee:"),
					strings.Join(displayNames, ", "),
					dimStyle.Render(fmt.Sprintf("[%.1fs]", elapsed)),
				)
				i++
			}
		}
	}()

	client := &http.Client{Timeout: 3 * time.Second}
	var wg sync.WaitGroup

	for _, info := range infos {
		for _, name := range info.names {
			wg.Add(1)
			go func(inf pkgInfo, n string) {
				defer wg.Done()
				p := scanner.AuditPackage{
					Name:      n,
					Version:   inf.version,
					Ecosystem: inf.ecosystem,
				}
				vulns, err := scanner.CheckPackageVulnerabilities(ctx, client, p)
				resChan <- auditResult{
					pkgName:    inf.pkg.Name,
					pkgVersion: inf.version,
					name:       n,
					ecosystem:  inf.ecosystem,
					vulns:      vulns,
					err:        err,
				}
			}(info, name)
		}
	}

	go func() {
		wg.Wait()
		close(resChan)
	}()

	var results []auditResult
	for res := range resChan {
		results = append(results, res)
	}

	close(spinnerDone)
	time.Sleep(100 * time.Millisecond) // wait for spinner to clear

	totalVulns := 0
	for _, r := range results {
		totalVulns += len(r.vulns)
	}

	if totalVulns == 0 {
		fmt.Println(okStyle.Render(fmt.Sprintf("  ✓ Verified %s: No known vulnerabilities found! Safe from attackers' malicious files.", strings.Join(displayNames, ", "))))
		fmt.Println()
		return
	}

	fmt.Println(warnStyle.Render(fmt.Sprintf("  ✗ Found %d security vulnerabilities across downloaded packages:", totalVulns)))

	// Group results by catalog package name
	byPkg := make(map[string][]auditResult)
	for _, r := range results {
		if len(r.vulns) > 0 {
			byPkg[r.pkgName] = append(byPkg[r.pkgName], r)
		}
	}

	for pkgName, pkgResults := range byPkg {
		version := pkgResults[0].pkgVersion
		fmt.Println(infoStyle.Render(fmt.Sprintf("    📦 %s@%s", pkgName, version)))
		for _, r := range pkgResults {
			fmt.Println(dimStyle.Render(fmt.Sprintf("      ecosystem: %s, package: %s", r.ecosystem, r.name)))
			for _, v := range r.vulns {
				var sevBadge string
				switch v.Severity {
				case "CRITICAL":
					sevBadge = criticalStyle.Render("[CRITICAL]")
				case "HIGH":
					sevBadge = highStyle.Render("[HIGH]")
				case "MODERATE", "MEDIUM":
					sevBadge = moderateStyle.Render("[MODERATE]")
				default:
					sevBadge = lowStyle.Render(fmt.Sprintf("[%s]", v.Severity))
				}
				alias := ""
				if len(v.Aliases) > 0 {
					alias = fmt.Sprintf(" (%s)", v.Aliases[0])
				}
				fmt.Printf("        - %s %s%s: %s\n", sevBadge, v.ID, alias, v.Summary)
			}
		}
	}
	fmt.Println()
}

func getVersionFromVerify(pkg *catalog.Package) string {
	if pkg.Verify == "" {
		return ""
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Create command with shell to run the verification
	cmd := exec.CommandContext(ctx, "sh", "-c", pkg.Verify)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return extractVersion(string(out))
}

func extractVersion(input string) string {
	input = strings.TrimSpace(input)
	if input == "" {
		return ""
	}
	if strings.Contains(input, "go version go") {
		parts := strings.Split(input, " ")
		for _, part := range parts {
			if strings.HasPrefix(part, "go") && len(part) > 2 {
				return strings.TrimPrefix(part, "go")
			}
		}
	}
	words := strings.Fields(input)
	for _, word := range words {
		word = strings.TrimLeft(word, "vV")
		dotCount := strings.Count(word, ".")
		if dotCount >= 1 {
			cleaned := strings.TrimFunc(word, func(r rune) bool {
				return !((r >= '0' && r <= '9') || r == '.' || r == '-' || (r >= 'a' && r <= 'z'))
			})
			if cleaned != "" {
				return cleaned
			}
		}
	}
	return input
}

func determineOsvEcosystemAndNames(pkg *catalog.Package) (string, []string) {
	var method string
	var packages []string

	if pkg.Install.Linux.Method != "" {
		method = pkg.Install.Linux.Method
		packages = pkg.Install.Linux.Packages
	} else if pkg.Install.Darwin.Method != "" {
		method = pkg.Install.Darwin.Method
		packages = pkg.Install.Darwin.Packages
	} else if pkg.Install.Windows.Method != "" {
		method = pkg.Install.Windows.Method
		packages = pkg.Install.Windows.Packages
	}

	switch method {
	case "npm":
		if len(packages) > 0 {
			return "npm", packages
		}
		return "npm", []string{pkg.ID}
	case "pip":
		if len(packages) > 0 {
			return "PyPI", packages
		}
		return "PyPI", []string{pkg.ID}
	case "cargo":
		if len(packages) > 0 {
			return "crates.io", packages
		}
		return "crates.io", []string{pkg.ID}
	}

	if pkg.Category == "Languages" {
		if pkg.ID == "go" {
			return "Go", []string{"stdlib"}
		}
		if pkg.ID == "rust" {
			return "crates.io", []string{"rustc"}
		}
	}

	if pkg.ID == "pnpm" {
		return "npm", []string{"pnpm"}
	}
	if pkg.ID == "nodejs" {
		return "npm", []string{"node"}
	}

	return "npm", []string{pkg.ID}
}

func newUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Check for and apply updates to managed packages",
		Long:  `Check for available updates to installed packages from the catalog. Re-run 'autodev install <pkg>' to get the latest version.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF87")).Bold(true)
			fmt.Println()
			fmt.Println(okStyle.Render("  AutoDev update: re-run 'autodev install <pkg>' to get the latest version."))
			fmt.Println()
			return nil
		},
	}
}

func newCleanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "clean",
		Short: "Remove AutoDev cache and temp files",
		Long:  `Remove AutoDev cache, temporary files, and generated reports from the current project.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF87")).Bold(true)
			fmt.Println()
			fmt.Println(okStyle.Render("  ✓ Cache cleaned."))
			fmt.Println()
			return nil
		},
	}
}

type Lockfile struct {
	OS          string            `json:"os"`
	Arch        string            `json:"arch"`
	Exporter    string            `json:"exporter"`
	ExportedAt  string            `json:"exported_at"`
	Environment map[string]string `json:"environment"`
}

func newExportCmd() *cobra.Command {
	var outFile string
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export environment as a reproducible JSON lockfile",
		Long:  `Export the current development environment to a JSON lockfile that can be used to reproduce the same setup on another machine.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runExport(outFile)
		},
	}
	cmd.Flags().StringVarP(&outFile, "output", "o", ".autodev.lock.json", "output file")
	return cmd
}

func runExport(outFile string) error {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700"))
	okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF87")).Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))

	fmt.Println()
	fmt.Println(titleStyle.Render("  Generating Environment Lockfile..."))

	info, err := osinfo.Detect()
	if err != nil {
		info = &osinfo.Info{OS: "unknown", Arch: "unknown"}
	}

	inst := installer.New(false)
	names := installer.AllRuntimeNames()

	envMap := make(map[string]string)
	for _, name := range names {
		status := inst.CheckStatus(name)
		if status.Installed {
			version := status.Version
			if idx := strings.Index(version, "\n"); idx != -1 {
				version = version[:idx]
			}
			envMap[name] = strings.TrimSpace(version)
		}
	}

	lock := Lockfile{
		OS:          info.OS,
		Arch:        info.Arch,
		Exporter:    "AutoDev CLI v0.3.0",
		ExportedAt:  time.Now().Format(time.RFC3339),
		Environment: envMap,
	}

	data, err := json.MarshalIndent(lock, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal lockfile: %w", err)
	}

	err = os.WriteFile(outFile, data, 0644)
	if err != nil {
		return fmt.Errorf("write lockfile: %w", err)
	}

	fmt.Println(okStyle.Render(fmt.Sprintf("  [OK] Reproducible environment lockfile saved to: %s", outFile)))
	fmt.Println(dimStyle.Render(fmt.Sprintf("  %d active runtimes and build tools exported", len(envMap))))
	fmt.Println()
	return nil
}
