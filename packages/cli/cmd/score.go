package cmd

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// scoreOptions holds the score command flags.
type scoreOptions struct {
	badge bool // print only the shields.io README markdown
	min   int  // fail (exit 1) when the score is below this threshold; 0 disables
}

func newScoreCmd() *cobra.Command {
	var opts scoreOptions
	cmd := &cobra.Command{
		Use:   "score",
		Short: "Rate your codebase and get a shareable README badge",
		Long: `Score the codebase from 0 to 100 using weighted doctor diagnostics —
secrets exposure and supply-chain health weigh double, hygiene checks weigh
in once — then map the result to a letter grade (A–F).

Prints a shields.io markdown badge you can paste straight into your README,
and --min turns it into a CI gate (exit 1 below the threshold).`,
		Example: `  autodev score
  autodev score --badge
  autodev score --json
  autodev score --min 80`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runScore(opts)
		},
	}
	cmd.Flags().BoolVar(&opts.badge, "badge", false,
		"print only the shields.io badge markdown (for piping into a README)")
	cmd.Flags().IntVar(&opts.min, "min", 0,
		"exit non-zero when the score is below this value (CI gate)")
	return cmd
}

// scoreWeights assigns relative importance per diagnostic check. Unknown
// checks default to 1 so new doctor checks are scored without edits here.
var scoreWeights = map[string]int{
	"Exposed Secrets Scanner":        3, // security-critical
	"Supply-Chain Vulnerabilities":   3, // security-critical
	"Git Configuration (.gitignore)": 2,
	"Environment Config (.env)":      2,
	"Dependency Lockfiles":           2,
	"Linter & Code Format Status":    1,
	"Environment Lockfile Mismatch":  1,
}

// scorePoints converts a check status into earned points per weight unit.
func scorePoints(status string) float64 {
	switch status {
	case "ok":
		return 1
	case "warning":
		return 0.5
	default: // error
		return 0
	}
}

// computeScore aggregates weighted check statuses into a 0-100 score.
func computeScore(results []checkResult) int {
	var earned, total float64
	for _, r := range results {
		w := float64(scoreWeights[r.Name])
		if w == 0 {
			w = 1
		}
		total += w
		earned += w * scorePoints(r.Status)
	}
	if total == 0 {
		return 0
	}
	return int(math.Round(100 * earned / total))
}

// scoreGrade maps a score to a letter grade.
func scoreGrade(score int) string {
	switch {
	case score >= 90:
		return "A"
	case score >= 75:
		return "B"
	case score >= 60:
		return "C"
	case score >= 40:
		return "D"
	default:
		return "F"
	}
}

// scoreBadgeColor picks the shields.io color for a grade — the same
// success/accent/warning/error semantic palette as the CLI.
func scoreBadgeColor(grade string) string {
	switch grade {
	case "A":
		return "00c853"
	case "B":
		return "0070f3"
	case "C", "D":
		return "f5a623"
	default:
		return "ee0000"
	}
}

// scoreBadgeURL builds the shields.io static badge endpoint for a score.
func scoreBadgeURL(score int, grade string) string {
	message := fmt.Sprintf("%d%%2F100_%s", score, grade) // "87/100 B"
	return "https://img.shields.io/badge/AutoDev_Score-" + message + "-" + scoreBadgeColor(grade)
}

// scoreBadgeMarkdown renders the paste-ready README image link.
func scoreBadgeMarkdown(score int, grade string) string {
	return fmt.Sprintf("[![AutoDev Score](%s)](https://github.com/HEETMEHTA18/autodev)",
		scoreBadgeURL(score, grade))
}

type scoreSummary struct {
	Passed   int `json:"passed"`
	Warnings int `json:"warnings"`
	Errors   int `json:"errors"`
}

type scoreReport struct {
	Version     string        `json:"version"`
	Path        string        `json:"path"`
	GeneratedAt string        `json:"generated_at"`
	Score       int           `json:"score"`
	Grade       string        `json:"grade"`
	Summary     scoreSummary  `json:"summary"`
	Badge       string        `json:"badge"`
	Checks      []checkResult `json:"checks"`
}

func summarizeChecks(results []checkResult) scoreSummary {
	var s scoreSummary
	for _, r := range results {
		switch r.Status {
		case "ok":
			s.Passed++
		case "warning":
			s.Warnings++
		default:
			s.Errors++
		}
	}
	return s
}

func runScore(opts scoreOptions) error {
	path := "."
	results, _ := evaluateDiagnostics(path, nil)
	score := computeScore(results)
	grade := scoreGrade(score)
	summary := summarizeChecks(results)
	badge := scoreBadgeMarkdown(score, grade)

	if jsonOut {
		report := scoreReport{
			Version:     rootCmd.Version,
			Path:        path,
			GeneratedAt: time.Now().UTC().Format(time.RFC3339),
			Score:       score,
			Grade:       grade,
			Summary:     summary,
			Badge:       badge,
			Checks:      results,
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(report); err != nil {
			return err
		}
		return checkScoreMinimum(score, opts.min)
	}

	if opts.badge {
		fmt.Println(badge)
		return checkScoreMinimum(score, opts.min)
	}

	// ── Text view ────────────────────────────────────────────────────────────
	fmt.Println()
	fmt.Println(styleHeading.Render("  AutoDev Score"))
	fmt.Println()

	scoreColor := styleOK
	switch grade {
	case "C", "D":
		scoreColor = styleWarn
	case "F":
		scoreColor = styleErr
	}
	content := fmt.Sprintf("%s   %s\n%s",
		scoreColor.Render(fmt.Sprintf("%d/100", score)),
		styleH3.Render("grade "+grade),
		ProgressBar(score),
	)
	fmt.Println(Panel("Score", content))
	fmt.Println()

	for _, r := range results {
		printCheckResult(r)
	}
	fmt.Println()
	fmt.Println("  " + StatusRow(
		fmt.Sprintf("%d checks", len(results)),
		styleOK.Render(fmt.Sprintf("%d %s", summary.Passed, iconSuccess)),
		styleWarn.Render(fmt.Sprintf("%d %s", summary.Warnings, iconWarn)),
		styleErr.Render(fmt.Sprintf("%d %s", summary.Errors, iconError)),
		styleH3.Render("grade "+grade),
	))
	fmt.Println()
	fmt.Println(Caption("  README badge:"))
	fmt.Println("  " + badge)
	fmt.Println(Caption("  Run 'autodev doctor --fix' to improve your score."))

	return checkScoreMinimum(score, opts.min)
}

func checkScoreMinimum(score, min int) error {
	if min > 0 && score < min {
		return fmt.Errorf("score %d is below the required minimum %d", score, min)
	}
	return nil
}
