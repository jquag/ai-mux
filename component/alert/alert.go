package alert

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jquag/ai-mux/component/modal"
	"github.com/jquag/ai-mux/theme"
)

type AlertType int

const (
	AlertTypeInfo AlertType = iota
	AlertTypeWarning
	AlertTypeError
)

type AlertMsg struct {
	Content string
	Type    AlertType
}

func Alert(content string, alertType AlertType) tea.Cmd {
	title := "Alert"
	switch alertType {
	case AlertTypeInfo:
		title = "Info"
	case AlertTypeWarning:
		title = "Warning"
	case AlertTypeError:
		title = "Error"
	}
	return modal.ShowModal(Model{
		Content: content,
		Type:    alertType,
	}, title)
}

type Model struct {
	Content string
	Type    AlertType
	width   int
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (modal.ModalContent, tea.Cmd) {
	return m, nil
}

func (m Model) WithWidth(width int) modal.ModalContent {
	m.width = width
	return m
}

func (m Model) WithHeight(height int) modal.ModalContent {
	return m
}

func (m Model) ShouldCloseOnEscape() bool {
	return true
}

func (m Model) View() string {
	style := lipgloss.NewStyle().Width(m.width).MaxWidth(m.width)
	icon := ""
	switch m.Type {
	case AlertTypeInfo:
		style = style.Foreground(theme.Colors.Info)
		icon = "  "
	case AlertTypeWarning:
		style = style.Foreground(theme.Colors.Primary)
		icon = "  "
	case AlertTypeError:
		style = style.Foreground(theme.Colors.Error)
		icon = "  "
	}
	return style.Render(icon + m.Content + "\n")
}
