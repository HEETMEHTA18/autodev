package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// ── Help ─────────────────────────────────────────────────────────────────────

func TestHelpBrandedAndGrouped(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"--help"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("--help returned error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{
		"autodev v" + rootCmd.Version,
		"Usage:",
		"Diagnose:",
		"Examples:",
		"Flags:",
		"Run 'autodev <command> --help'",
		"completion",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("help output missing %q", want)
		}
	}
}

// ── Color preference ─────────────────────────────────────────────────────────

func TestNoColorPreferenceDisablesColor(t *testing.T) {
	orig := lipgloss.ColorProfile()
	defer lipgloss.SetColorProfile(orig)

	lipgloss.SetColorProfile(termenv.ANSI256)

	noColor = true
	applyColorPreference(nil)
	noColor = false
	if lipgloss.ColorProfile() != termenv.Ascii {
		t.Errorf("--no-color flag: profile = %v, want Ascii", lipgloss.ColorProfile())
	}

	lipgloss.SetColorProfile(termenv.ANSI256)
	applyColorPreference([]string{"scan", "--no-color"})
	if lipgloss.ColorProfile() != termenv.Ascii {
		t.Errorf("--no-color in argv: profile = %v, want Ascii", lipgloss.ColorProfile())
	}

	lipgloss.SetColorProfile(termenv.ANSI256)
	origNoColor := os.Getenv("NO_COLOR")
	_ = os.Setenv("NO_COLOR", "1")
	applyColorPreference(nil)
	if lipgloss.ColorProfile() != termenv.Ascii {
		t.Errorf("NO_COLOR env: profile = %v, want Ascii", lipgloss.ColorProfile())
	}
	_ = os.Setenv("NO_COLOR", origNoColor)
}

// ── Root command without a TTY ───────────────────────────────────────────────

// resetRootLocalFlags clears root-level flag values left behind by earlier
// tests in this process: pflag does not reset a parsed bool flag on the next
// Parse call, so a stale --help/--version would make Execute short-circuit
// before RunE ever runs.
func resetRootLocalFlags(t *testing.T) {
	t.Helper()
	for _, name := range []string{"help", "version"} {
		f := rootCmd.Flags().Lookup(name)
		if f == nil {
			continue
		}
		if err := f.Value.Set("false"); err != nil {
			t.Fatalf("reset --%s flag: %v", name, err)
		}
		f.Changed = false
	}
}

func TestRootWithoutTTYPrintsHelp(t *testing.T) {
	resetRootLocalFlags(t)

	oldIn, oldOut := os.Stdin, os.Stdout
	devNull, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatalf("open /dev/null: %v", err)
	}
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}

	helpBuf := new(bytes.Buffer)
	rootCmd.SetOut(helpBuf)
	rootCmd.SetArgs([]string{})

	os.Stdin = devNull
	os.Stdout = w
	execErr := rootCmd.Execute()
	w.Close()
	os.Stdin, os.Stdout = oldIn, oldOut
	devNull.Close()

	out, _ := io.ReadAll(r)
	r.Close()

	if execErr != nil {
		t.Fatalf("bare root without TTY returned error: %v", execErr)
	}
	if !strings.Contains(string(out), "terminal") {
		t.Errorf("expected terminal notice on stdout, got: %s", string(out))
	}
	if !strings.Contains(helpBuf.String(), "Usage:") {
		t.Errorf("expected help output, got: %s", helpBuf.String()[:min(len(helpBuf.String()), 300)])
	}
}

// ── doctor --json / --strict ─────────────────────────────────────────────────

func TestDoctorJSONOutput(t *testing.T) {
	jsonOut = true
	defer func() { jsonOut = false }()

	var err error
	var raw string
	raw = captureStdoutStderr(func() {
		err = runDoctor(false)
	})
	if err != nil {
		t.Fatalf("runDoctor(json) returned error: %v", err)
	}

	var report doctorReport
	if dErr := json.Unmarshal([]byte(raw), &report); dErr != nil {
		t.Fatalf("doctor --json output is not valid JSON: %v\n%s", dErr, raw)
	}
	if report.Version != rootCmd.Version {
		t.Errorf("version = %q, want %q", report.Version, rootCmd.Version)
	}
	if len(report.Checks) == 0 {
		t.Error("expected at least one check in JSON report")
	}
	if report.Issues != len(report.Checks)-countStatus(report.Checks, "ok") {
		t.Errorf("issues = %d does not match non-ok checks", report.Issues)
	}
	if report.Healthy != (report.Issues == 0) {
		t.Errorf("healthy = %v inconsistent with issues = %d", report.Healthy, report.Issues)
	}
	for _, c := range report.Checks {
		switch c.Status {
		case "ok", "warning", "error":
		default:
			t.Errorf("check %q has invalid status %q", c.Name, c.Status)
		}
	}
}

func countStatus(checks []checkResult, status string) int {
	n := 0
	for _, c := range checks {
		if c.Status == status {
			n++
		}
	}
	return n
}

func TestDoctorStrictReportsIssues(t *testing.T) {
	dir := t.TempDir()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(cwd) }()

	var runErr error
	_ = captureStdoutStderr(func() {
		runErr = runDoctorStrict(false, true)
	})
	// An empty directory has no .gitignore, so strict mode must fail.
	if runErr == nil {
		t.Error("runDoctorStrict(false, true) = nil, want error for empty dir")
	}

	var permissive error
	_ = captureStdoutStderr(func() {
		permissive = runDoctorStrict(false, false)
	})
	if permissive != nil {
		t.Errorf("runDoctorStrict(false, false) = %v, want nil", permissive)
	}
}

// ── clean ────────────────────────────────────────────────────────────────────

func TestCleanDryRunAndRemove(t *testing.T) {
	dir := t.TempDir()
	reports := filepath.Join(dir, "autodev-reports")
	if err := os.MkdirAll(filepath.Join(reports, "nested"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(reports, "nested", "a.txt"), []byte("12345"), 0644); err != nil {
		t.Fatal(err)
	}
	targets := []string{reports, filepath.Join(dir, "autodev-missing")}

	jsonOut = true
	defer func() { jsonOut = false }()

	var dryErr error
	raw := captureStdoutStderr(func() {
		dryErr = runClean(true, targets)
	})
	if dryErr != nil {
		t.Fatalf("runClean(dry-run) error: %v", dryErr)
	}
	var dryReport cleanReport
	if err := json.Unmarshal([]byte(raw), &dryReport); err != nil {
		t.Fatalf("dry-run output not JSON: %v", err)
	}
	if !dryReport.DryRun {
		t.Error("expected dry_run=true")
	}
	if _, err := os.Stat(reports); err != nil {
		t.Errorf("dry-run must not delete files: %v", err)
	}

	var cleanErr error
	raw = captureStdoutStderr(func() {
		cleanErr = runClean(false, targets)
	})
	if cleanErr != nil {
		t.Fatalf("runClean() error: %v", cleanErr)
	}
	var cleanReportOut cleanReport
	if err := json.Unmarshal([]byte(raw), &cleanReportOut); err != nil {
		t.Fatalf("clean output not JSON: %v", err)
	}
	if cleanReportOut.DryRun {
		t.Error("expected dry_run=false")
	}
	if cleanReportOut.Bytes < 5 {
		t.Errorf("bytes_total = %d, want >= 5", cleanReportOut.Bytes)
	}
	if _, err := os.Stat(reports); !os.IsNotExist(err) {
		t.Errorf("expected autodev-reports to be removed, stat err = %v", err)
	}
}

func TestCleanSafety(t *testing.T) {
	home, _ := os.UserHomeDir()
	unsafe := []string{"/", "/tmp", "", home, "."}
	for _, p := range unsafe {
		if isSafeCleanTarget(p) {
			t.Errorf("isSafeCleanTarget(%q) = true, want false", p)
		}
	}
	safe := []string{"autodev-reports", filepath.Join(os.TempDir(), "autodev-123")}
	for _, p := range safe {
		if !isSafeCleanTarget(p) {
			t.Errorf("isSafeCleanTarget(%q) = false, want true", p)
		}
	}
}

// ── version ──────────────────────────────────────────────────────────────────

func TestVersionJSON(t *testing.T) {
	jsonOut = true
	defer func() { jsonOut = false }()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"version", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("version --json error: %v", err)
	}

	var info map[string]string
	if err := json.Unmarshal(buf.Bytes(), &info); err != nil {
		t.Fatalf("version output not JSON: %v", err)
	}
	if info["version"] != rootCmd.Version {
		t.Errorf("version = %q, want %q", info["version"], rootCmd.Version)
	}
	if info["os"] == "" || info["arch"] == "" {
		t.Errorf("missing os/arch in version output: %v", info)
	}
}
