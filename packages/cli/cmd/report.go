package cmd

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/autodev-sh/autodev/scanner"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func newReportCmd() *cobra.Command {
	var path string
	var format string
	var outDir string

	cmd := &cobra.Command{
		Use:   "report [path]",
		Short: "Generate a detailed environment report (HTML, JSON, Markdown)",
		Long:  `Scan the repository and generate a comprehensive environment report in the specified format.`,
		Example: `  autodev report
  autodev report --format html
  autodev report --format json --output ./reports
  autodev report --format markdown`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				path = args[0]
			}
			return runReport(path, format, outDir)
		},
	}

	cmd.Flags().StringVarP(&path, "path", "p", ".", "path to scan")
	cmd.Flags().StringVarP(&format, "format", "f", "html", "output format: html|json|markdown")
	cmd.Flags().StringVarP(&outDir, "output", "o", "./autodev-reports", "output directory")
	return cmd
}

func runReport(path, format, outDir string) error {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#0070f3"))
	okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00c853")).Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))

	fmt.Println()
	fmt.Println(titleStyle.Render("  Generating Report..."))

	s := scanner.New(path)
	result, err := s.Scan()
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	var outFile string

	switch format {
	case "json":
		outFile = filepath.Join(outDir, fmt.Sprintf("autodev-report-%s.json", timestamp))
		if err := writeJSONReport(result, outFile); err != nil {
			return err
		}
	case "markdown", "md":
		outFile = filepath.Join(outDir, fmt.Sprintf("autodev-report-%s.md", timestamp))
		if err := writeMarkdownReport(result, outFile); err != nil {
			return err
		}
	default: // html
		outFile = filepath.Join(outDir, fmt.Sprintf("autodev-report-%s.html", timestamp))
		if err := writeHTMLReport(result, outFile); err != nil {
			return err
		}
	}

	fmt.Println(okStyle.Render(fmt.Sprintf("  ✓ Report saved: %s", outFile)))
	fmt.Println(dimStyle.Render(fmt.Sprintf("  %d technologies detected", len(result.Technologies))))
	fmt.Println()
	return nil
}

func writeJSONReport(result *scanner.ScanResult, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(result); err != nil {
		return err
	}
	return f.Close()
}

func writeMarkdownReport(result *scanner.ScanResult, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	now := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(f, "# AutoDev Environment Report\n\n")
	fmt.Fprintf(f, "_Generated: %s_\n\n", now)
	fmt.Fprintf(f, "**Repository:** `%s`\n\n", result.Path)
	fmt.Fprintf(f, "---\n\n")

	if len(result.Languages) > 0 {
		fmt.Fprintf(f, "## Languages\n\n")
		for _, l := range result.Languages {
			fmt.Fprintf(f, "- %s\n", l)
		}
		fmt.Fprintf(f, "\n")
	}

	if len(result.Frameworks) > 0 {
		fmt.Fprintf(f, "## Frameworks\n\n")
		for _, fw := range result.Frameworks {
			fmt.Fprintf(f, "- %s\n", fw)
		}
		fmt.Fprintf(f, "\n")
	}

	if len(result.PackageManagers) > 0 {
		fmt.Fprintf(f, "## Package Managers\n\n")
		for _, pm := range result.PackageManagers {
			fmt.Fprintf(f, "- %s\n", pm)
		}
		fmt.Fprintf(f, "\n")
	}

	if len(result.Databases) > 0 {
		fmt.Fprintf(f, "## Databases & Services\n\n")
		for _, db := range result.Databases {
			fmt.Fprintf(f, "- %s\n", db)
		}
		fmt.Fprintf(f, "\n")
	}

	fmt.Fprintf(f, "## Setup Plan\n\n")
	for i, step := range result.RecommendedSetup {
		fmt.Fprintf(f, "%d. `%s`\n", i+1, step)
	}

	return f.Close()
}

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>AutoDev Report</title>
<style>
  @import url('https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@400;600;700&family=JetBrains+Mono:wght@400;600&display=swap');
  *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
  body { font-family: 'Space Grotesk', sans-serif; background: #0a0a0a; color: #e8e8e8; padding: 2rem; min-height: 100vh; }
  .container { max-width: 900px; margin: 0 auto; }
  header { border: 3px solid #0070f3; padding: 2rem; margin-bottom: 2rem; box-shadow: 6px 6px 0 #0070f3; }
  h1 { font-size: 2.5rem; font-weight: 700; color: #0070f3; }
  .meta { color: #888; margin-top: 0.5rem; font-family: 'JetBrains Mono', monospace; font-size: 0.85rem; }
  .section { border: 2px solid #333; padding: 1.5rem; margin-bottom: 1.5rem; }
  .section h2 { font-size: 1.2rem; font-weight: 700; color: #0070f3; margin-bottom: 1rem; text-transform: uppercase; letter-spacing: 0.1em; }
  .tag { display: inline-block; background: #1a1a1a; border: 2px solid #0070f3; color: #0070f3; padding: 0.25rem 0.75rem; margin: 0.25rem; font-weight: 600; font-size: 0.9rem; }
  .tag.green { border-color: #00c853; color: #00c853; }
  .tag.blue { border-color: #0070f3; color: #0070f3; }
  .tag.red { border-color: #ee0000; color: #ee0000; }
  .setup-step { font-family: 'JetBrains Mono', monospace; background: #111; border-left: 4px solid #0070f3; padding: 0.75rem 1rem; margin: 0.5rem 0; font-size: 0.9rem; }
  .badge { display: inline-block; background: #0070f3; color: #000; font-weight: 700; padding: 0.1rem 0.5rem; font-size: 0.8rem; margin-left: 0.5rem; }
  footer { text-align: center; color: #555; margin-top: 3rem; font-size: 0.85rem; }
</style>
</head>
<body>
<div class="container">
  <header>
    <h1>AutoDev Report</h1>
    <p class="meta">Generated: {{.Generated}} · Repository: {{.Path}}</p>
    <p class="meta">{{.Count}} technologies detected</p>
  </header>

  {{if .Languages}}
  <div class="section">
    <h2>Languages</h2>
    {{range .Languages}}<span class="tag green">{{.}}</span>{{end}}
  </div>
  {{end}}

  {{if .Frameworks}}
  <div class="section">
    <h2>Frameworks</h2>
    {{range .Frameworks}}<span class="tag">{{.}}</span>{{end}}
  </div>
  {{end}}

  {{if .PackageManagers}}
  <div class="section">
    <h2>Package Managers</h2>
    {{range .PackageManagers}}<span class="tag blue">{{.}}</span>{{end}}
  </div>
  {{end}}

  {{if .Databases}}
  <div class="section">
    <h2>Databases & Services</h2>
    {{range .Databases}}<span class="tag red">{{.}}</span>{{end}}
  </div>
  {{end}}

  {{if .Infra}}
  <div class="section">
    <h2>Infrastructure</h2>
    {{range .Infra}}<span class="tag">{{.}}</span>{{end}}
  </div>
  {{end}}

  {{if .SetupPlan}}
  <div class="section">
    <h2>Setup Plan</h2>
    {{range $i, $step := .SetupPlan}}
    <div class="setup-step"><span class="badge">{{add $i 1}}</span> {{$step}}</div>
    {{end}}
  </div>
  {{end}}

  <footer>Generated by AutoDev · https://github.com/HEETMEHTA18/autodev</footer>
</div>
</body>
</html>`

type htmlReportData struct {
	Generated       string
	Path            string
	Count           int
	Languages       []string
	Frameworks      []string
	PackageManagers []string
	Databases       []string
	Infra           []string
	SetupPlan       []string
}

func writeHTMLReport(result *scanner.ScanResult, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
	}
	tmpl, err := template.New("report").Funcs(funcMap).Parse(htmlTemplate)
	if err != nil {
		return err
	}

	data := htmlReportData{
		Generated:       time.Now().Format("2006-01-02 15:04:05"),
		Path:            result.Path,
		Count:           len(result.Technologies),
		Languages:       result.Languages,
		Frameworks:      result.Frameworks,
		PackageManagers: result.PackageManagers,
		Databases:       result.Databases,
		Infra:           result.Infra,
		SetupPlan:       result.RecommendedSetup,
	}

	if err := tmpl.Execute(f, data); err != nil {
		return err
	}

	return f.Close()
}
