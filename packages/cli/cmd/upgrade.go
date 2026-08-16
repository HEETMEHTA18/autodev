package cmd

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	owner = "HEETMEHTA18"
	repo  = "autodev"
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
	var installDir string
	cmd := &cobra.Command{
		Use:     "upgrade",
		Aliases: []string{"update-self", "self-update"},
		Short:   "Check for and apply AutoDev CLI updates",
		Long: `Check the latest release of AutoDev on GitHub and upgrade the CLI binary if a
newer version is available.

By default the binary that is currently running is replaced in place. Use
--install-dir to install to a new folder instead (the target is added to PATH
on Linux/macOS and the Windows user PATH automatically).`,
		Example: `  autodev upgrade
  autodev upgrade --check-only
  autodev upgrade --install-dir ~/.local/bin`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpgrade(checkOnly, installDir)
		},
	}
	cmd.Flags().BoolVar(&checkOnly, "check-only", false, "only check for updates, don't install")
	cmd.Flags().StringVar(&installDir, "install-dir", "", "install the updated binary to this folder (added to PATH if missing)")
	return cmd
}

func runUpgrade(checkOnly bool, installDir string) error {
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
		fmt.Println(dimStyle.Render("  Check manually: https://github.com/" + owner + "/" + repo + "/releases"))
		fmt.Println()
		return nil
	}

	currentTag := "v" + currentVersion()
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

	if err := downloadAndInstall(latest, installDir); err != nil {
		return fmt.Errorf("upgrade failed: %w", err)
	}

	// Make sure the target folder is on PATH for all platforms so the newly
	// installed binary is runnable after a fresh shell.
	if target, ok := upgradedBinaryDir(installDir); ok {
		if err := ensureOnPath(target); err != nil {
			fmt.Println(warnStyle.Render(fmt.Sprintf("  ⚠ Could not add %s to PATH: %v", target, err)))
		} else {
			fmt.Println(dimStyle.Render(fmt.Sprintf("  ✓ %s is on your PATH", target)))
		}
	}

	fmt.Println(successStyle.Render(fmt.Sprintf("  ✓ Upgraded to %s!", latest.TagName)))
	fmt.Println(dimStyle.Render("  Restart your shell or run 'autodev version' to verify."))
	fmt.Println()
	return nil
}

// upgradedBinaryDir returns the directory the freshly upgraded binary was
// installed into. It only returns a meaningful target when --install-dir was
// supplied or when the current binary lives in a user-writable location.
func upgradedBinaryDir(installDir string) (string, bool) {
	if installDir != "" {
		return installDir, true
	}
	exe, err := os.Executable()
	if err != nil {
		return "", false
	}
	dir := filepath.Dir(exe)
	home, _ := os.UserHomeDir()
	if home != "" && (dir == home+"/.local/bin" || dir == home+"/bin" || dir == home+"/.cargo/bin" || strings.Contains(dir, "/.local/bin")) {
		return dir, true
	}
	return "", false
}

// ensureOnPath adds dir to PATH across platforms:
//   - Linux/macOS: appends to ~/.bashrc, ~/.zshrc and ~/.profile if missing.
//   - Windows: appends to the user PATH via the registry (setx).
//
// It is best-effort and never fails hard.
func ensureOnPath(dir string) error {
	if dir == "" {
		return nil
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	switch runtime.GOOS {
	case "windows":
		return ensureOnPathWindows(abs)
	default:
		return ensureOnPathUnix(abs)
	}
}

func ensureOnPathUnix(dir string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	rcFiles := []string{home + "/.bashrc", home + "/.zshrc", home + "/.profile"}
	export := fmt.Sprintf("\n# AutoDev CLI\nexport PATH=%s:$PATH\n", quotePath(dir))

	for _, rc := range rcFiles {
		if _, err := os.Stat(rc); err != nil {
			continue // only touch files that already exist
		}
		data, err := os.ReadFile(rc)
		if err != nil {
			continue
		}
		if strings.Contains(string(data), dir) {
			continue
		}
		f, err := os.OpenFile(rc, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			continue
		}
		_, _ = f.WriteString(export)
		f.Close()
	}

	// Ensure a profile file exists so the PATH survives new shells even when
	// the user has no rc file yet.
	profile := home + "/.profile"
	if _, err := os.Stat(profile); err != nil {
		_ = os.WriteFile(profile, []byte(export), 0644)
		return nil
	}
	data, err := os.ReadFile(profile)
	if err != nil {
		return err
	}
	if strings.Contains(string(data), dir) {
		return nil
	}
	return os.WriteFile(profile, append(data, []byte(export)...), 0644)
}

func ensureOnPathWindows(dir string) error {
	out, err := exec.Command("reg", "query", "HKCU\\Environment", "/v", "Path").Output()
	if err != nil {
		// Fall back to setx directly.
		_ = exec.Command("setx", "PATH", dir).Run()
		return nil
	}
	current := strings.TrimSpace(string(out))
	if strings.Contains(current, dir) {
		return nil
	}
	newPath := dir
	if current != "" {
		lines := strings.Split(current, "\n")
		if len(lines) > 0 {
			last := strings.TrimSpace(lines[len(lines)-1])
			if idx := strings.Index(last, "Path"); idx != -1 {
				if val := strings.TrimSpace(last[idx+4:]); val != "" {
					newPath = dir + ";" + val
				}
			}
		}
	}
	_ = exec.Command("setx", "PATH", newPath).Run()
	return nil
}

func quotePath(dir string) string {
	if strings.ContainsAny(dir, " \t\"'") {
		return fmt.Sprintf("%q", dir)
	}
	return dir
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

func downloadAndInstall(release *GitHubRelease, installDir string) error {
	isWindows := runtime.GOOS == "windows"
	ext := ".tar.gz"
	if isWindows {
		ext = ".zip"
	}
	assetName := fmt.Sprintf("autodev_%s_%s_%s", strings.TrimPrefix(release.TagName, "v"), runtime.GOOS, runtime.GOARCH)
	exactName := assetName + ext

	var downloadURL string
	for _, asset := range release.Assets {
		if asset.Name == exactName {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}
	if downloadURL == "" {
		for _, asset := range release.Assets {
			if strings.HasPrefix(asset.Name, assetName) {
				downloadURL = asset.BrowserDownloadURL
				break
			}
		}
	}
	if downloadURL == "" {
		downloadURL = release.Assets[0].BrowserDownloadURL
	}

	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	fmt.Println(dimStyle.Render(fmt.Sprintf("  Downloading %s...", exactName)))

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: HTTP %d from %s", resp.StatusCode, downloadURL)
	}

	tmpFile, err := os.CreateTemp("", "autodev-*"+ext)
	if err != nil {
		return fmt.Errorf("temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		return err
	}
	tmpFile.Close()

	// Determine the destination: explicit install-dir, or replace in place.
	binaryPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot find current binary: %w", err)
	}
	if installDir != "" {
		abs, aerr := filepath.Abs(installDir)
		if aerr != nil {
			return aerr
		}
		if merr := os.MkdirAll(abs, 0755); merr != nil {
			return merr
		}
		binName := "autodev"
		if isWindows {
			binName = "autodev.exe"
		}
		binaryPath = filepath.Join(abs, binName)
	}

	// Extract into a temp dir, then locate the binary no matter what
	// layout the archive uses (GoReleaser nests it under a folder).
	extractDir, err := os.MkdirTemp("", "autodev-extract-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(extractDir)

	if isWindows {
		if err := extractZip(tmpFile.Name(), extractDir); err != nil {
			return fmt.Errorf("extraction failed: %w", err)
		}
	} else {
		if err := extractTarGz(tmpFile.Name(), extractDir); err != nil {
			return fmt.Errorf("extraction failed: %w", err)
		}
	}

	newBinary, err := findBinary(extractDir, isWindows)
	if err != nil {
		return err
	}

	if err := installBinary(newBinary, binaryPath); err != nil {
		return err
	}

	return nil
}

// currentVersion returns the runtime version of the running binary without a
// leading "v", preferring the ldflags-injected build version and falling back
// to the Cobra version constant.
func currentVersion() string {
	v := rootCmd.Version
	v = strings.TrimPrefix(strings.TrimSpace(v), "v")
	if v == "" {
		v = "0.5.2"
	}
	return v
}

// findBinary locates the autodev (or autodev.exe) executable anywhere under the
// extraction directory, handling both flat and nested archive layouts.
func findBinary(root string, isWindows bool) (string, error) {
	want := "autodev"
	if isWindows {
		want = "autodev.exe"
	}
	var found string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() == want {
			found = path
			return io.EOF
		}
		return nil
	})
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to scan extracted archive: %w", err)
	}
	if found == "" {
		return "", fmt.Errorf("open %s/autodev: no such file or directory", root)
	}
	return found, nil
}

// safeExtractPath resolves an archive entry name inside dir, returning the
// absolute target path. It returns ok=false for any entry that would escape
// dir via path traversal (e.g. "../", absolute paths, or symlink-writable
// paths), preventing "Zip Slip" style attacks.
func safeExtractPath(dir, name string) (string, bool) {
	if filepath.IsAbs(name) {
		return "", false
	}
	target := filepath.Join(dir, name)
	rel, err := filepath.Rel(dir, target)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", false
	}
	return target, true
}

// extractTarGz extracts a .tar.gz archive into dir, tolerating both flat and
// single-top-level-folder layouts.
func extractTarGz(archivePath, dir string) error {
	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		// Prevent path traversal from malicious archives.
		name := filepath.Clean(hdr.Name)
		target, ok := safeExtractPath(dir, name)
		if !ok {
			continue
		}
		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(out, tr); err != nil {
				out.Close()
				return err
			}
			out.Close()
		}
	}
	return nil
}

// extractZip extracts a .zip archive into dir.
func extractZip(archivePath, dir string) error {
	zr, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer zr.Close()

	for _, file := range zr.File {
		name := filepath.Clean(file.Name)
		target, ok := safeExtractPath(dir, name)
		if !ok {
			continue
		}
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}
		rc, err := file.Open()
		if err != nil {
			return err
		}
		out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, file.Mode())
		if err != nil {
			rc.Close()
			return err
		}
		if _, err := io.Copy(out, rc); err != nil {
			rc.Close()
			out.Close()
			return err
		}
		rc.Close()
		out.Close()
	}
	return nil
}

// installBinary replaces the running binary with the freshly extracted one,
// using an atomic rename where possible and falling back to a copy.
func installBinary(newBinary, binaryPath string) error {
	if err := os.Rename(newBinary, binaryPath); err == nil {
		_ = os.Chmod(binaryPath, 0755)
		return nil
	}

	// Fallback to copy (works across filesystems).
	input, err := os.ReadFile(newBinary)
	if err != nil {
		return fmt.Errorf("open %s: %w", newBinary, err)
	}
	if err := os.WriteFile(binaryPath, input, 0755); err != nil {
		return err
	}
	_ = os.Chmod(binaryPath, 0755)
	return nil
}
