package tui

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/autodev-sh/autodev/catalog"
	"github.com/autodev-sh/autodev/core/osinfo"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Professional Command Center is intentionally separate from the legacy catalog
// browser so the onboarding experience can evolve without destabilizing the
// existing installer implementation.

var (
	pcBg      = lipgloss.Color("#0B1020")
	pcSurface = lipgloss.Color("#11182B")
	pcBorder  = lipgloss.Color("#26324D")
	pcText    = lipgloss.Color("#F8FAFC")
	pcMuted   = lipgloss.Color("#94A3B8")
	pcAccent  = lipgloss.Color("#7C3AED")
	pcCyan    = lipgloss.Color("#22D3EE")
	pcGreen   = lipgloss.Color("#34D399")
	pcRed     = lipgloss.Color("#FB7185")

	pcTitle = lipgloss.NewStyle().Bold(true).Foreground(pcText)
	pcMutedStyle = lipgloss.NewStyle().Foreground(pcMuted)
	pcAccentStyle = lipgloss.NewStyle().Bold(true).Foreground(pcAccent)
	pcSelected = lipgloss.NewStyle().Bold(true).Foreground(pcText).Background(pcAccent).Padding(0, 1)
	pcCard = lipgloss.NewStyle().Background(pcSurface).Border(lipgloss.RoundedBorder()).BorderForeground(pcBorder).Padding(1, 2)
	pcPill = lipgloss.NewStyle().Background(lipgloss.Color("#17213A")).Foreground(pcCyan).Padding(0, 1)
)

type pcScreen int

const (
	pcHome pcScreen = iota
	pcCatalog
	pcProfiles
	pcCommands
	pcConfirm
	pcInstalling
	pcComplete
)

type pcMenuItem struct { label, description, action string }

var pcMenu = []pcMenuItem{
	{"Setup this machine", "Install languages, frameworks, databases and developer tools", "catalog"},
	{"Use a developer profile", "Bootstrap a focused environment in one step", "profiles"},
	{"Run diagnostics", "Check your project and development environment", "doctor"},
	{"Scan this project", "Understand dependencies, tooling and configuration", "scan"},
	{"Open AI agent", "Let AutoDev inspect and operate on your workspace", "agent"},
	{"Command center", "Discover every AutoDev capability without memorizing commands", "commands"},
	{"Quit", "Exit AutoDev", "quit"},
}

type professionalModel struct {
	catalog *catalog.Catalog
	sys *osinfo.Info
	screen pcScreen
	cursor int
	catCursor int
	profileCursor int
	commandCursor int
	catName string
	packages []*catalog.Package
	selected map[string]bool
	queue []*catalog.Package
	installIndex int
	success []string
	failed []string
	started time.Time
	width int
	errorText string
}

type pcInstallMsg struct { pkg *catalog.Package; err error }

type pcExecMsg struct { err error }

func newProfessionalModel(c *catalog.Catalog) professionalModel {
	sys, _ := osinfo.Detect()
	return professionalModel{catalog: c, sys: sys, screen: pcHome, selected: map[string]bool{}}
}

func (m professionalModel) Init() tea.Cmd { return nil }

func (m professionalModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case pcInstallMsg:
		if msg.err != nil { m.failed = append(m.failed, msg.pkg.Name) } else { m.success = append(m.success, msg.pkg.Name) }
		m.installIndex++
		if m.installIndex >= len(m.queue) { m.screen = pcComplete; return m, nil }
		return m, m.installNext()
	case pcExecMsg:
		m.screen = pcHome
		if msg.err != nil { m.errorText = msg.err.Error() }
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m professionalModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()
	if key == "ctrl+c" { return m, tea.Quit }
	switch m.screen {
	case pcHome:
		switch key {
		case "up", "k": m.cursor = (m.cursor + len(pcMenu) - 1) % len(pcMenu)
		case "down", "j": m.cursor = (m.cursor + 1) % len(pcMenu)
		case "enter", " ":
			switch pcMenu[m.cursor].action {
			case "catalog": m.screen = pcCatalog; m.catName = "Languages"; m.packages = m.catalog.ByCategory()[m.catName]; m.catCursor = 0
			case "profiles": m.screen = pcProfiles; m.profileCursor = 0
			case "commands": m.screen = pcCommands; m.commandCursor = 0
			case "doctor", "scan", "agent": return m, m.execAutoDev(pcMenu[m.cursor].action)
			case "quit": return m, tea.Quit
			}
		}
	case pcCatalog:
		switch key {
		case "esc": m.screen = pcHome
		case "left", "h": m.catName = pcPrevCategory(m.catName); m.packages = m.catalog.ByCategory()[m.catName]; m.catCursor = 0
		case "right", "l": m.catName = pcNextCategory(m.catName); m.packages = m.catalog.ByCategory()[m.catName]; m.catCursor = 0
		case "up", "k": if len(m.packages) > 0 { m.catCursor = (m.catCursor + len(m.packages) - 1) % len(m.packages) }
		case "down", "j": if len(m.packages) > 0 { m.catCursor = (m.catCursor + 1) % len(m.packages) }
		case "space": if len(m.packages) > 0 { p := m.packages[m.catCursor]; m.selected[p.ID] = !m.selected[p.ID] }
		case "a": for _, p := range m.packages { m.selected[p.ID] = true }
		case "n": for _, p := range m.packages { m.selected[p.ID] = false }
		case "enter":
			ids := selectedIDs(m.selected); if len(ids) == 0 { return m, nil }
			resolved, err := m.catalog.Resolve(ids); if err != nil { m.errorText = err.Error(); return m, nil }
			m.queue = resolved; m.screen = pcConfirm
		}
	case pcProfiles:
		switch key {
		case "esc": m.screen = pcHome
		case "up", "k": if len(m.catalog.Profiles) > 0 { m.profileCursor = (m.profileCursor + len(m.catalog.Profiles) - 1) % len(m.catalog.Profiles) }
		case "down", "j": if len(m.catalog.Profiles) > 0 { m.profileCursor = (m.profileCursor + 1) % len(m.catalog.Profiles) }
		case "enter":
			if len(m.catalog.Profiles) > 0 { p := m.catalog.Profiles[m.profileCursor]; resolved, err := m.catalog.Resolve(p.Packages); if err != nil { m.errorText = err.Error(); return m, nil }; m.queue = resolved; m.screen = pcConfirm }
		}
	case pcCommands:
		commands := pcCommandsList()
		switch key {
		case "esc": m.screen = pcHome
		case "up", "k": if len(commands) > 0 { m.commandCursor = (m.commandCursor + len(commands) - 1) % len(commands) }
		case "down", "j": if len(commands) > 0 { m.commandCursor = (m.commandCursor + 1) % len(commands) }
		case "enter": if len(commands) > 0 { return m, m.execAutoDev(commands[m.commandCursor].action) }
		}
	case pcConfirm:
		switch key { case "esc", "n": m.screen = pcHome; case "enter", "y": m.installIndex = 0; m.success = nil; m.failed = nil; m.started = time.Now(); m.screen = pcInstalling; return m, m.installNext() }
	case pcComplete:
		if key == "enter" || key == "esc" { m.screen = pcHome }
	case pcInstalling:
		// Installation is intentionally not cancellable while an installer owns the terminal.
	}
	return m, nil
}

func (m professionalModel) installNext() tea.Cmd {
	if m.installIndex >= len(m.queue) { return nil }
	pkg := m.queue[m.installIndex]
	if pkg.IsInstalled() { return func() tea.Msg { return pcInstallMsg{pkg: pkg} } }
	cmd := buildInstallCmd(pkg)
	return tea.ExecProcess(cmd, func(err error) tea.Msg { return pcInstallMsg{pkg: pkg, err: err} })
}

func (m professionalModel) execAutoDev(command string) tea.Cmd {
	return tea.ExecProcess(exec.Command(os.Args[0], command), func(err error) tea.Msg { return pcExecMsg{err: err} })
}

func (m professionalModel) View() string {
	if m.width < 70 { return lipgloss.NewStyle().Background(pcBg).Foreground(pcText).Padding(1, 2).Render("AutoDev\n\nPlease widen your terminal to at least 70 columns.") }
	var body string
	switch m.screen { case pcHome: body = m.viewHome(); case pcCatalog: body = m.viewCatalog(); case pcProfiles: body = m.viewProfiles(); case pcCommands: body = m.viewCommands(); case pcConfirm: body = m.viewConfirm(); case pcInstalling: body = m.viewInstalling(); case pcComplete: body = m.viewComplete() }
	return lipgloss.NewStyle().Background(pcBg).Foreground(pcText).Padding(1, 2).Width(m.width).Render(m.header() + "\n" + body)
}

func (m professionalModel) header() string {
	platform := runtime.GOOS + "/" + runtime.GOARCH
	if m.sys != nil { platform = fmt.Sprintf("%s • %s", m.sys.Version, m.sys.Arch) }
	return lipgloss.JoinHorizontal(lipgloss.Top, pcTitle.Render("◆ AutoDev"), "  ", pcPill.Render(platform), "  ", pcMutedStyle.Render("Developer Environment Control Center"))
}

func (m professionalModel) viewHome() string {
	selected := len(selectedIDs(m.selected))
	intro := pcCard.Render(pcTitle.Render("Welcome to AutoDev") + "\n" + pcMutedStyle.Render("Set up, inspect and operate your development environment from one place.") + "\n\n" + pcAccentStyle.Render("Detect  →  Plan  →  Execute  →  Verify"))
	items := make([]string, 0, len(pcMenu))
	for i, item := range pcMenu { line := item.label + "\n" + pcMutedStyle.Render(item.description); if i == m.cursor { line = pcSelected.Render(item.label) + "\n" + pcMutedStyle.Render(item.description) }; items = append(items, line) }
	status := pcMutedStyle.Render(fmt.Sprintf("%d selected  •  %s  •  ↑↓ navigate  enter select  q/ctrl-c quit", selected, platformLabel(m.sys)))
	if m.errorText != "" { status = lipgloss.NewStyle().Foreground(pcRed).Render("! " + m.errorText) }
	return intro + "\n\n" + strings.Join(items, "\n\n") + "\n\n" + status
}

func (m professionalModel) viewCatalog() string {
	cats := pcCategories(m.catalog)
	count := 0; for _, v := range m.selected { if v { count++ } }
	left := pcTitle.Render("Tool catalog") + "\n" + pcMutedStyle.Render("←/→ category  ↑/↓ package  space select  enter review") + "\n\n" + strings.Join(cats, "  ")
	var rows []string
	for i, p := range m.packages { mark := "○"; if m.selected[p.ID] { mark = lipgloss.NewStyle().Foreground(pcGreen).Render("●") }; line := fmt.Sprintf("%s  %-22s %s", mark, p.Name, pcMutedStyle.Render(p.Description)); if i == m.catCursor { line = pcSelected.Render(fmt.Sprintf("%s  %s", mark, p.Name)) + "  " + pcMutedStyle.Render(p.Description) }; rows = append(rows, line) }
	return pcCard.Render(left) + "\n\n" + strings.Join(rows, "\n")
}

func (m professionalModel) viewProfiles() string {
	var rows []string
	for i, p := range m.catalog.Profiles { line := fmt.Sprintf("%s\n%s", p.Name, pcMutedStyle.Render(p.Description)); if i == m.profileCursor { line = pcSelected.Render(p.Name) + "\n" + pcMutedStyle.Render(p.Description) }; rows = append(rows, line) }
	return pcCard.Render(pcTitle.Render("Developer profiles") + "\n" + pcMutedStyle.Render("Choose a known-good environment configuration.") + "\n\n" + strings.Join(rows, "\n\n"))
}

func (m professionalModel) viewCommands() string {
	commands := pcCommandsList(); var rows []string
	for i, c := range commands { line := fmt.Sprintf("%-18s %s", c.action, pcMutedStyle.Render(c.description)); if i == m.commandCursor { line = pcSelected.Render(c.action) + " " + pcMutedStyle.Render(c.description) }; rows = append(rows, line) }
	return pcCard.Render(pcTitle.Render("AutoDev command center") + "\n" + pcMutedStyle.Render("You do not need to memorize the CLI. Select a capability and AutoDev will launch it.") + "\n\n" + strings.Join(rows, "\n"))
}

func (m professionalModel) viewConfirm() string {
	var rows []string; for _, p := range m.queue { rows = append(rows, "• "+p.Name) }
	return pcCard.Render(pcTitle.Render(fmt.Sprintf("Review setup • %d actions", len(m.queue))) + "\n" + pcMutedStyle.Render("Nothing changes until you confirm.") + "\n\n" + strings.Join(rows, "\n") + "\n\n" + pcAccentStyle.Render("Press Enter to install") + "  " + pcMutedStyle.Render("Esc to cancel"))
}

func (m professionalModel) viewInstalling() string {
	done := m.installIndex; total := len(m.queue); if total == 0 { total = 1 }
	pct := done * 100 / total; width := 32; filled := pct * width / 100
	bar := strings.Repeat("█", filled) + strings.Repeat("·", width-filled)
	current := "Preparing..."; if done < len(m.queue) { current = "Installing " + m.queue[done].Name }
	return pcCard.Render(pcTitle.Render("Setting up your environment") + "\n\n" + current + "\n\n" + pcAccentStyle.Render(bar) + fmt.Sprintf("  %d%%", pct) + "\n" + pcMutedStyle.Render(fmt.Sprintf("Step %d of %d  •  %s", min(done+1,total), total, time.Since(m.started).Round(time.Second))))
}

func (m professionalModel) viewComplete() string {
	elapsed := time.Since(m.started).Round(time.Second)
	return pcCard.Render(pcTitle.Render("Environment setup complete") + "\n\n" + lipgloss.NewStyle().Foreground(pcGreen).Render(fmt.Sprintf("✓ %d succeeded", len(m.success))) + "   " + lipgloss.NewStyle().Foreground(pcRed).Render(fmt.Sprintf("✕ %d failed", len(m.failed))) + "\n" + pcMutedStyle.Render(fmt.Sprintf("Elapsed: %s", elapsed)) + "\n\n" + pcMutedStyle.Render("Next steps:") + "\n  autodev doctor\n  autodev scan\n  autodev agent\n\n" + pcMutedStyle.Render("Press Enter to return to the command center."))
}

func pcCategories(c *catalog.Catalog) []string { names := make([]string, 0); for k := range c.ByCategory() { names = append(names, k) }; sort.Strings(names); return names }
func pcPrevCategory(cur string) string { return cur }
func pcNextCategory(cur string) string { return cur }
func selectedIDs(m map[string]bool) []string { ids:=[]string{}; for id, ok := range m { if ok { ids=append(ids,id) } }; sort.Strings(ids); return ids }
func platformLabel(info *osinfo.Info) string { if info == nil { return runtime.GOOS + "/" + runtime.GOARCH }; return info.Version }
func min(a,b int) int { if a < b { return a }; return b }

type pcCommand struct { action, description string }
func pcCommandsList() []pcCommand { return []pcCommand{{"doctor","Diagnose the current environment"},{"scan","Analyze the current project"},{"setup","Open the setup workflow"},{"agent","Start the AI environment agent"},{"audit","Run a security audit"},{"upgrade","Upgrade AutoDev and tools"},{"skills","Browse installed skills"},{"mcp","Manage MCP integrations"},{"create","Create a new project"},{"containerize","Generate container configuration"}} }

// RunProfessional launches the redesigned first-run experience.
func RunProfessional(c *catalog.Catalog) error { _, err := tea.NewProgram(newProfessionalModel(c), tea.WithAltScreen()).Run(); return err }
