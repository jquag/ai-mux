package startinfo

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jquag/ai-mux/component/modal"
	"github.com/jquag/ai-mux/theme"
)

type Model struct {
	TmuxSessionMessage    string
	TmuxWindowMessage     string
	WorktreeFolderMessage string
	GitBranchMessage      string
	Error                 error
}

func (m Model) Update(msg tea.Msg) (modal.ModalContent, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	if m.Error != nil {
		return lipgloss.NewStyle().Foreground(theme.Colors.Error).Render(m.Error.Error())
	}

	labelStyle := lipgloss.NewStyle().Foreground(theme.Colors.Primary).Width(15).MaxWidth(15).MaxHeight(1)
	labels := []string{
		labelStyle.Render("tmux session"),
		labelStyle.Render("tmux window"),
		labelStyle.Render("worktree"),
		labelStyle.Render("branch"),
	}
	labelContent := lipgloss.JoinVertical(lipgloss.Left, labels...)

	valueStyle := lipgloss.NewStyle().Foreground(theme.Colors.Text).MaxHeight(1)
	values := []string{
		valueStyle.Render(m.TmuxSessionMessage),
		valueStyle.Render(m.TmuxWindowMessage),
		valueStyle.Render(m.WorktreeFolderMessage),
		valueStyle.Render(m.GitBranchMessage),
	}
	valueContent := lipgloss.JoinVertical(lipgloss.Left, values...)

	return lipgloss.JoinHorizontal(lipgloss.Top, labelContent, valueContent)
}

func (m Model) ShouldCloseOnEscape() bool {
	return true
}

func (m Model) WithWidth(width int) modal.ModalContent {
	return m
}

func (m Model) WithHeight(height int) modal.ModalContent {
	return m
}

func Alert(m Model) modal.ShowModalMsg {
	return modal.ShowModalMsg{Content: m, Title: "Session Started"}
}
