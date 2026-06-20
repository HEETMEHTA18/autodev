package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/autodev-sh/autodev/scanner"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

type devContainerBuild struct {
	Dockerfile string `json:"dockerfile"`
}

type devContainerConfig struct {
	Name           string                 `json:"name"`
	Build          *devContainerBuild     `json:"build,omitempty"`
	Image          string                 `json:"image,omitempty"`
	Features       map[string]interface{} `json:"features,omitempty"`
	Customizations map[string]interface{} `json:"customizations,omitempty"`
	PostCreateCmd  string                 `json:"postCreateCommand,omitempty"`
}

type extensionsConfig struct {
	Recommendations []string `json:"recommendations"`
}

func newContainerizeCmd() *cobra.Command {
	var path string

	cmd := &cobra.Command{
		Use:   "containerize [path]",
		Short: "Generate DevContainer and VSCode workspace configuration",
		Long: `Scan the workspace to detect languages and technologies, then generate a 
reproducible .devcontainer.json setup and recommended VSCode extension configurations.`,
		Example: `  autodev containerize
  autodev containerize ./my-project`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				path = args[0]
			}
			if path == "" {
				path = "."
			}
			return runContainerize(path)
		},
	}

	cmd.Flags().StringVarP(&path, "path", "p", ".", "path to project directory")
	return cmd
}

func runContainerize(path string) error {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700"))
	okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF87")).Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))

	fmt.Println()
	fmt.Println(titleStyle.Render("⚡ AutoDev Cloud IDE Containerizer"))
	fmt.Println(dimStyle.Render("  Scanning stack to configure DevContainer & VSCode..."))

	s := scanner.New(path)
	result, err := s.Scan()
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	// 1. Determine base features and extensions
	features := make(map[string]interface{})
	var extensions []string

	hasNode := false
	hasPython := false
	hasGo := false
	hasRust := false
	hasJava := false
	hasFlutter := false
	hasDocker := false

	for _, t := range result.Technologies {
		switch strings.ToLower(t.Name) {
		case "node.js", "nodejs", "typescript", "react", "next.js":
			hasNode = true
		case "python":
			hasPython = true
		case "go":
			hasGo = true
		case "rust":
			hasRust = true
		case "java", "kotlin", "maven", "gradle":
			hasJava = true
		case "flutter", "dart":
			hasFlutter = true
		case "docker", "docker compose":
			hasDocker = true
		}
	}

	// Add devcontainer features & vscode extensions
	if hasNode {
		features["ghcr.io/devcontainers/features/node:1"] = map[string]interface{}{
			"version": "lts",
		}
		extensions = append(extensions, "dbaeumer.vscode-eslint", "esbenp.prettier-vscode")
	}
	if hasPython {
		features["ghcr.io/devcontainers/features/python:1"] = map[string]interface{}{
			"version": "latest",
		}
		extensions = append(extensions, "ms-python.python", "ms-python.vscode-pylance")
	}
	if hasGo {
		features["ghcr.io/devcontainers/features/go:1"] = map[string]interface{}{
			"version": "latest",
		}
		extensions = append(extensions, "golang.go")
	}
	if hasRust {
		features["ghcr.io/devcontainers/features/rust:1"] = map[string]interface{}{}
		extensions = append(extensions, "rust-lang.rust-analyzer")
	}
	if hasJava {
		features["ghcr.io/devcontainers/features/java:1"] = map[string]interface{}{
			"version": "21",
		}
		extensions = append(extensions, "vscjava.vscode-java-pack")
	}
	if hasDocker {
		features["ghcr.io/devcontainers/features/docker-in-docker:2"] = map[string]interface{}{}
		extensions = append(extensions, "ms-azuretools.vscode-docker")
	}
	if hasFlutter {
		extensions = append(extensions, "dart-code.flutter", "dart-code.dart-code")
	}

	// Always recommend git extensions
	extensions = append(extensions, "eamodio.gitlens")

	// 2. Generate .devcontainer directory, Dockerfile, and devcontainer.json
	devcontainerDir := filepath.Join(path, ".devcontainer")
	if err := os.MkdirAll(devcontainerDir, 0755); err != nil {
		return fmt.Errorf("create .devcontainer directory: %w", err)
	}

	dockerfilePath := filepath.Join(devcontainerDir, "Dockerfile")
	dockerfileContent := `FROM mcr.microsoft.com/devcontainers/base:ubuntu

# Pre-install base dependencies for AutoDev operations
RUN apt-get update && apt-get install -y \
    curl \
    git \
    build-essential \
    && rm -rf /var/lib/apt/lists/*
`
	if err := os.WriteFile(dockerfilePath, []byte(dockerfileContent), 0644); err != nil {
		return fmt.Errorf("write devcontainer Dockerfile: %w", err)
	}
	fmt.Printf("  %s Generated DevContainer Dockerfile: %s\n", okStyle.Render("✓"), dockerfilePath)

	devContainerPath := filepath.Join(devcontainerDir, "devcontainer.json")
	dc := devContainerConfig{
		Name:  "AutoDev Dev Container",
		Build: &devContainerBuild{Dockerfile: "Dockerfile"},
		Features: features,
		Customizations: map[string]interface{}{
			"vscode": map[string]interface{}{
				"extensions": extensions,
			},
		},
		PostCreateCmd: "curl -fsSL https://raw.githubusercontent.com/HEETMEHTA18/autodev/main/install.sh | bash && autodev setup --yes",
	}

	dcData, err := json.MarshalIndent(dc, "", "  ")
	if err != nil {
		return fmt.Errorf("serialize devcontainer configuration: %w", err)
	}

	err = os.WriteFile(devContainerPath, dcData, 0644)
	if err != nil {
		return fmt.Errorf("write devcontainer.json: %w", err)
	}
	fmt.Printf("  %s Generated DevContainer config: %s\n", okStyle.Render("✓"), devContainerPath)

	// 3. Generate .vscode/extensions.json
	vscodeDir := filepath.Join(path, ".vscode")
	if err := os.MkdirAll(vscodeDir, 0755); err != nil {
		return fmt.Errorf("create .vscode directory: %w", err)
	}

	extConfigPath := filepath.Join(vscodeDir, "extensions.json")
	extCfg := extensionsConfig{
		Recommendations: extensions,
	}

	extData, err := json.MarshalIndent(extCfg, "", "  ")
	if err != nil {
		return fmt.Errorf("serialize vscode extensions configuration: %w", err)
	}

	err = os.WriteFile(extConfigPath, extData, 0644)
	if err != nil {
		return fmt.Errorf("write extensions.json: %w", err)
	}
	fmt.Printf("  %s Generated VSCode extension recommendations: %s\n", okStyle.Render("✓"), extConfigPath)

	fmt.Println()
	fmt.Println(okStyle.Render("  ✓ Containerization setup complete! Open this workspace in VSCode DevContainers to run in a secure sandbox."))
	fmt.Println()
	PrintGitHubCTA()
	return nil
}
