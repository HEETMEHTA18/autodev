package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

// chatOptions controls the AI behavior of the interactive chat session.
type chatOptions struct {
	thinkingBudget int  // cap on thinking tokens per turn (0 disables thinking)
	showThoughts   bool // display the model's thought summaries
	noRefine       bool // skip the two-phase prompt refinement step
	noGrounding    bool // disable Google Search grounding for answers
}

func newChatCmd() *cobra.Command {
	var opts chatOptions

	cmd := &cobra.Command{
		Use:   "chat",
		Short: "Start an interactive prompt session to capture and track prompts",
		Long: `Start an interactive AI coding session. AutoDev will capture all your prompts,
responses, commands run, and files modified, saving them to .autodevs/sessions/ and prompts.md.

Each message is first refined into a scoped, self-contained prompt (full context,
explicit direction), then answered with a capped thinking budget, a topic-scope
guard, and search grounding to prevent hallucinated answers.`,
		Example: `  autodev chat
  autodev chat --thinking-budget 2048 --show-thoughts
  autodev chat --no-refine --no-grounding`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runChat(opts)
		},
	}

	cmd.Flags().IntVar(&opts.thinkingBudget, "thinking-budget", 1024,
		"max thinking tokens per turn (0 disables thinking, max 24576)")
	cmd.Flags().BoolVar(&opts.showThoughts, "show-thoughts", false,
		"display the model's thought summaries")
	cmd.Flags().BoolVar(&opts.noRefine, "no-refine", false,
		"skip prompt refinement and answer the raw message directly")
	cmd.Flags().BoolVar(&opts.noGrounding, "no-grounding", false,
		"disable Google Search grounding for answers")
	return cmd
}

func runChat(opts chatOptions) error {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#0070f3"))
	accentStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00c853"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	warnStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#f5a623"))
	cmdStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#0070f3"))

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
	fmt.Println(titleStyle.Render("  AutoDev Prompt Capture Chat Engine v0.3.2"))
	fmt.Println(dimStyle.Render(fmt.Sprintf("  Session %s initialized.", engine.Session.SessionID)))
	fmt.Println(dimStyle.Render("  Capturing prompts, files generated, and commands executed."))

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		fmt.Println()
		fmt.Println(warnStyle.Render("  !  GEMINI_API_KEY not set. Running in simulated developer mode."))
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
			// Two-phase turn: refine the message, then answer with bounded thinking.
			responseText, queryErr = askGemini(input, engine.Session.Events, opts)
			if queryErr != nil {
				fmt.Println(warnStyle.Render(fmt.Sprintf("  Error querying Gemini: %v. Falling back to simulation.", queryErr)))
				responseText = getSimulatedResponse(input)
			}
		} else {
			responseText = getSimulatedResponse(input)
		}

		// Display response
		fmt.Println()
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#ededed")).Render(responseText))
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

// ── Gemini integration ──────────────────────────────────────────────────────

const (
	geminiModel       = "gemini-2.5-flash"
	maxThinkingBudget = 24576 // Gemini 2.5 Flash thinking limit
	refineBudget      = 512   // phase-1 refinement stays cheap and fast
	refineStartMarker = "<<REFINED>>"
	refineEndMarker   = "<<END_REFINED>>"
)

type geminiPart struct {
	Text    string `json:"text,omitempty"`
	Thought bool   `json:"thought,omitempty"`
}

type geminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []geminiPart `json:"parts"`
}

type geminiRequest struct {
	Contents          []geminiContent `json:"contents"`
	SystemInstruction struct {
		Parts []geminiPart `json:"parts"`
	} `json:"systemInstruction,omitempty"`
	GenerationConfig *generationConfig `json:"generationConfig,omitempty"`
	Tools            []map[string]any  `json:"tools,omitempty"`
}

type generationConfig struct {
	ThinkingConfig *thinkingConfig `json:"thinkingConfig,omitempty"`
}

type thinkingConfig struct {
	ThinkingBudget  int  `json:"thinkingBudget"`
	IncludeThoughts bool `json:"includeThoughts"`
}

// geminiReply is the normalized result of one generateContent call.
type geminiReply struct {
	Text     string
	Thoughts []string
	Sources  []string
}

func clampThinkingBudget(n int) int {
	if n < 0 {
		return 0
	}
	if n > maxThinkingBudget {
		return maxThinkingBudget
	}
	return n
}

// assistantSystemPrompt keeps the original lazy-senior-dev rules and adds
// bounded-thinking, topic-scope, and anti-hallucination discipline.
func assistantSystemPrompt() string {
	return "You are AutoDev CLI Assistant, a helper to the developer. Respond concisely. " +
		"You are a lazy senior dev. The best code is the code you never wrote. Follow these rules:\n" +
		"1. Does this need to exist? → no: skip it (YAGNI)\n" +
		"2. Stdlib does it? → use it\n" +
		"3. Native platform feature? → use it\n" +
		"4. Installed dependency? → use it\n" +
		"5. One line? → one line\n" +
		"6. Only then: the minimum that works\n\n" +
		"Scope discipline:\n" +
		"- Stay strictly on the topic of the user's question. Never drift into unrelated domains " +
		"(a finance question must not turn into medical, sports, or celebrity news).\n" +
		"- If the intent is ambiguous, ask ONE short clarifying question instead of guessing.\n" +
		"- For out-of-scope follow-ups, say they are out of scope and return to the task.\n\n" +
		"Anti-hallucination:\n" +
		"- Never invent facts, statistics, URLs, citations, quotes, or API names.\n" +
		"- If you are unsure or lack information, say so plainly instead of guessing.\n" +
		"- Separate verified facts from opinions; cite grounded sources when provided.\n\n" +
		"Bounded thinking:\n" +
		"- Think only as long as needed to choose the approach, then answer directly.\n" +
		"- Do not explore tangents or re-derive points already settled in this conversation.\n\n" +
		"When you suggest running a command, output: [RUN_COMMAND:your command here]\n" +
		"When you suggest writing/creating/modifying a file, output:\n" +
		"[WRITE_FILE:path/to/file]\n" +
		"file contents here\n" +
		"[/WRITE_FILE]\n"
}

// refineSystemPrompt instructs phase 1: rewrite the raw message into one
// scoped, self-contained prompt that carries the full context and direction.
func refineSystemPrompt() string {
	return "You are a prompt-refinement engine. You rewrite a user's raw message into ONE precise, " +
		"self-contained prompt for an AI coding assistant. Rules:\n" +
		"- Keep it strictly inside the original topic domain — never widen scope into unrelated subjects.\n" +
		"- Fold in the relevant earlier conversation context so the prompt stands alone.\n" +
		"- State the explicit constraints, the desired output, and the direction/next action.\n" +
		"- Instruct the assistant to answer only from verifiable knowledge and to admit when it does not know.\n" +
		"- Output ONLY the refined prompt wrapped in " + refineStartMarker + " and " + refineEndMarker +
		". No commentary, no explanations."
}

func buildContents(prompt string, history []promptcapture.PromptEvent) []geminiContent {
	var contents []geminiContent
	for _, ev := range history {
		if ev.Prompt == "" && ev.Response == "" {
			continue
		}
		contents = append(contents, geminiContent{
			Role:  "user",
			Parts: []geminiPart{{Text: ev.Prompt}},
		})
		contents = append(contents, geminiContent{
			Role:  "model",
			Parts: []geminiPart{{Text: ev.Response}},
		})
	}
	contents = append(contents, geminiContent{
		Role:  "user",
		Parts: []geminiPart{{Text: prompt}},
	})
	return contents
}

// askGemini runs one two-phase turn: refine the message into a scoped prompt
// (shown to the user), then answer with bounded thinking and grounding.
// Refinement failures fall back silently to the raw input.
func askGemini(input string, history []promptcapture.PromptEvent, opts chatOptions) (string, error) {
	prompt := input
	if !opts.noRefine {
		if refined, err := refinePrompt(input, history); err == nil {
			prompt = refined
			printRefinedPrompt(refined)
		}
	}

	reply, err := callGemini(
		buildContents(prompt, history),
		assistantSystemPrompt(),
		opts.thinkingBudget,
		opts.showThoughts,
		!opts.noGrounding,
	)
	if err != nil {
		return "", err
	}
	printThoughts(reply.Thoughts)
	printSources(reply.Sources)
	return reply.Text, nil
}

// refinePrompt runs phase 1: fast, tiny thinking budget, no grounding.
func refinePrompt(input string, history []promptcapture.PromptEvent) (string, error) {
	reply, err := callGemini(buildContents(input, history), refineSystemPrompt(), refineBudget, false, false)
	if err != nil {
		return "", err
	}
	return extractRefinedPrompt(reply.Text)
}

func extractRefinedPrompt(text string) (string, error) {
	start := strings.Index(text, refineStartMarker)
	if start < 0 {
		return "", fmt.Errorf("refinement start marker missing")
	}
	rest := text[start+len(refineStartMarker):]
	end := strings.Index(rest, refineEndMarker)
	if end < 0 {
		return "", fmt.Errorf("refinement end marker missing")
	}
	refined := strings.TrimSpace(rest[:end])
	if refined == "" {
		return "", fmt.Errorf("empty refined prompt")
	}
	return refined, nil
}

// callGemini performs one generateContent request. When Google Search
// grounding is rejected by the API, it retries once without tools.
func callGemini(contents []geminiContent, system string, budget int, thoughts, grounding bool) (geminiReply, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return geminiReply{}, fmt.Errorf("GEMINI_API_KEY not set")
	}

	var req geminiRequest
	req.Contents = contents
	req.SystemInstruction.Parts = []geminiPart{{Text: system}}
	req.GenerationConfig = &generationConfig{
		ThinkingConfig: &thinkingConfig{
			ThinkingBudget:  clampThinkingBudget(budget),
			IncludeThoughts: thoughts,
		},
	}
	if grounding {
		req.Tools = []map[string]any{{"google_search": map[string]any{}}}
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return geminiReply{}, err
	}

	url := "https://generativelanguage.googleapis.com/v1beta/models/" +
		geminiModel + ":generateContent?key=" + apiKey

	status, body, err := postGemini(url, jsonData)
	if err != nil {
		return geminiReply{}, err
	}

	if status != http.StatusOK && grounding && status == http.StatusBadRequest {
		// Search grounding is unavailable for some requests/accounts — retry once without it.
		req.Tools = nil
		jsonData, err = json.Marshal(req)
		if err != nil {
			return geminiReply{}, err
		}
		status, body, err = postGemini(url, jsonData)
		if err != nil {
			return geminiReply{}, err
		}
	}
	if status != http.StatusOK {
		return geminiReply{}, fmt.Errorf("Gemini API returned status %d: %s", status, trimAPIError(body))
	}

	return parseGeminiResponse(body, thoughts)
}

func postGemini(url string, payload []byte) (int, []byte, error) {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}
	return resp.StatusCode, body, nil
}

func trimAPIError(body []byte) string {
	s := strings.Join(strings.Fields(string(body)), " ")
	if len(s) > 200 {
		s = s[:200] + "…"
	}
	return s
}

func parseGeminiResponse(body []byte, wantThoughts bool) (geminiReply, error) {
	var parsed struct {
		Candidates []struct {
			Content struct {
				Parts []geminiPart `json:"parts"`
			} `json:"content"`
			GroundingMetadata struct {
				GroundingChunks []struct {
					Web struct {
						URI   string `json:"uri"`
						Title string `json:"title"`
					} `json:"web"`
				} `json:"groundingChunks"`
			} `json:"groundingMetadata"`
		} `json:"candidates"`
		PromptFeedback struct {
			BlockReason string `json:"blockReason"`
		} `json:"promptFeedback"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return geminiReply{}, err
	}
	if parsed.PromptFeedback.BlockReason != "" {
		return geminiReply{}, fmt.Errorf("request blocked by safety filter: %s", parsed.PromptFeedback.BlockReason)
	}
	if len(parsed.Candidates) == 0 {
		return geminiReply{}, fmt.Errorf("no response candidate found")
	}

	cand := parsed.Candidates[0]
	var reply geminiReply
	var textParts []string
	for _, p := range cand.Content.Parts {
		if p.Thought {
			if wantThoughts && p.Text != "" {
				reply.Thoughts = append(reply.Thoughts, p.Text)
			}
			continue
		}
		if p.Text != "" {
			textParts = append(textParts, p.Text)
		}
	}
	reply.Text = strings.Join(textParts, "\n")

	seen := map[string]bool{}
	for _, chunk := range cand.GroundingMetadata.GroundingChunks {
		uri := chunk.Web.URI
		if uri == "" || seen[uri] {
			continue
		}
		seen[uri] = true
		title := chunk.Web.Title
		if title == "" {
			title = uri
		}
		reply.Sources = append(reply.Sources, title+" — "+uri)
		if len(reply.Sources) >= 5 {
			break
		}
	}

	if reply.Text == "" {
		return geminiReply{}, fmt.Errorf("empty response from Gemini")
	}
	return reply, nil
}

func printRefinedPrompt(p string) {
	dim := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	accent := lipgloss.NewStyle().Foreground(lipgloss.Color("#0070f3"))

	const maxDisplay = 600
	shown := p
	if len(shown) > maxDisplay {
		shown = shown[:maxDisplay] + " …"
	}

	fmt.Println()
	fmt.Println(accent.Render("  ▸ refined prompt (scoped, full context)"))
	for _, line := range strings.Split(shown, "\n") {
		fmt.Println(dim.Render("    " + line))
	}
}

func printThoughts(thoughts []string) {
	if len(thoughts) == 0 {
		return
	}
	dim := lipgloss.NewStyle().Foreground(lipgloss.Color("#555555"))
	for _, t := range thoughts {
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}
		if len(t) > 400 {
			t = t[:400] + " …"
		}
		fmt.Println(dim.Render("  · " + t))
	}
}

func printSources(sources []string) {
	if len(sources) == 0 {
		return
	}
	dim := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	fmt.Println(dim.Render("  · grounded sources:"))
	for _, s := range sources {
		fmt.Println(dim.Render("    ·" + s))
	}
}
