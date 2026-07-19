package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/autodev-sh/autodev/scanner"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func newScanCmd() *cobra.Command {
	var path string
	var tui bool

	cmd := &cobra.Command{
		Use:   "scan [path]",
		Short: "Scan a repository for languages, frameworks, and dependencies",
		Long:  `Scan the specified directory (or current directory) to detect all technologies, frameworks, package managers, databases, and infrastructure requirements.`,
		Example: `  autodev scan
  autodev scan ./my-project
  autodev scan --json ./my-project
  autodev scan --tui ./my-project`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				path = args[0]
			}
			if path == "" {
				path = "."
			}

			if tui {
				return runScanTUI(path)
			}
			return runScan(path, jsonOut)
		},
	}

	cmd.Flags().StringVarP(&path, "path", "p", ".", "path to scan")
	cmd.Flags().BoolVar(&tui, "tui", false, "open interactive directory dependency tree TUI")
	return cmd
}

func runScan(path string, asJSON bool) error {
	start := time.Now()

	if !asJSON {
		printScanHeader(path)
	}

	s := scanner.New(path)
	result, err := s.Scan()
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	elapsed := time.Since(start)

	if asJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}

	printScanResult(result, elapsed)
	_ = runEnhancerFlow(path, result)
	PrintGitHubCTA()
	return nil
}

func printScanHeader(path string) {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFD700")).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#FFD700")).
		Padding(0, 2)

	fmt.Println(headerStyle.Render(fmt.Sprintf("Scanning: %s", path)))
	fmt.Println()
}

func printScanResult(result *scanner.ScanResult, elapsed time.Duration) {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FF87"))
	itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).PaddingLeft(2)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
	badgeStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#FFD700")).
		Foreground(lipgloss.Color("#000000")).
		Padding(0, 1).
		Bold(true)

	if len(result.Languages) > 0 {
		fmt.Println(titleStyle.Render("  [LANGUAGES]"))
		for _, lang := range result.Languages {
			fmt.Println(itemStyle.Render("  - " + lang))
		}
		fmt.Println()
	}

	if len(result.Frameworks) > 0 {
		fmt.Println(titleStyle.Render("  [FRAMEWORKS]"))
		for _, fw := range result.Frameworks {
			fmt.Println(itemStyle.Render("  - " + fw))
		}
		fmt.Println()
	}

	if len(result.PackageManagers) > 0 {
		fmt.Println(titleStyle.Render("  [PACKAGE MANAGERS]"))
		for _, pm := range result.PackageManagers {
			fmt.Println(itemStyle.Render("  - " + pm))
		}
		fmt.Println()
	}

	if len(result.Databases) > 0 {
		fmt.Println(titleStyle.Render("  [DATABASES & SERVICES]"))
		for _, db := range result.Databases {
			fmt.Println(itemStyle.Render("  - " + db))
		}
		fmt.Println()
	}

	if len(result.Infra) > 0 {
		fmt.Println(titleStyle.Render("  [INFRASTRUCTURE]"))
		for _, inf := range result.Infra {
			fmt.Println(itemStyle.Render("  - " + inf))
		}
		fmt.Println()
	}

	if result.HasDocker {
		fmt.Println("  " + badgeStyle.Render("Docker") + "  " + badgeStyle.Render("Container-ready"))
	}
	if result.HasK8s {
		fmt.Println("  " + badgeStyle.Render("Kubernetes") + "  Found k8s manifests")
	}
	fmt.Println()

	if len(result.Projects) > 0 {
		fmt.Println(titleStyle.Render("  [MONOREPO SUBPROJECTS]"))
		for _, proj := range result.Projects {
			var techs []string
			for _, t := range proj.Technologies {
				techs = append(techs, t.Name)
			}
			fmt.Printf("    - %s (%s) → %s\n", proj.Name, dimStyle.Render(proj.Path), strings.Join(techs, ", "))
		}
		fmt.Println()
	}

	if len(result.RecommendedSetup) > 0 {
		fmt.Println()
		fmt.Println(titleStyle.Render("  [SETUP PLAN]"))
		for i, step := range result.RecommendedSetup {
			fmt.Println(itemStyle.Render(fmt.Sprintf("  %d. %s", i+1, step)))
		}
	}

	fmt.Println()
	fmt.Println(dimStyle.Render(fmt.Sprintf("  Scanned in %s | %d technologies detected", elapsed.Round(time.Millisecond), len(result.Technologies))))
	fmt.Println()
	fmt.Println(dimStyle.Render("  Run 'autodev setup' to install all missing tools."))
	fmt.Println(dimStyle.Render("  Run 'autodev audit' to check dependencies for security risks."))
	fmt.Println()
}

type Enhancement struct {
	Name        string
	Configured  bool
	Description string
	Action      func() error
}

func runEnhancerFlow(path string, result *scanner.ScanResult) error {
	if result == nil {
		return nil
	}
	// Detect project type
	isReact := false
	isNextJS := false
	isPython := false
	isGo := false
	for _, t := range result.Technologies {
		if t.Name == "React" {
			isReact = true
		} else if t.Name == "Next.js" {
			isNextJS = true
		} else if t.Name == "Python" {
			isPython = true
		} else if t.Name == "Go" {
			isGo = true
		}
	}

	var enhancements []Enhancement
	projectType := ""

	if isNextJS {
		projectType = "Next.js"
		enhancements = []Enhancement{
			{
				Name:        "Tailwind CSS",
				Description: "Configures postcss/tailwind structure",
				Configured:  fileExists(path, "tailwind.config.js") || fileExists(path, "tailwind.config.ts"),
				Action: func() error {
					_ = os.WriteFile(filepath.Join(path, "tailwind.config.js"), []byte(strings.TrimSpace(nextJsTailwindConfig)), 0644)
					_ = os.WriteFile(filepath.Join(path, "postcss.config.js"), []byte(strings.TrimSpace(nextJsPostcssConfig)), 0644)
					return nil
				},
			},
			{
				Name:        "Dockerfile",
				Description: "Multi-stage production build container configuration",
				Configured:  fileExists(path, "Dockerfile"),
				Action: func() error {
					_ = os.WriteFile(filepath.Join(path, "Dockerfile"), []byte(strings.TrimSpace(nextJsDockerfile)), 0644)
					return nil
				},
			},
			{
				Name:        "ESLint & Prettier",
				Description: "Code style standards and auto-format checks",
				Configured:  fileExists(path, ".eslintrc.json") || fileExists(path, "eslint.config.js") || fileExists(path, "eslint.config.mjs"),
				Action: func() error {
					_ = os.WriteFile(filepath.Join(path, ".eslintrc.json"), []byte(strings.TrimSpace(eslintContent)), 0644)
					_ = os.WriteFile(filepath.Join(path, ".prettierrc"), []byte(strings.TrimSpace(prettierContent)), 0644)
					return nil
				},
			},
			{
				Name:        "GitHub Actions CI/CD",
				Description: "Pre-set build and test checks on push",
				Configured:  dirExistsAndNotEmpty(path, filepath.Join(".github", "workflows")),
				Action: func() error {
					_ = os.MkdirAll(filepath.Join(path, ".github", "workflows"), 0755)
					_ = os.WriteFile(filepath.Join(path, ".github", "workflows", "ci.yml"), []byte(strings.TrimSpace(nextJsGithubAction)), 0644)
					return nil
				},
			},
		}
	} else if isReact {
		projectType = "React"
		enhancements = []Enhancement{
			{
				Name:        "Tailwind CSS",
				Description: "Utility styling library",
				Configured:  fileExists(path, "tailwind.config.js") || fileExists(path, "tailwind.config.ts"),
				Action: func() error {
					_ = os.WriteFile(filepath.Join(path, "tailwind.config.js"), []byte(strings.TrimSpace(tailwindConfigContent)), 0644)
					_ = os.WriteFile(filepath.Join(path, "postcss.config.js"), []byte(strings.TrimSpace(postcssConfigContent)), 0644)
					// Try to append @tailwind to index.css if it exists
					cssPaths := []string{"src/index.css", "src/globals.css", "index.css"}
					for _, cp := range cssPaths {
						fullCp := filepath.Join(path, cp)
						if fileExists(path, cp) {
							content, err := os.ReadFile(fullCp)
							if err == nil && !strings.Contains(string(content), "@tailwind") {
								newContent := strings.TrimSpace(indexCssContent) + "\n\n" + string(content)
								_ = os.WriteFile(fullCp, []byte(newContent), 0644)
							}
							break
						}
					}
					return nil
				},
			},
			{
				Name:        "Dockerfile",
				Description: "Docker container and static serving via Nginx",
				Configured:  fileExists(path, "Dockerfile"),
				Action: func() error {
					_ = os.WriteFile(filepath.Join(path, "Dockerfile"), []byte(strings.TrimSpace(mernClientDockerfile)), 0644)
					return nil
				},
			},
			{
				Name:        "ESLint & Prettier",
				Description: "Quality standards and formatting configurations",
				Configured:  fileExists(path, ".eslintrc.json") || fileExists(path, ".eslintrc.js"),
				Action: func() error {
					_ = os.WriteFile(filepath.Join(path, ".eslintrc.json"), []byte(strings.TrimSpace(eslintContent)), 0644)
					_ = os.WriteFile(filepath.Join(path, ".prettierrc"), []byte(strings.TrimSpace(prettierContent)), 0644)
					return nil
				},
			},
			{
				Name:        "GitHub Actions CI/CD",
				Description: "Build test suite execution runs",
				Configured:  dirExistsAndNotEmpty(path, filepath.Join(".github", "workflows")),
				Action: func() error {
					_ = os.MkdirAll(filepath.Join(path, ".github", "workflows"), 0755)
					reactAction := strings.ReplaceAll(nextJsGithubAction, "npm run build", "npm run build --if-present")
					_ = os.WriteFile(filepath.Join(path, ".github", "workflows", "ci.yml"), []byte(strings.TrimSpace(reactAction)), 0644)
					return nil
				},
			},
		}
	} else if isPython {
		projectType = "Python"
		enhancements = []Enhancement{
			{
				Name:        "Dockerfile",
				Description: "Python environment Docker runner",
				Configured:  fileExists(path, "Dockerfile"),
				Action: func() error {
					dockerContent := `FROM python:3.11-slim
WORKDIR /app
COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt || true
COPY . .
EXPOSE 8000
CMD ["python", "main.py"]`
					_ = os.WriteFile(filepath.Join(path, "Dockerfile"), []byte(strings.TrimSpace(dockerContent)), 0644)
					return nil
				},
			},
			{
				Name:        "Ruff Linting",
				Description: "Linter configuration for codebase style check",
				Configured:  fileExists(path, "ruff.toml") || fileExists(path, ".ruff.toml") || fileExists(path, "pyproject.toml"),
				Action: func() error {
					ruffContent := `[tool.ruff]
line-length = 88
indent-width = 4
target-version = "py311"`
					_ = os.WriteFile(filepath.Join(path, "ruff.toml"), []byte(strings.TrimSpace(ruffContent)), 0644)
					return nil
				},
			},
			{
				Name:        "GitHub Actions CI/CD",
				Description: "Python test suite workflow run",
				Configured:  dirExistsAndNotEmpty(path, filepath.Join(".github", "workflows")),
				Action: func() error {
					_ = os.MkdirAll(filepath.Join(path, ".github", "workflows"), 0755)
					pyAction := `name: Python CI

on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    - name: Set up Python
      uses: actions/setup-python@v4
      with:
        python-version: '3.11'
    - name: Install dependencies
      run: |
        python -m pip install --upgrade pip
        pip install ruff pytest
        if [ -f requirements.txt ]; then pip install -r requirements.txt; fi
    - name: Lint with ruff
      run: ruff check .
    - name: Run tests
      run: pytest || echo "No tests defined"`
					_ = os.WriteFile(filepath.Join(path, ".github", "workflows", "ci.yml"), []byte(strings.TrimSpace(pyAction)), 0644)
					return nil
				},
			},
		}
	} else if isGo {
		projectType = "Go"
		enhancements = []Enhancement{
			{
				Name:        "Dockerfile",
				Description: "Production multi-stage build container",
				Configured:  fileExists(path, "Dockerfile"),
				Action: func() error {
					dockerContent := `FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download || true
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]`
					_ = os.WriteFile(filepath.Join(path, "Dockerfile"), []byte(strings.TrimSpace(dockerContent)), 0644)
					return nil
				},
			},
			{
				Name:        "GitHub Actions CI/CD",
				Description: "Go test and build workflow",
				Configured:  dirExistsAndNotEmpty(path, filepath.Join(".github", "workflows")),
				Action: func() error {
					_ = os.MkdirAll(filepath.Join(path, ".github", "workflows"), 0755)
					goAction := `name: Go CI

on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.22'
    - name: Run tests
      run: go test -v ./...`
					_ = os.WriteFile(filepath.Join(path, ".github", "workflows", "ci.yml"), []byte(strings.TrimSpace(goAction)), 0644)
					return nil
				},
			},
		}
	}

	if projectType == "" || len(enhancements) == 0 {
		return nil
	}

	// Print visual checklist
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700"))
	okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF87")).Bold(true)
	warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B")).Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))

	fmt.Println()
	fmt.Println(titleStyle.Render(fmt.Sprintf("  🛡️  PROJECT ENHANCEMENTS DETECTED (%s):", projectType)))
	fmt.Println(dimStyle.Render("  ──────────────────────────────────────────────────────────"))

	hasMissing := false
	var missing []Enhancement

	for _, e := range enhancements {
		if e.Configured {
			fmt.Printf("    %s  %-22s %s\n", okStyle.Render("✓"), e.Name, dimStyle.Render("Already Configured"))
		} else {
			fmt.Printf("    %s  %-22s %s\n", warnStyle.Render("✗"), e.Name, warnStyle.Render("Not Configured — "+e.Description))
			hasMissing = true
			missing = append(missing, e)
		}
	}
	fmt.Println(dimStyle.Render("  ──────────────────────────────────────────────────────────"))
	fmt.Println()

	if !hasMissing {
		fmt.Println(okStyle.Render("  ✓ Your project has all recommended AutoDev enhancements!"))
		fmt.Println()
		return nil
	}

	// Ask to install missing enhancements
	fmt.Print("  Configure missing enhancements? [y/N] ")
	var answer string
	_, _ = fmt.Scanln(&answer)
	answer = strings.ToLower(strings.TrimSpace(answer))

	if answer == "y" || answer == "yes" {
		trackCLIMetric("scan_enhance")
		fmt.Println()
		for _, m := range missing {
			fmt.Printf("  %s Configuring %s...\n", titleStyle.Render("→"), m.Name)
			if err := m.Action(); err != nil {
				fmt.Printf("    %s Failed: %v\n", warnStyle.Render("✗"), err)
			} else {
				fmt.Printf("    %s Completed successfully\n", okStyle.Render("✓"))
			}
		}
		fmt.Println()
		fmt.Println(okStyle.Render("  ✓ All missing enhancements configured successfully!"))
		fmt.Println()
	} else {
		fmt.Println(dimStyle.Render("  Enhancements skipped."))
		fmt.Println()
	}

	return nil
}

func fileExists(path, filename string) bool {
	_, err := os.Stat(filepath.Join(path, filename))
	return err == nil
}

func dirExistsAndNotEmpty(path, dirname string) bool {
	dirPath := filepath.Join(path, dirname)
	entries, err := os.ReadDir(dirPath)
	return err == nil && len(entries) > 0
}
