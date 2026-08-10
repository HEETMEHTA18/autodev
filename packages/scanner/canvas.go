package scanner

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// CanvasGenerator builds a Canvas from a project root.
type CanvasGenerator struct {
	RootPath   string
	ScanResult *ScanResult
	fileCount  int
	locCount   int
}

// NewCanvasGenerator creates a generator for the given path.
func NewCanvasGenerator(rootPath string) *CanvasGenerator {
	abs, _ := filepath.Abs(rootPath)
	return &CanvasGenerator{RootPath: abs}
}

// Generate produces the full Canvas. Returns canvas and any non-fatal errors.
func (cg *CanvasGenerator) Generate() (*Canvas, error) {
	// Run the technology scan first
	s := New(cg.RootPath)
	sr, err := s.Scan()
	if err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}
	cg.ScanResult = sr

	canvas := &Canvas{
		SchemaVersion: CanvasSchemaVersion,
		GeneratedAt:   time.Now().UTC(),
	}

	canvas.Project = cg.buildProjectOverview(sr)
	canvas.Files = cg.buildFileCards()
	canvas.Architecture = cg.buildArchitecture()
	canvas.Dependencies = cg.buildDependencyGraph(canvas.Files)
	canvas.API = cg.buildAPIView(canvas.Files)
	canvas.Database = cg.buildDatabaseView(canvas.Files)
	canvas.Security = cg.buildSecurityView()
	canvas.Git = cg.buildGitSummary()
	canvas.Notes = loadExistingNotes(cg.RootPath)
	canvas.Groups = loadExistingGroups(cg.RootPath)
	canvas.Connections = loadExistingConnections(cg.RootPath)
	cg.buildAIContext(canvas)
	canvas.Lenses = buildLenses(canvas)

	return canvas, nil
}

// Save writes the canvas to the .autodevs/canvas/ directory.
func (cg *CanvasGenerator) Save(canvas *Canvas) error {
	dir := filepath.Join(cg.RootPath, ".autodevs", "canvas")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create canvas dir: %w", err)
	}
	data, err := json.MarshalIndent(canvas, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal canvas: %w", err)
	}
	path := filepath.Join(dir, "canvas.json")
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write canvas: %w", err)
	}
	return nil
}

// ── Project Overview ────────────────────────────────────────────────

func (cg *CanvasGenerator) buildProjectOverview(sr *ScanResult) ProjectOverview {
	name := filepath.Base(cg.RootPath)
	cg.fileCount = 0
	cg.locCount = 0
	_ = filepath.Walk(cg.RootPath, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		base := info.Name()
		if info.IsDir() {
			if strings.HasPrefix(base, ".") || base == "node_modules" || base == "vendor" || base == ".git" || base == "dist" || base == "build" || base == ".next" || base == ".turbo" {
				return filepath.SkipDir
			}
			return nil
		}
		cg.fileCount++
		if ext := filepath.Ext(base); ext == ".go" || ext == ".js" || ext == ".ts" || ext == ".tsx" || ext == ".jsx" || ext == ".py" || ext == ".rs" || ext == ".java" || ext == ".rb" || ext == ".php" {
			lines, _ := countLines(p)
			cg.locCount += lines
		}
		return nil
	})

	return ProjectOverview{
		Name:            name,
		RootPath:        cg.RootPath,
		Languages:       sr.Languages,
		Frameworks:      sr.Frameworks,
		PackageManagers: sr.PackageManagers,
		Databases:       sr.Databases,
		Infra:           sr.Infra,
		TotalFiles:      cg.fileCount,
		TotalLOC:        cg.locCount,
		Monorepo:        len(sr.Projects) > 0,
		Subprojects:     len(sr.Projects),
	}
}

// ── File Cards ──────────────────────────────────────────────────────

func (cg *CanvasGenerator) buildFileCards() []FileCard {
	var cards []FileCard
	langMap := buildLanguageMap()

	_ = filepath.Walk(cg.RootPath, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		base := info.Name()
		if info.IsDir() {
			if strings.HasPrefix(base, ".") || base == "node_modules" || base == "vendor" || base == ".git" || base == "dist" || base == "build" || base == "__pycache__" || base == ".next" || base == ".turbo" {
				return filepath.SkipDir
			}
			return nil
		}

		ext := filepath.Ext(base)
		lang := langMap[ext]
		if lang == "" {
			return nil
		}

		rel, _ := filepath.Rel(cg.RootPath, p)
		loc, _ := countLines(p)
		content, _ := os.ReadFile(p)

		card := FileCard{
			Path:         rel,
			Language:     lang,
			LOC:          loc,
			FileType:     classifyFile(rel, lang, content),
			LastModified: info.ModTime().UTC().Format(time.RFC3339),
			IsEntryPoint: isEntryPoint(rel, content),
			HasTodos:     hasTodos(content),
			TODOList:     extractTodos(content),
			Functions:    extractFunctions(content, lang),
		}

		card.Imports = extractImports(content, lang)
		card.Exports = extractExports(content, lang)
		card.Tests = findTestFile(rel)

		// Security scan on the file
		card.SecurityIssues = scanFileSecurity(p, content)

		cards = append(cards, card)
		return nil
	})

	// Build dependency edges
	cg.buildDependencyEdges(cards)
	cg.detectDeadCode(cards)

	return cards
}

func (cg *CanvasGenerator) buildDependencyEdges(cards []FileCard) {
	fileMap := make(map[string]int)
	for i, c := range cards {
		fileMap[c.Path] = i
	}
	for i, c := range cards {
		for _, imp := range c.Imports {
			resolved := resolveImport(c.Path, imp, cards)
			if resolved != "" {
				cards[i].Dependencies = append(cards[i].Dependencies, resolved)
				if j, ok := fileMap[resolved]; ok {
					cards[j].Dependents = append(cards[j].Dependents, c.Path)
				}
			}
		}
	}
}

func (cg *CanvasGenerator) detectDeadCode(cards []FileCard) {
	for i := range cards {
		if len(cards[i].Dependents) == 0 && !cards[i].IsEntryPoint &&
			!strings.Contains(cards[i].FileType, "test") &&
			!strings.Contains(cards[i].FileType, "config") &&
			!strings.Contains(cards[i].FileType, "entry") {
			cards[i].IsDeadCode = true
		}
	}
}

// ── Architecture ────────────────────────────────────────────────────

func (cg *CanvasGenerator) buildArchitecture() ArchitectureView {
	av := ArchitectureView{}

	// Build directory tree
	av.DirectoryTree = cg.buildDirTree(0)

	// Detect layers from directory structure
	av.Layers = cg.detectLayers()
	av.EntryPoints = cg.detectEntryPoints()
	av.DataFlow = cg.detectDataFlow()

	return av
}

func (cg *CanvasGenerator) buildDirTree(depth int) []DirNode {
	if depth > 3 {
		return nil
	}
	var nodes []DirNode
	entries, err := os.ReadDir(cg.RootPath)
	if err != nil {
		return nil
	}
	for _, e := range entries {
		name := e.Name()
		if strings.HasPrefix(name, ".") || name == "node_modules" || name == "vendor" || name == "dist" || name == "build" || name == "__pycache__" {
			continue
		}
		node := DirNode{Name: name, IsDir: e.IsDir()}
		if e.IsDir() {
			node.Path = name
			subPath := filepath.Join(cg.RootPath, name)
			old := cg.RootPath
			cg.RootPath = subPath
			node.Children = cg.buildDirTree(depth + 1)
			cg.RootPath = old
		} else {
			node.Path = name
		}
		nodes = append(nodes, node)
	}
	return nodes
}

func (cg *CanvasGenerator) detectLayers() []ArchLayer {
	layers := []ArchLayer{
		{Name: "Presentation", Description: "UI components, views, templates", Technologies: []string{}},
		{Name: "API/Route", Description: "HTTP handlers, API routes, controllers", Technologies: []string{}},
		{Name: "Service/Business", Description: "Core business logic, services", Technologies: []string{}},
		{Name: "Data/Repository", Description: "Data access, ORM, database layer", Technologies: []string{}},
		{Name: "Infrastructure", Description: "Docker, CI/CD, deployment", Technologies: []string{}},
		{Name: "Configuration", Description: "Config files, environment, settings", Technologies: []string{}},
		{Name: "Tests", Description: "Unit tests, integration tests", Technologies: []string{}},
	}

	pathMap := map[string][]string{
		"Presentation":   {"components", "views", "templates", "pages", "ui", "frontend", "app"},
		"API/Route":      {"api", "routes", "handlers", "controllers", "endpoints", "server", "cmd"},
		"Service/Business": {"services", "lib", "core", "domain", "business", "usecases"},
		"Data/Repository": {"db", "database", "models", "repositories", "dao", "orm", "migrations", "data", "store"},
		"Infrastructure": {"infra", "infrastructure", "deploy", "docker", "ci", "cd", "terraform", "helm", "scripts"},
		"Configuration":  {"config", "cfg", "settings", "env"},
		"Tests":          {"test", "tests", "__tests__", "spec", "specs"},
	}

	_ = filepath.Walk(cg.RootPath, func(p string, info os.FileInfo, err error) error {
		if err != nil || !info.IsDir() {
			return nil
		}
		base := strings.ToLower(info.Name())
		if strings.HasPrefix(base, ".") || base == "node_modules" || base == "vendor" || base == ".git" || base == "dist" || base == "build" || base == "__pycache__" || base == ".next" || base == ".turbo" {
			return filepath.SkipDir
		}
		for i := range layers {
			for _, pattern := range pathMap[layers[i].Name] {
				if strings.Contains(base, pattern) {
					rel, _ := filepath.Rel(cg.RootPath, p)
					layers[i].Paths = appendUniq(layers[i].Paths, rel)
					break
				}
			}
		}
		return nil
	})

	return layers
}

func (cg *CanvasGenerator) detectEntryPoints() []EntryPoint {
	var eps []EntryPoint
	entryPatterns := []struct {
		Pattern string
		Type    string
	}{
		{"main.go", "cli"},
		{"main.py", "cli"},
		{"index.ts", "web"},
		{"index.tsx", "web"},
		{"index.js", "web"},
		{"App.tsx", "web"},
		{"App.jsx", "web"},
		{"app.py", "web"},
		{"server.go", "server"},
		{"server.ts", "server"},
		{"server.js", "server"},
		{"cmd/root.go", "cli"},
		{"manage.py", "cli"},
	}

	_ = filepath.Walk(cg.RootPath, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		base := info.Name()
		if info.IsDir() {
			if strings.HasPrefix(base, ".") || base == "node_modules" || base == "vendor" || base == ".git" || base == "dist" || base == "build" || base == "__pycache__" || base == ".next" || base == ".turbo" {
				return filepath.SkipDir
			}
			return nil
		}
		for _, ep := range entryPatterns {
			if base == ep.Pattern || strings.HasSuffix("/"+p, "/"+ep.Pattern) {
				rel, _ := filepath.Rel(cg.RootPath, p)
				eps = append(eps, EntryPoint{Path: rel, Type: ep.Type})
				break
			}
		}
		return nil
	})

	return eps
}

func (cg *CanvasGenerator) detectDataFlow() []DataFlow {
	var flows []DataFlow
	// Detect common patterns from file relationships
	flows = append(flows, DataFlow{From: "Client", To: "API", Via: "HTTP", Type: "request"})
	flows = append(flows, DataFlow{From: "API", To: "Service", Via: "function calls", Type: "request"})
	flows = append(flows, DataFlow{From: "Service", To: "Database", Via: "ORM/SQL", Type: "data"})
	return flows
}

// ── Dependency Graph ────────────────────────────────────────────────

func (cg *CanvasGenerator) buildDependencyGraph(files []FileCard) DependencyGraph {
	dg := DependencyGraph{
		Files:     make([]DepNode, len(files)),
		Edges:     []DepEdge{},
		ImpactMap: make(map[string][]string),
	}

	fileMap := make(map[string]int)
	for i, f := range files {
		fileMap[f.Path] = i
		dg.Files[i] = DepNode{
			Path:    f.Path,
			Imports: len(f.Imports),
			Exports: len(f.Exports),
			Layer:   f.FileType,
		}
	}

	// Build edges from dependencies
	seen := make(map[string]bool)
	for _, f := range files {
		for _, dep := range f.Dependencies {
			key := f.Path + "->" + dep
			if !seen[key] {
				seen[key] = true
				dg.Edges = append(dg.Edges, DepEdge{From: f.Path, To: dep, Type: "import"})
			}
		}
	}

	// Detect circular dependencies
	dg.CircularGroups = detectCircularDeps(files)

	// Detect dead files
	for _, f := range files {
		if f.IsDeadCode {
			dg.DeadFiles = append(dg.DeadFiles, f.Path)
		}
	}

	// Build impact map
	for _, f := range files {
		dg.ImpactMap[f.Path] = append(dg.ImpactMap[f.Path], f.Dependents...)
	}

	return dg
}

// ── API View ────────────────────────────────────────────────────────

func (cg *CanvasGenerator) buildAPIView(files []FileCard) APIView {
	av := APIView{}
	routePatterns := map[string][]string{
		"express":    {"router\\.", "app\\.(get|post|put|delete|patch)", "app\\.use"},
		"gin":        {"r\\.(GET|POST|PUT|DELETE|PATCH)", "router\\.(GET|POST|PUT|DELETE|PATCH)"},
		"echo":       {"e\\.(GET|POST|PUT|DELETE|PATCH|Add)", "api\\.(GET|POST|PUT|DELETE|PATCH)"},
		"fiber":      {"app\\.(Get|Post|Put|Delete|Patch)"},
		"fasthttp":   {"fasthttp\\.ListenAndServe"},
		"nextjs":     {"export (default function|const) (GET|POST|PUT|DELETE|PATCH)"},
		"django":     {"path\\(|url\\(|@api_view|@app\\.route"},
		"flask":      {"@app\\.route|@blueprint\\.route"},
		"go":         {"http\\.Handle|http\\.ListenAndServe|mux\\.Handle"},
	}

	for _, f := range files {
		content, _ := os.ReadFile(filepath.Join(cg.RootPath, f.Path))
		if content == nil {
			continue
		}
		contentStr := string(content)

		for fw, patterns := range routePatterns {
			for _, pat := range patterns {
				re := regexp.MustCompile(pat)
				matches := re.FindAllStringIndex(contentStr, -1)
				if len(matches) > 0 {
					route := APIRoute{
						Method: extractHTTPMethod(contentStr, matches),
						Path:   extractRoutePath(contentStr, matches),
						Handler: fw,
						File:   f.Path,
					}
					av.Routes = append(av.Routes, route)
				}
			}
		}
	}

	return av
}

func extractHTTPMethod(content string, matches [][]int) string {
	if len(matches) == 0 {
		return "GET"
	}
	// Look around the match for method hints
	start := matches[0][0]
	if start > 0 {
		before := content[:start]
		methods := []string{"POST", "PUT", "DELETE", "PATCH", "GET"}
		for _, m := range methods {
			if strings.Contains(before, m) {
				return m
			}
		}
	}
	return "GET"
}

func extractRoutePath(content string, matches [][]int) string {
	if len(matches) == 0 {
		return "/"
	}
	// Try to extract path from surrounding context
	start := matches[0][0]
	end := matches[0][1]
	chunk := content[start:min(end+80, len(content))]
	re := regexp.MustCompile(`["'` + "`" + `](/[^"'` + "`" + `\s]+)["'` + "`" + `]`)
	if m := re.FindStringSubmatch(chunk); len(m) > 1 {
		return m[1]
	}
	return "/"
}

// ── Database View ───────────────────────────────────────────────────

func (cg *CanvasGenerator) buildDatabaseView(files []FileCard) DatabaseView {
	dv := DatabaseView{}
	modelPatterns := []struct {
		Re   *regexp.Regexp
		Type string
	}{
		{regexp.MustCompile(`type\s+(\w+)\s+struct`), "go_struct"},
		{regexp.MustCompile(`class\s+(\w+)\(models\.Model\)`), "django_model"},
		{regexp.MustCompile(`@Entity`), "jpa_entity"},
		{regexp.MustCompile(`@Table`), "jpa_table"},
		{regexp.MustCompile(`Schema\.define`), "sequelize_schema"},
		{regexp.MustCompile(`mongoose\.Schema`), "mongoose_schema"},
		{regexp.MustCompile(`CREATE\s+TABLE\s+(\w+)`), "sql_table"},
	}

	for _, f := range files {
		content, _ := os.ReadFile(filepath.Join(cg.RootPath, f.Path))
		if content == nil {
			continue
		}
		contentStr := string(content)

		for _, mp := range modelPatterns {
			matches := mp.Re.FindAllStringSubmatch(contentStr, -1)
			for _, m := range matches {
				name := ""
				if len(m) > 1 {
					name = m[1]
				}
				if name == "" {
					continue
				}
				dv.Models = append(dv.Models, DBModel{
					Name: name,
					File: f.Path,
				})
			}
		}
	}

	return dv
}

// ── Security View ───────────────────────────────────────────────────

func (cg *CanvasGenerator) buildSecurityView() SecurityView {
	sv := SecurityView{}
	sv.Findings = []SecurityFinding{}
	sv.Score = 100

	_ = filepath.Walk(cg.RootPath, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		base := info.Name()
		if info.IsDir() {
			if strings.HasPrefix(base, ".") || base == "node_modules" || base == "vendor" || base == ".git" || base == "dist" || base == "build" || base == "__pycache__" || base == ".next" || base == ".turbo" {
				return filepath.SkipDir
			}
			return nil
		}

		ext := filepath.Ext(base)
		if ext != ".go" && ext != ".js" && ext != ".ts" && ext != ".tsx" && ext != ".jsx" &&
			ext != ".py" && ext != ".java" && ext != ".rb" {
			return nil
		}

		content, err := os.ReadFile(p)
		if err != nil {
			return nil
		}
		rel, _ := filepath.Rel(cg.RootPath, p)
		sv.Findings = append(sv.Findings, scanFileSecurity(rel, content)...)
		return nil
	})

	for _, f := range sv.Findings {
		switch f.Severity {
		case "CRITICAL":
			sv.CriticalCnt++
			sv.Score -= 15
		case "HIGH":
			sv.HighCnt++
			sv.Score -= 5
		case "MEDIUM":
			sv.Score -= 2
		}
	}
	if sv.Score < 0 {
		sv.Score = 0
	}

	return sv
}

func scanFileSecurity(path string, content []byte) []SecurityFinding {
	var findings []SecurityFinding
	if len(content) == 0 {
		return findings
	}
	s := string(content)

	patterns := []struct {
		Re       *regexp.Regexp
		Severity string
		Type     string
		Message  string
	}{
		{regexp.MustCompile(`eval\s*\(`), "HIGH", "code_injection", "eval() usage detected — potential code injection"},
		{regexp.MustCompile(`exec\s*\(`), "HIGH", "code_injection", "exec() usage detected — potential command injection"},
		{regexp.MustCompile(`innerHTML\s*=`), "MEDIUM", "xss", "innerHTML assignment — potential XSS"},
		{regexp.MustCompile(`dangerouslySetInnerHTML`), "MEDIUM", "xss", "dangerouslySetInnerHTML — potential XSS"},
		{regexp.MustCompile(`password.*=.*["']`), "HIGH", "hardcoded_secret", "Hardcoded password detected"},
		{regexp.MustCompile(`api[_-]?key.*=.*["']`), "HIGH", "hardcoded_secret", "Hardcoded API key detected"},
		{regexp.MustCompile(`secret.*=.*["']`), "HIGH", "hardcoded_secret", "Hardcoded secret detected"},
		{regexp.MustCompile(`token.*=.*["']`), "MEDIUM", "hardcoded_secret", "Hardcoded token detected"},
		{regexp.MustCompile(`SELECT.*FROM.*WHERE.*\+`), "MEDIUM", "sql_injection", "String concatenation in SQL — potential SQL injection"},
		{regexp.MustCompile(`fmt\.Sprintf.*SELECT`), "MEDIUM", "sql_injection", "Sprintf in SQL query — potential SQL injection"},
		{regexp.MustCompile(`os\.Open.*\+`), "LOW", "path_traversal", "Dynamic file path — potential path traversal"},
		{regexp.MustCompile(`http\.ListenAndServe\(.*:.*`), "LOW", "security_config", "HTTP server on non-standard port"},
		{regexp.MustCompile(`TLSClientConfig`), "MEDIUM", "crypto", "Custom TLS config — verify certificate validation"},
		{regexp.MustCompile(`InsecureSkipVerify\s*:\s*true`), "HIGH", "crypto", "TLS verification disabled"},
	}

	lines := strings.Split(s, "\n")
	for i, line := range lines {
		for _, pat := range patterns {
			if pat.Re.MatchString(line) {
				// Skip comments
				trimmed := strings.TrimSpace(line)
				if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "/*") {
					continue
				}
				findings = append(findings, SecurityFinding{
					File:     path,
					Line:     i + 1,
					Severity: pat.Severity,
					Type:     pat.Type,
					Message:  pat.Message,
				})
			}
		}
	}

	return findings
}

// ── Git Summary ─────────────────────────────────────────────────────

func (cg *CanvasGenerator) buildGitSummary() GitSummary {
	gs := GitSummary{
		ChurnMap: make(map[string]int),
	}

	// Check if git is available
	if err := exec.Command("git", "-C", cg.RootPath, "rev-parse", "--git-dir").Run(); err != nil {
		return gs
	}

	// Total commits
	if out, err := exec.Command("git", "-C", cg.RootPath, "rev-list", "--count", "HEAD").Output(); err == nil {
		gs.TotalCommits, _ = strconv.Atoi(strings.TrimSpace(string(out)))
	}

	// Last commit
	if out, err := exec.Command("git", "-C", cg.RootPath, "log", "-1", "--format=%H %s").Output(); err == nil {
		gs.LastCommit = strings.TrimSpace(string(out))
	}

	// Contributors
	if out, err := exec.Command("git", "-C", cg.RootPath, "shortlog", "-sn", "--all").Output(); err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			parts := strings.SplitN(line, "\t", 2)
			if len(parts) == 2 {
				count, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
				gs.Contributors = append(gs.Contributors, Contributor{
					Name:  strings.TrimSpace(parts[1]),
					Count: count,
				})
			}
		}
	}

	// Branches
	if out, err := exec.Command("git", "-C", cg.RootPath, "branch", "--format=%(refname:short)").Output(); err == nil {
		for _, b := range strings.Split(strings.TrimSpace(string(out)), "\n") {
			b = strings.TrimSpace(b)
			if b != "" {
				gs.Branches = append(gs.Branches, b)
			}
		}
	}

	// Hot files (most changed) — limit to recent commits for speed
	if out, err := exec.Command("git", "-C", cg.RootPath, "log", "--format=", "--numstat", "-50", "--all").Output(); err == nil {
		fileChanges := make(map[string]int)
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			parts := strings.Fields(line)
			if len(parts) == 3 && parts[0] != "-" && parts[1] != "-" {
				fileChanges[parts[2]]++
			}
		}
		type fc struct {
			Path   string
			Changes int
		}
		var sorted []fc
		for p, c := range fileChanges {
			sorted = append(sorted, fc{p, c})
		}
		sort.Slice(sorted, func(i, j int) bool { return sorted[i].Changes > sorted[j].Changes })
		for i, s := range sorted {
			if i >= 20 {
				break
			}
			gs.HotFiles = append(gs.HotFiles, HotFile(s))
			gs.ChurnMap[s.Path] = s.Changes
		}
	}

	return gs
}

// ── AI Context Builder ──────────────────────────────────────────────

func (cg *CanvasGenerator) buildAIContext(canvas *Canvas) {
	ai := AIContext{}

	// One-liner
	if len(canvas.Project.Languages) > 0 {
		langs := strings.Join(canvas.Project.Languages, ", ")
		fw := strings.Join(canvas.Project.Frameworks, ", ")
		if fw != "" {
			ai.OneLiner = fmt.Sprintf("%s project using %s (%d files, %d LOC)",
				fw, langs, canvas.Project.TotalFiles, canvas.Project.TotalLOC)
		} else {
			ai.OneLiner = fmt.Sprintf("%s project (%d files, %d LOC)",
				langs, canvas.Project.TotalFiles, canvas.Project.TotalLOC)
		}
	} else {
		ai.OneLiner = fmt.Sprintf("Project with %d files", canvas.Project.TotalFiles)
	}

	// Summary
	var summaryParts []string
	if canvas.Project.Monorepo {
		summaryParts = append(summaryParts, fmt.Sprintf("Monorepo with %d subprojects", canvas.Project.Subprojects))
	}
	if len(canvas.Architecture.EntryPoints) > 0 {
		summaryParts = append(summaryParts, fmt.Sprintf("%d entry points", len(canvas.Architecture.EntryPoints)))
	}
	if len(canvas.API.Routes) > 0 {
		summaryParts = append(summaryParts, fmt.Sprintf("%d API routes", len(canvas.API.Routes)))
	}
	if len(canvas.Database.Models) > 0 {
		summaryParts = append(summaryParts, fmt.Sprintf("%d database models", len(canvas.Database.Models)))
	}
	if canvas.Security.CriticalCnt > 0 {
		summaryParts = append(summaryParts, fmt.Sprintf("%d critical security findings", canvas.Security.CriticalCnt))
	}
	ai.Summary = strings.Join(summaryParts, "; ")

	// Key insights
	if len(canvas.Git.HotFiles) > 0 {
		ai.FrequentChanges = []string{}
		for i, hf := range canvas.Git.HotFiles {
			if i >= 10 {
				break
			}
			ai.FrequentChanges = append(ai.FrequentChanges, hf.Path)
		}
		ai.KeyInsights = append(ai.KeyInsights,
			fmt.Sprintf("Most frequently changed files: %s", strings.Join(ai.FrequentChanges[:min(3, len(ai.FrequentChanges))], ", ")))
	}

	if len(canvas.Dependencies.DeadFiles) > 0 {
		ai.KeyInsights = append(ai.KeyInsights,
			fmt.Sprintf("%d potentially dead/unused files detected", len(canvas.Dependencies.DeadFiles)))
		ai.DangerZones = append(ai.DangerZones,
			fmt.Sprintf("Dead code candidates: %s", strings.Join(canvas.Dependencies.DeadFiles[:min(3, len(canvas.Dependencies.DeadFiles))], ", ")))
	}

	if len(canvas.Dependencies.CircularGroups) > 0 {
		ai.KeyInsights = append(ai.KeyInsights,
			fmt.Sprintf("%d circular dependency groups detected", len(canvas.Dependencies.CircularGroups)))
		ai.DangerZones = append(ai.DangerZones, "Circular dependencies found — refactoring candidates")
	}

	if canvas.Security.Score < 100 {
		ai.KeyInsights = append(ai.KeyInsights,
			fmt.Sprintf("Security score: %d/100", canvas.Security.Score))
	}

	// Tech stack
	ai.TechStack = append(ai.TechStack, canvas.Project.Languages...)
	ai.TechStack = append(ai.TechStack, canvas.Project.Frameworks...)
	ai.TechStack = append(ai.TechStack, canvas.Project.Databases...)

	// Critical files
	ai.CriticalFiles = []string{}
	for _, ep := range canvas.Architecture.EntryPoints {
		ai.CriticalFiles = append(ai.CriticalFiles, ep.Path)
	}
	if len(canvas.API.Routes) > 0 {
		ai.CriticalFiles = append(ai.CriticalFiles, canvas.API.Routes[0].File)
	}

	// Conventions (inferred)
	ai.Conventions = cg.inferConventions(canvas)

	// Patterns
	ai.Patterns = cg.inferPatterns(canvas)

	// Recommended reading order
	ai.RecommendedReading = cg.buildReadingOrder(canvas)

	// FAQs
	ai.FAQs = cg.generateFAQs(canvas)

	// Glossary
	ai.Glossary = cg.generateGlossary(canvas)

	// Architecture description
	ai.ArchitectureDescription = cg.describeArchitecture(canvas)

	// Gotchas
	ai.Gotchas = cg.findGotchas(canvas)

	canvas.AI = ai
}

func (cg *CanvasGenerator) inferConventions(canvas *Canvas) []string {
	var conv []string

	// Detect naming conventions
	conv = append(conv, fmt.Sprintf("Primary languages: %s", strings.Join(canvas.Project.Languages, ", ")))

	// Check file organization
	hasTests := false
	hasConfig := false
	for _, f := range canvas.Files {
		if strings.Contains(f.FileType, "test") {
			hasTests = true
		}
		if strings.Contains(f.FileType, "config") {
			hasConfig = true
		}
	}
	if hasTests {
		conv = append(conv, "Tests are co-located with source files")
	}
	if hasConfig {
		conv = append(conv, "Configuration files present")
	}

	return conv
}

func (cg *CanvasGenerator) inferPatterns(canvas *Canvas) []string {
	var patterns []string

	if canvas.Project.Monorepo {
		patterns = append(patterns, "Monorepo architecture")
	}

	if len(canvas.Architecture.Layers) > 0 {
		patterns = append(patterns, "Layered architecture detected")
	}

	if len(canvas.API.Routes) > 0 {
		patterns = append(patterns, "REST API pattern")
	}

	if len(canvas.Database.Models) > 0 {
		patterns = append(patterns, "ORM/Active Record pattern")
	}

	return patterns
}

func (cg *CanvasGenerator) buildReadingOrder(canvas *Canvas) []string {
	var order []string

	// Start with config and entry points
	for _, ep := range canvas.Architecture.EntryPoints {
		order = append(order, ep.Path)
	}

	// Add high-connectivity files
	type fileConn struct {
		Path  string
		Count int
	}
	var conns []fileConn
	depCount := make(map[string]int)
	for _, f := range canvas.Files {
		depCount[f.Path] = len(f.Dependencies) + len(f.Dependents)
	}
	for p, c := range depCount {
		conns = append(conns, fileConn{p, c})
	}
	sort.Slice(conns, func(i, j int) bool { return conns[i].Count > conns[j].Count })
	for i, fc := range conns {
		if i >= 20 {
			break
		}
		if fc.Count > 2 {
			order = append(order, fc.Path)
		}
	}

	// Add API routes
	for _, r := range canvas.API.Routes {
		order = append(order, r.File)
	}

	return order
}

func (cg *CanvasGenerator) generateFAQs(canvas *Canvas) []FAQ {
	var faqs []FAQ

	faqs = append(faqs, FAQ{
		Question: "What is the project?",
		Answer:   canvas.AI.OneLiner,
	})

	if len(canvas.Architecture.EntryPoints) > 0 {
		faqs = append(faqs, FAQ{
			Question: "Where does execution start?",
			Answer:   fmt.Sprintf("Entry points: %s", joinPaths(canvas.Architecture.EntryPoints)),
		})
	}

	if canvas.Project.TotalFiles > 0 {
		faqs = append(faqs, FAQ{
			Question: "How big is the codebase?",
			Answer:   fmt.Sprintf("%d files, %d lines of code", canvas.Project.TotalFiles, canvas.Project.TotalLOC),
		})
	}

	if len(canvas.API.Routes) > 0 {
		faqs = append(faqs, FAQ{
			Question: "How many API endpoints are there?",
			Answer:   fmt.Sprintf("%d routes detected", len(canvas.API.Routes)),
		})
	}

	return faqs
}

func (cg *CanvasGenerator) generateGlossary(canvas *Canvas) []GlossaryTerm {
	var glossary []GlossaryTerm

	terms := map[string]string{
		"LOC":   "Lines of Code",
		"DI":    "Dependency Injection",
		"ORM":   "Object-Relational Mapping",
		"IaC":   "Infrastructure as Code",
		"CI":    "Continuous Integration",
		"CD":    "Continuous Delivery",
	}
	for term, def := range terms {
		glossary = append(glossary, GlossaryTerm{Term: term, Definition: def})
	}

	return glossary
}

func (cg *CanvasGenerator) describeArchitecture(canvas *Canvas) string {
	var parts []string
	parts = append(parts, fmt.Sprintf("The project is a %s with %d layers.",
		strings.Join(canvas.Project.Languages, "/"), len(canvas.Architecture.Layers)))

	if len(canvas.Architecture.EntryPoints) > 0 {
		parts = append(parts, fmt.Sprintf("It starts from %d entry points.", len(canvas.Architecture.EntryPoints)))
	}

	if len(canvas.API.Routes) > 0 {
		parts = append(parts, fmt.Sprintf("The API layer has %d routes.", len(canvas.API.Routes)))
	}

	if len(canvas.Database.Models) > 0 {
		parts = append(parts, fmt.Sprintf("The data layer has %d models.", len(canvas.Database.Models)))
	}

	return strings.Join(parts, " ")
}

func (cg *CanvasGenerator) findGotchas(canvas *Canvas) []string {
	var gotchas []string

	if len(canvas.Dependencies.CircularGroups) > 0 {
		gotchas = append(gotchas, "Circular dependencies detected — can cause initialization order issues")
	}

	if len(canvas.Dependencies.DeadFiles) > 0 {
		gotchas = append(gotchas, "Dead code detected — may cause confusion during refactoring")
	}

	if canvas.Security.CriticalCnt > 0 {
		gotchas = append(gotchas, fmt.Sprintf("%d critical security findings — address before production", canvas.Security.CriticalCnt))
	}

	if canvas.Project.TotalLOC > 100000 {
		gotchas = append(gotchas, "Large codebase (100K+ LOC) — consider splitting modules")
	}

	return gotchas
}

// ── Lenses ──────────────────────────────────────────────────────────

func buildLenses(canvas *Canvas) LensesView {
	lv := LensesView{}

	// Technology Lens
	lv.Technology = &TechnologyLens{
		Categories: map[string][]string{
			"Languages":  canvas.Project.Languages,
			"Frameworks": canvas.Project.Frameworks,
			"Databases":  canvas.Project.Databases,
			"Infra":      canvas.Project.Infra,
		},
	}

	// File Lens
	byLayer := make(map[string][]string)
	byType := make(map[string][]string)
	for _, f := range canvas.Files {
		if f.FileType != "" {
			byType[f.FileType] = append(byType[f.FileType], f.Path)
		}
	}
	lv.File = &FileLens{ByLayer: byLayer, ByType: byType}

	// Function Lens
	callGraph := make(map[string][]string)
	for _, f := range canvas.Files {
		for _, fn := range f.Functions {
			if len(fn.Calls) > 0 {
				callGraph[f.Path+":"+fn.Name] = fn.Calls
			}
		}
	}
	lv.Function = &FunctionLens{CallGraph: callGraph}

	// Security Lens
	severityMap := make(map[string][]SecurityFinding)
	for _, f := range canvas.Security.Findings {
		severityMap[f.Severity] = append(severityMap[f.Severity], f)
	}
	lv.Security = &SecurityLens{SeverityMap: severityMap}

	// Git Lens
	lv.Git = &GitLens{
		ContributorMap: make(map[string]int),
		ChurnMap:       canvas.Git.ChurnMap,
	}
	for _, c := range canvas.Git.Contributors {
		lv.Git.ContributorMap[c.Name] = c.Count
	}

	// Runtime Lens (placeholder)
	lv.Runtime = &RuntimeLens{}

	// Architecture Lens
	lv.Architecture = &ArchitectureLens{
		AIArchitecture: canvas.Architecture.Layers,
	}

	return lv
}

// ── Helper functions ────────────────────────────────────────────────

func loadExistingNotes(rootPath string) []CanvasNote {
	path := filepath.Join(rootPath, ".autodevs", "canvas", "notes.json")
	return loadJSON[[]CanvasNote](path)
}

func loadExistingGroups(rootPath string) []CanvasGroup {
	path := filepath.Join(rootPath, ".autodevs", "canvas", "groups.json")
	return loadJSON[[]CanvasGroup](path)
}

func loadExistingConnections(rootPath string) []CanvasEdge {
	path := filepath.Join(rootPath, ".autodevs", "canvas", "connections.json")
	return loadJSON[[]CanvasEdge](path)
}

func loadJSON[T any](path string) T {
	var result T
	data, err := os.ReadFile(path)
	if err != nil {
		return result
	}
	_ = json.Unmarshal(data, &result)
	return result
}

func countLines(path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	return strings.Count(string(data), "\n") + 1, nil
}

func buildLanguageMap() map[string]string {
	return map[string]string{
		".go":    "Go",
		".js":    "JavaScript",
		".ts":    "TypeScript",
		".tsx":   "TypeScript React",
		".jsx":   "JavaScript React",
		".py":    "Python",
		".rs":    "Rust",
		".java":  "Java",
		".rb":    "Ruby",
		".php":   "PHP",
		".c":     "C",
		".cpp":   "C++",
		".h":     "C/C++ Header",
		".cs":    "C#",
		".swift": "Swift",
		".kt":    "Kotlin",
		".dart":  "Dart",
		".lua":   "Lua",
		".sh":    "Shell",
		".yaml":  "YAML",
		".yml":   "YAML",
		".json":  "JSON",
		".toml":  "TOML",
		".sql":   "SQL",
		".css":   "CSS",
		".scss":  "SCSS",
		".html":  "HTML",
		".vue":   "Vue",
		".svelte": "Svelte",
		".tf":    "Terraform",
		".proto": "Protocol Buffers",
	}
}

func classifyFile(relPath string, lang string, content []byte) string {
	base := filepath.Base(relPath)
	dir := filepath.Dir(relPath)
	lower := strings.ToLower(base)

	if strings.HasSuffix(lower, "_test.go") || strings.HasSuffix(lower, ".test.js") || strings.HasSuffix(lower, ".test.ts") || strings.HasSuffix(lower, ".spec.js") || strings.HasSuffix(lower, ".spec.ts") || strings.Contains(lower, "__test__") {
		return "test"
	}
	if strings.Contains(lower, "config") || strings.Contains(lower, ".env") || strings.HasSuffix(lower, "rc") || strings.HasSuffix(lower, ".config.js") || strings.HasSuffix(lower, ".config.ts") {
		return "config"
	}
	if base == "main.go" || base == "main.py" || base == "index.ts" || base == "index.tsx" || base == "index.js" || base == "App.tsx" || base == "App.jsx" {
		return "entry"
	}
	if strings.Contains(dir, "route") || strings.Contains(dir, "api") || strings.Contains(dir, "handler") || strings.Contains(dir, "controller") {
		return "route"
	}
	if strings.Contains(dir, "component") || strings.Contains(dir, "ui") || strings.Contains(dir, "view") {
		return "component"
	}
	if strings.Contains(dir, "util") || strings.Contains(dir, "helper") || strings.Contains(dir, "lib") {
		return "util"
	}
	if strings.Contains(dir, "model") || strings.Contains(dir, "entity") || strings.Contains(dir, "schema") {
		return "model"
	}
	if strings.Contains(dir, "service") || strings.Contains(dir, "domain") {
		return "service"
	}
	if strings.Contains(dir, "store") || strings.Contains(dir, "repo") || strings.Contains(dir, "dao") {
		return "repository"
	}
	if strings.Contains(dir, "middleware") || strings.Contains(dir, "interceptor") {
		return "middleware"
	}
	if strings.Contains(dir, "hook") {
		return "hook"
	}
	if strings.Contains(dir, "context") && lang == "Go" {
		return "context"
	}
	if strings.Contains(dir, "test") || strings.Contains(dir, "spec") {
		return "test"
	}
	if strings.Contains(dir, "migration") || strings.Contains(dir, "migrate") {
		return "migration"
	}

	return "source"
}

func isEntryPoint(relPath string, content []byte) bool {
	base := filepath.Base(relPath)
	entryNames := []string{"main.go", "main.py", "index.ts", "index.tsx", "index.js", "App.tsx", "App.jsx", "cmd/root.go", "manage.py", "server.go", "app.py"}
	for _, e := range entryNames {
		if base == e {
			return true
		}
	}
	return false
}

func hasTodos(content []byte) bool {
	s := string(content)
	return strings.Contains(s, "TODO") || strings.Contains(s, "FIXME") || strings.Contains(s, "HACK")
}

func extractTodos(content []byte) []string {
	var todos []string
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.Contains(line, "TODO:") || strings.Contains(line, "FIXME:") || strings.Contains(line, "HACK:") {
			todos = append(todos, strings.TrimSpace(line))
		}
	}
	return todos
}

func extractFunctions(content []byte, lang string) []Function {
	var funcs []Function
	s := string(content)

	patterns := []struct {
		Re *regexp.Regexp
		Langs []string
	}{
		{regexp.MustCompile(`func\s+(?:\([^)]+\)\s+)?(\w+)\s*\(`), []string{"Go"}},
		{regexp.MustCompile(`(?:export\s+)?(?:async\s+)?function\s+(\w+)`), []string{"JavaScript", "TypeScript", "TypeScript React", "JavaScript React"}},
		{regexp.MustCompile(`(?:export\s+)?(?:const|let|var)\s+(\w+)\s*=\s*(?:async\s*)?\(`), []string{"JavaScript", "TypeScript", "TypeScript React", "JavaScript React"}},
		{regexp.MustCompile(`def\s+(\w+)\s*\(`), []string{"Python"}},
		{regexp.MustCompile(`fn\s+(\w+)`), []string{"Rust"}},
		{regexp.MustCompile(`func\s+(\w+)`), []string{"Go"}},
	}

	for _, p := range patterns {
		supported := false
		for _, l := range p.Langs {
			if l == lang {
				supported = true
				break
			}
		}
		if !supported {
			continue
		}

		matches := p.Re.FindAllStringSubmatchIndex(s, -1)
		for _, m := range matches {
			if len(m) >= 4 {
				start := m[2]
				end := m[3]
				name := s[start:end]
				lineNum := strings.Count(s[:start], "\n") + 1

				isExported := strings.Contains(s[max(0, start-10):start], "export")
				isTest := strings.HasPrefix(name, "Test") || strings.HasPrefix(name, "test_")

				funcs = append(funcs, Function{
					Name:       name,
					Line:       lineNum,
					IsExported: isExported,
					IsTest:     isTest,
				})
			}
		}
	}

	return funcs
}

func extractImports(content []byte, lang string) []string {
	var imports []string
	s := string(content)

	var patterns []*regexp.Regexp
	switch lang {
	case "Go":
		patterns = []*regexp.Regexp{
			regexp.MustCompile(`"([^"]+)"`), // import "path"
		}
		// Also check for grouped imports
		if strings.Contains(s, "import (") {
			lines := strings.Split(s, "\n")
			inImport := false
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "import (" {
					inImport = true
					continue
				}
				if line == ")" {
					inImport = false
					continue
				}
				if inImport {
					m := regexp.MustCompile(`"([^"]+)"`).FindStringSubmatch(line)
					if len(m) > 1 {
						imports = append(imports, m[1])
					}
				}
			}
			return imports
		}
	case "JavaScript", "TypeScript", "TypeScript React", "JavaScript React":
		patterns = []*regexp.Regexp{
			regexp.MustCompile(`from\s+['"]([^'"]+)['"]`),
			regexp.MustCompile(`import\s+['"]([^'"]+)['"]`),
			regexp.MustCompile(`require\s*\(\s*['"]([^'"]+)['"]\s*\)`),
		}
	case "Python":
		patterns = []*regexp.Regexp{
			regexp.MustCompile(`(?:from|import)\s+(\S+)`),
		}
	case "Rust":
		patterns = []*regexp.Regexp{
			regexp.MustCompile(`use\s+([\w:]+)`),
		}
	default:
		patterns = []*regexp.Regexp{
			regexp.MustCompile(`import\s+['"]([^'"]+)['"]`),
		}
	}

	for _, p := range patterns {
		matches := p.FindAllStringSubmatch(s, -1)
		for _, m := range matches {
			if len(m) > 1 {
				imp := m[1]
				if imp != "" && len(imp) > 1 {
					imports = append(imports, imp)
				}
			}
		}
	}

	return imports
}

func extractExports(content []byte, lang string) []string {
	var exports []string
	s := string(content)

	var patterns []*regexp.Regexp
	switch lang {
	case "Go":
		patterns = []*regexp.Regexp{
			regexp.MustCompile(`^func\s+([A-Z]\w*)`),    // exported functions
			regexp.MustCompile(`^type\s+([A-Z]\w*)`),    // exported types
			regexp.MustCompile(`^var\s+([A-Z]\w*)`),     // exported vars
			regexp.MustCompile(`^const\s+([A-Z]\w*)`),   // exported consts
		}
	case "JavaScript", "TypeScript", "TypeScript React", "JavaScript React":
		patterns = []*regexp.Regexp{
			regexp.MustCompile(`export\s+(?:default\s+)?(?:function|const|let|var|class)\s+(\w+)`),
			regexp.MustCompile(`module\.exports\s*=\s*(\w+)`),
		}
	case "Python":
		patterns = []*regexp.Regexp{
			regexp.MustCompile(`^def\s+(\w+)`),    // public functions
			regexp.MustCompile(`^class\s+(\w+)`),  // public classes
		}
	case "Rust":
		patterns = []*regexp.Regexp{
			regexp.MustCompile(`pub\s+(?:fn|struct|enum|trait)\s+(\w+)`),
		}
	}

	for _, p := range patterns {
		matches := p.FindAllStringSubmatch(s, -1)
		for _, m := range matches {
			if len(m) > 1 {
				exports = append(exports, m[1])
			}
		}
	}

	return exports
}

func findTestFile(relPath string) []string {
	base := filepath.Base(relPath)
	ext := filepath.Ext(relPath)
	name := strings.TrimSuffix(base, ext)

	// Check for paired test file
	testExtensions := []string{"_test.go", ".test.js", ".test.ts", ".test.tsx", ".test.jsx", ".spec.js", ".spec.ts", ".spec.tsx", ".spec.jsx", "_test.py"}
	var tests []string
	for _, te := range testExtensions {
		testPath := filepath.Join(filepath.Dir(relPath), name+te)
		tests = append(tests, testPath)
	}

	return tests
}

func resolveImport(importerPath, imp string, cards []FileCard) string {
	// Try to find the import target in the cards
	impBase := filepath.Base(imp)
	for _, c := range cards {
		cBase := filepath.Base(c.Path)
		if strings.HasPrefix(cBase, impBase) || strings.HasSuffix(imp, c.Path) {
			return c.Path
		}
	}
	return ""
}

func detectCircularDeps(files []FileCard) [][]string {
	// Build adjacency list
	graph := make(map[string][]string)
	for _, f := range files {
		graph[f.Path] = f.Dependencies
	}

	// DFS-based cycle detection
	var cycles [][]string
	visited := make(map[string]bool)
	inStack := make(map[string]bool)
	path := []string{}

	var dfs func(string)
	dfs = func(node string) {
		if inStack[node] {
			// Found cycle
			cycleStart := -1
			for i, p := range path {
				if p == node {
					cycleStart = i
					break
				}
			}
			if cycleStart >= 0 {
				cycle := append([]string{}, path[cycleStart:]...)
				cycle = append(cycle, node)
				cycles = append(cycles, cycle)
			}
			return
		}
		if visited[node] {
			return
		}
		visited[node] = true
		inStack[node] = true
		path = append(path, node)

		for _, dep := range graph[node] {
			dfs(dep)
		}

		path = path[:len(path)-1]
		inStack[node] = false
	}

	for _, f := range files {
		dfs(f.Path)
	}

	return cycles
}

func joinPaths(eps []EntryPoint) string {
	var paths []string
	for _, ep := range eps {
		paths = append(paths, fmt.Sprintf("%s (%s)", ep.Path, ep.Type))
	}
	return strings.Join(paths, ", ")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
