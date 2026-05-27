package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/eljaska/claustomize/internal/slot"
	"github.com/eljaska/claustomize/internal/statusline"
	"github.com/eljaska/claustomize/internal/statusline/blocks"
)

const (
	focusPreview = iota
	focusPalette
)

const emptyMarker = "·"

var (
	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("213"))
	dimStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	mutedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	activeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("231"))
	emptyStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
	selectStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("231")).
			Background(lipgloss.Color("57"))
	panelStyle           = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1)
	focusBorderColor     = lipgloss.Color("213")
	unfocusedBorderColor = lipgloss.Color("238")
	statusOKStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("84"))
	statusErrStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))
)

type Model struct {
	slots        slot.List
	cursor       int
	focus        int
	palette      []blocks.Block
	paletteType  int
	paletteStyle int
	status       string
	isError      bool
}

func New() Model {
	return Model{
		slots:   slot.New(),
		palette: blocks.All(),
		focus:   focusPreview,
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	if m.focus == focusPalette {
		return m.updatePalette(keyMsg)
	}
	return m.updatePreview(keyMsg)
}

func (m Model) updatePreview(k tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch k.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "left":
		if m.cursor > 0 {
			m.cursor--
		}
	case "right":
		if m.cursor < len(m.slots)-1 {
			m.cursor++
		}
	case "enter":
		m.enterPalette()
	case "delete", "backspace":
		if !m.slots[m.cursor].IsEmpty() {
			m.slots, m.cursor = m.slots.Empty(m.cursor)
			m.status = ""
		}
	case "i":
		if err := statusline.Install(m.slots); err != nil {
			m.status = fmt.Sprintf("install failed: %v", err)
			m.isError = true
		} else {
			m.status = "installed to ~/.claude/settings.json"
			m.isError = false
		}
	}
	return m, nil
}

func (m Model) updatePalette(k tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch k.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "esc":
		m.focus = focusPreview
	case "up":
		if m.paletteType > 0 {
			m.paletteType--
			m.paletteStyle = 0
		}
	case "down":
		if m.paletteType < len(m.palette)-1 {
			m.paletteType++
			m.paletteStyle = 0
		}
	case "left":
		if m.paletteStyle > 0 {
			m.paletteStyle--
		}
	case "right":
		if m.paletteStyle < len(m.palette[m.paletteType].Styles)-1 {
			m.paletteStyle++
		}
	case "enter":
		b := &m.palette[m.paletteType]
		if m.slots[m.cursor].IsEmpty() {
			m.slots, m.cursor = m.slots.Fill(m.cursor, b, m.paletteStyle)
		} else {
			m.slots = m.slots.Replace(m.cursor, b, m.paletteStyle)
		}
		m.focus = focusPreview
		m.status = ""
	}
	return m, nil
}

// enterPalette transfers focus to the palette and seeds the palette cursors
// from the active slot (so the current selection is visible).
func (m *Model) enterPalette() {
	m.focus = focusPalette
	if s := m.slots[m.cursor]; !s.IsEmpty() {
		for i := range m.palette {
			if m.palette[i].ID == s.Block.ID {
				m.paletteType = i
				m.paletteStyle = s.StyleIdx
				return
			}
		}
	}
	m.paletteType = 0
	m.paletteStyle = 0
}

func (m Model) View() string {
	header := titleStyle.Render("claustomize — statusline builder")

	body := lipgloss.JoinHorizontal(lipgloss.Top,
		m.renderPalette(),
		m.renderPreview(),
	)

	footer := m.renderFooter()

	parts := []string{header, "", body, "", footer}
	if m.status != "" {
		style := statusOKStyle
		if m.isError {
			style = statusErrStyle
		}
		parts = append(parts, style.Render(m.status))
	}
	return strings.Join(parts, "\n") + "\n"
}

func (m Model) renderPalette() string {
	paletteFocused := m.focus == focusPalette
	rowStyle := mutedStyle
	if paletteFocused {
		rowStyle = activeStyle
	}

	var lines []string
	lines = append(lines, titleStyle.Render("Blocks"), "")
	for i, b := range m.palette {
		prefix := "  "
		if i == m.paletteType {
			prefix = "▸ "
		}
		line := prefix + b.Name
		if i == m.paletteType {
			line = rowStyle.Bold(true).Render(line)
		} else {
			line = rowStyle.Render(line)
		}
		lines = append(lines, line)
	}
	return m.panel(paletteFocused).Width(24).Render(strings.Join(lines, "\n"))
}

func (m Model) panel(focused bool) lipgloss.Style {
	if focused {
		return panelStyle.BorderForeground(focusBorderColor)
	}
	return panelStyle.BorderForeground(unfocusedBorderColor)
}

func (m Model) renderPreview() string {
	previewFocused := m.focus == focusPreview

	var b strings.Builder
	b.WriteString(titleStyle.Render("Preview"))
	b.WriteString("\n\n")
	b.WriteString(m.renderStatusline(previewFocused))
	b.WriteString("\n")

	if m.focus == focusPalette && m.paletteType < len(m.palette) {
		b.WriteString("\n")
		b.WriteString(m.renderStyleStrip())
	}

	return m.panel(previewFocused).Width(60).Render(b.String())
}

// renderStatusline composes the per-slot rendered preview, highlighting the
// active slot. When previewFocused is true, the highlight is bright; when
// not focused (palette has focus), the highlight is shown more subtly so
// the user still sees which slot they are editing.
func (m Model) renderStatusline(previewFocused bool) string {
	var b strings.Builder
	for i, s := range m.slots {
		text := statusline.RenderSlot(s)
		if s.IsEmpty() {
			text = emptyMarker
		}
		switch {
		case i == m.cursor && previewFocused:
			b.WriteString(selectStyle.Render(text))
		case i == m.cursor:
			b.WriteString(activeStyle.Underline(true).Render(text))
		case s.IsEmpty():
			b.WriteString(emptyStyle.Render(text))
		default:
			b.WriteString(text)
		}
	}
	return b.String()
}

func (m Model) renderStyleStrip() string {
	styles := m.palette[m.paletteType].Styles
	parts := make([]string, len(styles))
	for i, s := range styles {
		if i == m.paletteStyle {
			parts[i] = activeStyle.Bold(true).Render(s.Name)
		} else {
			parts[i] = mutedStyle.Render(s.Name)
		}
	}
	return dimStyle.Render("Style: < ") + strings.Join(parts, dimStyle.Render(" · ")) + dimStyle.Render(" >")
}

func (m Model) renderFooter() string {
	if m.focus == focusPalette {
		return dimStyle.Render("↑/↓ block · ←/→ style · Enter apply · Esc cancel · q quit")
	}
	return dimStyle.Render("←/→ move · Enter edit · Del empty · i install · q quit")
}

// Run starts the TUI.
func Run() error {
	_, err := tea.NewProgram(New(), tea.WithAltScreen()).Run()
	return err
}
