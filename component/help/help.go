package help

import (
	"strings"
	
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jquag/ai-mux/component/modal"
	"github.com/jquag/ai-mux/theme"
)

type Model struct {
	viewport viewport.Model
	width    int
	height   int
}

func New() *Model {
	vp := viewport.New(0, 0)
	return &Model{
		viewport: vp,
	}
}

func (m *Model) Update(msg tea.Msg) (modal.ModalContent, tea.Cmd) {
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	content := m.buildContent()
	m.viewport.SetContent(content)
	return m.viewport.View()
}

func (m *Model) buildContent() string {
	headerStyle := lipgloss.NewStyle().
		Foreground(theme.Colors.Primary).
		Bold(true)
	
	keyStyle := lipgloss.NewStyle().
		Foreground(theme.Colors.Success).
		Bold(true).
		Width(10)
	
	descStyle := lipgloss.NewStyle().
		Width(m.width - 12).
		Foreground(theme.Colors.Text)
	
	// Build key bindings sections
	var sections []string
	
	// General section
	sections = append(sections, headerStyle.Render("General"))
	sections = append(sections,
		lipgloss.JoinHorizontal(lipgloss.Top, keyStyle.Render("q"), descStyle.Render("Quit application")),
		lipgloss.JoinHorizontal(lipgloss.Top, keyStyle.Render("?"), descStyle.Render("Show this help")),
		lipgloss.JoinHorizontal(lipgloss.Top, keyStyle.Render("j/↓"), descStyle.Render("Move selection down")),
		lipgloss.JoinHorizontal(lipgloss.Top, keyStyle.Render("k/↑"), descStyle.Render("Move selection up")),
		lipgloss.JoinHorizontal(lipgloss.Top, keyStyle.Render("Esc"), descStyle.Render("Close modal/dialog")),
		"", // Empty line for spacing
	)
	
	// Work items management
	sections = append(sections, headerStyle.Render("Work Items"))
	sections = append(sections,
		lipgloss.JoinHorizontal(lipgloss.Top, keyStyle.Render("a"), descStyle.Render("Add new work item")),
		lipgloss.JoinHorizontal(lipgloss.Top, keyStyle.Render("e"), descStyle.Render("Edit work item")),
		lipgloss.JoinHorizontal(lipgloss.Top, keyStyle.Render("Enter"), descStyle.Render("Show work item details inlcuding code changes made")),
		lipgloss.JoinHorizontal(lipgloss.Top, keyStyle.Render("c"), descStyle.Render("Close work item")),
		"", // Empty line for spacing
	)
	
	// Session management
	sections = append(sections, headerStyle.Render("Session Management"))
	sections = append(sections,
		lipgloss.JoinHorizontal(lipgloss.Top, keyStyle.Render("s"), descStyle.Render("Start session in default mode (manual accept)")),
		lipgloss.JoinHorizontal(lipgloss.Top, keyStyle.Render("p"), descStyle.Render("Start session in plan mode")),
		lipgloss.JoinHorizontal(lipgloss.Top, keyStyle.Render("v"), descStyle.Render("Start session in vibe/accept-edits mode")),
		lipgloss.JoinHorizontal(lipgloss.Top, keyStyle.Render("r"), descStyle.Render("Resume existing session (in case a claude session was interrupted)")),
		lipgloss.JoinHorizontal(lipgloss.Top, keyStyle.Render("o"), descStyle.Render("Open/switch to tmux window")),
	)
	
	return strings.Join(sections, "\n")
}

func (m *Model) ShouldCloseOnEscape() bool {
	return true
}

func (m *Model) WithWidth(width int) modal.ModalContent {
	m.width = width
	m.viewport.Width = width
	return m
}

func (m *Model) WithHeight(height int) modal.ModalContent {
	m.height = height
	m.viewport.Height = height - 4
	return m
}
