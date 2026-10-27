package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/autodev-sh/autodev/cli/tui"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"golang.org/x/term"
)

// Minimal — TUI Design System, blue / galaxy edition.
//
// Background #0a0a0a · Foreground #ededed · Primary #ffffff · Secondary #888888
// Accent #0070f3 · Success #00c853 · Warning #f5a623 · Error #ee0000
// Muted #555555 · Surface #1a1a1a · Border #2a2a2a
//
// Rules: single-line borders, bold-only emphasis, single-width icons (no
// emoji), at most one accent color per view.
const (
	paletteBg        = "#0a0a0a"
	paletteFg        = "#ededed"
	palettePrimary   = "#ffffff"
	paletteSecondary = "#888888"
	paletteAccent    = "#0070f3"
	paletteSuccess   = "#00c853"
	paletteWarning   = "#f5a623"
	paletteError     = "#ee0000"
	paletteMuted     = "#555555"
	paletteSurface   = "#1a1a1a"
	paletteBorder    = "#2a2a2a"
	paletteInk       = "#000000"
)

// Icon set — single-width characters only, never emoji.
const (
	iconSuccess = "✓"
	iconError   = "✗"
	iconWarn    = "!"
	iconInfo    = "·"
	iconPending = "○"
	iconRunning = "▶"
	iconArrow   = "→"
	iconBullet  = "·"
	iconSelect  = "▸"
)

var (
	// H2 section headers carry the galaxy identity: BOLD + Accent.
	styleHeading = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(paletteAccent))
	styleH2      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(paletteFg))
	styleH3      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(paletteSecondary))
	styleOK      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(paletteSuccess))
	styleErr     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(paletteError))
	styleWarn    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(paletteWarning))
	styleCyan    = lipgloss.NewStyle().Foreground(lipgloss.Color(paletteAccent))
	styleAccent  = lipgloss.NewStyle().Foreground(lipgloss.Color(paletteAccent))
	styleDim     = lipgloss.NewStyle().Foreground(lipgloss.Color(paletteSecondary))
	styleCaption = lipgloss.NewStyle().Foreground(lipgloss.Color(paletteMuted))
	styleBody    = lipgloss.NewStyle().Foreground(lipgloss.Color(paletteFg))
	styleLabel   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(paletteSecondary))
	styleBold    = lipgloss.NewStyle().Bold(true)
	styleBadge   = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(palettePrimary)).
			Background(lipgloss.Color(paletteAccent)).
			Padding(0, 1)
)

// applyColorPreference disables color when requested via the --no-color flag,
// the NO_COLOR environment variable, or the AUTODEV_NO_COLOR env var. It is
// called before any rendering so both help output and command output obey it.
func applyColorPreference(argv []string) {
	if os.Getenv("AUTODEV_NO_COLOR") != "" || os.Getenv("NO_COLOR") != "" {
		lipgloss.SetColorProfile(termenv.Ascii)
		return
	}
	if noColor {
		lipgloss.SetColorProfile(termenv.Ascii)
		return
	}
	for _, arg := range argv {
		if arg == "--no-color" || arg == "--no-color=true" {
			lipgloss.SetColorProfile(termenv.Ascii)
			return
		}
	}
}

// isInteractive reports whether both stdin and stdout are terminals.
func isInteractive() bool {
	return term.IsTerminal(int(os.Stdin.Fd())) && term.IsTerminal(int(os.Stdout.Fd()))
}

// printError writes a single styled error line to stderr. Cobra's own error
// printing is silenced (SilenceErrors) so users never see the message twice.
func printError(err error) {
	fmt.Fprintf(os.Stderr, "%s %s\n", styleErr.Render(iconError), err.Error())
	fmt.Fprintf(os.Stderr, "  %s\n", styleCaption.Render("Run 'autodev --help' for usage."))
}

// ── Minimal design-system helpers ────────────────────────────────────────────

// H2 renders a section header (BOLD + Foreground).
func H2(s string) string { return styleH2.Render(s) }

// H3 renders a subsection header (BOLD + Secondary).
func H3(s string) string { return styleH3.Render(s) }

// Caption renders muted help text.
func Caption(s string) string { return styleCaption.Render(s) }

// Label renders a bold secondary label.
func Label(s string) string { return styleLabel.Render(s) }

// ListItem renders a list row: selected rows get "▸" + BOLD, the rest a
// plain 4-space indent.
func ListItem(s string, selected bool) string {
	if selected {
		return styleH2.Render(iconSelect + " " + s)
	}
	return "    " + styleBody.Render(s)
}

// Panel renders content inside a single-line box with the title embedded in
// the top border, 1-space padding, and a muted secondary border.
//
//	┌─ Deploy Status ──────────┐
//	│  ✓ Ready                 │
//	└──────────────────────────┘
func Panel(title, content string) string {
	width := 0
	for _, line := range strings.Split(content, "\n") {
		if w := lipgloss.Width(line); w > width {
			width = w
		}
	}
	inner := width + 2 // 1-space padding each side
	top := "┌─ " + title + " "
	fill := inner - lipgloss.Width(top)
	if fill < 1 {
		fill = 1
	}
	top += strings.Repeat("─", fill) + "┐"

	border := lipgloss.NewStyle().Foreground(lipgloss.Color(paletteMuted))
	var body []string
	for _, line := range strings.Split(content, "\n") {
		body = append(body, "│ "+line+strings.Repeat(" ", inner-lipgloss.Width(line)-1)+"│")
	}
	bottom := "└" + strings.Repeat("─", inner) + "┘"

	out := []string{top}
	out = append(out, body...)
	out = append(out, bottom)
	return border.Render(strings.Join(out, "\n"))
}

// StatusRow renders a single-line status bar whose items are separated by
// " ─ " — the Minimal status-bar pattern.
func StatusRow(items ...string) string {
	sep := styleCaption.Render(" ─ ")
	joined := make([]string, 0, len(items)*2-1)
	for i, it := range items {
		if i > 0 {
			joined = append(joined, sep)
		}
		joined = append(joined, it)
	}
	return strings.Join(joined, "")
}

// ProgressBar renders ▕████░░▏ <pct>% — accent-filled, dim-empty.
func ProgressBar(pct int) string {
	if pct < 0 {
		pct = 0
	}
	if pct > 100 {
		pct = 100
	}
	const width = 24
	filled := width * pct / 100
	bar := "▕" +
		styleAccent.Render(strings.Repeat("█", filled)) +
		styleCaption.Render(strings.Repeat("░", width-filled)) +
		"▏"
	return fmt.Sprintf("%s %d%%", bar, pct)
}

// ── Galaxy wordmark ──────────────────────────────────────────────────────────

// Galaxy renders the wordmark with a blue-family gradient — deep space blue
// into light starlight blue, a single hue family that keeps the "one accent"
// rule. With --no-color the renderer flattens it to plain text. The gradient
// itself lives in package tui so the CLI banner and TUI headers share it.
func Galaxy(s string) string { return tui.Galaxy(s) }
