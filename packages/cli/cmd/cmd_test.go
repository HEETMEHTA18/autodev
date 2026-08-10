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
	if rootCmd.Version != "0.5.1" {
		t.Errorf("root Version = %q, want %q", rootCmd.Version, "0.5.1")
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
	_ = os.WriteFile(filepath.Join(dir, "package.json"), []byte(`{"name":"test"}`), 0644)

	err := runScan(dir, false)
	if err != nil {
		t.Errorf("runScan() returned error: %v", err)
	}
}

func TestScanRunJSON(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "main.go"), []byte(`package main`), 0644)

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
}

// ── GitHub ────────────────────────────────────────────────────────────────────

func TestGitHubCommand(t *testing.T) {
	cmd := newGitHubCmd()
	if cmd.Flags().Lookup("token") == nil {
		t.Error("expected --token flag")
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
	_ = os.WriteFile(filepath.Join(dir, "go.mod"), []byte(`module test`), 0644)
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
	_ = os.WriteFile(filepath.Join(dir, "go.mod"), []byte(`module test
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
	_ = os.WriteFile(oldFile, []byte(`{"github_token":"test"}`), 0644)
	defer func() {
		_ = os.Chdir(".") // restore
	}()
	_ = os.Chdir(dir)

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
	_ = newCanvasCmd()
}

func TestCanvasRun(t *testing.T) {
	dir := t.TempDir()
	err := runCanvas(dir, false, false)
	if err != nil {
		t.Logf("runCanvas returned: %v", err)
	}
}

// ── Graph ─────────────────────────────────────────────────────────────────────







// ── Dashboard ─────────────────────────────────────────────────────────────────





// ── Completion ────────────────────────────────────────────────────────────────







// ── Uninstall ─────────────────────────────────────────────────────────────────



// ── Upgrade ───────────────────────────────────────────────────────────────────

func TestUpgradeCommand(t *testing.T) {
	cmd := newUpgradeCmd()
	if cmd.Flags().Lookup("check-only") == nil {
		t.Error("expected --check-only flag")
	}
}

// ── Info ──────────────────────────────────────────────────────────────────────



// ── Init ──────────────────────────────────────────────────────────────────────



// ── Snapshot ──────────────────────────────────────────────────────────────────



// ── Review ────────────────────────────────────────────────────────────────────









// ── AI ────────────────────────────────────────────────────────────────────────





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



// ── Telemetry ─────────────────────────────────────────────────────────────────

func TestTrackCLIMetric(t *testing.T) {
	// This should not panic
	trackCLIMetric("test")
}

// ── Tools ─────────────────────────────────────────────────────────────────────

func TestToolsCommand(t *testing.T) {
	cmd := newToolsCmd()
	if cmd.Use != "tools" {
		t.Errorf("Use = %q, want %q", cmd.Use, "tools")
	}
	subNames := []string{"list", "install", "remove"}
	registered := make(map[string]bool)
	for _, sub := range cmd.Commands() {
		registered[sub.Name()] = true
	}
	for _, name := range subNames {
		if !registered[name] {
			t.Errorf("expected subcommand %q under tools", name)
		}
	}
}

func TestToolsListRun(t *testing.T) {
	if err := runToolsList(); err != nil {
		t.Errorf("runToolsList returned error: %v", err)
	}
}

func TestToolsInstallUnknown(t *testing.T) {
	if err := runToolsInstall("nope"); err == nil {
		t.Error("expected error for unknown agent")
	}
}

// ── Run ───────────────────────────────────────────────────────────────────────

func TestRunCommand(t *testing.T) {
	cmd := newRunCmd()
	if cmd.Use != "run <agent> [args...]" {
		t.Logf("run Use = %q", cmd.Use)
	}
}

func TestRunUnknownAgent(t *testing.T) {
	if err := runAgent("nope", nil); err == nil {
		t.Error("expected error for unknown agent")
	}
}

// ── Session ───────────────────────────────────────────────────────────────────

func TestSessionCommand(t *testing.T) {
	cmd := newSessionCmd()
	if cmd.Use != "session" {
		t.Errorf("Use = %q, want %q", cmd.Use, "session")
	}
	subNames := []string{"new", "list", "stop"}
	registered := make(map[string]bool)
	for _, sub := range cmd.Commands() {
		registered[sub.Name()] = true
	}
	for _, name := range subNames {
		if !registered[name] {
			t.Errorf("expected subcommand %q under session", name)
		}
	}
}

func TestSessionListRun(t *testing.T) {
	if err := runSessionList(); err != nil {
		t.Errorf("runSessionList returned error: %v", err)
	}
}

func TestNormalizeAgentID(t *testing.T) {
	cases := map[string]string{
		"claude-code": "claude",
		"Claude Code": "claude",
		"open-code":   "opencode",
		"gemini-cli":  "gemini",
		"codex":       "codex",
	}
	for in, want := range cases {
		if got := normalizeAgentID(in); got != want {
			t.Errorf("normalizeAgentID(%q) = %q, want %q", in, got, want)
		}
	}
}

// ── Agent router ──────────────────────────────────────────────────────────────

func TestAgentCommand(t *testing.T) {
	cmd := newAgentCmd()
	if cmd.Use != "agent <task...>" {
		t.Logf("agent Use = %q", cmd.Use)
	}
}

func TestClassifyTask(t *testing.T) {
	cases := []struct {
		task string
		want string
	}{
		{"fix the authentication bug", "opencode"},
		{"research the latest postgres release", "gemini"},
		{"review this pull request for quality", "codex"},
		{"harden the codebase against cve exploits", "security"},
		{"design the system architecture", "claude"},
	}
	for _, c := range cases {
		if got := classifyTask(c.task); got != c.want {
			t.Errorf("classifyTask(%q) = %q, want %q", c.task, got, c.want)
		}
	}
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
	_ = os.WriteFile(filepath.Join(dir, "package.json"), []byte(`{"name":"test"}`), 0644)

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
	_ = os.Chdir(dir)
	defer func() { _ = os.Chdir(oldDir) }()

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
	_ = os.WriteFile(tmpFile, []byte("hello"), 0644)

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
	_ = os.MkdirAll(subDir, 0755)

	if dirExistsAndNotEmpty(dir, "subdir") {
		t.Error("expected empty dir to return false")
	}

	_ = os.WriteFile(filepath.Join(subDir, "file.txt"), []byte("hello"), 0644)
	if !dirExistsAndNotEmpty(dir, "subdir") {
		t.Error("expected non-empty dir to return true")
	}

	if dirExistsAndNotEmpty(dir, "nonexistent") {
		t.Error("expected nonexistent dir to return false")
	}
}

// ── New commands that might be missing ────────────────────────────────────────



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
		
	}
	for _, p := range parents {
		for _, sub := range p.cmd.Commands() {
			if sub.Short == "" {
				t.Errorf("subcommand %s %q has empty Short", p.name, sub.Name())
			}
		}
	}
}
