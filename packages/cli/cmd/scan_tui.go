package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TreeNode represents a file or directory node in our dependency tree.
type TreeNode struct {
	Path     string
	Name     string
	IsDir    bool
	IsOpen   bool
	Level    int
	Size     int64
	Tech     string
	Deps     []string
	Children []*TreeNode
}

// scanTUIModel holds state for the TUI rendering loop.
type scanTUIModel struct {
	rootPath     string
	rootNode     *TreeNode
	flatNodes    []*TreeNode
	cursor       int
	width        int
	height       int
	scrollOffset int
}

func detectTechForFile(name string) string {
	switch name {
	case "package.json":
		return "Node.js Project"
	case "pnpm-lock.yaml":
		return "pnpm Lockfile"
	case "package-lock.json":
		return "npm Lockfile"
	case "yarn.lock":
		return "yarn Lockfile"
	case "go.mod":
		return "Go Module"
	case "go.sum":
		return "Go Checksums"
	case "Cargo.toml":
		return "Rust Cargo Manifest"
	case "Cargo.lock":
		return "Rust Cargo Lockfile"
	case "requirements.txt", "pyproject.toml":
		return "Python Environment"
	case "pom.xml":
		return "Java Maven Project"
	case "build.gradle", "build.gradle.kts":
		return "Java Gradle Project"
	case "composer.json":
		return "PHP Composer Project"
	case "Gemfile":
		return "Ruby Bundler Project"
	case "Dockerfile":
		return "Docker Container Build"
	case "docker-compose.yml", "docker-compose.yaml":
		return "Docker Compose Spec"
	case "main.tf":
		return "Terraform Configuration"
	case "pubspec.yaml":
		return "Flutter / Dart Project"
	case "next.config.js", "next.config.ts", "next.config.mjs":
		return "Next.js Framework config"
	case "vite.config.ts", "vite.config.js":
		return "Vite Builder config"
	}

	ext := filepath.Ext(name)
	switch ext {
	case ".go":
		return "Go Source"
	case ".ts", ".tsx":
		return "TypeScript Source"
	case ".js", ".jsx":
		return "JavaScript Source"
	case ".py":
		return "Python Source"
	case ".rs":
		return "Rust Source"
	case ".java":
		return "Java Source"
	case ".kt":
		return "Kotlin Source"
	case ".php":
		return "PHP Source"
	case ".rb":
		return "Ruby Source"
	case ".cpp", ".hpp", ".cc", ".c", ".h":
		return "C/C++ Source"
	case ".cs":
		return "C# Source"
	case ".html", ".htm":
		return "HTML Document"
	case ".css", ".scss":
		return "CSS Stylesheet"
	case ".md":
		return "Markdown File"
	case ".json":
		return "JSON Configuration"
	case ".yaml", ".yml":
		return "YAML Configuration"
	case ".tf":
		return "Terraform Source"
	case ".sql":
		return "SQL Script"
	case ".sh":
		return "Shell Script"
	}
	return ""
}

func parseDepsForFile(path string) []string {
	name := filepath.Base(path)
	switch name {
	case "package.json":
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		var pkg map[string]interface{}
		if err := json.Unmarshal(data, &pkg); err != nil {
			return nil
		}
		var deps []string
		for _, key := range []string{"dependencies", "devDependencies"} {
			if m, ok := pkg[key].(map[string]interface{}); ok {
				for k, v := range m {
					deps = append(deps, fmt.Sprintf("%s: %v", k, v))
				}
			}
		}
		sort.Strings(deps)
		return deps
	case "go.mod":
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		var deps []string
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "require ") {
				dep := strings.TrimPrefix(line, "require ")
				deps = append(deps, strings.Trim(dep, "()"))
			} else if (strings.Contains(line, "github.com/") || strings.Contains(line, "golang.org/")) && !strings.HasPrefix(line, "module ") && !strings.HasPrefix(line, "go ") {
				deps = append(deps, line)
			}
		}
		sort.Strings(deps)
		return deps
	case "Cargo.toml":
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		var deps []string
		lines := strings.Split(string(data), "\n")
		inDeps := false
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "[") {
				inDeps = strings.Contains(line, "dependencies")
				continue
			}
			if inDeps && line != "" && !strings.HasPrefix(line, "#") {
				deps = append(deps, line)
			}
		}
		sort.Strings(deps)
		return deps
	case "requirements.txt":
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		var deps []string
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "#") {
				deps = append(deps, line)
			}
		}
		sort.Strings(deps)
		return deps
	}
	return nil
}

func buildTree(root string, maxDepth int) (*TreeNode, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	return buildNode(abs, filepath.Base(abs), 0, maxDepth)
}

func buildNode(path, name string, level, maxDepth int) (*TreeNode, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	node := &TreeNode{
		Path:  path,
		Name:  name,
		IsDir: info.IsDir(),
		Level: level,
		Size:  info.Size(),
	}

	if node.IsDir {
		node.IsOpen = level < 2
		if level >= maxDepth {
			return node, nil
		}

		entries, err := os.ReadDir(path)
		if err != nil {
			return node, nil
		}

		for _, entry := range entries {
			n := entry.Name()
			if strings.HasPrefix(n, ".") || n == "node_modules" || n == "vendor" || n == "dist" || n == "build" || n == ".git" || n == ".turbo" || n == ".next" {
				continue
			}
			childPath := filepath.Join(path, n)
			child, err := buildNode(childPath, n, level+1, maxDepth)
			if err == nil {
				node.Children = append(node.Children, child)
			}
		}

		// Sort: directories first, then files
		sort.Slice(node.Children, func(i, j int) bool {
			if node.Children[i].IsDir && !node.Children[j].IsDir {
				return true
			}
			if !node.Children[i].IsDir && node.Children[j].IsDir {
				return false
			}
			return node.Children[i].Name < node.Children[j].Name
		})
	} else {
		node.Tech = detectTechForFile(name)
		if node.Tech != "" {
			node.Deps = parseDepsForFile(path)
		}
	}

	return node, nil
}

func (m *scanTUIModel) flatten(node *TreeNode) {
	m.flatNodes = append(m.flatNodes, node)
	if node.IsDir && node.IsOpen {
		for _, child := range node.Children {
			m.flatten(child)
		}
	}
}

func (m *scanTUIModel) refreshFlattened() {
	m.flatNodes = nil
	m.flatten(m.rootNode)
	if m.cursor >= len(m.flatNodes) {
		m.cursor = len(m.flatNodes) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

func runScanTUI(path string) error {
	rootNode, err := buildTree(path, 5)
	if err != nil {
		return err
	}

	model := scanTUIModel{
		rootPath: path,
		rootNode: rootNode,
	}
	model.refreshFlattened()

	p := tea.NewProgram(model, tea.WithAltScreen())
	_, err = p.Run()
	return err
}

func (m scanTUIModel) Init() tea.Cmd {
	return nil
}

func (m scanTUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.flatNodes)-1 {
				m.cursor++
			}
		case "enter", "space":
			selected := m.flatNodes[m.cursor]
			if selected.IsDir {
				selected.IsOpen = !selected.IsOpen
				m.refreshFlattened()
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	// Dynamic scroll view logic
	maxVisibleHeight := m.height - 7 // leaving room for header and footer
	if maxVisibleHeight < 5 {
		maxVisibleHeight = 5
	}
	if m.cursor < m.scrollOffset {
		m.scrollOffset = m.cursor
	} else if m.cursor >= m.scrollOffset+maxVisibleHeight {
		m.scrollOffset = m.cursor - maxVisibleHeight + 1
	}

	return m, nil
}

func (m scanTUIModel) View() string {
	if m.width < 20 || m.height < 10 {
		return "Terminal too small."
	}

	// ── Lipgloss Styles ────────────────────────────────────────────────────────
	colorGold := lipgloss.Color("#0070f3")
	colorGreen := lipgloss.Color("#00c853")
	colorCyan := lipgloss.Color("#0070f3")
	colorGray := lipgloss.Color("#555555")
	colorWhite := lipgloss.Color("#FFFFFF")

	styleTitle := lipgloss.NewStyle().Bold(true).Foreground(colorGold).Padding(0, 1).BorderStyle(lipgloss.NormalBorder()).BorderForeground(colorGold)
	styleHeaderDesc := lipgloss.NewStyle().Foreground(colorCyan)
	stylePane := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(colorGray).Padding(1, 2)
	styleSelectedNode := lipgloss.NewStyle().Background(lipgloss.Color("#2a2a2a")).Foreground(colorGold).Bold(true)
	styleNormalNode := lipgloss.NewStyle().Foreground(colorWhite)
	styleDetailTitle := lipgloss.NewStyle().Bold(true).Foreground(colorCyan).Underline(true)
	styleDetailLabel := lipgloss.NewStyle().Foreground(colorGold).Bold(true)
	styleDetailVal := lipgloss.NewStyle().Foreground(colorWhite)
	styleBadge := lipgloss.NewStyle().Background(colorGreen).Foreground(lipgloss.Color("#000000")).Bold(true).Padding(0, 1)
	styleFooter := lipgloss.NewStyle().Foreground(colorGray)

	// Render Header
	headerStr := styleTitle.Render("AUTODEV DEPENDENCY TUI SCANNER")
	headerDesc := styleHeaderDesc.Render("  Interactively explore directory hierarchy and configuration dependencies")
	header := lipgloss.JoinVertical(lipgloss.Left, headerStr, headerDesc, "")

	// Pane Heights
	paneHeight := m.height - 8
	if paneHeight < 5 {
		paneHeight = 5
	}

	// Left Pane: Interactive Dependency Tree
	var treeLines []string
	maxVisibleHeight := paneHeight - 2

	visibleEnd := m.scrollOffset + maxVisibleHeight
	if visibleEnd > len(m.flatNodes) {
		visibleEnd = len(m.flatNodes)
	}

	for idx := m.scrollOffset; idx < visibleEnd; idx++ {
		node := m.flatNodes[idx]

		// Construct prefix guidelines representing tree depth
		indent := strings.Repeat("  ", node.Level)
		prefix := ""
		if node.IsDir {
			if node.IsOpen {
				prefix = ""
			} else {
				prefix = ""
			}
		}

		lineContent := fmt.Sprintf("%s%s%s", indent, prefix, node.Name)
		if node.Tech != "" {
			lineContent += fmt.Sprintf(" [%s]", node.Tech)
		}

		if idx == m.cursor {
			treeLines = append(treeLines, styleSelectedNode.Render(lineContent))
		} else {
			treeLines = append(treeLines, styleNormalNode.Render(lineContent))
		}
	}

	leftPaneContent := strings.Join(treeLines, "\n")
	leftPaneWidth := (m.width - 6) * 55 / 100
	if leftPaneWidth < 25 {
		leftPaneWidth = 25
	}
	leftPaneStyled := stylePane.Width(leftPaneWidth).Height(paneHeight).Render(leftPaneContent)

	// Right Pane: Details & Dependencies
	selectedNode := m.flatNodes[m.cursor]
	var detailLines []string

	detailLines = append(detailLines, styleDetailTitle.Render("SELECTED ELEMENT INFO"), "")
	detailLines = append(detailLines, fmt.Sprintf("%s %s", styleDetailLabel.Render("Name:"), styleDetailVal.Render(selectedNode.Name)))

	nodeTypeStr := "File"
	if selectedNode.IsDir {
		nodeTypeStr = "Directory"
	}
	detailLines = append(detailLines, fmt.Sprintf("%s %s", styleDetailLabel.Render("Type:"), styleDetailVal.Render(nodeTypeStr)))

	// File Size formatting
	sizeStr := fmt.Sprintf("%d B", selectedNode.Size)
	if selectedNode.Size > 1024*1024 {
		sizeStr = fmt.Sprintf("%.2f MB", float64(selectedNode.Size)/(1024*1024))
	} else if selectedNode.Size > 1024 {
		sizeStr = fmt.Sprintf("%.2f KB", float64(selectedNode.Size)/1024)
	}
	detailLines = append(detailLines, fmt.Sprintf("%s %s", styleDetailLabel.Render("Size:"), styleDetailVal.Render(sizeStr)))

	relPath, _ := filepath.Rel(m.rootPath, selectedNode.Path)
	detailLines = append(detailLines, fmt.Sprintf("%s %s", styleDetailLabel.Render("Path:"), styleDetailVal.Render(relPath)))

	if selectedNode.Tech != "" {
		detailLines = append(detailLines, fmt.Sprintf("%s %s", styleDetailLabel.Render("Detected Tech:"), styleBadge.Render(selectedNode.Tech)))
	}

	detailLines = append(detailLines, "")

	if selectedNode.IsDir {
		detailLines = append(detailLines, styleDetailLabel.Render(fmt.Sprintf("Contains %d items", len(selectedNode.Children))))
		for i, child := range selectedNode.Children {
			if i >= 10 {
				detailLines = append(detailLines, styleDetailVal.Render(fmt.Sprintf("  ... and %d more", len(selectedNode.Children)-10)))
				break
			}
			childType := ""
			if child.IsDir {
				childType = ""
			}
			detailLines = append(detailLines, styleDetailVal.Render(fmt.Sprintf("  %s %s", childType, child.Name)))
		}
	} else {
		if len(selectedNode.Deps) > 0 {
			detailLines = append(detailLines, styleDetailTitle.Render("RESOLVED DEPENDENCIES"), "")
			for i, dep := range selectedNode.Deps {
				if i >= (paneHeight - len(detailLines) - 5) {
					detailLines = append(detailLines, styleDetailVal.Render(fmt.Sprintf("  ... and %d more", len(selectedNode.Deps)-i)))
					break
				}
				detailLines = append(detailLines, styleDetailVal.Render("  - "+dep))
			}
		} else {
			detailLines = append(detailLines, styleDetailVal.Render("No configuration dependencies found in this file."))
		}
	}

	rightPaneContent := strings.Join(detailLines, "\n")
	rightPaneWidth := m.width - leftPaneWidth - 8
	if rightPaneWidth < 20 {
		rightPaneWidth = 20
	}
	rightPaneStyled := stylePane.Width(rightPaneWidth).Height(paneHeight).Render(rightPaneContent)

	// Combine Left & Right panes horizontally
	mainBody := lipgloss.JoinHorizontal(lipgloss.Top, leftPaneStyled, rightPaneStyled)

	// Render Footer
	footer := styleFooter.Render("  ↑/↓: Navigate | Enter/Space: Expand/Collapse folder | q: Exit TUI")

	return lipgloss.JoinVertical(lipgloss.Left, header, mainBody, "", footer)
}
