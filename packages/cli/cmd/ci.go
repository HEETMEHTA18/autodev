package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/autodev-sh/autodev/scanner"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func newCICmd() *cobra.Command {
	var path string

	cmd := &cobra.Command{
		Use:   "ci [path]",
		Short: "Generate CI/CD configuration template for the detected stack",
		Long: `Scan the workspace to detect languages and technologies, then generate a 
reproducible .github/workflows/ci.yml configuration.`,
		Example: `  autodev ci
  autodev ci ./my-project`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				path = args[0]
			}
			if path == "" {
				path = "."
			}
			return runCI(path)
		},
	}

	cmd.Flags().StringVarP(&path, "path", "p", ".", "path to project directory")
	return cmd
}

func runCI(path string) error {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700"))
	okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF87")).Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))

	fmt.Println()
	fmt.Println(titleStyle.Render("⚡ AutoDev CI/CD Template Generator"))
	fmt.Println(dimStyle.Render("  Scanning stack to generate GitHub Actions workflow..."))

	s := scanner.New(path)
	result, err := s.Scan()
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	// Determine primary language/tech
	primaryTech := "node"
	for _, t := range result.Technologies {
		name := strings.ToLower(t.Name)
		if name == "go" {
			primaryTech = "go"
			break
		} else if name == "rust" {
			primaryTech = "rust"
			break
		} else if name == "python" {
			primaryTech = "python"
			break
		} else if name == "flutter" || name == "dart" {
			primaryTech = "flutter"
			break
		} else if name == "java" || name == "kotlin" {
			primaryTech = "java"
			break
		}
	}

	// Workflow files directory
	workflowDir := filepath.Join(path, ".github", "workflows")
	if err := os.MkdirAll(workflowDir, 0755); err != nil {
		return fmt.Errorf("create .github/workflows directory: %w", err)
	}

	ciFilePath := filepath.Join(workflowDir, "ci.yml")
	var ciContent string

	switch primaryTech {
	case "go":
		ciContent = `name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    name: Go Lint & Test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.22"
          cache: true
      - name: Run Tests
        run: go test ./... -v -race
      - name: Build Binary
        run: go build -v ./...
`
	case "rust":
		ciContent = `name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    name: Rust Check & Test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: dtolnay/rust-toolchain@stable
        with:
          components: clippy, rustfmt
      - name: Check Format
        run: cargo fmt --check
      - name: Clippy Lint
        run: cargo clippy -- -D warnings
      - name: Run Tests
        run: cargo test --verbose
`
	case "python":
		ciContent = `name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    name: Python Lint & Test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-python@v5
        with:
          python-version: "3.11"
          cache: "pip"
      - name: Install Dependencies
        run: |
          python -m pip install --upgrade pip
          pip install flake8 pytest
          if [ -f requirements.txt ]; then pip install -r requirements.txt; fi
      - name: Lint with flake8
        run: flake8 . --count --select=E9,F63,F7,F82 --show-source --statistics
      - name: Run Tests
        run: pytest
`
	case "flutter":
		ciContent = `name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    name: Flutter Lint & Test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: subosito/flutter-action@v2
        with:
          channel: "stable"
          cache: true
      - name: Install Dependencies
        run: flutter pub get
      - name: Static Analysis
        run: flutter analyze
      - name: Run Tests
        run: flutter test
`
	case "java":
		ciContent = `name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    name: Java Build & Test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-java@v4
        with:
          distribution: "temurin"
          java-version: "17"
          cache: "maven"
      - name: Build with Maven
        run: mvn -B package --file pom.xml
`
	default: // Node.js / TypeScript / React / Next.js
		ciContent = `name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    name: Node Lint & Build
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: pnpm/action-setup@v3
        with:
          version: 9
      - uses: actions/setup-node@v4
        with:
          node-version: "20"
          cache: pnpm
      - name: Install Dependencies
        run: pnpm install
      - name: Lint
        run: pnpm run lint
      - name: Build
        run: pnpm run build
`
	}

	err = os.WriteFile(ciFilePath, []byte(ciContent), 0644)
	if err != nil {
		return fmt.Errorf("write ci.yml: %w", err)
	}

	fmt.Printf("  %s Generated CI/CD configuration template: %s\n", okStyle.Render("✓"), ciFilePath)
	fmt.Println()
	fmt.Println(okStyle.Render("  ✓ CI/CD template generation complete! You are ready to ship."))
	fmt.Println()
	PrintGitHubCTA()
	return nil
}
