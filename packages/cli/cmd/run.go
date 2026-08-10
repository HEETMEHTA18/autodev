package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"

	"github.com/autodev-sh/autodev/registry"
	"github.com/creack/pty"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// newRunCmd launches an AI agent CLI inside a managed pseudo-terminal.
func newRunCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run <agent> [args...]",
		Short: "Launch an AI agent in a managed terminal",
		Long: `Launch an AI agent CLI (OpenCode, Claude, Codex, Gemini, Aider) inside a
pseudo-terminal. Arguments after the agent name are passed through.

  autodev run opencode
  autodev run codex exec "fix the auth bug"
  autodev run gemini -- "research prompt..."`,
		Args:               cobra.MinimumNArgs(1),
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAgent(args[0], args[1:])
		},
	}
}

// runAgent launches agent id with the given args, preferring a PTY and
// falling back to plain stdio passthrough when stdin is not a terminal.
func runAgent(id string, args []string) error {
	agent, err := registry.Get(id)
	if err != nil {
		return err
	}

	path, err := exec.LookPath(agent.Command)
	if err != nil {
		return fmt.Errorf("%s is not installed. Run: autodev tools install %s", agent.Name, agent.ID)
	}

	cmd := exec.Command(path, args...)
	cmd.Env = os.Environ()

	// PTY is a Unix feature; on Windows (or non-TTY stdin) use stdio passthrough.
	if runtime.GOOS == "windows" || !term.IsTerminal(int(os.Stdin.Fd())) {
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		return cmd.Run()
	}

	ptmx, err := pty.Start(cmd)
	if err != nil {
		return fmt.Errorf("failed to allocate pty: %w", err)
	}
	defer func() { _ = ptmx.Close() }()

	// Put the parent terminal in raw mode so arrow keys etc. pass through.
	state, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err == nil {
		defer func() { _ = term.Restore(int(os.Stdin.Fd()), state) }()
	}

	// Forward window resize events to the child pty.
	ch := make(chan os.Signal, 1)
	notifyResize(ch)
	defer stopResize(ch)
	go func() {
		for range ch {
			_ = pty.InheritSize(os.Stdin, ptmx)
		}
	}()
	ch <- resizeSignal()

	go func() { _, _ = io.Copy(ptmx, os.Stdin) }()
	_, _ = io.Copy(os.Stdout, ptmx)

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("%s exited: %w", agent.ID, err)
	}
	return nil
}
