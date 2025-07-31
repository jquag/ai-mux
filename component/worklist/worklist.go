package worklist

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jquag/ai-mux/component/modal"
	"github.com/jquag/ai-mux/component/workform"
	workitem "github.com/jquag/ai-mux/data"
	"github.com/jquag/ai-mux/theme"
)

type Model struct {
	width     int
	height    int
	viewport  viewport.Model
	workItems []*workitem.WorkItem
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			form := workform.New()
			initCmd := form.Init()
			return m, tea.Batch(initCmd, modal.ShowModal(form, "Add Work Item"))
		}
	case workitem.NewWorkItemMsg:
		m.workItems = append(m.workItems, msg.WorkItem)
	}

	return m, nil
}

func (m *Model) View() string {
	title := lipgloss.NewStyle().Foreground(theme.Colors.Title).Render("Work Items")
	body := ""

	if len(m.workItems) == 0 {
		body = m.emptyBody()
	} else {
		body = m.listBody()
	}

	m.viewport.SetContent(fmt.Sprintf("%s\n\n%s", title, body))

	var style = lipgloss.NewStyle().
		Width(m.width).
		Height(m.height)

	return style.Render(m.viewport.View())
}

func (m *Model) emptyBody() string {
	body := lipgloss.NewStyle().
		Foreground(theme.Colors.Muted).
		Italic(true).Render("--None--")

	body += "\n\n[Press " +
		lipgloss.NewStyle().Foreground(theme.Colors.Primary).Render("a") +
		" to add a work item.]"

	return body
}

func (m *Model) listBody() string {
	body := ""
	for _, item := range m.workItems {
		body += lipgloss.NewStyle().
			Render(fmt.Sprintf("- %s: %s", item.BranchName, item.Description)) + "\n"
	}

	return body
}

func (m *Model) SetWidth(width int) {
	m.viewport.Width = width
	m.width = width
}

func (m *Model) SetHeight(height int) {
	m.viewport.Height = height
	m.height = height
}

func New(width, height int) *Model {
	return &Model{
		width:    width,
		height:   height,
		viewport: viewport.New(width, height),
	}
}
