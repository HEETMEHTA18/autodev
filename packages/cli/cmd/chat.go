package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/autodev-sh/autodev/core/promptcapture"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func newChatCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chat",
		Short: "Start an interactive prompt session to capture and track prompts",
		Long: `Start an interactive AI coding session. AutoDev will capture all your prompts,
responses, commands run, and files modified, saving them to .autodevs/sessions/ and prompts.md.`,
		Example: `  autodev chat`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runChat()
		},
	}
	return cmd
}

func runChat() error {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700"))
	accentStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FF87"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	warnStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF8700"))
	cmdStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF"))

	engine, err := promptcapture.NewEngine("")
	if err != nil {
		return fmt.Errorf("failed to start capture engine: %w", err)
	}

	err = engine.StartSession()
	if err != nil {
		return fmt.Errorf("failed to initialize session: %w", err)
	}
	defer func() { _ = engine.EndSession() }()

	fmt.Println()
	fmt.Println(titleStyle.Render("  ⚡ AutoDev Prompt Capture Chat Engine v0.3.2"))
	fmt.Println(dimStyle.Render(fmt.Sprintf("  Session %s initialized.", engine.Session.SessionID)))
	fmt.Println(dimStyle.Render("  Capturing prompts, files generated, and commands executed."))

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		fmt.Println()
		fmt.Println(warnStyle.Render("  ⚠️  GEMINI_API_KEY not set. Running in simulated developer mode."))
		fmt.Println(dimStyle.Render("     Set GEMINI_API_KEY to query real AI models."))
	} else {
		fmt.Println()
		fmt.Println(accentStyle.Render("  ✓ Real Gemini AI integration active."))
	}
	fmt.Println(dimStyle.Render("  Type 'exit' or 'quit' to end session."))
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(accentStyle.Render("agy> "))
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		if input == "exit" || input == "quit" {
			fmt.Println(dimStyle.Render("  Ending session and exporting logs..."))
			break
		}

		var responseText string
		var queryErr error

		if apiKey != "" {
			// Query Gemini
			responseText, queryErr = queryGemini(input, engine.Session.Events)
			if queryErr != nil {
				fmt.Println(warnStyle.Render(fmt.Sprintf("  Error querying Gemini: %v. Falling back to simulation.", queryErr)))
				responseText = getSimulatedResponse(input)
			}
		} else {
			responseText = getSimulatedResponse(input)
		}

		// Display response
		fmt.Println()
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#F0F0F0")).Render(responseText))
		fmt.Println()

		// Parse commands and files suggested in response
		suggestedCommands, suggestedFiles := parseActions(responseText)

		var runCmds []promptcapture.ExecutedCommand
		var genFiles []promptcapture.GeneratedFile

		// Handle suggested files
		for _, f := range suggestedFiles {
			fmt.Printf("  Create/Write file %s? [y/N] ", cmdStyle.Render(f.Path))
			if askConfirm() {
				// Make sure dir exists
				dir := filepath.Dir(f.Path)
				if dir != "." {
					_ = os.MkdirAll(filepath.Join(engine.Root, dir), 0755)
				}
				err := os.WriteFile(filepath.Join(engine.Root, f.Path), []byte(f.Content), 0644)
				if err != nil {
					fmt.Printf("    %s Failed to write file: %v\n", warnStyle.Render("✗"), err)
				} else {
					fmt.Printf("    %s Saved %s (%d bytes)\n", accentStyle.Render("✓"), f.Path, len(f.Content))
					genFiles = append(genFiles, promptcapture.GeneratedFile{
						FilePath:  f.Path,
						SizeBytes: int64(len(f.Content)),
						Action:    "created",
						Timestamp: time.Now(),
					})
				}
			}
		}

		// Handle suggested commands
		for _, c := range suggestedCommands {
			fmt.Printf("  Run command: %s? [y/N] ", cmdStyle.Render(c))
			if askConfirm() {
				fmt.Printf("    Running `%s`...\n", c)
				execCmd, err := runCommandAndCapture(engine.Root, c)
				if err != nil {
					fmt.Printf("    %s Run failed: %v\n", warnStyle.Render("✗"), err)
				} else {
					fmt.Printf("    %s Finished (Exit Code: %d, Duration: %.2fs)\n",
						accentStyle.Render("✓"), execCmd.ExitCode, float64(execCmd.DurationMs)/1000.0)
					runCmds = append(runCmds, execCmd)
				}
			}
		}

		// Add event to capture engine
		err = engine.AddEvent(input, responseText, runCmds, genFiles)
		if err != nil {
			fmt.Println(warnStyle.Render(fmt.Sprintf("  Warning: failed to log prompt event: %v", err)))
		} else {
			fmt.Println(dimStyle.Render("  [captured prompt, files, and commands]"))
		}
		fmt.Println()
	}

	return nil
}

func askConfirm() bool {
	var ans string
	_, _ = fmt.Scanln(&ans)
	ans = strings.ToLower(strings.TrimSpace(ans))
	return ans == "y" || ans == "yes"
}

func parseActions(response string) ([]string, []struct{ Path, Content string }) {
	var commands []string
	var files []struct{ Path, Content string }

	// Parse RUN_COMMAND tags
	lines := strings.Split(response, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "[RUN_COMMAND:") && strings.HasSuffix(line, "]") {
			cmd := strings.TrimSuffix(strings.TrimPrefix(line, "[RUN_COMMAND:"), "]")
			commands = append(commands, cmd)
		}
	}

	// Parse WRITE_FILE tags
	startTag := "[WRITE_FILE:"
	endTag := "[/WRITE_FILE]"

	temp := response
	for {
		startIdx := strings.Index(temp, startTag)
		if startIdx == -1 {
			break
		}

		closeBracketIdx := strings.Index(temp[startIdx:], "]")
		if closeBracketIdx == -1 {
			break
		}
		closeBracketIdx += startIdx

		filePath := temp[startIdx+len(startTag) : closeBracketIdx]

		endIdx := strings.Index(temp[closeBracketIdx:], endTag)
		if endIdx == -1 {
			break
		}
		endIdx += closeBracketIdx

		fileContent := temp[closeBracketIdx+1 : endIdx]
		files = append(files, struct{ Path, Content string }{Path: filePath, Content: strings.TrimSpace(fileContent)})

		temp = temp[endIdx+len(endTag):]
	}

	return commands, files
}

func runCommandAndCapture(root, cmdStr string) (promptcapture.ExecutedCommand, error) {
	parts := strings.Fields(cmdStr)
	if len(parts) == 0 {
		return promptcapture.ExecutedCommand{}, fmt.Errorf("empty command")
	}
	command := parts[0]
	args := parts[1:]

	start := time.Now()
	cmd := exec.Command(command, args...)
	cmd.Dir = root

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	duration := time.Since(start)

	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			exitCode = -1
		}
	}

	return promptcapture.ExecutedCommand{
		Command:    command,
		Args:       args,
		ExitCode:   exitCode,
		Stdout:     stdout.String(),
		Stderr:     stderr.String(),
		DurationMs: duration.Milliseconds(),
		Timestamp:  start,
	}, nil
}

func getSimulatedResponse(prompt string) string {
	promptLower := strings.ToLower(prompt)

	if strings.Contains(promptLower, "react") || strings.Contains(promptLower, "dashboard") {
		return `I will create a standard React dashboard page and run a production build check.

[WRITE_FILE:src/components/Dashboard.tsx]
import React from 'react';

export default function Dashboard() {
  return (
    <div className="min-h-screen bg-slate-900 text-slate-100 p-8">
      <header className="mb-8">
        <h1 className="text-3xl font-bold text-emerald-400">AutoDev Dashboard</h1>
        <p className="text-slate-400">Bioluminescent Ocean Farming Control Panel</p>
      </header>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="bg-slate-800 p-6 rounded-lg border border-slate-700">
          <h2 className="text-lg font-medium text-slate-300">Sonar Sweeper</h2>
          <p className="text-2xl font-bold text-emerald-400 mt-2">Active</p>
        </div>
        <div className="bg-slate-800 p-6 rounded-lg border border-slate-700">
          <h2 className="text-lg font-medium text-slate-300">Ballast Pods</h2>
          <p className="text-2xl font-bold text-emerald-400 mt-2">74% Depth OK</p>
        </div>
        <div className="bg-slate-800 p-6 rounded-lg border border-slate-700">
          <h2 className="text-lg font-medium text-slate-300">Tomato Yield</h2>
          <p className="text-2xl font-bold text-amber-400 mt-2">250 kg/day</p>
        </div>
      </div>
    </div>
  );
}
[/WRITE_FILE]

To compile and verify the React setup, let's run the build check:
[RUN_COMMAND:npm run build]
`
	}

	if strings.Contains(promptLower, "docker") {
		return `I will generate a Dockerfile to package the Ocean Farming simulator app.

[WRITE_FILE:Dockerfile]
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o autodev-simulator .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/autodev-simulator .
EXPOSE 8080
CMD ["./autodev-simulator"]
[/WRITE_FILE]

Let's test building the image locally:
[RUN_COMMAND:docker build -t autodev-simulator-app .]
`
	}

	if strings.Contains(promptLower, "auth") {
		return `I will add basic JSON Web Token (JWT) authentication logic to protect the control panel APIs.

[WRITE_FILE:src/auth.ts]
export interface Session {
  userId: string;
  role: string;
  exp: number;
}

export function verifySession(token: string): Session | null {
  if (token === "mock-valid-token") {
    return {
      userId: "dolphin-node-1",
      role: "operator",
      exp: Date.now() + 3600000
    };
  }
  return nil;
}
[/WRITE_FILE]
`
	}

	return `I've received your prompt: "` + prompt + `".
Let's run a doctor toolchain health diagnostics scan to make sure the environment is fully compliant.

[RUN_COMMAND:autodev doctor]
`
}

func queryGemini(prompt string, history []promptcapture.PromptEvent) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY not set")
	}

	type Part struct {
		Text string `json:"text"`
	}
	type Content struct {
		Role  string `json:"role,omitempty"`
		Parts []Part `json:"parts"`
	}
	type Request struct {
		Contents          []Content `json:"contents"`
		SystemInstruction struct {
			Parts []Part `json:"parts"`
		} `json:"systemInstruction,omitempty"`
	}

	var contents []Content

	for _, ev := range history {
		contents = append(contents, Content{
			Role:  "user",
			Parts: []Part{{Text: ev.Prompt}},
		})
		contents = append(contents, Content{
			Role:  "model",
			Parts: []Part{{Text: ev.Response}},
		})
	}

	contents = append(contents, Content{
		Role:  "user",
		Parts: []Part{{Text: prompt}},
	})

	reqBody := Request{
		Contents: contents,
	}

	reqBody.SystemInstruction.Parts = []Part{{
		Text: "You are AutoDev CLI Assistant, a helper to the developer. Respond concisely. " +
			"You are a lazy senior dev. The best code is the code you never wrote. Follow these rules:\n" +
			"1. Does this need to exist? → no: skip it (YAGNI)\n" +
			"2. Stdlib does it? → use it\n" +
			"3. Native platform feature? → use it\n" +
			"4. Installed dependency? → use it\n" +
			"5. One line? → one line\n" +
			"6. Only then: the minimum that works\n\n" +
			"When you suggest running a command, output: [RUN_COMMAND:your command here]\n" +
			"When you suggest writing/creating/modifying a file, output:\n" +
			"[WRITE_FILE:path/to/file]\n" +
			"file contents here\n" +
			"[/WRITE_FILE]\n",
	}}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=" + apiKey
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Gemini API returned status %d", resp.StatusCode)
	}

	var geminiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return "", err
	}

	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		return geminiResp.Candidates[0].Content.Parts[0].Text, nil
	}

	return "", fmt.Errorf("no response candidate found")
}
