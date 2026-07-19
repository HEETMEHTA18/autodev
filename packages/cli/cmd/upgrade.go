package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	currentVersion = "0.4.2"
	owner          = "HEETMEHTA18"
	repo           = "autodev"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
	Body    string `json:"body"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func newUpgradeCmd() *cobra.Command {
	var checkOnly bool
	cmd := &cobra.Command{
		Use:     "upgrade",
		Aliases: []string{"update-self", "self-update"},
		Short:   "Check for and apply AutoDev CLI updates",
		Long:    `Check the latest release of AutoDev on GitHub and upgrade the CLI binary if a newer version is available.`,
		Example: `  autodev upgrade
  autodev upgrade --check-only`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpgrade(checkOnly)
		},
	}
	cmd.Flags().BoolVar(&checkOnly, "check-only", false, "only check for updates, don't install")
	return cmd
}

func runUpgrade(checkOnly bool) error {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700"))
	successStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FF87"))
	warnStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF6B6B"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))

	fmt.Println()
	fmt.Println(titleStyle.Render("  🔄 AutoDev Upgrade Check"))
	fmt.Println()

	latest, err := fetchLatestRelease()
	if err != nil {
		fmt.Println(warnStyle.Render(fmt.Sprintf("  ✗ Failed to check for updates: %v", err)))
		fmt.Println(dimStyle.Render("  Check manually: https://github.com/"+owner+"/"+repo+"/releases"))
		fmt.Println()
		return nil
	}

	currentTag := "v" + currentVersion
	if latest.TagName == currentTag {
		fmt.Println(successStyle.Render(fmt.Sprintf("  ✓ You're running the latest version (%s)", currentTag)))
		fmt.Println()
		return nil
	}

	fmt.Println(dimStyle.Render(fmt.Sprintf("  Current version: %s", warnStyle.Render(currentTag))))
	fmt.Println(dimStyle.Render(fmt.Sprintf("  Latest version:  %s", successStyle.Render(latest.TagName))))
	fmt.Println()

	if checkOnly {
		fmt.Println(dimStyle.Render("  Run 'autodev upgrade' without --check-only to update."))
		fmt.Println()
		return nil
	}

	fmt.Print(dimStyle.Render("  Download and install ") + successStyle.Render(latest.TagName) + dimStyle.Render("? [y/N] "))
	var answer string
	_, _ = fmt.Scanln(&answer)
	answer = strings.ToLower(strings.TrimSpace(answer))
	if answer != "y" && answer != "yes" {
		fmt.Println(warnStyle.Render("  Upgrade cancelled."))
		fmt.Println()
		return nil
	}

	if err := downloadAndInstall(latest); err != nil {
		return fmt.Errorf("upgrade failed: %w", err)
	}

	fmt.Println(successStyle.Render(fmt.Sprintf("  ✓ Upgraded to %s!", latest.TagName)))
	fmt.Println(dimStyle.Render("  Restart your shell or run 'autodev version' to verify."))
	fmt.Println()
	return nil
}

func fetchLatestRelease() (*GitHubRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "autodev-upgrader")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned %d: %s", resp.StatusCode, string(body))
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}
	return &release, nil
}

func downloadAndInstall(release *GitHubRelease) error {
	assetName := fmt.Sprintf("autodev_%s_%s_%s.tar.gz", strings.TrimPrefix(release.TagName, "v"), runtime.GOOS, runtime.GOARCH)

	var downloadURL string
	for _, asset := range release.Assets {
		if asset.Name == assetName {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}
	if downloadURL == "" {
		downloadURL = release.Assets[0].BrowserDownloadURL
	}

	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	fmt.Println(dimStyle.Render(fmt.Sprintf("  Downloading %s...", assetName)))

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	tmpFile, err := os.CreateTemp("", "autodev-*.tar.gz")
	if err != nil {
		return fmt.Errorf("temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		return err
	}
	tmpFile.Close()

	// Find current binary path
	binaryPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot find current binary: %w", err)
	}

	// Extract and replace
	extractDir, err := os.MkdirTemp("", "autodev-extract-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(extractDir)

	cmd := exec.Command("tar", "-xzf", tmpFile.Name(), "-C", extractDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("extraction failed: %w", err)
	}

	newBinary := extractDir + "/autodev"
	if runtime.GOOS == "windows" {
		newBinary = extractDir + "/autodev.exe"
	}

	if err := os.Rename(newBinary, binaryPath); err != nil {
		// Fallback to copy
		input, err := os.ReadFile(newBinary)
		if err != nil {
			return err
		}
		if err := os.WriteFile(binaryPath, input, 0755); err != nil {
			return err
		}
	}

	return nil
}
