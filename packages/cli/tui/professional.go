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

// The professional TUI is the first-run product surface. The legacy catalog
// browser remains available in the package for compatibility, while this model
// provides a guided detect -> plan -> execute -> verify workflow.
var (
	pcBg      = lipgloss.Color("#000000")
	pcSurface = lipgloss.Color("#0C0C0C")
	pcBorder  = lipgloss.Color("#222222")
	pcText    = lipgloss.Color("#FFFFFF")
	pcMuted   = lipgloss.Color("#888888")
	pcAccent  = lipgloss.Color("#FFD700")
	pcGreen   = lipgloss.Color("#00FF87")
	pcRed     = lipgloss.Color("#FF5F56")

	pcTitle       = lipgloss.NewStyle().Bold(true).Foreground(pcText)
	pcMutedStyle  = lipgloss.NewStyle().Foreground(pcMuted)
	pcAccentStyle = lipgloss.NewStyle().Bold(true).Foreground(pcAccent)
	pcSelected    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#000000")).Background(pcAccent).Padding(0, 1)
	pcCard        = lipgloss.NewStyle().Background(pcSurface).Border(lipgloss.RoundedBorder()).BorderForeground(pcBorder).Padding(1, 2)
	pcPill        = lipgloss.NewStyle().Background(lipgloss.Color("#1A1A1A")).Foreground(pcAccent).Padding(0, 1)
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

type pcMenuItem struct{ label, description, action string }

var pcMenu = []pcMenuItem{
	{"Setup this machine", "Install languages, frameworks, databases and developer tools", "catalog"},
	{"Use a developer profile", "Bootstrap a focused environment in one step", "profiles"},
	{"Run diagnostics", "Check your project and development environment", "doctor"},
	{"Scan this project", "Understand dependencies, tooling and configuration", "scan"},
	{"Open AI agent", "Start a new OpenCode session aware of the whole codebase", "opencode"},
	{"Command center", "Discover every AutoDev capability without memorizing commands", "commands"},
	{"Quit", "Exit AutoDev", "quit"},
}

type professionalModel struct {
	catalog                                         *catalog.Catalog
	sys                                             *osinfo.Info
	screen                                          pcScreen
	cursor, catCursor, profileCursor, commandCursor int
	catName                                         string
	packages                                        []*catalog.Package
	selected                                        map[string]bool
	queue                                           []*catalog.Package
	installIndex                                    int
	success, failed                                 []string
	started                                         time.Time
	width                                           int
	errorText                                       string
}

type pcInstallMsg struct {
	pkg *catalog.Package
	err error
}
type pcExecMsg struct{ err error }

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
		if msg.err != nil {
			m.failed = append(m.failed, msg.pkg.Name)
		} else {
			m.success = append(m.success, msg.pkg.Name)
		}
		m.installIndex++
		if m.installIndex >= len(m.queue) {
			m.screen = pcComplete
			return m, nil
		}
		return m, m.installNext()
	case pcExecMsg:
		m.screen = pcHome
		if msg.err != nil {
			m.errorText = friendlyExecError(msg.err)
		}
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m professionalModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()
	if key == "ctrl+c" || (m.screen == pcHome && key == "q") {
		return m, tea.Quit
	}
	switch m.screen {
	case pcHome:
		switch key {
		case "up", "k":
			m.cursor = (m.cursor + len(pcMenu) - 1) % len(pcMenu)
		case "down", "j":
			m.cursor = (m.cursor + 1) % len(pcMenu)
		case "enter", " ":
			switch pcMenu[m.cursor].action {
			case "catalog":
				m.openCategory("Languages")
			case "profiles":
				m.screen = pcProfiles
				m.profileCursor = 0
			case "commands":
				m.screen = pcCommands
				m.commandCursor = 0
			case "doctor", "scan", "opencode":
				return m, m.execAutoDev(pcMenu[m.cursor].action)
			case "quit":
				return m, tea.Quit
			}
		}
	case pcCatalog:
		switch key {
		case "esc":
			m.screen = pcHome
		case "left", "h":
			m.openCategory(m.prevCategory())
		case "right", "l":
			m.openCategory(m.nextCategory())
		case "up", "k":
			if len(m.packages) > 0 {
				m.catCursor = (m.catCursor + len(m.packages) - 1) % len(m.packages)
			}
		case "down", "j":
			if len(m.packages) > 0 {
				m.catCursor = (m.catCursor + 1) % len(m.packages)
			}
		case " ":
			if len(m.packages) > 0 {
				p := m.packages[m.catCursor]
				m.selected[p.ID] = !m.selected[p.ID]
			}
		case "a":
			for _, p := range m.packages {
				m.selected[p.ID] = true
			}
		case "n":
			for _, p := range m.packages {
				m.selected[p.ID] = false
			}
		case "enter":
			ids := selectedIDs(m.selected)
			if len(ids) == 0 {
				m.errorText = "Select at least one package before continuing."
				return m, nil
			}
			resolved, err := m.catalog.Resolve(ids)
			if err != nil {
				m.errorText = err.Error()
				return m, nil
			}
			m.queue, m.errorText, m.screen = resolved, "", pcConfirm
		}
	case pcProfiles:
		switch key {
		case "esc":
			m.screen = pcHome
		case "up", "k":
			if len(m.catalog.Profiles) > 0 {
				m.profileCursor = (m.profileCursor + len(m.catalog.Profiles) - 1) % len(m.catalog.Profiles)
			}
		case "down", "j":
			if len(m.catalog.Profiles) > 0 {
				m.profileCursor = (m.profileCursor + 1) % len(m.catalog.Profiles)
			}
		case "enter":
			if len(m.catalog.Profiles) > 0 {
				p := m.catalog.Profiles[m.profileCursor]
				resolved, err := m.catalog.Resolve(p.Packages)
				if err != nil {
					m.errorText = err.Error()
					return m, nil
				}
				m.queue, m.errorText, m.screen = resolved, "", pcConfirm
			}
		}
	case pcCommands:
		commands := pcCommandsList()
		switch key {
		case "esc":
			m.screen = pcHome
		case "up", "k":
			if len(commands) > 0 {
				m.commandCursor = (m.commandCursor + len(commands) - 1) % len(commands)
			}
		case "down", "j":
			if len(commands) > 0 {
				m.commandCursor = (m.commandCursor + 1) % len(commands)
			}
		case "enter":
			if len(commands) > 0 {
				return m, m.execAutoDev(commands[m.commandCursor].action)
			}
		}
	case pcConfirm:
		switch key {
		case "esc", "n":
			m.screen = pcHome
		case "enter", "y":
			m.installIndex, m.success, m.failed, m.started, m.screen = 0, nil, nil, time.Now(), pcInstalling
			return m, m.installNext()
		}
	case pcComplete:
		if key == "enter" || key == "esc" {
			m.screen = pcHome
		}
	}
	return m, nil
}

func (m *professionalModel) openCategory(name string) {
	cats := m.categories()
	if len(cats) == 0 {
		m.packages = nil
		return
	}
	m.catName = name
	if _, ok := m.catalog.ByCategory()[name]; !ok {
		m.catName = cats[0]
	}
	m.packages = m.catalog.ByCategory()[m.catName]
	m.catCursor = 0
	m.errorText = ""
	m.screen = pcCatalog
}

func (m professionalModel) categories() []string {
	cats := make([]string, 0)
	for name := range m.catalog.ByCategory() {
		cats = append(cats, name)
	}
	sort.Strings(cats)
	return cats
}

func (m professionalModel) prevCategory() string { return m.shiftCategory(-1) }
func (m professionalModel) nextCategory() string { return m.shiftCategory(1) }
func (m professionalModel) shiftCategory(delta int) string {
	cats := m.categories()
	if len(cats) == 0 {
		return m.catName
	}
	idx := sort.SearchStrings(cats, m.catName)
	if idx >= len(cats) || cats[idx] != m.catName {
		idx = 0
	}
	idx = (idx + delta + len(cats)) % len(cats)
	return cats[idx]
}

func (m professionalModel) installNext() tea.Cmd {
	if m.installIndex >= len(m.queue) {
		return nil
	}
	pkg := m.queue[m.installIndex]
	if pkg.IsInstalled() {
		return func() tea.Msg { return pcInstallMsg{pkg: pkg} }
	}
	return tea.ExecProcess(buildInstallCmd(pkg), func(err error) tea.Msg { return pcInstallMsg{pkg: pkg, err: err} })
}

func (m professionalModel) execAutoDev(command string) tea.Cmd {
	return tea.ExecProcess(exec.Command(os.Args[0], command), func(err error) tea.Msg { return pcExecMsg{err: err} })
}

// friendlyExecError turns a raw subprocess error (e.g. "exit status 1") into a
// readable message that points the user at the failing command.
func friendlyExecError(err error) string {
	if err == nil {
		return ""
	}
	msg := strings.TrimSpace(err.Error())
	if strings.HasPrefix(msg, "exit status") {
		return "The command exited with an error. Run it directly for the full output."
	}
	return msg
}

func (m professionalModel) View() string {
	if m.width < 70 {
		return lipgloss.NewStyle().Background(pcBg).Foreground(pcText).Padding(1, 2).Render("AutoDev\n\nPlease widen your terminal to at least 70 columns.")
	}
	var body string
	switch m.screen {
	case pcHome:
		body = m.viewHome()
	case pcCatalog:
		body = m.viewCatalog()
	case pcProfiles:
		body = m.viewProfiles()
	case pcCommands:
		body = m.viewCommands()
	case pcConfirm:
		body = m.viewConfirm()
	case pcInstalling:
		body = m.viewInstalling()
	case pcComplete:
		body = m.viewComplete()
	}
	return lipgloss.NewStyle().Background(pcBg).Foreground(pcText).Padding(1, 2).Width(m.width).Render(m.header() + "\n\n" + body)
}

func (m professionalModel) header() string {
	platform := runtime.GOOS + "/" + runtime.GOARCH
	if m.sys != nil {
		platform = fmt.Sprintf("%s • %s", m.sys.Version, m.sys.Arch)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, pcTitle.Render("◆ AutoDev"), "  ", pcPill.Render(platform), "  ", pcMutedStyle.Render("Developer Environment Control Center"))
}

func (m professionalModel) viewHome() string {
	selected := len(selectedIDs(m.selected))
	intro := pcCard.Render(pcTitle.Render("Welcome to AutoDev") + "\n" + pcMutedStyle.Render("Set up, inspect and operate your development environment from one place.") + "\n\n" + pcAccentStyle.Render("Detect  →  Plan  →  Execute  →  Verify"))
	items := make([]string, 0, len(pcMenu))
	for i, item := range pcMenu {
		line := item.label + "\n" + pcMutedStyle.Render(item.description)
		if i == m.cursor {
			line = pcSelected.Render(item.label) + "\n" + pcMutedStyle.Render(item.description)
		}
		items = append(items, line)
	}
	status := pcMutedStyle.Render(fmt.Sprintf("%d selected • %s • ↑↓ navigate • enter select • q quit", selected, platformLabel(m.sys)))
	if m.errorText != "" {
		status = lipgloss.NewStyle().Foreground(pcRed).Render("! " + m.errorText)
	}
	return intro + "\n\n" + strings.Join(items, "\n\n") + "\n\n" + status
}

func (m professionalModel) viewCatalog() string {
	count := len(selectedIDs(m.selected))
	catLine := pcAccentStyle.Render("Category: "+m.catName) + "  " + pcMutedStyle.Render(fmt.Sprintf("%d selected", count))
	var rows []string
	for i, p := range m.packages {
		mark := "○"
		if m.selected[p.ID] {
			mark = lipgloss.NewStyle().Foreground(pcGreen).Render("●")
		}
		line := fmt.Sprintf("%s  %-22s %s", mark, p.Name, pcMutedStyle.Render(p.Description))
		if i == m.catCursor {
			line = pcSelected.Render(fmt.Sprintf("%s  %s", mark, p.Name)) + "  " + pcMutedStyle.Render(p.Description)
		}
		rows = append(rows, line)
	}
	return pcCard.Render(pcTitle.Render("Tool catalog") + "\n" + catLine + "\n" + pcMutedStyle.Render("←/→ category • ↑/↓ package • space select • a all • n none • enter review • esc back") + "\n\n" + strings.Join(rows, "\n"))
}

func (m professionalModel) viewProfiles() string {
	var rows []string
	for i, p := range m.catalog.Profiles {
		line := fmt.Sprintf("%s\n%s", p.Name, pcMutedStyle.Render(p.Description))
		if i == m.profileCursor {
			line = pcSelected.Render(p.Name) + "\n" + pcMutedStyle.Render(p.Description)
		}
		rows = append(rows, line)
	}
	return pcCard.Render(pcTitle.Render("Developer profiles") + "\n" + pcMutedStyle.Render("Choose a known-good environment configuration. Enter to review before installing.") + "\n\n" + strings.Join(rows, "\n\n"))
}

func (m professionalModel) viewCommands() string {
	commands := pcCommandsList()
	var rows []string
	for i, c := range commands {
		line = fmt.Sprintf("%-18s %s", c.action, pcMutedStyle.Render(c.description))
		if i == m.commandCursor {
			line = pcSelected.Render(c.action) + " " + pcMutedStyle.Render(c.description)
		}
		rows = append(rows, line)
	}
	return pcCard.Render(pcTitle.Render("AutoDev command center") + "\n" + pcMutedStyle.Render("Select a capability and AutoDev launches the real command for you.") + "\n\n" + strings.Join(rows, "\n"))
}

func (m professionalModel) viewConfirm() string {
	var rows []string
	for _, p := range m.queue {
		state := "install"
		if p.IsInstalled() {
			state = "already installed — verify"
		}
		rows = append(rows, fmt.Sprintf("• %-24s %s", p.Name, pcMutedStyle.Render(state)))
	}
	return pcCard.Render(pcTitle.Render(fmt.Sprintf("Review setup • %d actions", len(m.queue))) + "\n" + pcMutedStyle.Render("Nothing changes until you confirm. Dependencies have already been resolved.") + "\n\n" + strings.Join(rows, "\n") + "\n\n" + pcAccentStyle.Render("Enter / y: install") + "  " + pcMutedStyle.Render("Esc / n: cancel"))
}

func (m professionalModel) viewInstalling() string {
	total := len(m.queue)
	if total == 0 {
		total = 1
	}
	done := m.installIndex
	pct := done * 100 / total
	width := 36
	filled := pct * width / 100
	bar := strings.Repeat("█", filled) + strings.Repeat("·", width-filled)
	current := "Preparing…"
	if done < len(m.queue) {
		current = m.queue[done].Name
	}
	return pcCard.Render(pcTitle.Render("Setting up your environment") + "\n\n" + pcAccentStyle.Render("Step "+fmt.Sprint(min(done+1, total))+" of "+fmt.Sprint(total)) + "  " + current + "\n\n" + bar + fmt.Sprintf("  %d%%", pct) + "\n" + pcMutedStyle.Render("Installer output is shown directly in the terminal when required (including sudo prompts)."))
}

func (m professionalModel) viewComplete() string {
	return pcCard.Render(pcTitle.Render("Environment setup complete") + "\n\n" + lipgloss.NewStyle().Foreground(pcGreen).Render(fmt.Sprintf("✓ %d succeeded", len(m.success))) + "   " + lipgloss.NewStyle().Foreground(pcRed).Render(fmt.Sprintf("✕ %d failed", len(m.failed))) + "\n" + pcMutedStyle.Render(fmt.Sprintf("Elapsed: %s", time.Since(m.started).Round(time.Second))) + "\n\n" + pcMutedStyle.Render("Recommended next commands:") + "\n  autodev doctor   — verify the environment\n  autodev scan     — understand the current project\n  autodev agent    — start the AI workflow\n\n" + pcMutedStyle.Render("Enter / Esc to return to the command center"))
}

func selectedIDs(selected map[string]bool) []string {
	ids := []string{}
	for id, ok := range selected {
		if ok {
			ids = append(ids, id)
		}
	}
	sort.Strings(ids)
	return ids
}
func platformLabel(info *osinfo.Info) string {
	if info == nil {
		return runtime.GOOS + "/" + runtime.GOARCH
	}
	return info.Version
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type pcCommand struct{ action, description string }

func pcCommandsList() []pcCommand {
	return []pcCommand{
		{"doctor", "Diagnose the current environment"},
		{"scan", "Analyze the current project"},
		{"setup", "Open the environment setup workflow"},
		{"opencode", "Start an OpenCode session with shared codebase knowledge"},
		{"audit", "Run a security audit"},
		{"upgrade", "Upgrade AutoDev and managed tools"},
		{"skills", "Browse installed skills"},
		{"mcp", "Manage MCP integrations"},
		{"create", "Create a new project"},
		{"containerize", "Generate container configuration"},
	}
}

func RunProfessional(c *catalog.Catalog) error {
	_, err := tea.NewProgram(newProfessionalModel(c), tea.WithAltScreen()).Run()
	return err
}
