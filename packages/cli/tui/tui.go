// Package tui implements the AutoDev interactive terminal UI using BubbleTea.
// Running `autodev` with no args opens this UI — the "App Store for Developers".
package tui

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/autodev-sh/autodev/catalog"
	"github.com/autodev-sh/autodev/core/osinfo"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ── Styles ───────────────────────────────────────────────────────────────────

var (
	colorAccent = lipgloss.Color("#0070f3")
	colorGreen  = lipgloss.Color("#00c853")
	colorRed    = lipgloss.Color("#ee0000")
	colorBlue   = lipgloss.Color("#0070f3")
	colorGray   = lipgloss.Color("#555555")
	colorWhite  = lipgloss.Color("#ededed")

	styleSelected = lipgloss.NewStyle().Bold(true).Foreground(colorAccent)
	styleNormal   = lipgloss.NewStyle().Foreground(colorWhite)
	styleDim      = lipgloss.NewStyle().Foreground(colorGray)
	styleCheck    = lipgloss.NewStyle().Foreground(colorGreen).Bold(true)
	styleCategory = lipgloss.NewStyle().Bold(true).Foreground(colorBlue).MarginTop(1)
	styleBadge    = lipgloss.NewStyle().Background(colorAccent).Foreground(lipgloss.Color("#ffffff")).Padding(0, 1).Bold(true)
	styleSuccess  = lipgloss.NewStyle().Foreground(colorGreen).Bold(true)
	styleError    = lipgloss.NewStyle().Foreground(colorRed).Bold(true)
	styleHint     = lipgloss.NewStyle().Foreground(colorGray)
)

// ── State machine ─────────────────────────────────────────────────────────────

type screen int

const (
	screenMenu screen = iota
	screenCategory
	screenProfile
	screenConfirm
	screenDone
)

// ── Messages ─────────────────────────────────────────────────────────────────

type installDoneMsg struct {
	pkg *catalog.Package
	err error
}

// ── Model ─────────────────────────────────────────────────────────────────────

type Model struct {
	catalog *catalog.Catalog

	screen      screen
	menuCursor  int
	catCursor   int
	curCategory string
	catPkgs     []*catalog.Package

	selected map[string]bool

	// Install queue
	toInstall  []*catalog.Package
	installIdx int
	installLog []string

	// Track success/failure
	installedSuccess []string
	installedFailed  []string

	sysInfo *osinfo.Info
	width   int
}

var menuItems = []struct {
	label string
	icon  string
	value string
}{
	{"Languages", "", "cat:Languages"},
	{"Frameworks", "", "cat:Frameworks"},
	{"Databases", "", "cat:Databases"},
	{"DevOps", "", "cat:DevOps"},
	{"Mobile Development", "", "cat:Mobile"},
	{"AI / ML", "", "cat:AI/ML"},
	{"Tools", "", "cat:Tools"},
	{"", "", "---"},
	{"Install by Profile", "", "profiles"},
	{"", "", "---"},
	{"Start Installation", "", "install"},
	{"Quit", "", "quit"},
}

func New(c *catalog.Catalog) Model {
	sysInfo, _ := osinfo.Detect()
	return Model{
		catalog:  c,
		screen:   screenMenu,
		selected: map[string]bool{},
		sysInfo:  sysInfo,
	}
}

func (m Model) Init() tea.Cmd { return nil }

// ── Update ────────────────────────────────────────────────────────────────────

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width

	case tea.KeyMsg:
		switch m.screen {
		case screenMenu:
			return m.updateMenu(msg)
		case screenCategory:
			return m.updateCategory(msg)
		case screenProfile:
			return m.updateProfile(msg)
		case screenConfirm:
			return m.updateConfirm(msg)
		case screenDone:
			if msg.String() == "esc" || msg.String() == "enter" {
				return m, tea.Quit
			}
		}

	// ── KEY FIX: installDoneMsg is sent AFTER tea.ExecProcess finishes ──────
	case installDoneMsg:
		if msg.err != nil {
			m.installedFailed = append(m.installedFailed, msg.pkg.Name)
			m.installLog = append(m.installLog,
				styleError.Render(fmt.Sprintf("✗ %s: %v", msg.pkg.Name, msg.err)))
		} else {
			m.installedSuccess = append(m.installedSuccess, msg.pkg.Name)
			m.installLog = append(m.installLog,
				styleSuccess.Render(fmt.Sprintf("✓ %s installed", msg.pkg.Name)))
		}
		m.installIdx++
		if m.installIdx >= len(m.toInstall) {
			m.screen = screenDone
			return m, nil
		}
		// Run the next package via ExecProcess (suspends TUI → full stdin)
		return m, m.execNextInstall()
	}

	return m, nil
}

func (m Model) updateMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		m.menuCursor--
		if m.menuCursor < 0 {
			m.menuCursor = len(menuItems) - 1
		}
		for menuItems[m.menuCursor].value == "---" {
			m.menuCursor--
			if m.menuCursor < 0 {
				m.menuCursor = len(menuItems) - 1
			}
		}
	case "down", "j":
		m.menuCursor++
		if m.menuCursor >= len(menuItems) {
			m.menuCursor = 0
		}
		for menuItems[m.menuCursor].value == "---" {
			m.menuCursor++
			if m.menuCursor >= len(menuItems) {
				m.menuCursor = 0
			}
		}
	case "enter", " ":
		item := menuItems[m.menuCursor]
		switch {
		case strings.HasPrefix(item.value, "cat:"):
			m.curCategory = strings.TrimPrefix(item.value, "cat:")
			m.catPkgs = m.catalog.ByCategory()[m.curCategory]
			m.catCursor = 0
			m.screen = screenCategory
		case item.value == "profiles":
			m.catCursor = 0
			m.screen = screenProfile
		case item.value == "install":
			var ids []string
			for id, checked := range m.selected {
				if checked {
					ids = append(ids, id)
				}
			}
			if len(ids) == 0 {
				return m, nil
			}
			resolved, err := m.catalog.Resolve(ids)
			if err != nil {
				m.installLog = append(m.installLog,
					styleError.Render("Dependency error: "+err.Error()))
				return m, nil
			}
			m.toInstall = resolved
			m.screen = screenConfirm
		case item.value == "quit":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) updateCategory(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.screen = screenMenu
	case "up", "k":
		if m.catCursor > 0 {
			m.catCursor--
		}
	case "down", "j":
		if m.catCursor < len(m.catPkgs)-1 {
			m.catCursor++
		}
	case " ", "enter":
		if len(m.catPkgs) > 0 {
			pkg := m.catPkgs[m.catCursor]
			m.selected[pkg.ID] = !m.selected[pkg.ID]
		}
	case "a":
		for _, p := range m.catPkgs {
			m.selected[p.ID] = true
		}
	case "n":
		for _, p := range m.catPkgs {
			m.selected[p.ID] = false
		}
	}
	return m, nil
}

func (m Model) updateProfile(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.screen = screenMenu
	case "up", "k":
		if m.catCursor > 0 {
			m.catCursor--
		}
	case "down", "j":
		if m.catCursor < len(m.catalog.Profiles)-1 {
			m.catCursor++
		}
	case "enter", " ":
		prof := m.catalog.Profiles[m.catCursor]
		for _, pkgID := range prof.Packages {
			m.selected[pkgID] = true
		}
		var ids []string
		for id, checked := range m.selected {
			if checked {
				ids = append(ids, id)
			}
		}
		resolved, err := m.catalog.Resolve(ids)
		if err != nil {
			m.screen = screenMenu
			return m, nil
		}
		m.toInstall = resolved
		m.screen = screenConfirm
	}
	return m, nil
}

func (m Model) updateConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y", "enter":
		m.installIdx = 0
		m.installLog = nil
		// Start the first install — TUI suspends, full terminal given to installer
		return m, m.execNextInstall()
	case "n", "N", "esc":
		m.screen = screenMenu
	}
	return m, nil
}

// execNextInstall uses tea.ExecProcess to SUSPEND the TUI and run the installer
// with full terminal access — this allows sudo to prompt for the password normally.
func (m Model) execNextInstall() tea.Cmd {
	if m.installIdx >= len(m.toInstall) {
		return nil
	}
	pkg := m.toInstall[m.installIdx]

	if pkg.IsInstalled() {
		return func() tea.Msg {
			return installDoneMsg{pkg: pkg, err: nil}
		}
	}

	cmd := buildInstallCmd(pkg)

	// Print a header so the user knows what's installing while TUI is suspended
	fmt.Printf("\n\033[1;33m  Installing %s...\033[0m\n\n", pkg.Name)

	// tea.ExecProcess suspends BubbleTea, gives full stdin/stdout/stderr to cmd,
	// then resumes and sends the callback msg when done.
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return installDoneMsg{pkg: pkg, err: err}
	})
}

// ── View ──────────────────────────────────────────────────────────────────────

func (m Model) View() string {
	switch m.screen {
	case screenMenu:
		return m.viewMenu()
	case screenCategory:
		return m.viewCategory()
	case screenProfile:
		return m.viewProfile()
	case screenConfirm:
		return m.viewConfirm()
	case screenDone:
		return m.viewDone()
	}
	return ""
}

func (m Model) viewMenu() string {
	var b strings.Builder
	b.WriteString(m.renderHeader())
	b.WriteString("\n\n")

	count := 0
	for _, v := range m.selected {
		if v {
			count++
		}
	}
	if count > 0 {
		b.WriteString(styleBadge.Render(fmt.Sprintf(" %d selected ", count)))
		b.WriteString("\n\n")
	}

	b.WriteString(styleCategory.Render("  Select Categories to Install:"))
	b.WriteString("\n\n")

	for i, item := range menuItems {
		if item.value == "---" {
			b.WriteString(styleDim.Render("  ────────────────────────────"))
			b.WriteString("\n")
			continue
		}

		cursor := "  "
		style := styleNormal
		if i == m.menuCursor {
			cursor = "▶ "
			style = styleSelected
		}

		line := fmt.Sprintf("%s%s", cursor, item.label)

		if strings.HasPrefix(item.value, "cat:") {
			cat := strings.TrimPrefix(item.value, "cat:")
			c := 0
			for _, p := range m.catalog.ByCategory()[cat] {
				if m.selected[p.ID] {
					c++
				}
			}
			if c > 0 {
				line += "  " + styleBadge.Render(fmt.Sprintf("%d", c))
			}
		}

		b.WriteString(style.Render(line))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(styleHint.Render("  navigate: arrow keys | select: space/enter | quit: q"))
	b.WriteString("\n")
	return b.String()
}

func (m Model) viewCategory() string {
	var b strings.Builder
	b.WriteString(m.renderHeader())
	b.WriteString("\n\n")
	b.WriteString(styleCategory.Render("  " + m.curCategory))
	b.WriteString("\n\n")

	for i, pkg := range m.catPkgs {
		cursor := "  "
		style := styleNormal
		if i == m.catCursor {
			cursor = "▶ "
			style = styleSelected
		}
		checkbox := "[ ]"
		if m.selected[pkg.ID] {
			checkbox = styleCheck.Render("[X]")
		}
		desc := styleDim.Render(pkg.Description)
		if len(pkg.Deps) > 0 {
			desc += styleDim.Render(fmt.Sprintf(" (needs: %s)", strings.Join(pkg.Deps, ", ")))
		}
		b.WriteString(fmt.Sprintf("%s%s %s  %s\n",
			cursor, checkbox, style.Render(pkg.Name), desc))
	}

	b.WriteString("\n")
	b.WriteString(styleHint.Render("  space toggle   a select all   n deselect all   esc back"))
	b.WriteString("\n")
	return b.String()
}

func (m Model) viewProfile() string {
	var b strings.Builder
	b.WriteString(m.renderHeader())
	b.WriteString("\n\n")
	b.WriteString(styleCategory.Render("  Developer Profiles"))
	b.WriteString("\n")
	b.WriteString(styleDim.Render("  Select a developer profile configuration to auto-setup:"))
	b.WriteString("\n\n")

	for i, prof := range m.catalog.Profiles {
		cursor := "  "
		style := styleNormal
		if i == m.catCursor {
			cursor = "▶ "
			style = styleSelected
		}
		b.WriteString(fmt.Sprintf("%s%s\n", cursor, style.Render(prof.Name)))
		b.WriteString(fmt.Sprintf("     %s\n", styleDim.Render(prof.Description)))
		b.WriteString(fmt.Sprintf("     %s\n\n", styleDim.Render(strings.Join(prof.Packages, " · "))))
	}

	b.WriteString(styleHint.Render("  enter install profile | esc back"))
	b.WriteString("\n")
	return b.String()
}

func (m Model) viewConfirm() string {
	var b strings.Builder
	b.WriteString(m.renderHeader())
	b.WriteString("\n\n")
	b.WriteString(styleCategory.Render(fmt.Sprintf("  %d Packages to Install (resolved with dependencies):", len(m.toInstall))))
	b.WriteString("\n\n")
	for _, pkg := range m.toInstall {
		b.WriteString(fmt.Sprintf("  %-20s  %s\n",
			styleNormal.Render(pkg.Name), styleDim.Render(pkg.Description)))
	}
	b.WriteString("\n")
	b.WriteString(styleSelected.Render("  Proceed? [y/n]  "))
	b.WriteString("\n")
	return b.String()
}

func (m Model) viewDone() string {
	var b strings.Builder
	b.WriteString(m.renderHeader())
	b.WriteString("\n\n")

	b.WriteString(lipgloss.NewStyle().Foreground(colorAccent).Bold(true).Render("  INSTALLATION SUMMARY"))
	b.WriteString("\n")
	b.WriteString(styleDim.Render("  ──────────────────────────────────────────────────────────"))
	b.WriteString("\n\n")

	// 1. Success List
	if len(m.installedSuccess) > 0 {
		b.WriteString(styleSuccess.Render("  [SUCCESS] DOWNLOADED & INSTALLED (") + fmt.Sprintf("%d", len(m.installedSuccess)) + "):\n")
		for _, name := range m.installedSuccess {
			b.WriteString(fmt.Sprintf("     - %s\n", name))
		}
	} else {
		b.WriteString(styleSuccess.Render("  [SUCCESS] DOWNLOADED & INSTALLED (0):") + "\n     - None\n")
	}
	b.WriteString("\n")

	// 2. Failed List
	if len(m.installedFailed) > 0 {
		b.WriteString(styleError.Render("  [FAILED] FAILED TO DOWNLOAD/INSTALL (") + fmt.Sprintf("%d", len(m.installedFailed)) + "):\n")
		for _, name := range m.installedFailed {
			b.WriteString(fmt.Sprintf("     - %s\n", name))
		}
	} else {
		b.WriteString(styleError.Render("  [FAILED] FAILED TO DOWNLOAD/INSTALL (0):") + "\n     - None\n")
	}
	b.WriteString("\n")

	// 3. Not Selected / Not Downloaded List
	installedMap := make(map[string]bool)
	for _, name := range m.installedSuccess {
		installedMap[name] = true
	}
	for _, name := range m.installedFailed {
		installedMap[name] = true
	}

	var notDownloaded []string
	for _, pkg := range m.catalog.Packages {
		if !installedMap[pkg.Name] {
			notDownloaded = append(notDownloaded, pkg.Name)
		}
	}

	b.WriteString(lipgloss.NewStyle().Foreground(colorBlue).Bold(true).Render("  [NOT INSTALLED] NOT DOWNLOADED / NOT SELECTED (") + fmt.Sprintf("%d", len(notDownloaded)) + "):\n")
	var notDownloadedStr string
	if len(notDownloaded) > 8 {
		notDownloadedStr = strings.Join(notDownloaded[:8], ", ") + fmt.Sprintf("... (+%d more)", len(notDownloaded)-8)
	} else if len(notDownloaded) > 0 {
		notDownloadedStr = strings.Join(notDownloaded, ", ")
	} else {
		notDownloadedStr = "None"
	}
	b.WriteString("     " + styleDim.Render(notDownloadedStr) + "\n\n")

	b.WriteString(styleDim.Render("  Run 'autodev doctor' to verify your environment."))
	b.WriteString("\n")
	b.WriteString(styleDim.Render("  Run 'autodev skills' to see your learning roadmap."))
	b.WriteString("\n\n")
	b.WriteString(styleHint.Render("  press esc or enter to exit"))
	b.WriteString("\n")
	return b.String()
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func (m Model) renderHeader() string {
	var b strings.Builder
	ascii := `  █████╗ ██╗   ██╗████████╗ ██████╗ ██████╗ ███████╗██╗   ██╗
  ██╔══██╗██║   ██║╚══██╔══╝██╔═══██╗██╔══██╗██╔════╝██║   ██║
  ███████║██║   ██║   ██║   ██║   ██║██║  ██║█████╗  ██║   ██║
  ██╔══██║██║   ██║   ██║   ██║   ██║██║  ██║██╔══╝  ╚██╗ ██╔╝
  ██║  ██║╚██████╔╝   ██║   ╚██████╔╝██████╔╝███████╗ ╚████╔╝ 
  ╚═╝  ╚═╝ ╚═════╝    ╚═╝    ╚══════╝ ╚═════╝ ╚══════╝  ╚═══╝ `
	b.WriteString(GalaxyBlock(ascii))
	b.WriteString("\n\n")
	if m.sysInfo != nil {
		specs := fmt.Sprintf("  OS: %s (%s) | RAM: %s | PM: %s",
			m.sysInfo.Version,
			m.sysInfo.Arch,
			osinfo.FormatRAM(m.sysInfo.RAMBytes),
			m.sysInfo.PackageManager,
		)
		b.WriteString(styleDim.Render(specs))
	}
	return b.String()
}

// buildInstallCmd builds an *exec.Cmd for a package, with stdin/stdout/stderr
// connected to the real terminal (so sudo password prompts work).
func buildInstallCmd(pkg *catalog.Package) *exec.Cmd {
	var steps catalog.InstallSteps
	switch runtime.GOOS {
	case "linux":
		steps = pkg.Install.Linux
	case "darwin":
		steps = pkg.Install.Darwin
	case "windows":
		steps = pkg.Install.Windows
	}

	var cmd *exec.Cmd

	switch steps.Method {
	case "apt":
		if commandExists("apt-get") {
			cmd = exec.Command("sudo", append([]string{"apt-get", "install", "-y"}, steps.Packages...)...)
		} else if commandExists("dnf") {
			cmd = exec.Command("sudo", append([]string{"dnf", "install", "-y"}, steps.Packages...)...)
		} else if commandExists("pacman") {
			cmd = exec.Command("sudo", append([]string{"pacman", "-S", "--noconfirm"}, steps.Packages...)...)
		} else if commandExists("yum") {
			cmd = exec.Command("sudo", append([]string{"yum", "install", "-y"}, steps.Packages...)...)
		} else if commandExists("apk") {
			cmd = exec.Command("sudo", append([]string{"apk", "add"}, steps.Packages...)...)
		} else {
			cmd = exec.Command("sudo", append([]string{"apt-get", "install", "-y"}, steps.Packages...)...)
		}

	case "brew":
		cmd = exec.Command("brew", append([]string{"install"}, steps.Packages...)...)

	case "winget":
		// winget: install sequentially — wrap in a shell
		cmds := make([]string, len(steps.Packages))
		for i, p := range steps.Packages {
			cmds[i] = "winget install --silent " + p
		}
		cmd = exec.Command("cmd", "/C", strings.Join(cmds, " && "))

	case "npm":
		cmd = exec.Command("npm", append([]string{"install", "-g"}, steps.Packages...)...)

	case "pip":
		args := []string{"install", "--upgrade"}
		if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
			args = append(args, "--break-system-packages")
		}
		args = append(args, steps.Packages...)
		cmd = exec.Command("pip3", args...)

	case "cargo":
		cmd = exec.Command("cargo", append([]string{"install"}, steps.Packages...)...)

	case "script", "manual":
		script := strings.Join(steps.Script, "\n")
		cmd = exec.Command("sh", "-c", script)

	default:
		// Fallback: echo unsupported
		cmd = exec.Command("sh", "-c",
			fmt.Sprintf("echo 'No install method defined for %s on %s'", pkg.Name, runtime.GOOS))
	}

	// Connect to real terminal so sudo password prompts work
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

// Run launches the TUI.
func Run(c *catalog.Catalog) error {
	m := New(c)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
