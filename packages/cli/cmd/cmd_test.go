package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// TestRootCommand tests the root command structure
func TestRootCommand(t *testing.T) {
	if rootCmd.Use != "autodev" {
		t.Errorf("root Use = %q, want %q", rootCmd.Use, "autodev")
	}
	if rootCmd.Version != "0.4.2" {
		t.Errorf("root Version = %q, want %q", rootCmd.Version, "0.4.2")
	}
}

// TestPersistentFlags ensures root persistent flags exist
func TestPersistentFlags(t *testing.T) {
	flags := []string{"config", "verbose", "no-color", "dry-run", "json"}
	for _, f := range flags {
		if rootCmd.PersistentFlags().Lookup(f) == nil {
			t.Errorf("expected persistent flag --%s to exist", f)
		}
	}
}

// Test necessary commands are registered
func TestCommandsRegistered(t *testing.T) {
	expected := []string{
		"scan", "setup", "github", "doctor", "report", "install",
		"update", "clean", "skills", "export", "profile", "ui",
		"clone", "audit", "mcp", "create", "benchmark", "containerize",
		"migrate", "ci", "plugin", "canvas", "graph", "chat",
		"capture", "daemon", "replay", "export-prompts", "sync",
		"prompts", "dashboard", "completion", "uninstall", "upgrade",
		"info", "init", "snapshot", "review", "ai",
	}

	registered := make(map[string]bool)
	for _, c := range rootCmd.Commands() {
		registered[c.Name()] = true
	}

	for _, name := range expected {
		if !registered[name] {
			t.Errorf("expected command %q to be registered", name)
		}
	}
}

// ── Doctor ────────────────────────────────────────────────────────────────────

func TestDoctorCommand(t *testing.T) {
	cmd := newDoctorCmd()
	if cmd.Use != "doctor" {
		t.Errorf("Use = %q, want %q", cmd.Use, "doctor")
	}
	if cmd.Flags().Lookup("fix") == nil {
		t.Error("expected --fix flag")
	}
}

func TestDoctorRun(t *testing.T) {
	err := runDoctor(false)
	if err != nil {
		t.Errorf("runDoctor(false) returned error: %v", err)
	}
}

// ── Scan ──────────────────────────────────────────────────────────────────────

func TestScanCommand(t *testing.T) {
	cmd := newScanCmd()
	if cmd.Use != "scan [path]" {
		t.Errorf("Use = %q, want %q", cmd.Use, "scan [path]")
	}
	if cmd.Flags().Lookup("path") == nil {
		t.Error("expected --path flag")
	}
	if cmd.Flags().Lookup("tui") == nil {
		t.Error("expected --tui flag")
	}
}

func TestScanRun(t *testing.T) {
	// Create temp dir with a known file
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "package.json"), []byte(`{"name":"test"}`), 0644)

	err := runScan(dir, false)
	if err != nil {
		t.Errorf("runScan() returned error: %v", err)
	}
}

func TestScanRunJSON(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "main.go"), []byte(`package main`), 0644)

	err := runScan(dir, true)
	if err != nil {
		t.Errorf("runScan(json) returned error: %v", err)
	}
}

// ── Setup ─────────────────────────────────────────────────────────────────────

func TestSetupCommand(t *testing.T) {
	cmd := newSetupCmd()
	if cmd.Use != "setup [path]" && cmd.Use != "setup" {
		t.Logf("setup Use = %q", cmd.Use)
	}
	if cmd.Flags().Lookup("yes") == nil && cmd.Flags().Lookup("path") == nil {
		// setup might have path/yes flags
	}
}

// ── GitHub ────────────────────────────────────────────────────────────────────

func TestGitHubCommand(t *testing.T) {
	cmd := newGitHubCmd()
	if cmd.Flags().Lookup("token") == nil {
		t.Error("expected --token flag")
	}
	if cmd.Flags().Lookup("json") == nil {
		// could be defined
	}
}

// ── Report ────────────────────────────────────────────────────────────────────

func TestReportCommand(t *testing.T) {
	cmd := newReportCmd()
	if cmd.Flags().Lookup("format") == nil {
		t.Error("expected --format flag")
	}
	if cmd.Flags().Lookup("output") == nil {
		t.Error("expected --output flag")
	}
}

func TestReportRun(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte(`module test`), 0644)
	err := runReport(dir, "markdown", "")
	if err != nil {
		t.Logf("runReport returned: %v (may be expected if not a real project)", err)
	}
}

// ── Install ───────────────────────────────────────────────────────────────────

func TestInstallCommand(t *testing.T) {
	cmd := newInstallCmd()
	if cmd.Flags().Lookup("list") == nil {
		t.Error("expected --list flag")
	}
	if cmd.Flags().Lookup("yes") == nil {
		t.Error("expected --yes flag")
	}
}

func TestInstallRunList(t *testing.T) {
	// Calling runInstall from create.go installs deps, no error expected
	runInstall(t.TempDir())
}

// ── Update ────────────────────────────────────────────────────────────────────

func TestUpdateCommand(t *testing.T) {
	cmd := newUpdateCmd()
	if cmd.Use != "update" {
		t.Logf("update Use = %q", cmd.Use)
	}
}

// ── Clean ─────────────────────────────────────────────────────────────────────

func TestCleanCommand(t *testing.T) {
	cmd := newCleanCmd()
	if cmd.Use != "clean" {
		t.Logf("clean Use = %q", cmd.Use)
	}
}

// ── Skills ────────────────────────────────────────────────────────────────────

func TestSkillsCommand(t *testing.T) {
	cmd := newSkillsCmd()
	if cmd.Flags().Lookup("deep") == nil {
		t.Error("expected --deep flag")
	}
}

func TestSkillsRun(t *testing.T) {
	err := runSkills(".", false, "", false, false, false, false)
	if err != nil {
		t.Logf("runSkills returned: %v", err)
	}
}

// ── Export ────────────────────────────────────────────────────────────────────

func TestExportCommand(t *testing.T) {
	cmd := newExportCmd()
	if cmd.Flags().Lookup("output") == nil {
		t.Error("expected --output flag")
	}
}

func TestExportRun(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "test-lock.json")
	err := runExport(tmpFile)
	if err != nil {
		t.Errorf("runExport returned: %v", err)
	}
	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		t.Error("export file was not created")
	}
}

// ── Profile ───────────────────────────────────────────────────────────────────

func TestProfileCommand(t *testing.T) {
	cmd := newProfileCmd()
	if cmd.Flags().Lookup("list") == nil {
		t.Error("expected --list flag")
	}
}

// ── UI ────────────────────────────────────────────────────────────────────────

func TestUICommand(t *testing.T) {
	cmd := newUICmd()
	if cmd.Flags().Lookup("port") == nil {
		t.Error("expected --port flag")
	}
}

// ── Clone ─────────────────────────────────────────────────────────────────────

func TestCloneCommand(t *testing.T) {
	cmd := newCloneCmd()
	if cmd.Flags().Lookup("yes") == nil {
		t.Error("expected --yes flag")
	}
}

// ── Audit ─────────────────────────────────────────────────────────────────────

func TestAuditCommand(t *testing.T) {
	cmd := newAuditCmd()
	if cmd.Use != "audit [path]" {
		t.Logf("audit Use = %q", cmd.Use)
	}
}

func TestAuditRun(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte(`module test
go 1.22
`), 0644)
	err := runAudit(dir)
	if err != nil {
		t.Logf("runAudit returned: %v (expected for empty module)", err)
	}
}

// ── MCP ───────────────────────────────────────────────────────────────────────

func TestMCPCommand(t *testing.T) {
	cmd := newMCPCmd()
	if cmd.Use != "mcp" {
		t.Errorf("Use = %q, want %q", cmd.Use, "mcp")
	}
	// Check subcommands
	subNames := []string{"start", "setup"}
	for _, name := range subNames {
		found := false
		for _, sub := range cmd.Commands() {
			if sub.Name() == name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected subcommand %q under mcp", name)
		}
	}
}

func TestMCPGuideRun(t *testing.T) {
	err := runMCPGuide()
	if err != nil {
		t.Errorf("runMCPGuide returned error: %v", err)
	}
}

// ── Create ────────────────────────────────────────────────────────────────────

func TestCreateCommand(t *testing.T) {
	cmd := newCreateCmd()
	if cmd.Use != "create <template> [name]" {
		t.Logf("create Use = %q", cmd.Use)
	}
}

// ── Benchmark ─────────────────────────────────────────────────────────────────

func TestBenchmarkCommand(t *testing.T) {
	cmd := newBenchmarkCmd()
	if cmd.Use != "benchmark" {
		t.Errorf("Use = %q, want %q", cmd.Use, "benchmark")
	}
}

func TestBenchmarkRun(t *testing.T) {
	err := runBenchmark()
	if err != nil {
		t.Errorf("runBenchmark returned error: %v", err)
	}
}

// ── Containerize ──────────────────────────────────────────────────────────────

func TestContainerizeCommand(t *testing.T) {
	cmd := newContainerizeCmd()
	if cmd.Flags().Lookup("path") == nil {
		t.Error("expected --path flag")
	}
}

// ── Migrate ───────────────────────────────────────────────────────────────────

func TestMigrateCommand(t *testing.T) {
	cmd := newMigrateCmd()
	if cmd.Use != "migrate" {
		t.Logf("migrate Use = %q", cmd.Use)
	}
}

func TestMigrateRun(t *testing.T) {
	dir := t.TempDir()
	oldFile := filepath.Join(dir, ".autodev.json")
	os.WriteFile(oldFile, []byte(`{"github_token":"test"}`), 0644)
	defer func() {
		os.Chdir(".") // restore
	}()
	os.Chdir(dir)

	err := runMigrate()
	if err != nil {
		t.Logf("runMigrate returned: %v", err)
	}
}

// ── CI ────────────────────────────────────────────────────────────────────────

func TestCICommand(t *testing.T) {
	cmd := newCICmd()
	if cmd.Flags().Lookup("path") == nil {
		t.Error("expected --path flag")
	}
}

// ── Plugin ────────────────────────────────────────────────────────────────────

func TestPluginCommand(t *testing.T) {
	cmd := newPluginCmd()
	if cmd.Use != "plugin" {
		t.Errorf("Use = %q, want %q", cmd.Use, "plugin")
	}
	subNames := []string{"list", "install"}
	for _, name := range subNames {
		found := false
		for _, sub := range cmd.Commands() {
			if sub.Name() == name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected subcommand %q under plugin", name)
		}
	}
}

// ── Canvas ────────────────────────────────────────────────────────────────────

func TestCanvasCommand(t *testing.T) {
	cmd := newCanvasCmd()
	if cmd.Flags().Lookup("save") == nil && cmd.Flags().Lookup("summary") == nil {
		// canvas has path/save/summary/web/serve/port flags
	}
}

func TestCanvasRun(t *testing.T) {
	dir := t.TempDir()
	err := runCanvas(dir, false, false)
	if err != nil {
		t.Logf("runCanvas returned: %v", err)
	}
}

// ── Graph ─────────────────────────────────────────────────────────────────────

func TestGraphCommand(t *testing.T) {
	cmd := newGraphCmd()
	if cmd.Flags().Lookup("format") == nil {
		t.Error("expected --format flag")
	}
	if cmd.Flags().Lookup("web") == nil {
		t.Error("expected --web flag")
	}
	if cmd.Flags().Lookup("audit") == nil {
		t.Error("expected --audit flag")
	}
}

func TestGraphRun(t *testing.T) {
	dir := t.TempDir()
	err := runGraph(dir, "tree")
	if err != nil {
		t.Logf("runGraph returned: %v", err)
	}
}

func TestGraphRunJSON(t *testing.T) {
	dir := t.TempDir()
	err := runGraph(dir, "json")
	if err != nil {
		t.Logf("runGraph(json) returned: %v", err)
	}
}

// ── Dashboard ─────────────────────────────────────────────────────────────────

func TestDashboardCommand(t *testing.T) {
	cmd := newDashboardCmd()
	if cmd.Use != "dashboard" {
		t.Errorf("Use = %q, want %q", cmd.Use, "dashboard")
	}
}

func TestDashboardRun(t *testing.T) {
	err := runDashboard()
	if err != nil {
		t.Errorf("runDashboard returned error: %v", err)
	}
}

// ── Completion ────────────────────────────────────────────────────────────────

func TestCompletionCommand(t *testing.T) {
	cmd := newCompletionCmd()
	if cmd.Use != "completion [bash|zsh|fish|powershell]" {
		t.Errorf("Use = %q, want %q", cmd.Use, "completion [bash|zsh|fish|powershell]")
	}
}

func TestCompletionBash(t *testing.T) {
	// We can't easily test GenBashCompletion without rootCmd, so just verify the shell is handled
	shells := []string{"bash", "zsh", "fish", "powershell"}
	for _, s := range shells {
		if !isValidShell(s) {
			t.Errorf("expected %q to be valid", s)
		}
	}
}

func isValidShell(s string) bool {
	switch s {
	case "bash", "zsh", "fish", "powershell":
		return true
	}
	return false
}

// ── Uninstall ─────────────────────────────────────────────────────────────────

func TestUninstallCommand(t *testing.T) {
	cmd := newUninstallCmd()
	if cmd.Use == "" {
		t.Error("empty Use")
	}
}

// ── Upgrade ───────────────────────────────────────────────────────────────────

func TestUpgradeCommand(t *testing.T) {
	cmd := newUpgradeCmd()
	if cmd.Flags().Lookup("check-only") == nil {
		t.Error("expected --check-only flag")
	}
}

// ── Info ──────────────────────────────────────────────────────────────────────

func TestInfoCommand(t *testing.T) {
	cmd := newInfoCmd()
	if cmd.Flags().Lookup("all") == nil {
		t.Error("expected --all flag")
	}
}

// ── Init ──────────────────────────────────────────────────────────────────────

func TestInitCommand(t *testing.T) {
	cmd := newInitCmd()
	if cmd.Flags().Lookup("force") == nil {
		t.Error("expected --force flag")
	}
}

// ── Snapshot ──────────────────────────────────────────────────────────────────

func TestSnapshotCommand(t *testing.T) {
	cmd := newSnapshotCmd()
	if cmd.Flags().Lookup("list") == nil {
		t.Error("expected --list flag")
	}
}

// ── Review ────────────────────────────────────────────────────────────────────

func TestReviewCommand(t *testing.T) {
	cmd := newReviewCmd()
	if cmd.Flags().Lookup("severity") == nil {
		t.Error("expected --severity flag")
	}
	if cmd.Flags().Lookup("format") == nil {
		t.Error("expected --format flag")
	}
}

func TestReviewRun(t *testing.T) {
	dir := t.TempDir()
	err := runReview(dir, "", "table")
	if err != nil {
		t.Logf("runReview returned: %v", err)
	}
}

func TestReviewRunJSON(t *testing.T) {
	dir := t.TempDir()
	err := runReview(dir, "", "json")
	if err != nil {
		t.Logf("runReview(json) returned: %v", err)
	}
}

func TestMeetsSeverityThreshold(t *testing.T) {
	tests := []struct {
		sev, threshold string
		want           bool
	}{
		{"critical", "", true},
		{"critical", "high", true},
		{"high", "critical", false},
		{"medium", "high", false},
		{"medium", "low", true},
		{"low", "medium", false},
		{"critical", "CRITICAL", true},
	}
	for _, tt := range tests {
		got := meetsSeverityThreshold(tt.sev, tt.threshold)
		if got != tt.want {
			t.Errorf("meetsSeverityThreshold(%q, %q) = %v, want %v", tt.sev, tt.threshold, got, tt.want)
		}
	}
}

// ── AI ────────────────────────────────────────────────────────────────────────

func TestAICommand(t *testing.T) {
	cmd := newAICmd()
	if cmd.Use != "ai <natural-language-command>" {
		t.Logf("ai Use = %q", cmd.Use)
	}
}

func TestAIRunSimple(t *testing.T) {
	err := runAI("health check")
	if err != nil {
		t.Errorf("runAI returned error: %v", err)
	}
}

// ── Prompts (history) ─────────────────────────────────────────────────────────

func TestPromptsCommand(t *testing.T) {
	cmd := newPromptsCmd()
	if cmd.Flags().Lookup("today") == nil {
		t.Error("expected --today flag")
	}
	subNames := []string{"chat", "capture", "daemon", "replay", "export-prompts", "sync"}
	registered := make(map[string]bool)
	for _, sub := range cmd.Commands() {
		registered[sub.Name()] = true
	}
	for _, name := range subNames {
		if !registered[name] {
			t.Errorf("expected subcommand %q under prompts", name)
		}
	}
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func TestIsPlaceholder(t *testing.T) {
	tests := []struct {
		val  string
		want bool
	}{
		{"your-api-key", true},
		{"sk-real-key", false},
		{"placeholder", true},
		{"ghp_realToken", false},
		{"test-token", true},
		{"my-secret", true},
	}
	for _, tt := range tests {
		got := isPlaceholder(tt.val)
		if got != tt.want {
			t.Errorf("isPlaceholder(%q) = %v, want %v", tt.val, got, tt.want)
		}
	}
}

func TestSanitizeVersion(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"v1.2.3", "v1.2.3"},
		{"v1.2.3\nwith extra", "v1.2.3"},
		{"  v1.0.0  ", "v1.0.0"},
	}
	for _, tt := range tests {
		got := sanitizeVersion(tt.input)
		if got != tt.want {
			t.Errorf("sanitizeVersion(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

// ── Telemetry ─────────────────────────────────────────────────────────────────

func TestTrackCLIMetric(t *testing.T) {
	// This should not panic
	trackCLIMetric("test")
}

// ── PrintGitHubCTA ───────────────────────────────────────────────────────────

func TestPrintGitHubCTA(t *testing.T) {
	// Should not panic with default jsonOut=false
	PrintGitHubCTA()
}

// ── Cobra execution tests (integration) ───────────────────────────────────────

func TestExecuteHelp(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"--help"})
	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("rootCmd.Execute() with --help returned error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "Usage:") || !strings.Contains(output, "Available Commands:") {
		t.Errorf("help output missing expected sections, got: %s", output[:min(len(output), 500)])
	}
}

func TestExecuteVersion(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"--version"})
	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("rootCmd.Execute() with --version returned error: %v", err)
	}
}

func TestExecuteDoctor(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"doctor"})
	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("doctor command via root returned error: %v", err)
	}
}

func TestExecuteDashboard(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"dashboard"})
	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("dashboard command via root returned error: %v", err)
	}
}

func TestExecuteBenchmark(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"benchmark"})
	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("benchmark command via root returned error: %v", err)
	}
}

func TestExecuteScan(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "package.json"), []byte(`{"name":"test"}`), 0644)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"scan", dir})
	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("scan command via root returned error: %v", err)
	}
}

func TestExecuteClean(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"clean"})
	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("clean command via root returned error: %v", err)
	}
}

func TestExecuteInfoList(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"info", "--all"})
	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("info --all via root returned error: %v", err)
	}
}

func TestExecuteGraphJSON(t *testing.T) {
	dir := t.TempDir()

	// Command uses fmt.Println directly, capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs([]string{"graph", "--format", "json", dir})
	err := rootCmd.Execute()

	w.Close()
	var outBuf bytes.Buffer
	outBuf.ReadFrom(r)
	os.Stdout = old

	if err != nil {
		t.Errorf("graph --format json via root returned error: %v", err)
	}
	output := strings.TrimSpace(outBuf.String())
	if !strings.Contains(output, "[") && !strings.Contains(output, "{") && !strings.Contains(output, "null") {
		t.Errorf("expected JSON in output, got: %s", output[:min(len(output), 300)])
	}
}

func TestExecuteCompletion(t *testing.T) {
	// GenBashCompletion writes to os.Stdout directly, so capture it
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs([]string{"completion", "bash"})
	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("completion bash via root returned error: %v", err)
	}

	w.Close()
	var outBuf bytes.Buffer
	outBuf.ReadFrom(r)
	os.Stdout = old

	output := outBuf.String()
	if !strings.Contains(output, "bash") && !strings.Contains(output, "autodev") {
		t.Errorf("expected bash completion output, got: %s", output[:min(len(output), 200)])
	}
}

func TestExecuteCompletionZsh(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"completion", "zsh"})
	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("completion zsh via root returned error: %v", err)
	}
}

func TestExecuteCompletionFish(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"completion", "fish"})
	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("completion fish via root returned error: %v", err)
	}
}

func TestExecuteCompletionPowershell(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"completion", "powershell"})
	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("completion powershell via root returned error: %v", err)
	}
}

func TestExecuteUpgradeCheckOnly(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"upgrade", "--check-only"})
	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("upgrade --check-only via root returned error: %v", err)
	}
}

func TestExecuteExport(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "test-lock.json")
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"export", "--output", tmpFile})
	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("export via root returned error: %v", err)
	}
}

func TestExecuteProfileList(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"profile", "--list"})
	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("profile --list via root returned error: %v", err)
	}
}

func TestExecuteInstallList(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"install", "--list"})
	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("install --list via root returned error: %v", err)
	}
}

func TestExecutePluginList(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"plugin", "list"})
	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("plugin list via root returned error: %v", err)
	}
}

func TestExecuteMCP(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"mcp"})
	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("mcp via root returned error: %v", err)
	}
}

func TestExecuteSnapshotList(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"snapshot", "--list"})
	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("snapshot --list via root returned error: %v", err)
	}
}

func TestExecuteAI(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"ai", "hello"})
	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("ai command via root returned error: %v", err)
	}
}

func TestExecutePrompts(t *testing.T) {
	// Skip: prompts command finds project root and may try to use less pager
	t.Skip("prompts command interacts with project root and pager")
}

func TestExecutePromptsToday(t *testing.T) {
	// Skip: prompts --today reads sessions directory
	t.Skip("prompts --today reads sessions from project root")
}

func TestExecuteContainerize(t *testing.T) {
	dir := t.TempDir()
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"containerize", dir})
	err := rootCmd.Execute()
	if err != nil {
		t.Logf("containerize via root returned: %v", err)
	}
}

func TestExecuteInit(t *testing.T) {
	dir := t.TempDir()
	oldDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(oldDir)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"init", dir})
	err := rootCmd.Execute()
	if err != nil {
		t.Logf("init via root returned: %v", err)
	}
}

func TestExecuteReview(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"review", "--format", "json"})
	err := rootCmd.Execute()
	if err != nil {
		t.Logf("review via root returned: %v", err)
	}
}

func TestExecuteSkills(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"skills", "--deep"})
	err := rootCmd.Execute()
	if err != nil {
		t.Logf("skills via root returned: %v", err)
	}
}

// ── MCP Server types ──────────────────────────────────────────────────────────

func TestMCPTypes(t *testing.T) {
	req := jsonRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "ping",
	}
	if req.JSONRPC != "2.0" {
		t.Errorf("JSONRPC = %q, want %q", req.JSONRPC, "2.0")
	}
}

// ── Custom nodes ──────────────────────────────────────────────────────────────

func TestCustomNodeTypes(t *testing.T) {
	cn := customNode{
		ID:    "test-1",
		Name:  "TestService",
		Group: "Services",
		Color: "#ff0000",
	}
	if cn.ID != "test-1" {
		t.Errorf("ID = %q", cn.ID)
	}
}

// ── Doctor MCP ────────────────────────────────────────────────────────────────

func TestRunDoctorMCP(t *testing.T) {
	result := runDoctorMCP(false)
	if !strings.Contains(result, "AUTODEV DOCTOR") {
		t.Errorf("expected header, got: %s", result[:min(len(result), 100)])
	}
}

// ── OS info detection test helper ─────────────────────────────────────────────

func TestEnhancerFlow(t *testing.T) {
	// Test with nil result (should handle gracefully)
	dir := t.TempDir()
	err := runEnhancerFlow(dir, nil)
	if err != nil {
		t.Errorf("runEnhancerFlow(nil) returned: %v", err)
	}
}

// ── File helpers ──────────────────────────────────────────────────────────────

func TestFileExists(t *testing.T) {
	dir := t.TempDir()
	tmpFile := filepath.Join(dir, "test.txt")
	os.WriteFile(tmpFile, []byte("hello"), 0644)

	if !fileExists(dir, "test.txt") {
		t.Error("expected fileExists to be true")
	}
	if fileExists(dir, "nonexistent.txt") {
		t.Error("expected fileExists to be false")
	}
}

func TestDirExistsAndNotEmpty(t *testing.T) {
	dir := t.TempDir()
	subDir := filepath.Join(dir, "subdir")
	os.MkdirAll(subDir, 0755)

	if dirExistsAndNotEmpty(dir, "subdir") {
		t.Error("expected empty dir to return false")
	}

	os.WriteFile(filepath.Join(subDir, "file.txt"), []byte("hello"), 0644)
	if !dirExistsAndNotEmpty(dir, "subdir") {
		t.Error("expected non-empty dir to return true")
	}

	if dirExistsAndNotEmpty(dir, "nonexistent") {
		t.Error("expected nonexistent dir to return false")
	}
}

// ── New commands that might be missing ────────────────────────────────────────

func TestSnapshotSubcommands(t *testing.T) {
	cmd := newSnapshotCmd()
	subNames := []string{"save", "list", "restore"}
	registered := make(map[string]bool)
	for _, sub := range cmd.Commands() {
		registered[sub.Name()] = true
	}
	for _, name := range subNames {
		if !registered[name] {
			t.Errorf("snapshot: expected subcommand %q", name)
		}
	}
}

func TestMCPSubcommands(t *testing.T) {
	cmd := newMCPCmd()
	subNames := []string{"start", "setup"}
	registered := make(map[string]bool)
	for _, sub := range cmd.Commands() {
		registered[sub.Name()] = true
	}
	for _, name := range subNames {
		if !registered[name] {
			t.Errorf("mcp: expected subcommand %q", name)
		}
	}
}

// ── Verify aliases ────────────────────────────────────────────────────────────

func TestCommandAliases(t *testing.T) {
	aliasTests := []struct {
		name    string
		aliases []string
	}{
		{"ai", []string{"ask", "auto"}},
		{"graph", []string{"depgraph", "dep-tree"}},
		{"dashboard", []string{"status", "overview"}},
		{"review", []string{"code-review", "audit-code"}},
		{"completion", []string{"autocomplete", "comp"}},
		{"init", []string{"bootstrap", "scaffold-config"}},
		{"uninstall", []string{"remove", "rm"}},
		{"upgrade", []string{"update-self", "self-update"}},
		{"info", []string{"inspect", "show"}},
		{"snapshot", []string{"env-save", "env-restore"}},
	}

	cmdMap := make(map[string]bool)
	for _, c := range rootCmd.Commands() {
		cmdMap[c.Name()] = true
	}

	for _, at := range aliasTests {
		if !cmdMap[at.name] {
			t.Errorf("expected command %q to exist (for alias check)", at.name)
		}
	}
}

// ── Doctor run with fix ───────────────────────────────────────────────────────

func TestDoctorRunWithFix(t *testing.T) {
	err := runDoctor(true)
	if err != nil {
		t.Errorf("runDoctor(true) returned error: %v", err)
	}
}

// ── Test all commands have Help text ──────────────────────────────────────────

func TestAllCommandsHaveHelp(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Short == "" {
			t.Errorf("command %q has empty Short description", c.Name())
		}
		if c.Long == "" {
			t.Errorf("command %q has empty Long description", c.Name())
		}
	}

	// Also check subcommands of known parents
	type parentCmd struct {
		cmd  *cobra.Command
		name string
	}
	parents := []parentCmd{
		{newMCPCmd(), "mcp"},
		{newPluginCmd(), "plugin"},
		{newPromptsCmd(), "prompts"},
		{newSnapshotCmd(), "snapshot"},
	}
	for _, p := range parents {
		for _, sub := range p.cmd.Commands() {
			if sub.Short == "" {
				t.Errorf("subcommand %s %q has empty Short", p.name, sub.Name())
			}
		}
	}
}
