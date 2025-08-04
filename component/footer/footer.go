package footer

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jquag/ai-mux/theme"
)

type mapping struct {
	label string
	key   string
}

type Model struct {
	mappings  []mapping
	overlayed bool
	width     int
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	labelColor := theme.Colors.Text
	keyColor := theme.Colors.Primary
	borderColor := theme.Colors.Border
	if m.overlayed {
		labelColor = theme.Colors.Muted
		keyColor = theme.Colors.Muted
		borderColor = theme.Colors.Muted
	}

	views := []string{}
	for _, mapping := range m.mappings {
		v := fmt.Sprintf("%s: %s  ",
			lipgloss.NewStyle().Foreground(labelColor).Render(mapping.label),
			lipgloss.NewStyle().Foreground(keyColor).Render(mapping.key),
		)
		views = append(views, v)
	}
	style := lipgloss.NewStyle().
		Width(m.width).
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(borderColor)
	return style.Render(lipgloss.JoinHorizontal(lipgloss.Top, views...))
}

func (m Model) WithWidth(width int) Model {
	m.width = width
	return m
}

func (m Model) WithOverlayed(overlayed bool) Model {
	m.overlayed = overlayed
	return m
}

func New() Model {
	mappings := []mapping{
		{"Quit", "q"},
		{"Add", "a"},
		{"Remove", "d"},
		{"(Re)Start", "s"},
		{"Open", "o"},
		{"Info", "enter"},
	}

	return Model{
		mappings:  mappings,
		overlayed: false,
	}
}
