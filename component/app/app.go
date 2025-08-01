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
		m.width = min(msg.Width, 100)
		m.height = msg.Height
		m.updateLayout()
		return m, nil

	case modal.ShowModalMsg:
		m.currentModal = modal.New(m.width, m.height, msg.Content, msg.Title, theme.Colors.Border)
		m.currentModal.Show = true
		m.workListModel.Overlayed = true
		return m, nil

	case modal.CloseMsg:
		m.currentModal.Show = false
		m.workListModel.Overlayed = false
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
	v := lipgloss.NewStyle().
		Padding(0, 1).
		Render(m.workListModel.View())
	if m.currentModal.Show {
		m.currentModal.BackgroundView = v
		return m.currentModal.View()
	} else {
		return v
	}
}

func (m *Model) updateLayout() {
	m.workListModel.SetHeight(m.height - 2)
	m.workListModel.SetWidth(m.width - 2)

	m.currentModal = m.currentModal.WithWidth(m.width)
	m.currentModal = m.currentModal.WithHeight(m.height)
}
