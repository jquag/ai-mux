package workform

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/google/uuid"
	"github.com/jquag/ai-mux/component/modal"
	workitem "github.com/jquag/ai-mux/data"
)

type Model struct {
	form      *huh.Form
	submitted bool
	width     int
	height    int

	branchName    string
	planMode bool
	description   string
}

func (m Model) Init() tea.Cmd {
	return m.form.Init()
}

func (m Model) Update(msg tea.Msg) (modal.ModalContent, tea.Cmd) {
	if m.submitted {
		return m, nil
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f

		if m.form.State == huh.StateCompleted {
			m.submitted = true
			return m, tea.Batch(cmd, m.submitCmd())
		}
	}

	return m, cmd
}

func (m Model) View() string {
	return m.form.View()
}

func (m Model) WithWidth(width int) modal.ModalContent {
	m.width = min(width, 90)
	m.form = m.form.WithWidth(m.width)
	return m
}

func (m Model) WithHeight(height int) modal.ModalContent {
	m.height = min(height, 40)
	return m
}

func (m Model) ShouldCloseOnEscape() bool {
	return true
}

func New() Model {
	m := Model{
		width:  0,
		height: 0,
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("branchName").
				Title("Branch name"),
			huh.NewText().
				Key("description").
				Title("Description"),
			huh.NewConfirm().
				Key("startWithPlan").
				Title("Use plan mode"),
		),
	).WithWidth(0).WithHeight(0)

	m.form = form
	return m
}

func (m Model) submitCmd() tea.Cmd {
	// TODO: save the work item to file

	workItem := &workitem.WorkItem{
		Id:          uuid.New().String(),
		BranchName:  m.form.GetString("branchName"),
		Description: m.form.GetString("description"),
		PlanMode:    m.form.GetBool("startWithPlan"),
	}

	newWorkItemCmd := func() tea.Msg {
		return workitem.NewWorkItemMsg{
			WorkItem: workItem,
		}
	}
	return tea.Batch(modal.CloseCmd, newWorkItemCmd)
}
