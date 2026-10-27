package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/autodev-sh/autodev/core/osinfo"
	"github.com/autodev-sh/autodev/installer"
	"github.com/autodev-sh/autodev/registry"
	"github.com/autodev-sh/autodev/scanner"
	"github.com/spf13/cobra"
)

func newDoctorCmd() *cobra.Command {
	var fix bool
	var strict bool
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Check the health and security of your codebase",
		Long:  `Verify that your codebase is healthy, secure, has correct git configurations, is free of exposed secrets, and matches standard linting and environment variables.`,
		Example: `  autodev doctor
  autodev doctor --fix
  autodev doctor --json
  autodev doctor --json --fix > report.json
  autodev doctor --strict`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDoctorStrict(fix, strict)
		},
	}
	cmd.Flags().BoolVar(&fix, "fix", false, "automatically attempt to fix detected codebase issues")
	cmd.Flags().BoolVar(&strict, "strict", false, "exit with a non-zero status when issues are found")
	return cmd
}

type codebaseDiagnostic struct {
	name        string
	description string
	checkFn     func(path string) (bool, string, error)
	fixFn       func(path string) (bool, error)
}

// checkResult is the machine-readable outcome of a single diagnostic check.
type checkResult struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"` // ok | warning | error
	Message     string `json:"message,omitempty"`
}

// fixResult is the machine-readable outcome of a single auto-fix attempt.
type fixResult struct {
	Name   string `json:"name"`
	Status string `json:"status"` // fixed | manual | failed
	Error  string `json:"error,omitempty"`
}

type systemReport struct {
	OS             string `json:"os"`
	Arch           string `json:"arch"`
	Version        string `json:"version"`
	PackageManager string `json:"package_manager"`
	CPUCores       int    `json:"cpu_cores"`
}

type agentReport struct {
	Name      string `json:"name"`
	Installed bool   `json:"installed"`
	Version   string `json:"version,omitempty"`
}

type doctorReport struct {
	Version     string        `json:"version"`
	Path        string        `json:"path"`
	GeneratedAt string        `json:"generated_at"`
	Healthy     bool          `json:"healthy"`
	Issues      int           `json:"issues"`
	System      *systemReport `json:"system,omitempty"`
	Agents      []agentReport `json:"agents"`
	Checks      []checkResult `json:"checks"`
	Fix         fixReport     `json:"fix"`
}

type fixReport struct {
	Attempted bool        `json:"attempted"`
	Results   []fixResult `json:"results,omitempty"`
	Log       string      `json:"log,omitempty"`
}

// runDoctor keeps the historical signature used by tests and call sites.
func runDoctor(fix bool) error {
	return runDoctorStrict(fix, false)
}

func runDoctorStrict(fix, strict bool) error {
	path := "."
	info, infoErr := osinfo.Detect()
	agents := registry.DetectAll()

	if jsonOut {
		return runDoctorJSON(path, fix, strict, info, infoErr, agents)
	}

	fmt.Println()
	fmt.Println(styleHeading.Render("AUTODEV CODEBASE HEALTH & SECURITY DOCTOR"))
	fmt.Println()

	if infoErr == nil {
		fmt.Println(styleHeading.Render("SYSTEM SPECIFICATIONS"))
		fmt.Printf("  %-20s %s\n", styleBold.Render("OS"), info.Version)
		fmt.Printf("  %-20s %s\n", styleBold.Render("Architecture"), info.Arch)
		fmt.Printf("  %-20s %s\n", styleBold.Render("Package Manager"), info.PackageManager)
		fmt.Println()
	}

	fmt.Println(styleHeading.Render("AI AGENTS"))
	for _, s := range agents {
		mark := styleWarn.Render("✗")
		status := styleDim.Render("not installed")
		if s.Installed {
			mark = styleOK.Render("✓")
			status = styleOK.Render(s.Version)
		}
		fmt.Printf("  %-6s %-16s %s\n", mark, s.Agent.Name, status)
	}
	fmt.Println()

	fmt.Println(styleHeading.Render("CODEBASE DIAGNOSTICS SCAN"))
	fmt.Println(styleDim.Render("  Scanning for secrets, configuration mismatches, and code errors..."))
	fmt.Println()

	_, toFix := evaluateDiagnostics(path, printCheckResult)

	fmt.Println()
	if len(toFix) == 0 {
		fmt.Println(styleOK.Render("  ✓ Codebase is completely healthy, secure, and ready for production!"))
		return nil
	}

	fmt.Printf("  %s issues detected in the codebase.\n", styleWarn.Render(fmt.Sprintf("%d", len(toFix))))

	if fix {
		fmt.Println()
		fmt.Println(styleHeading.Render("Attempting auto-remediation (--fix)..."))
		fmt.Println()
		applyFixes(path, toFix, func(name string) {
			fmt.Printf("  Healing %s...\n", name)
		}, func(fr fixResult) {
			switch fr.Status {
			case "fixed":
				fmt.Printf("    %s Successfully healed %s!\n", styleOK.Render("✓"), fr.Name)
			case "failed":
				fmt.Printf("    %s Failed to fix: %s\n", styleWarn.Render("✗"), fr.Error)
			default:
				fmt.Printf("    %s Requires manual action (refer to warning logs)\n", styleWarn.Render("!"))
			}
		})
		fmt.Println()
	} else {
		fmt.Println()
		fmt.Println(styleDim.Render("  Run 'autodev doctor --fix' to automatically resolve fixable issues."))
	}

	PrintGitHubCTA()
	if strict {
		return fmt.Errorf("%d issue(s) found", len(toFix))
	}
	return nil
}

func printCheckResult(r checkResult) {
	switch r.Status {
	case "ok":
		fmt.Printf("  %s %-25s %s\n", styleOK.Render(fmt.Sprintf("%-3s", iconSuccess)), r.Name, styleDim.Render("Clean & Healthy"))
	case "warning":
		fmt.Printf("  %s %-25s %s\n", styleWarn.Render(fmt.Sprintf("%-3s", iconWarn)), r.Name, styleCyan.Render(r.Message))
	default:
		fmt.Printf("  %s %-25s %s\n", styleErr.Render(fmt.Sprintf("%-3s", iconError)), r.Name, styleDim.Render(r.Message))
	}
	if verbose && r.Status != "ok" && r.Description != "" {
		fmt.Printf("  %-3s %s %s\n", "", "", styleDim.Render(r.Description))
	}
}

// evaluateDiagnostics runs every diagnostic check. When report is non-nil it
// is invoked as each check completes, which keeps text output streaming while
// JSON output stays silent.
func evaluateDiagnostics(path string, report func(checkResult)) ([]checkResult, []*codebaseDiagnostic) {
	diagnostics := getCodebaseDiagnostics()
	results := make([]checkResult, 0, len(diagnostics))
	var toFix []*codebaseDiagnostic

	for i := range diagnostics {
		d := &diagnostics[i]
		ok, msg, err := d.checkFn(path)
		r := checkResult{Name: d.name, Description: d.description}
		switch {
		case err != nil:
			r.Status = "error"
			r.Message = err.Error()
			toFix = append(toFix, d)
		case !ok:
			r.Status = "warning"
			r.Message = msg
			toFix = append(toFix, d)
		default:
			r.Status = "ok"
		}
		results = append(results, r)
		if report != nil {
			report(r)
		}
	}
	return results, toFix
}

// applyFixes attempts remediation for every failing diagnostic.
func applyFixes(path string, toFix []*codebaseDiagnostic, onStart func(string), report func(fixResult)) []fixResult {
	results := make([]fixResult, 0, len(toFix))
	for _, d := range toFix {
		if onStart != nil {
			onStart(d.name)
		}
		fixed, err := d.fixFn(path)
		fr := fixResult{Name: d.name}
		switch {
		case err != nil:
			fr.Status = "failed"
			fr.Error = err.Error()
		case fixed:
			fr.Status = "fixed"
		default:
			fr.Status = "manual"
		}
		results = append(results, fr)
		if report != nil {
			report(fr)
		}
	}
	return results
}

func runDoctorJSON(path string, fix, strict bool, info *osinfo.Info, infoErr error, agents []registry.Status) error {
	results, toFix := evaluateDiagnostics(path, nil)

	report := doctorReport{
		Version:     rootCmd.Version,
		Path:        path,
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Healthy:     len(toFix) == 0,
		Issues:      len(toFix),
		Agents:      make([]agentReport, 0, len(agents)),
		Checks:      results,
	}
	if infoErr == nil {
		report.System = &systemReport{
			OS:             info.OS,
			Arch:           info.Arch,
			Version:        info.Version,
			PackageManager: info.PackageManager,
			CPUCores:       info.CPUCores,
		}
	}
	for _, s := range agents {
		report.Agents = append(report.Agents, agentReport{
			Name:      s.Agent.Name,
			Installed: s.Installed,
			Version:   s.Version,
		})
	}

	if fix && len(toFix) > 0 {
		report.Fix.Attempted = true
		// fix functions print progress to stdout; capture it so the JSON
		// document on stdout stays valid.
		log := captureStdoutStderr(func() {
			report.Fix.Results = applyFixes(path, toFix, nil, nil)
		})
		report.Fix.Log = strings.TrimSpace(log)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(report); err != nil {
		return err
	}

	if strict && !report.Healthy {
		return fmt.Errorf("%d issue(s) found", report.Issues)
	}
	return nil
}

func getCodebaseDiagnostics() []codebaseDiagnostic {
	return []codebaseDiagnostic{
		{
			name:        "Git Configuration (.gitignore)",
			description: "Verify that sensitive environment or package files are not tracked in git",
			checkFn: func(path string) (bool, string, error) {
				gitIgnorePath := filepath.Join(path, ".gitignore")
				if _, err := os.Stat(gitIgnorePath); os.IsNotExist(err) {
					return false, "No .gitignore file detected. Environment files and builds could be committed.", nil
				}
				data, err := os.ReadFile(gitIgnorePath)
				if err != nil {
					return false, "", err
				}
				content := string(data)
				hasEnv := strings.Contains(content, ".env")
				hasNode := strings.Contains(content, "node_modules")
				if !hasEnv || !hasNode {
					return false, ".gitignore does not ignore sensitive .env files or node_modules.", nil
				}
				return true, "", nil
			},
			fixFn: func(path string) (success bool, err error) {
				gitIgnorePath := filepath.Join(path, ".gitignore")
				if _, statErr := os.Stat(gitIgnorePath); os.IsNotExist(statErr) {
					defaultGitignore := `# Standard AutoDev Gitignore
node_modules/
dist/
build/
.next/
.turbo/
.env
.env.local
.env.*.local
*.log
packages/cli/bin/
`
					writeErr := os.WriteFile(gitIgnorePath, []byte(defaultGitignore), 0644)
					return writeErr == nil, writeErr
				}

				// Append to existing
				data, readErr := os.ReadFile(gitIgnorePath)
				if readErr != nil {
					return false, readErr
				}
				content := string(data)
				var toAppend []string
				if !strings.Contains(content, ".env") {
					toAppend = append(toAppend, ".env", ".env.local", ".env.*.local")
				}
				if !strings.Contains(content, "node_modules") {
					toAppend = append(toAppend, "node_modules/")
				}

				if len(toAppend) > 0 {
					f, openErr := os.OpenFile(gitIgnorePath, os.O_APPEND|os.O_WRONLY, 0644)
					if openErr != nil {
						return false, openErr
					}
					defer func() {
						if closeErr := f.Close(); closeErr != nil && err == nil {
							err = closeErr
						}
					}()
					if _, writeErr := f.WriteString("\n# Added by AutoDev Doctor\n" + strings.Join(toAppend, "\n") + "\n"); writeErr != nil {
						return false, writeErr
					}
				}
				return true, nil
			},
		},
		{
			name:        "Exposed Secrets Scanner",
			description: "Scan the codebase to verify no hardcoded secrets or API keys are present",
			checkFn: func(path string) (bool, string, error) {
				var matches []string
				secretRegexes := []*regexp.Regexp{
					// GitHub PAT
					regexp.MustCompile(`\b(ghp_[a-zA-Z0-9]{36})\b`),
					regexp.MustCompile(`\b(github_pat_[a-zA-Z0-9_]{82})\b`),
					// AWS Access Key ID
					regexp.MustCompile(`\b(AKIA[0-9A-Z]{16})\b`),
					// AWS Secret Key
					regexp.MustCompile(`(?i)aws_secret_access_key\s*[:=]\s*['"]([a-zA-Z0-9/+=]{40})['"]`),
					// Google API Key
					regexp.MustCompile(`\b(AIza[0-9A-Za-z\-_]{35})\b`),
					// Stripe API Key
					regexp.MustCompile(`\b(sk_live_[0-9a-zA-Z]{24})\b`),
					// Slack Webhook URL
					regexp.MustCompile(`^\s*(https://hooks\.slack\.com/services/T[a-zA-Z0-9_]+/B[a-zA-Z0-9_]+/[a-zA-Z0-9_]+)\s*$`),
					// General credentials pattern
					regexp.MustCompile(`(?i)(api_key|secret|password|private_key|token|auth_token)\s*[:=]\s*['"]([a-zA-Z0-9_\-\.]{16,})['"]`),
				}

				err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
					if err != nil {
						return nil
					}
					// skip noise
					base := d.Name()
					if d.IsDir() {
						if p != path && (strings.HasPrefix(base, ".") ||
							base == "node_modules" || base == "vendor" ||
							base == "dist" || base == "build" || base == ".next" || base == ".turbo") {
							return filepath.SkipDir
						}
						return nil
					}

					// Only scan code files
					ext := filepath.Ext(p)
					if ext != ".js" && ext != ".ts" && ext != ".tsx" && ext != ".py" && ext != ".go" && ext != ".json" && ext != ".yml" && ext != ".env" {
						return nil
					}

					// Skip lockfiles and .env.example
					if base == "package-lock.json" || base == "pnpm-lock.yaml" || base == "yarn.lock" || strings.HasSuffix(base, ".example") {
						return nil
					}

					data, err := os.ReadFile(p)
					if err != nil {
						return nil
					}

					lines := strings.Split(string(data), "\n")
					for lineNum, line := range lines {
						for _, rx := range secretRegexes {
							sub := rx.FindStringSubmatch(line)
							if len(sub) > 0 {
								secretVal := ""
								if len(sub) > 2 {
									secretVal = sub[2]
								} else if len(sub) > 1 {
									secretVal = sub[1]
								} else {
									secretVal = sub[0]
								}
								if secretVal != "" && !isPlaceholder(secretVal) {
									matches = append(matches, fmt.Sprintf("%s (Line %d)", filepath.Base(p), lineNum+1))
									break
								}
							}
						}
					}

					return nil
				})

				if err != nil {
					return false, "", err
				}

				if len(matches) > 0 {
					return false, fmt.Sprintf("Exposed credentials detected: %s", strings.Join(matches, ", ")), nil
				}
				return true, "", nil
			},
			fixFn: func(path string) (bool, error) {
				// Secrets are sensitive and hard to replace without breaking codes.
				// We prompt the user on how to fix instead of destructive rewrite.
				return false, fmt.Errorf("exposed secrets must be manually migrated to .env and referenced dynamically")
			},
		},
		{
			name:        "Environment Config (.env)",
			description: "Verify that environmental profiles are properly configured",
			checkFn: func(path string) (bool, string, error) {
				envExample := filepath.Join(path, ".env.example")
				envFile := filepath.Join(path, ".env")

				if _, err := os.Stat(envExample); err == nil {
					if _, err := os.Stat(envFile); os.IsNotExist(err) {
						return false, "Found .env.example but no active .env file was configured.", nil
					}
				}
				return true, "", nil
			},
			fixFn: func(path string) (bool, error) {
				envExample := filepath.Join(path, ".env.example")
				envFile := filepath.Join(path, ".env")
				if _, err := os.Stat(envExample); err == nil {
					data, err := os.ReadFile(envExample)
					if err != nil {
						return false, err
					}
					err = os.WriteFile(envFile, data, 0644)
					return err == nil, err
				}
				return true, nil
			},
		},
		{
			name:        "Dependency Lockfiles",
			description: "Verify dependencies configurations are locked and synchronized",
			checkFn: func(path string) (bool, string, error) {
				pkgJSON := filepath.Join(path, "package.json")
				if _, err := os.Stat(pkgJSON); err == nil {
					lockfiles := []string{"package-lock.json", "pnpm-lock.yaml", "yarn.lock"}
					found := false
					for _, lf := range lockfiles {
						if _, err := os.Stat(filepath.Join(path, lf)); err == nil {
							found = true
							break
						}
					}
					if !found {
						return false, "package.json found but no matching package lockfile exists.", nil
					}
				}
				return true, "", nil
			},
			fixFn: func(path string) (bool, error) {
				// If pnpm is available, run pnpm install. Otherwise try npm install.
				cmdName := "npm"
				if _, err := exec.LookPath("pnpm"); err == nil {
					cmdName = "pnpm"
				}

				cmd := exec.Command(cmdName, "install")
				cmd.Dir = path
				err := cmd.Run()
				return err == nil, err
			},
		},
		{
			name:        "Linter & Code Format Status",
			description: "Scan codebase for syntax compilation or structural formatting warnings",
			checkFn: func(path string) (bool, string, error) {
				pkgJSON := filepath.Join(path, "package.json")
				if _, err := os.Stat(pkgJSON); err == nil {
					// Check if lint script exists
					data, err := os.ReadFile(pkgJSON)
					if err != nil {
						return true, "", nil
					}
					if strings.Contains(string(data), `"lint"`) {
						cmdName := "npm"
						if _, err := exec.LookPath("pnpm"); err == nil {
							cmdName = "pnpm"
						}
						var stdout, stderr bytes.Buffer
						cmd := exec.Command(cmdName, "run", "lint")
						cmd.Dir = path
						cmd.Stdout = &stdout
						cmd.Stderr = &stderr
						if err := cmd.Run(); err != nil {
							return false, "Linter run reported syntax or formatting problems.", nil
						}
					}
				}
				return true, "", nil
			},
			fixFn: func(path string) (bool, error) {
				cmdName := "npm"
				if _, err := exec.LookPath("pnpm"); err == nil {
					cmdName = "pnpm"
				}

				// Attempt auto-linting fixes
				var cmd *exec.Cmd
				if cmdName == "pnpm" {
					cmd = exec.Command("pnpm", "run", "lint", "--fix")
				} else {
					cmd = exec.Command("npm", "run", "lint", "--", "--fix")
				}
				cmd.Dir = path
				_ = cmd.Run() // run best effort

				// Also run prettier if available
				if _, err := exec.LookPath("prettier"); err == nil {
					pcmd := exec.Command("npx", "prettier", "--write", "**/*.{js,ts,tsx,json,css,md}")
					pcmd.Dir = path
					_ = pcmd.Run()
				}
				return true, nil
			},
		},
		{
			name:        "Environment Lockfile Mismatch",
			description: "Verify active tools match the versions declared in .autodev.lock.json",
			checkFn: func(path string) (bool, string, error) {
				lockPath := filepath.Join(path, ".autodev.lock.json")
				if _, err := os.Stat(lockPath); os.IsNotExist(err) {
					return true, "", nil // lockfile is optional, skip check
				}
				data, err := os.ReadFile(lockPath)
				if err != nil {
					return false, "", err
				}
				var lock struct {
					Environment map[string]string `json:"environment"`
				}
				if err := json.Unmarshal(data, &lock); err != nil {
					return false, "Lockfile is corrupted or invalid JSON", nil
				}

				inst := installer.New(false)
				var mismatches []string
				for name, targetVer := range lock.Environment {
					status := inst.CheckStatus(name)
					if !status.Installed {
						mismatches = append(mismatches, fmt.Sprintf("%s (not installed, expected %s)", name, targetVer))
					} else {
						activeVer := strings.TrimSpace(status.Version)
						if idx := strings.Index(activeVer, "\n"); idx != -1 {
							activeVer = activeVer[:idx]
						}
						if !strings.Contains(strings.ToLower(activeVer), strings.ToLower(targetVer)) &&
							!strings.Contains(strings.ToLower(targetVer), strings.ToLower(activeVer)) {
							mismatches = append(mismatches, fmt.Sprintf("%s (active: %s, expected: %s)", name, activeVer, targetVer))
						}
					}
				}

				if len(mismatches) > 0 {
					return false, fmt.Sprintf("Mismatched environment versions: %s", strings.Join(mismatches, ", ")), nil
				}
				return true, "", nil
			},
			fixFn: func(path string) (bool, error) {
				lockPath := filepath.Join(path, ".autodev.lock.json")
				if _, err := os.Stat(lockPath); os.IsNotExist(err) {
					return true, nil
				}
				data, err := os.ReadFile(lockPath)
				if err != nil {
					return false, err
				}
				var lock struct {
					Environment map[string]string `json:"environment"`
				}
				if err := json.Unmarshal(data, &lock); err != nil {
					return false, err
				}

				inst := installer.New(false)
				fixedAny := false
				for name, targetVer := range lock.Environment {
					status := inst.CheckStatus(name)
					mismatch := false
					if !status.Installed {
						mismatch = true
					} else {
						activeVer := strings.TrimSpace(status.Version)
						if idx := strings.Index(activeVer, "\n"); idx != -1 {
							activeVer = activeVer[:idx]
						}
						if !strings.Contains(strings.ToLower(activeVer), strings.ToLower(targetVer)) &&
							!strings.Contains(strings.ToLower(targetVer), strings.ToLower(activeVer)) {
							mismatch = true
						}
					}
					if mismatch {
						fmt.Printf("    Re-installing/restoring matching runtime: %s...\n", name)
						if err := inst.Install(name); err != nil {
							return false, fmt.Errorf("failed to restore %s: %w", name, err)
						}
						fixedAny = true
					}
				}
				return fixedAny, nil
			},
		},
		{
			name:        "Supply-Chain Vulnerabilities",
			description: "Verify dependencies are free of known vulnerabilities",
			checkFn: func(path string) (bool, string, error) {
				results, err := scanner.AuditRepository(path)
				if err != nil {
					return false, "", err
				}
				if len(results) > 0 {
					var vulnPacks []string
					for _, res := range results {
						vulnPacks = append(vulnPacks, fmt.Sprintf("%s@%s", res.Package.Name, res.Package.Version))
					}
					return false, fmt.Sprintf("Known vulnerabilities found in packages: %s", strings.Join(vulnPacks, ", ")), nil
				}
				return true, "", nil
			},
			fixFn: func(path string) (bool, error) {
				pkgJSON := filepath.Join(path, "package.json")
				if _, err := os.Stat(pkgJSON); err == nil {
					cmdName := "npm"
					if _, err := exec.LookPath("pnpm"); err == nil {
						cmdName = "pnpm"
					}

					var cmd *exec.Cmd
					if cmdName == "pnpm" {
						cmd = exec.Command("pnpm", "update")
					} else {
						cmd = exec.Command("npm", "audit", "fix")
					}
					cmd.Dir = path
					err := cmd.Run()
					return err == nil, err
				}
				return true, nil
			},
		},
	}
}

func isPlaceholder(val string) bool {
	val = strings.ToLower(val)
	placeholders := []string{"your", "placeholder", "key_here", "my-secret", "dummy", "example", "test", "token_here", "config"}
	for _, p := range placeholders {
		if strings.Contains(val, p) {
			return true
		}
	}
	return false
}
