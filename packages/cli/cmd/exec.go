package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/autodev-sh/autodev/catalog"
)

// execInstall executes the platform-specific install steps for a catalog package.
// This is the shared execution engine used by both the CLI commands and TUI.
func execInstall(pkg *catalog.Package) error {
	var steps catalog.InstallSteps
	switch runtime.GOOS {
	case "linux":
		steps = pkg.Install.Linux
	case "darwin":
		steps = pkg.Install.Darwin
	case "windows":
		steps = pkg.Install.Windows
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	if steps.Method == "" {
		return fmt.Errorf("no install method defined for %s on %s", pkg.Name, runtime.GOOS)
	}

	switch steps.Method {
	case "apt":
		return runLinuxInstall(steps.Packages)

	case "brew":
		return shellRun("brew", append([]string{"install"}, steps.Packages...)...)

	case "winget":
		for _, p := range steps.Packages {
			if err := shellRun("winget", "install", "--silent", p); err != nil {
				return err
			}
		}

	case "npm":
		return shellRun("npm", append([]string{"install", "-g"}, steps.Packages...)...)

	case "pip":
		args := []string{"install", "--upgrade"}
		if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
			args = append(args, "--break-system-packages")
		}
		args = append(args, steps.Packages...)
		return shellRun("pip3", args...)

	case "cargo":
		return shellRun("cargo", append([]string{"install"}, steps.Packages...)...)

	case "script", "manual":
		for _, script := range steps.Script {
			if dryRun {
				fmt.Printf("  [dry-run] %s\n", script)
				continue
			}
			if err := shellScript(script); err != nil {
				return err
			}
		}

	default:
		return fmt.Errorf("unknown install method: %s", steps.Method)
	}

	// Run post-install steps
	for _, post := range pkg.PostInstall {
		_ = shellScript(post) // best-effort
	}

	return nil
}

func shellRun(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func shellScript(script string) error {
	cmd := exec.Command("sh", "-c", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func isValidPackageName(name string) bool {
	if len(name) == 0 || name[0] == '-' {
		return false
	}
	for _, c := range name {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_' || c == '.' || c == '+') {
			return false
		}
	}
	return true
}

func runLinuxInstall(packages []string) error {
	for _, p := range packages {
		if !isValidPackageName(p) {
			return fmt.Errorf("unsafe or invalid package name detected: %s", p)
		}
	}

	if commandExists("apt-get") {
		return shellRun("sudo", append([]string{"apt-get", "install", "-y"}, packages...)...)
	} else if commandExists("dnf") {
		return shellRun("sudo", append([]string{"dnf", "install", "-y"}, packages...)...)
	} else if commandExists("pacman") {
		return shellRun("sudo", append([]string{"pacman", "-S", "--noconfirm"}, packages...)...)
	} else if commandExists("yum") {
		return shellRun("sudo", append([]string{"yum", "install", "-y"}, packages...)...)
	} else if commandExists("apk") {
		return shellRun("sudo", append([]string{"apk", "add"}, packages...)...)
	}
	return shellRun("sudo", append([]string{"apt-get", "install", "-y"}, packages...)...)
}

func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
