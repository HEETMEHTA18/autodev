package promptcapture

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// DevMentorSyncResponse represents the response from DevMentor prompt intelligence API
type DevMentorSyncResponse struct {
	OriginalPrompt string   `json:"original_prompt"`
	RefinedPrompt  string   `json:"refined_prompt"`
	Score          int      `json:"score"`
	Workflow       string   `json:"workflow"`
	Technologies   []string `json:"technologies"`
}

// MetadataJSON represents .autodevs/metadata.json
type MetadataJSON struct {
	ProjectName  string    `json:"project_name"`
	Path         string    `json:"path"`
	Technologies []string  `json:"technologies"`
	LastSync     time.Time `json:"last_sync"`
}

// InitializeAutodevsDir ensures the .autodevs directory exists
func InitializeAutodevsDir(root string) error {
	dir := filepath.Join(root, ".autodevs")
	return os.MkdirAll(dir, 0755)
}

// IsSensitive returns true if the prompt likely contains sensitive data
func IsSensitive(prompt string) bool {
	sensitiveKeywords := []string{
		"password", "passwd", "secret", "private_key", "api_key",
		"token", "credentials", "auth_token", "mysql_pwd", "postgres_password",
	}
	lower := strings.ToLower(prompt)
	for _, kw := range sensitiveKeywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

// AppendToPromptsMD writes the prompt to .autodevs/prompts.md
func AppendToPromptsMD(root, prompt string) error {
	return AppendToPromptsMDWithTime(root, prompt, time.Now().Format("2006-01-02 15:04:05"))
}

// AppendToPromptsMDWithTime writes the prompt to .autodevs/prompts.md with a specific timestamp
func AppendToPromptsMDWithTime(root, prompt, timestamp string) error {
	if err := InitializeAutodevsDir(root); err != nil {
		return err
	}

	path := filepath.Join(root, ".autodevs", "prompts.md")

	// Read existing content to check for duplicates
	if data, err := os.ReadFile(path); err == nil {
		if isDuplicate(string(data), prompt) {
			return nil // Duplicate, skip
		}
	}

	// Format prompt entry
	entry := fmt.Sprintf("## %s\n\n%s\n\n---\n\n", timestamp, prompt)

	// Open file in append mode or create it
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// If file is new, write header
	info, err := f.Stat()
	if err == nil && info.Size() == 0 {
		header := "# Prompt History\n\n"
		if _, err := f.WriteString(header); err != nil {
			return err
		}
	}

	_, err = f.WriteString(entry)
	return err
}

func isDuplicate(content, prompt string) bool {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	prompt = strings.TrimSpace(strings.ReplaceAll(prompt, "\r\n", "\n"))

	// Split by entry delimiter
	blocks := strings.Split(content, "\n---\n")
	for _, block := range blocks {
		block = strings.TrimSpace(block)
		if block == "" {
			continue
		}
		// If it's the first block, it might contain the title "# Prompt History"
		if strings.HasPrefix(block, "# Prompt History") {
			block = strings.TrimPrefix(block, "# Prompt History")
			block = strings.TrimSpace(block)
		}

		// Now it should start with "## <timestamp>"
		if !strings.HasPrefix(block, "## ") {
			continue
		}

		// Split into title line and prompt body
		lines := strings.SplitN(block, "\n", 2)
		if len(lines) < 2 {
			continue
		}
		existingPrompt := strings.TrimSpace(lines[1])
		if existingPrompt == prompt {
			return true
		}
	}
	return false
}

// SyncWithDevMentor sends the prompt to the DevMentor intelligence API and updates local files
func SyncWithDevMentor(root, prompt string) error {
	token := os.Getenv("DEVMENTOR_TOKEN")
	if token == "" {
		// Fallback/Silent mode if token is not available
		return nil
	}

	apiURL := "https://devmentor-jmjh.onrender.com/api/v1/prompts/event"
	projectName := filepath.Base(root)

	payload := map[string]string{
		"original_prompt": prompt,
		"project_name":    projectName,
		"file_context":    "AutoDev Prompt Tracking System active",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("DevMentor API returned status %d", resp.StatusCode)
	}

	var syncRes DevMentorSyncResponse
	if err := json.NewDecoder(resp.Body).Decode(&syncRes); err != nil {
		return err
	}

	// Update local .autodevs files with DevMentor intelligence
	if err := UpdateRefinedPromptsMD(root, syncRes); err != nil {
		return err
	}
	if err := UpdateWorkflowsMD(root, syncRes); err != nil {
		return err
	}
	if err := UpdateMetadataJSON(root, syncRes.Technologies); err != nil {
		return err
	}

	return nil
}

// UpdateRefinedPromptsMD appends details to .autodevs/refined-prompts.md
func UpdateRefinedPromptsMD(root string, res DevMentorSyncResponse) error {
	path := filepath.Join(root, ".autodevs", "refined-prompts.md")
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	entry := fmt.Sprintf("## %s\n- **Original:** %s\n- **Refined:** %s\n- **Score:** %d/100\n- **Technologies:** %s\n\n---\n\n",
		timestamp,
		res.OriginalPrompt,
		res.RefinedPrompt,
		res.Score,
		strings.Join(res.Technologies, ", "),
	)

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	info, err := f.Stat()
	if err == nil && info.Size() == 0 {
		header := "# Refined Prompts History\n\n"
		if _, err := f.WriteString(header); err != nil {
			return err
		}
	}

	_, err = f.WriteString(entry)
	return err
}

// UpdateWorkflowsMD appends workflow information to .autodevs/workflows.md
func UpdateWorkflowsMD(root string, res DevMentorSyncResponse) (retErr error) {
	if res.Workflow == "" {
		return nil
	}

	path := filepath.Join(root, ".autodevs", "workflows.md")
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	entry := fmt.Sprintf("## %s\n- **Workflow:** %s\n- **Prompt:** %s\n\n---\n\n",
		timestamp,
		res.Workflow,
		res.OriginalPrompt,
	)

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); retErr == nil && cerr != nil {
			retErr = cerr
		}
	}()

	info, err := f.Stat()
	if err == nil && info.Size() == 0 {
		header := "# Detected Workflows\n\n"
		if _, err := f.WriteString(header); err != nil {
			return err
		}
	}

	_, err = f.WriteString(entry)
	return err
}

// UpdateMetadataJSON merges new technologies into .autodevs/metadata.json
func UpdateMetadataJSON(root string, newTechs []string) error {
	path := filepath.Join(root, ".autodevs", "metadata.json")
	var meta MetadataJSON

	if data, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(data, &meta)
	}

	meta.ProjectName = filepath.Base(root)
	meta.Path = root
	meta.LastSync = time.Now()

	// Merge tech arrays uniquely
	techMap := make(map[string]bool)
	for _, t := range meta.Technologies {
		techMap[strings.ToLower(t)] = true
	}
	for _, t := range newTechs {
		techMap[strings.ToLower(t)] = true
		// Keep original case
		found := false
		for _, existing := range meta.Technologies {
			if strings.EqualFold(existing, t) {
				found = true
				break
			}
		}
		if !found {
			meta.Technologies = append(meta.Technologies, t)
		}
	}

	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// StdinProxy interceptor structure
type StdinProxy struct {
	io.Reader
	buf    []byte
	onLine func(string)
}

func NewStdinProxy(r io.Reader, onLine func(string)) *StdinProxy {
	return &StdinProxy{
		Reader: r,
		onLine: onLine,
	}
}

func (sp *StdinProxy) Read(p []byte) (int, error) {
	n, err := sp.Reader.Read(p)
	if n > 0 {
		sp.buf = append(sp.buf, p[:n]...)
		for {
			idx := bytes.IndexByte(sp.buf, '\n')
			if idx == -1 {
				break
			}
			line := string(bytes.TrimSpace(sp.buf[:idx]))
			sp.buf = sp.buf[idx+1:]

			if len(line) > 3 && !IsSensitive(line) {
				// Avoid logging small control words
				words := strings.Fields(line)
				if len(words) > 1 {
					sp.onLine(line)
				}
			}
		}
	}
	return n, err
}

// AIProcess represents a detected running AI CLI process
type AIProcess struct {
	PID     int
	Name    string
	Cmdline []string
}

// FindActiveAISessions scans /proc to find active AI CLI processes
func FindActiveAISessions() ([]AIProcess, error) {
	files, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}

	var processes []AIProcess
	aiNames := map[string]bool{
		"gemini":  true,
		"claude":  true,
		"copilot": true,
		"agy":     true,
		"codex":   true,
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		// Check if folder name is a number (PID)
		name := file.Name()
		if len(name) == 0 || name[0] < '0' || name[0] > '9' {
			continue
		}

		pidPath := filepath.Join("/proc", name)
		cmdlineData, err := os.ReadFile(filepath.Join(pidPath, "cmdline"))
		if err != nil {
			continue
		}

		parts := bytes.Split(cmdlineData, []byte{0})
		if len(parts) == 0 || len(parts[0]) == 0 {
			continue
		}

		exePath := string(parts[0])
		exeName := filepath.Base(exePath)

		if aiNames[strings.ToLower(exeName)] {
			var cmdline []string
			for _, part := range parts {
				if len(part) > 0 {
					cmdline = append(cmdline, string(part))
				}
			}
			pid := 0
			_, _ = fmt.Sscanf(name, "%d", &pid)

			processes = append(processes, AIProcess{
				PID:     pid,
				Name:    exeName,
				Cmdline: cmdline,
			})
		}
	}

	return processes, nil
}

// ExtractPromptFromCmdline pulls the prompt string from process arguments
func ExtractPromptFromCmdline(cmdline []string) string {
	if len(cmdline) < 2 {
		return ""
	}
	// The prompt is typically the last non-flag argument or any non-flag argument that contains words
	for i := len(cmdline) - 1; i >= 1; i-- {
		arg := cmdline[i]
		if strings.HasPrefix(arg, "-") {
			continue
		}
		words := strings.Fields(arg)
		if len(words) >= 2 {
			return arg
		}
	}
	return ""
}

// StartDaemon runs a loop to monitor running AI CLI sessions and capture their prompts
func StartDaemon(root string, stopChan chan struct{}) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// Track PIDs we have already processed to avoid duplicate logging of the same execution
	processedPIDs := make(map[int]bool)

	for {
		select {
		case <-ticker.C:
			// Sync prompts from active Antigravity/agy session log
			_ = ImportAntigravityPrompts(root)

			procs, err := FindActiveAISessions()
			if err != nil {
				continue
			}

			// Clean up processed PIDs that are no longer running
			currentPIDs := make(map[int]bool)
			for _, p := range procs {
				currentPIDs[p.PID] = true
			}
			for pid := range processedPIDs {
				if !currentPIDs[pid] {
					delete(processedPIDs, pid)
				}
			}

			for _, proc := range procs {
				if processedPIDs[proc.PID] {
					continue
				}

				prompt := ExtractPromptFromCmdline(proc.Cmdline)
				if prompt != "" && !IsSensitive(prompt) {
					// Log the captured prompt
					_ = AppendToPromptsMD(root, prompt)
					_ = SyncWithDevMentor(root, prompt)
				}
				processedPIDs[proc.PID] = true
			}
		case <-stopChan:
			return
		}
	}
}

// GetLatestOverviewPath finds the active Antigravity CLI overview log path
func GetLatestOverviewPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	brainDir := filepath.Join(home, ".gemini", "antigravity", "brain")
	entries, err := os.ReadDir(brainDir)
	if err != nil {
		return "", err
	}

	var latestPath string
	var latestTime time.Time

	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == "tempmediaStorage" {
			continue
		}
		overviewPath := filepath.Join(brainDir, entry.Name(), ".system_generated", "logs", "overview.txt")
		info, err := os.Stat(overviewPath)
		if err != nil {
			continue
		}
		if info.ModTime().After(latestTime) {
			latestTime = info.ModTime()
			latestPath = overviewPath
		}
	}

	if latestPath == "" {
		return "", fmt.Errorf("no conversation brain overview logs found")
	}

	return latestPath, nil
}

// cleanPrompt extracts the plain user request content, removing internal XML structures
func cleanPrompt(content string) string {
	content = strings.TrimSpace(content)

	// If there's a USER_REQUEST tag, extract only the request content
	if strings.Contains(content, "<USER_REQUEST>") {
		start := strings.Index(content, "<USER_REQUEST>") + len("<USER_REQUEST>")
		if strings.Contains(content, "</USER_REQUEST>") {
			end := strings.Index(content, "</USER_REQUEST>")
			if end > start {
				content = content[start:end]
			} else {
				content = content[start:]
			}
		} else {
			content = content[start:]
		}
	}

	// Helper to strip tag blocks and their inner content
	stripBlock := func(s, openTag, closeTag string) string {
		for {
			start := strings.Index(s, openTag)
			if start == -1 {
				break
			}
			if closeIdx := strings.Index(s, closeTag); closeIdx != -1 && closeIdx > start {
				s = s[:start] + s[closeIdx+len(closeTag):]
			} else {
				s = s[:start]
				break
			}
		}
		return s
	}

	content = stripBlock(content, "<USER_INFORMATION>", "</USER_INFORMATION>")
	content = stripBlock(content, "<ADDITIONAL_METADATA>", "</ADDITIONAL_METADATA>")
	content = stripBlock(content, "<EPHEMERAL_MESSAGE>", "</EPHEMERAL_MESSAGE>")

	// Clean up any stray closing tags
	content = strings.ReplaceAll(content, "</USER_REQUEST>", "")
	content = strings.ReplaceAll(content, "</USER_INFORMATION>", "")
	content = strings.ReplaceAll(content, "</ADDITIONAL_METADATA>", "")
	content = strings.ReplaceAll(content, "</EPHEMERAL_MESSAGE>", "")

	return strings.TrimSpace(content)
}

// ImportAntigravityPrompts parses the active Antigravity session and writes prompts to prompts.md
func ImportAntigravityPrompts(workspaceRoot string) error {
	overviewPath, err := GetLatestOverviewPath()
	if err != nil {
		return err
	}

	file, err := os.Open(overviewPath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var oLine struct {
			Source    string `json:"source"`
			Type      string `json:"type"`
			Content   string `json:"content"`
			CreatedAt string `json:"created_at"`
		}
		if err := json.Unmarshal(line, &oLine); err != nil {
			continue
		}

		if oLine.Source == "MODEL" {
			continue
		}

		if oLine.Type == "USER_INPUT" && oLine.Source == "USER_EXPLICIT" {
			prompt := cleanPrompt(oLine.Content)
			if prompt == "" || len(prompt) < 4 || IsSensitive(prompt) {
				continue
			}
			timestamp := time.Now().Format("2006-01-02 15:04:05")
			if oLine.CreatedAt != "" {
				if t, err := time.Parse(time.RFC3339, oLine.CreatedAt); err == nil {
					timestamp = t.Local().Format("2006-01-02 15:04:05")
				} else if t, err := time.Parse(time.RFC3339Nano, oLine.CreatedAt); err == nil {
					timestamp = t.Local().Format("2006-01-02 15:04:05")
				}
			}
			// Deduplication and saving is handled by AppendToPromptsMDWithTime
			_ = AppendToPromptsMDWithTime(workspaceRoot, prompt, timestamp)
		}
	}

	return nil
}
