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
	Overlayed bool
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
	titleColor := theme.Colors.Title
	borderColor := theme.Colors.Border

	if m.Overlayed {
		titleColor = theme.Colors.Muted
		borderColor = theme.Colors.Muted
	}

	title := lipgloss.NewStyle().
		Foreground(titleColor).
		Border(lipgloss.NormalBorder(), false, false, true).BorderForeground(borderColor).
		Width(m.width).
		Render("Work Items")
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

	if !m.Overlayed {
		body += "\n\n[Press " +
			lipgloss.NewStyle().Foreground(theme.Colors.Primary).Render("a") +
			" to add a work item.]"
	}

	return body
}

func (m *Model) listBody() string {
	items := make([]string, len(m.workItems))
	for i, item := range m.workItems {
		items[i] = m.itemView(item)
	}
	return lipgloss.JoinVertical(lipgloss.Left, items...)
}

func (m *Model) itemView(item *workitem.WorkItem) string {
	lineStyle := lipgloss.NewStyle().Foreground(theme.Colors.Muted)
	left := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Foreground(theme.Colors.Muted).Render("● "),
		lineStyle.Render("│"),
		lineStyle.Render("│"),
		lineStyle.Render("╰"),
	)

	nameColor := theme.Colors.Primary
	descriptionColor := theme.Colors.Text

	if m.Overlayed {
		nameColor = theme.Colors.Muted
		descriptionColor = theme.Colors.Muted
	}

	centerWidth := m.width - lipgloss.Width(left) - 1
	name := lipgloss.NewStyle().Width(centerWidth).MaxWidth(centerWidth).MaxHeight(1).Foreground(nameColor).Render(item.BranchName)
	descr := lipgloss.NewStyle().
		Height(2).MaxHeight(2).Width(centerWidth).
		Foreground(descriptionColor).
		Render(item.Description)
	status := lipgloss.NewStyle().Foreground(theme.Colors.Muted).Render("[Not Started]")

	right := ""
	// Check if name was truncated
	if lipgloss.Width(item.BranchName) > centerWidth {
		right = lipgloss.NewStyle().Foreground(theme.Colors.Muted).Render("…")
	} else {
		right = lipgloss.NewStyle().Foreground(theme.Colors.Muted).Render(" ")
	}
	// Check if description exceeds 2 lines when wrapped
	descrHeight := lipgloss.Height(lipgloss.NewStyle().Width(centerWidth).Render(item.Description))
	if descrHeight > 2 {
		right += lipgloss.NewStyle().Foreground(theme.Colors.Muted).Render("\n\n…")
	}

	info := lipgloss.JoinVertical(lipgloss.Left, name, descr, status)
	return lipgloss.JoinHorizontal(lipgloss.Top, left, info, right) + "\n"
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
