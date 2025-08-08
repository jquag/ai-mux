package workform

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

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

func New() Model {
	m := Model{
		width:  0,
		height: 0,
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("shortName").
				Title("Short name"),
			huh.NewText().
				Key("description").
				Title("Description"),
		),
	).WithWidth(0).WithHeight(0)

	m.form = form
	return m
}

func (m Model) submitCmd() tea.Cmd {
	workItem := &workitem.WorkItem{
		Id:          uuid.New().String(),
		ShortName:  m.form.GetString("shortName"),
		Description: m.form.GetString("description"),
	}

	// Save the work item to file
	if err := saveWorkItem(workItem); err != nil {
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

func saveWorkItem(item *workitem.WorkItem) error {
	// Ensure .ai-mux directory exists
	if err := util.EnsureAiMuxDir(); err != nil {
		return err
	}

	// Create directory for this item using its UUID
	itemDir := filepath.Join(util.AiMuxDir, item.Id)
	if err := os.MkdirAll(itemDir, 0755); err != nil {
		return fmt.Errorf("failed to create item directory: %w", err)
	}

	// Save item as JSON
	itemPath := filepath.Join(itemDir, "item.json")
	itemData, err := json.MarshalIndent(item, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}
	if err := os.WriteFile(itemPath, itemData, 0644); err != nil {
		return fmt.Errorf("failed to write item.json: %w", err)
	}

	// Create state log with simple "created" entry
	stateLogPath := filepath.Join(itemDir, "state-log.txt")
	if err := os.WriteFile(stateLogPath, []byte("created\n"), 0644); err != nil {
		return fmt.Errorf("failed to write state-log.txt: %w", err)
	}

	return nil
}
