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

func newCanvasCmd() *cobra.Command {
	var path string
	var save bool
	var summary bool

	cmd := &cobra.Command{
		Use:   "canvas [path]",
		Aliases: []string{"cw", "knowledge"},
		Short: "Generate a Developer Knowledge Canvas for the project",
		Long: `Generate a comprehensive Developer Knowledge Canvas that captures:
file relationships, architecture, dependencies, API routes, database models,
security findings, git history, and AI-generated project context.

This canvas is saved to .autodevs/canvas/canvas.json and is designed to be
read by AI agents for rapid project understanding — reducing token consumption
and improving context accuracy.

Use --save to persist the canvas for future AI sessions.
Use --summary to print a quick AI-readable project summary.
Use --json to output raw JSON.`,
		Example: `  autodev canvas
  autodev canvas ./my-project
  autodev canvas --save
  autodev canvas --summary
  autodev canvas --json`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				path = args[0]
			}
			if path == "" {
				path = "."
			}
			return runCanvas(path, save, summary)
		},
	}

	cmd.Flags().StringVarP(&path, "path", "p", ".", "path to analyze")
	cmd.Flags().BoolVarP(&save, "save", "s", false, "save canvas to .autodevs/canvas/canvas.json")
	cmd.Flags().BoolVar(&summary, "summary", false, "print AI-readable project summary")
	return cmd
}

func runCanvas(path string, save, summary bool) error {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
	successStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FF87"))

	if !summary && !jsonOut {
		fmt.Println()
		fmt.Println(titleStyle.Render("  🎨 Developer Knowledge Canvas"))
		fmt.Println(dimStyle.Render(fmt.Sprintf("  Analyzing: %s", path)))
		fmt.Println()
	}

	gen := scanner.NewCanvasGenerator(path)
	canvas, err := gen.Generate()
	if err != nil {
		return fmt.Errorf("canvas generation failed: %w", err)
	}

	if save {
		if err := gen.Save(canvas); err != nil {
			return fmt.Errorf("save canvas: %w", err)
		}
		if !summary && !jsonOut {
			savePath := filepath.Join(path, ".autodevs", "canvas", "canvas.json")
			fmt.Println(successStyle.Render(fmt.Sprintf("  ✓ Canvas saved to: %s", savePath)))
			fmt.Println()
		}
	}

	if jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(canvas)
	}

	if summary {
		printCanvasSummary(canvas)
		return nil
	}

	printCanvasReport(canvas)
	return nil
}

func printCanvasSummary(canvas *scanner.Canvas) {
	ai := canvas.AI

	// Compact one-liner
	fmt.Println(ai.OneLiner)
	fmt.Println()

	// Key info
	if ai.ArchitectureDescription != "" {
		fmt.Println(ai.ArchitectureDescription)
		fmt.Println()
	}

	// Critical files
	if len(ai.CriticalFiles) > 0 {
		fmt.Println("CRITICAL FILES:")
		for _, f := range ai.CriticalFiles {
			fmt.Printf("  %s\n", f)
		}
		fmt.Println()
	}

	// Key insights
	if len(ai.KeyInsights) > 0 {
		fmt.Println("KEY INSIGHTS:")
		for _, i := range ai.KeyInsights {
			fmt.Printf("  - %s\n", i)
		}
		fmt.Println()
	}

	// Patterns
	if len(ai.Patterns) > 0 {
		fmt.Println("PATTERNS:")
		for _, p := range ai.Patterns {
			fmt.Printf("  - %s\n", p)
		}
		fmt.Println()
	}

	// Gotchas
	if len(ai.Gotchas) > 0 {
		fmt.Println("GOTCHAS:")
		for _, g := range ai.Gotchas {
			fmt.Printf("  ⚠ %s\n", g)
		}
		fmt.Println()
	}

	// Reading order
	if len(ai.RecommendedReading) > 0 {
		fmt.Println("RECOMMENDED READING ORDER:")
		for i, r := range ai.RecommendedReading {
			fmt.Printf("  %d. %s\n", i+1, r)
		}
	}
}

func printCanvasReport(canvas *scanner.Canvas) {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700"))
	sectionStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#4A90E2"))
	itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).PaddingLeft(2)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
	dangerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B"))

	fmt.Println(titleStyle.Render("  📊 Project Overview"))
	fmt.Printf("  %s\n", canvas.AI.OneLiner)
	fmt.Printf("  Files: %d | LOC: %d | Languages: %s\n",
		canvas.Project.TotalFiles, canvas.Project.TotalLOC,
		strings.Join(canvas.Project.Languages, ", "))
	fmt.Println()

	// Architecture
	if len(canvas.Architecture.Layers) > 0 {
		fmt.Println(sectionStyle.Render("  🏗️  Architecture Layers"))
		for _, l := range canvas.Architecture.Layers {
			if len(l.Paths) > 0 {
				fmt.Printf("    %s (%d paths)\n", l.Name, len(l.Paths))
			}
		}
		fmt.Println()
	}

	// Entry Points
	if len(canvas.Architecture.EntryPoints) > 0 {
		fmt.Println(sectionStyle.Render("  🚀 Entry Points"))
		for _, ep := range canvas.Architecture.EntryPoints {
			fmt.Println(itemStyle.Render(fmt.Sprintf("  %s [%s]", ep.Path, ep.Type)))
		}
		fmt.Println()
	}

	// Files by type
	if len(canvas.Files) > 0 {
		fmt.Println(sectionStyle.Render(fmt.Sprintf("  📁 Files (%d total)", len(canvas.Files))))
		typeCount := make(map[string]int)
		for _, f := range canvas.Files {
			typeCount[f.FileType]++
		}
		for ft, count := range typeCount {
			fmt.Printf("    %-15s %d\n", ft, count)
		}
		fmt.Println()
	}

	// API Routes
	if len(canvas.API.Routes) > 0 {
		fmt.Println(sectionStyle.Render(fmt.Sprintf("  🌐 API Routes (%d)", len(canvas.API.Routes))))
		for _, r := range canvas.API.Routes[:min(len(canvas.API.Routes), 15)] {
			fmt.Printf("    %-7s %-30s %s\n", r.Method, r.Path, dimStyle.Render(r.File))
		}
		if len(canvas.API.Routes) > 15 {
			fmt.Printf("    ... and %d more\n", len(canvas.API.Routes)-15)
		}
		fmt.Println()
	}

	// Database
	if len(canvas.Database.Models) > 0 {
		fmt.Println(sectionStyle.Render(fmt.Sprintf("  🗄️  Database Models (%d)", len(canvas.Database.Models))))
		for _, m := range canvas.Database.Models {
			fmt.Println(itemStyle.Render(fmt.Sprintf("  %s (%s)", m.Name, m.File)))
		}
		fmt.Println()
	}

	// Security
	if canvas.Security.CriticalCnt > 0 || canvas.Security.HighCnt > 0 {
		fmt.Println(sectionStyle.Render("  🔒 Security"))
		fmt.Printf("    Score: %d/100\n", canvas.Security.Score)
		fmt.Printf("    Critical: %d | High: %d\n", canvas.Security.CriticalCnt, canvas.Security.HighCnt)
		for _, f := range canvas.Security.Findings[:min(len(canvas.Security.Findings), 5)] {
			fmt.Println(dangerStyle.Render(fmt.Sprintf("    [%s] %s:%d — %s", f.Severity, f.File, f.Line, f.Message)))
		}
		if len(canvas.Security.Findings) > 5 {
			fmt.Printf("    ... and %d more findings\n", len(canvas.Security.Findings)-5)
		}
		fmt.Println()
	}

	// Dependencies
	if len(canvas.Dependencies.DeadFiles) > 0 {
		fmt.Println(sectionStyle.Render("  🗑️  Dead Code"))
		for _, df := range canvas.Dependencies.DeadFiles[:min(len(canvas.Dependencies.DeadFiles), 5)] {
			fmt.Println(dangerStyle.Render(fmt.Sprintf("    %s", df)))
		}
		fmt.Println()
	}

	// AI Context
	fmt.Println(sectionStyle.Render("  🤖 AI Context"))
	fmt.Printf("    %s\n", canvas.AI.Summary)
	if len(canvas.AI.KeyInsights) > 0 {
		for _, i := range canvas.AI.KeyInsights {
			fmt.Println(dimStyle.Render(fmt.Sprintf("    • %s", i)))
		}
	}
	fmt.Println()
	fmt.Println(dimStyle.Render("  Use --save to persist canvas for AI sessions."))
	fmt.Println(dimStyle.Render("  Use --summary for compact AI-readable output."))
	fmt.Println()
}
