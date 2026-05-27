package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/eljaska/claustomize/internal/block"
	"github.com/eljaska/claustomize/internal/statusline"
)

var (
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("213"))
	cursorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	dimStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	panelStyle    = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1)
	previewStyle  = panelStyle.Width(60)
	listStyle     = panelStyle.Width(36)
	statusOKStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("84"))
	statusErrStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))
)

type Model struct {
	blocks   []block.Block
	selected map[string]bool
	cursor   int
	status   string
	isError  bool
}

func New() Model {
	return Model{
		blocks:   block.All(),
		selected: map[string]bool{},
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch keyMsg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.blocks)-1 {
			m.cursor++
		}
	case " ":
		id := m.blocks[m.cursor].ID
		m.selected[id] = !m.selected[id]
		m.status = ""
	case "i":
		if err := statusline.Install(m.activeBlocks()); err != nil {
			m.status = fmt.Sprintf("install failed: %v", err)
			m.isError = true
		} else {
			m.status = "installed to ~/.claude/settings.json"
			m.isError = false
		}
	}
	return m, nil
}

func (m Model) activeBlocks() []block.Block {
	out := make([]block.Block, 0, len(m.blocks))
	for _, b := range m.blocks {
		if m.selected[b.ID] {
			out = append(out, b)
		}
	}
	return out
}

func (m Model) View() string {
	var list strings.Builder
	for i, b := range m.blocks {
		mark := "[ ]"
		if m.selected[b.ID] {
			mark = "[x]"
		}
		line := fmt.Sprintf("%s %s", mark, b.Name)
		if i == m.cursor {
			line = cursorStyle.Render("> " + line)
		} else {
			line = "  " + line
		}
		list.WriteString(line)
		list.WriteString("\n")
		list.WriteString(dimStyle.Render("    " + b.Description))
		list.WriteString("\n")
	}

	preview := statusline.Preview(m.activeBlocks())
	previewPanel := previewStyle.Render(
		titleStyle.Render("Preview") + "\n\n" + preview,
	)
	listPanel := listStyle.Render(
		titleStyle.Render("Blocks") + "\n\n" + strings.TrimRight(list.String(), "\n"),
	)

	body := lipgloss.JoinHorizontal(lipgloss.Top, listPanel, previewPanel)

	footer := dimStyle.Render("↑/↓ move  space toggle  i install  q quit")
	statusLine := ""
	if m.status != "" {
		if m.isError {
			statusLine = statusErrStyle.Render(m.status)
		} else {
			statusLine = statusOKStyle.Render(m.status)
		}
	}

	header := titleStyle.Render("claustomize — statusline builder")
	parts := []string{header, "", body, "", footer}
	if statusLine != "" {
		parts = append(parts, statusLine)
	}
	return strings.Join(parts, "\n") + "\n"
}

// Run starts the TUI.
func Run() error {
	_, err := tea.NewProgram(New(), tea.WithAltScreen()).Run()
	return err
}
