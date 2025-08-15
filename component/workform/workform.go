package workform

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/google/uuid"
	"github.com/jquag/ai-mux/component/modal"
	workitem "github.com/jquag/ai-mux/data"
	"github.com/jquag/ai-mux/util"
)

type Model struct {
	form      *huh.Form
	submitted bool
	width     int
	height    int
	editMode  bool
	existingItem *workitem.WorkItem

	shortName    string
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
	m.width = width
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

func New(item *workitem.WorkItem) Model {
	m := Model{
		width:  0,
		height: 0,
		editMode: item != nil && item.Id != "",
		existingItem: item,
	}
	
	// Set initial values for editing
	shortNameValue := ""
	descriptionValue := ""
	confirmValue := true // Default to Submit
	if item != nil {
		shortNameValue = item.ShortName
		descriptionValue = item.Description
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("shortName").
				Title("Short name").
				Value(&shortNameValue),
			huh.NewText().
				Key("description").
				Title("Description").
				Value(&descriptionValue),
			huh.NewConfirm().
				Key("done").
				Value(&confirmValue).
				Affirmative("Submit (s)").
				Negative("Cancel (c)"),
		),
	).WithWidth(0).WithHeight(0)

	m.form = form
	return m
}

func (m Model) submitCmd() tea.Cmd {
	// Check if user clicked Cancel
	if !m.form.GetBool("done") {
		// User cancelled, just close the modal
		return modal.CloseCmd
	}
	
	workItem := m.existingItem
	workItem.ShortName = m.form.GetString("shortName")
	workItem.Description = m.form.GetString("description")
	
	if m.editMode && m.existingItem != nil {
		// Update the work item file
		if err := util.UpdateWorkItem(workItem); err != nil {
			fmt.Fprintf(os.Stderr, "Error updating work item: %v\n", err)
			os.Exit(1)
		}
		
		updateWorkItemCmd := func() tea.Msg {
			return workitem.UpdateWorkItemMsg{
				WorkItem: workItem,
			}
		}
		return tea.Batch(modal.CloseCmd, updateWorkItemCmd)
	} else {
		workItem.Id = uuid.New().String()
		
		// Save the new work item to file
		if err := util.SaveWorkItem(workItem); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving work item: %v\n", err)
			os.Exit(1)
		}
		
		newWorkItemCmd := func() tea.Msg {
			return workitem.NewWorkItemMsg{
				WorkItem: workItem,
			}
		}
		return tea.Batch(modal.CloseCmd, newWorkItemCmd)
	}
}

