package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	gh "github.com/autodev-sh/autodev/github"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newGitHubCmd() *cobra.Command {
	var token string
	var setup bool

	cmd := &cobra.Command{
		Use:   "github <username>",
		Short: "Scan all public repositories for a GitHub user",
		Long: `Fetch all public repositories for a GitHub username, detect languages,
frameworks, and generate a recommended development environment setup plan.`,
		Example: `  autodev github HEETMEHTA18
  autodev github torvalds --json
  autodev github microsoft --setup`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			username := args[0]
			if token == "" {
				token = viper.GetString("github_token")
			}
			if token == "" {
				token = os.Getenv("GITHUB_TOKEN")
			}
			return runGitHub(username, token, setup)
		},
	}

	cmd.Flags().StringVarP(&token, "token", "t", "", "GitHub personal access token (or set AUTODEV_GITHUB_TOKEN)")
	cmd.Flags().BoolVarP(&setup, "setup", "s", false, "auto-install recommended environment after scan")
	return cmd
}

func runGitHub(username, token string, doSetup bool) error {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#0070f3"))
	okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00c853")).Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ee0000"))
	badgeStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#0070f3")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 1)

	if !jsonOut {
		fmt.Println()
		fmt.Println(titleStyle.Render(fmt.Sprintf("  GitHub Scanner — @%s", username)))
		fmt.Println(dimStyle.Render("  Fetching public repositories..."))
		fmt.Println()
	}

	client := gh.New(token)
	result, err := client.ScanUser(username)
	if err != nil {
		return fmt.Errorf("github scan failed: %w", err)
	}

	if jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}

	// Print results
	fmt.Printf("  %s  %s\n",
		okStyle.Render(fmt.Sprintf("@%s", result.Username)),
		dimStyle.Render(fmt.Sprintf("%d public repos analyzed", result.TotalRepos)),
	)
	fmt.Println()

	if len(result.Languages) > 0 {
		fmt.Println(titleStyle.Render("  [LANGUAGES USED]"))
		// Sort by count
		type langCount struct {
			lang  string
			count int
		}
		var langs []langCount
		for l, c := range result.Languages {
			langs = append(langs, langCount{l, c})
		}
		// Simple bubble sort for display
		for i := range langs {
			for j := i + 1; j < len(langs); j++ {
				if langs[j].count > langs[i].count {
					langs[i], langs[j] = langs[j], langs[i]
				}
			}
		}
		for _, lc := range langs {
			bar := buildBar(lc.count, result.TotalRepos, 20)
			fmt.Printf("  %-20s %s %s\n", lc.lang, bar, dimStyle.Render(fmt.Sprintf("%d repos", lc.count)))
		}
		fmt.Println()
	}

	if len(result.TopRepos) > 0 {
		fmt.Println(titleStyle.Render("  [TOP REPOSITORIES]"))
		for i, repo := range result.TopRepos {
			if i >= 5 {
				break
			}
			stars := fmt.Sprintf("Stars: %d", repo.StargazersCount)
			desc := repo.Description
			if len(desc) > 60 {
				desc = desc[:57] + "..."
			}
			fmt.Printf("  %s %-30s %s\n",
				badgeStyle.Render(repo.Language),
				repo.Name,
				dimStyle.Render(stars),
			)
			if desc != "" {
				fmt.Printf("    %s\n", dimStyle.Render(desc))
			}
		}
		fmt.Println()
	}

	if len(result.Recommended) > 0 {
		fmt.Println(titleStyle.Render("  [RECOMMENDED ENVIRONMENT]"))
		for _, env := range result.Recommended {
			fmt.Printf("  - %s\n", env)
		}
		fmt.Println()
	}

	if len(result.SkillGaps) > 0 {
		fmt.Println(titleStyle.Render("  [SUGGESTED SKILLS]"))
		for _, skill := range result.SkillGaps {
			fmt.Printf("  - %s\n", skill)
		}
		fmt.Println()
	}

	if doSetup {
		fmt.Println(titleStyle.Render("  Setting up recommended environment..."))
		fmt.Println(warnStyle.Render("  (run 'autodev setup' manually for fine-grained control)"))
	}

	return nil
}

func buildBar(count, total, width int) string {
	if total == 0 {
		return ""
	}
	filled := (count * width) / total
	bar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	return bar
}
