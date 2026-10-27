package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Galaxy renders a wordmark with a blue-family gradient — deep space blue into
// light starlight blue, a single hue family so views keep their one-accent
// rule. With --no-color the lipgloss renderer flattens it to plain text.
//
// Used by both TUI surfaces and the CLI banner (cmd/theme.go delegates here),
// so the gradient lives in exactly one place.
func Galaxy(s string) string {
	runes := []rune(s)
	if len(runes) == 0 {
		return s
	}
	const (
		deep     = "#003d91"
		bright   = "#7fc0ff"
		minAlpha = 0.35 // keep the tail readable instead of washed out
	)
	var b strings.Builder
	for i, r := range runes {
		if r == '\n' {
			b.WriteRune(r)
			continue
		}
		t := 0.0
		if len(runes) > 1 {
			t = float64(i) / float64(len(runes)-1)
		}
		t = minAlpha + t*(1-minAlpha)
		b.WriteString(lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(lerpHex(deep, bright, t))).
			Render(string(r)))
	}
	return b.String()
}

// GalaxyBlock renders multi-line ASCII art by applying the gradient across
// each line (spaces pass through unstyled).
func GalaxyBlock(art string) string {
	lines := strings.Split(art, "\n")
	for i, line := range lines {
		lines[i] = Galaxy(line)
	}
	return strings.Join(lines, "\n")
}

// hexToRGB parses "#RRGGBB" into 0-255 components.
func hexToRGB(hex string) (r, g, b int) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return 0, 0, 0
	}
	for i := 0; i < 6; i += 2 {
		c := 0
		for j := 0; j < 2; j++ {
			d := hex[i+j]
			var n int
			switch {
			case d >= '0' && d <= '9':
				n = int(d - '0')
			case d >= 'a' && d <= 'f':
				n = int(d-'a') + 10
			case d >= 'A' && d <= 'F':
				n = int(d-'A') + 10
			default:
				return 0, 0, 0
			}
			c = c*16 + n
		}
		switch i {
		case 0:
			r = c
		case 2:
			g = c
		case 4:
			b = c
		}
	}
	return r, g, b
}

// lerpHex blends two hex colors, t in [0,1].
func lerpHex(from, to string, t float64) string {
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	r1, g1, b1 := hexToRGB(from)
	r2, g2, b2 := hexToRGB(to)
	mix := func(a, b int) int { return a + int(float64(b-a)*t+.5) }
	return fmt.Sprintf("#%02x%02x%02x", mix(r1, r2), mix(g1, g2), mix(b1, b2))
}
