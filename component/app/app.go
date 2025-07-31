package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jquag/ai-mux/component/modal"
	"github.com/jquag/ai-mux/component/worklist"
	"github.com/jquag/ai-mux/theme"
)

type Model struct {
	width         int
	height        int
	workListModel *worklist.Model
	currentModal  modal.Model
}

func New() Model {
	return Model{
		workListModel: worklist.New(0, 0),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			if !m.currentModal.Show {
				return m, tea.Quit
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateLayout(msg.Width, msg.Height)
		return m, nil

	case modal.ShowModalMsg:
		m.currentModal = modal.New(m.width, m.height, msg.Content, msg.Title, theme.Colors.Border)
		m.currentModal.Show = true
		return m, nil

	case modal.CloseMsg:
		m.currentModal.Show = false
		return m, nil
	}

	// Only send key messages to the active window
	if _, ok := msg.(tea.KeyMsg); ok {
		if m.currentModal.Show {
			newModal, cmd := m.currentModal.Update(msg)
			m.currentModal = newModal
			return m, cmd
		} else {
			// Let pane handle other key messages
			_, cmd := m.workListModel.Update(msg)
			return m, cmd
		}
	}

	cmds := []tea.Cmd{}

	if m.currentModal.Show {
		newModal, cmd := m.currentModal.Update(msg)
		m.currentModal = newModal
		cmds = append(cmds, cmd)
	}

	_, cmd := m.workListModel.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	borderColor := theme.Colors.Border
	if m.currentModal.Show {
		borderColor = theme.Colors.Muted
	}

	v := lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.NormalBorder()).
		BorderForeground(borderColor).
		Render(m.workListModel.View())
	if m.currentModal.Show {
		m.currentModal.BackgroundView = v
		return m.currentModal.View()
	} else {
		return v
	}
}

func (m *Model) updateLayout(width, height int) {
	m.workListModel.SetHeight(height - 4)
	m.workListModel.SetWidth(width - 4)

	m.currentModal = m.currentModal.WithWidth(width)
	m.currentModal = m.currentModal.WithHeight(height)
}
