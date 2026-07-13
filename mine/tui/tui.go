// Package tui implements the Marathon-styled terminal interface for minenotyours.
package tui

import (
	"minenotyours/fileio"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

type state int

const (
	stateMenu state = iota
	stateBrowse
	statePassword
	stateWorking
	stateResult
)

type operation int

const (
	opEncrypt operation = iota
	opDecrypt
)

func (o operation) String() string {
	if o == opDecrypt {
		return "DECRYPT"
	}
	return "ENCRYPT"
}

type cryptoDoneMsg struct{ err error }

type model struct {
	state         state
	width, height int

	menuCursor int // 0 = ENCRYPT, 1 = DECRYPT
	op         operation

	fp       filepicker.Model
	filePath string

	passInputs []textinput.Model
	passFocus  int
	passErr    string

	spinner spinner.Model

	resultErr error // nil == success
}

func newModel() model {
	return model{state: stateMenu}
}

func Run() error {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.fp.Height = clamp(msg.Height-10, 3, 12)
		return m, nil
	case cryptoDoneMsg:
		m.state = stateResult
		m.resultErr = msg.err
		return m, nil
	case spinner.TickMsg:
		if m.state == stateWorking {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	var cmd tea.Cmd
	switch m.state {
	case stateBrowse:
		m.fp, cmd = m.fp.Update(msg)
	case statePassword:
		if len(m.passInputs) > 0 {
			m.passInputs[m.passFocus], cmd = m.passInputs[m.passFocus].Update(msg)
		}
	}
	return m, cmd
}

func (m model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "ctrl+c" && m.state != stateWorking {
		return m, tea.Quit
	}

	switch m.state {
	case stateMenu:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "up", "k":
			if m.menuCursor > 0 {
				m.menuCursor--
			}
		case "down", "j":
			if m.menuCursor < 1 {
				m.menuCursor++
			}
		case "enter":
			m.op = opEncrypt
			if m.menuCursor == 1 {
				m.op = opDecrypt
			}
			cmd := m.enterBrowse()
			return m, cmd
		}

	case stateBrowse:
		switch msg.String() {
		case "esc":
			m.state = stateMenu
			return m, nil
		case "q":
			return m, tea.Quit
		}
		var cmd tea.Cmd
		m.fp, cmd = m.fp.Update(msg)
		if ok, path := m.fp.DidSelectFile(msg); ok {
			m.filePath = path
			m.enterPassword()
			return m, tea.Batch(cmd, textinput.Blink)
		}
		return m, cmd

	case statePassword:
		switch msg.String() {
		case "esc":
			m.state = stateBrowse
			return m, nil
		case "tab", "down":
			m.passFocus = (m.passFocus + 1) % len(m.passInputs)
			m.focusPass()
			return m, nil
		case "shift+tab", "up":
			m.passFocus = (m.passFocus - 1 + len(m.passInputs)) % len(m.passInputs)
			m.focusPass()
			return m, nil
		case "enter":
			return m.submitPass()
		}
		var cmd tea.Cmd
		m.passInputs[m.passFocus], cmd = m.passInputs[m.passFocus].Update(msg)
		return m, cmd

	case stateWorking:
		return m, nil

	case stateResult:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "enter", "esc":
			m.reset()
			return m, nil
		}
	}
	return m, nil
}

func (m *model) enterBrowse() tea.Cmd {
	fp := filepicker.New()
	if home, err := os.UserHomeDir(); err == nil {
		fp.CurrentDirectory = home
	}
	fp.DirAllowed = true
	fp.FileAllowed = true
	fp.ShowHidden = false
	fp.ShowPermissions = false
	fp.AutoHeight = false
	fp.Height = clamp(m.height-10, 3, 12)
	fp.Styles = c2FilepickerStyles()
	m.fp = fp
	m.state = stateBrowse
	return fp.Init()
}

func (m *model) enterPassword() {
	n := 1
	if m.op == opEncrypt {
		n = 2
	}
	labels := []string{"passphrase : ", "confirm    : "}
	inputs := make([]textinput.Model, n)
	for i := 0; i < n; i++ {
		ti := textinput.New()
		ti.EchoMode = textinput.EchoPassword
		ti.EchoCharacter = '•'
		ti.Prompt = labels[i]
		ti.PromptStyle = rowStyle
		ti.TextStyle = tealStyle
		ti.Cursor.Style = tealStyle
		inputs[i] = ti
	}
	inputs[0].Focus()
	m.passInputs = inputs
	m.passFocus = 0
	m.passErr = ""
	m.state = statePassword
}

func (m *model) focusPass() {
	for i := range m.passInputs {
		if i == m.passFocus {
			m.passInputs[i].Focus()
		} else {
			m.passInputs[i].Blur()
		}
	}
}

func (m model) submitPass() (tea.Model, tea.Cmd) {
	pass := m.passInputs[0].Value()
	if pass == "" {
		m.passErr = "passphrase cannot be empty"
		return m, nil
	}
	if m.op == opEncrypt && pass != m.passInputs[1].Value() {
		m.passErr = "passphrases do not match"
		return m, nil
	}
	m.passErr = ""
	m.state = stateWorking
	m.spinner = newSpinner()
	return m, tea.Batch(m.spinner.Tick, runCrypto(m.op, m.filePath, pass))
}

func (m *model) reset() {
	m.state = stateMenu
	m.filePath = ""
	m.passInputs = nil
	m.passFocus = 0
	m.passErr = ""
	m.resultErr = nil
}

func runCrypto(op operation, path, pass string) tea.Cmd {
	return func() tea.Msg {
		var err error
		if op == opEncrypt {
			err = fileio.CallEncryption(pass, path)
		} else {
			err = fileio.CallDecryption(pass, path)
		}
		return cryptoDoneMsg{err: err}
	}
}

func newSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = tealStyle
	return s
}

// View aqui --

func (m model) View() string {
	if m.width > 0 && (m.width < contentWidth+6 || m.height < 12) {
		msg := tealStyle.Render("terminal too small - please enlarge the window")
		return lipgloss.Place(m.width, max(m.height, 1), lipgloss.Center, lipgloss.Center, msg,
			lipgloss.WithWhitespaceBackground(colBg))
	}

	header := frameHeader("minenotyours", m.statusRight(), contentWidth)
	rule := ruleStyle.Render(strings.Repeat("─", contentWidth))
	body := m.body()
	footer := footerStyle.Width(contentWidth).Render(m.footer())

	content := frameContent([]string{header, rule, "", body, "", footer})
	frame := frameStyle.Render(content)

	if m.width > 0 && m.height > 0 {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, frame,
			lipgloss.WithWhitespaceBackground(colBg))
	}
	return frame
}

func (m model) statusRight() string {
	if m.state == stateMenu {
		return "JSSA · AES-256/ARGON2"
	}
	return m.op.String()
}

func (m model) body() string {
	switch m.state {
	case stateMenu:
		return m.menuBody()
	case stateBrowse:
		return promptStyle.Render("> SELECT TARGET FILE") + "\n\n" + fillGaps(m.fp.View())
	case statePassword:
		return m.passwordBody()
	case stateWorking:
		return m.workingBody()
	case stateResult:
		return m.resultBody()
	}
	return ""
}

func (m model) menuBody() string {
	var b strings.Builder
	b.WriteString(promptStyle.Render("> SELECT OPERATION"))
	b.WriteString("\n\n")

	items := []string{"ENCRYPT", "DECRYPT"}
	for i, it := range items {
		if i == m.menuCursor {
			b.WriteString(rowSelStyle.Width(contentWidth).Render(" ▸ " + it))
		} else {
			b.WriteString(rowStyle.Width(contentWidth).Render("   " + it))
		}
		if i < len(items)-1 {
			b.WriteString("\n")
		}
	}
	return b.String()
}

func (m model) passwordBody() string {
	head := "> SET PASSPHRASE"
	if m.op == opDecrypt {
		head = "> ENTER PASSPHRASE"
	}
	var b strings.Builder
	b.WriteString(promptStyle.Render(head))
	b.WriteString("\n")
	b.WriteString(dimStyle.Render("  "+filepath.Base(m.filePath)) + "\n\n")
	for i := range m.passInputs {
		b.WriteString("  " + m.passInputs[i].View() + "\n")
	}
	if m.passErr != "" {
		b.WriteString("\n" + rowStyle.Render("  ! "+m.passErr))
	}
	return b.String()
}

func (m model) workingBody() string {
	verb := "encrypting"
	if m.op == opDecrypt {
		verb = "decrypting"
	}
	var b strings.Builder
	b.WriteString(tealStyle.Render(m.spinner.View()+" "+verb+" ") + rowStyle.Render(filepath.Base(m.filePath)))
	b.WriteString("\n\n")
	b.WriteString(dimStyle.Render("  deriving key · argon2id — this can take a moment"))
	return b.String()
}

func (m model) resultBody() string {
	var b strings.Builder
	if m.resultErr == nil {
		title := "✓ SEALED"
		detail := "file is now encrypted."
		if m.op == opDecrypt {
			title = "✓ OPENED"
			detail = "file is now decrypted."
		}
		b.WriteString(okStyle.Render(title) + "\n\n")
		b.WriteString(tealStyle.Render("  "+detail) + "\n")
		b.WriteString(dimStyle.Render("  permissions preserved"))
	} else {
		b.WriteString(rowStyle.Bold(true).Render("✗ FAILED") + "\n\n")
		b.WriteString(rowStyle.Render("  "+m.resultErr.Error()) + "\n")
		b.WriteString(dimStyle.Render("  file left untouched."))
	}
	return b.String()
}

func (m model) footer() string {
	switch m.state {
	case stateMenu:
		return "[↑↓] move   [enter] select   [q] quit"
	case stateBrowse:
		return "[↑↓] move   [→] open   [enter] choose   [esc] back   [q] quit"
	case statePassword:
		return "[enter] confirm   [tab] next field   [esc] back   [ctrl+c] quit"
	case stateWorking:
		return "please wait — do not close"
	case stateResult:
		return "[enter] return to menu   [q] quit"
	}
	return ""
}
func fillGaps(s string) string {
	return strings.ReplaceAll(s, "\x1b[0m ", "\x1b[0m"+bgFill.Render(" "))
}

func frameContent(parts []string) string {
	var lines []string
	for _, p := range parts {
		for _, ln := range strings.Split(p, "\n") {
			if lipgloss.Width(ln) > contentWidth {
				ln = ansi.Truncate(ln, contentWidth, "…")
			}
			if lipgloss.Width(ln) < contentWidth {
				ln = bgFill.Width(contentWidth).Render(ln)
			}
			lines = append(lines, ln)
		}
	}
	return strings.Join(lines, "\n")
}

func frameHeader(left, right string, width int) string {
	l := headerStyle.Render(left)
	r := headerStyle.Render(right)
	gap := width - lipgloss.Width(l) - lipgloss.Width(r)
	if gap < 1 {
		gap = 1
	}
	spacer := lipgloss.NewStyle().Background(colBg).Render(strings.Repeat(" ", gap))
	return l + spacer + r
}

func c2FilepickerStyles() filepicker.Styles {
	s := filepicker.DefaultStyles()
	s.Cursor = tealStyle
	s.Directory = tealStyle
	s.File = rowStyle
	s.Selected = lipgloss.NewStyle().Foreground(colTeal).Background(colBg).Bold(true)
	s.DisabledFile = dimStyle
	s.EmptyDirectory = dimStyle
	s.FileSize = s.FileSize.Foreground(colTeal).Background(colBg)
	s.Permission = s.Permission.Foreground(colTeal).Background(colBg).Faint(true)
	return s
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
