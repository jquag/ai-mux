package workitemdetails

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jquag/ai-mux/component/modal"
	"github.com/jquag/ai-mux/data"
	"github.com/jquag/ai-mux/theme"
	"github.com/jquag/ai-mux/util"
)

type Model struct {
	workItem *data.WorkItem
	viewport viewport.Model
	width    int
	height   int
}

func New(workItem *data.WorkItem) *Model {
	vp := viewport.New(0, 0)
	return &Model{
		workItem: workItem,
		viewport: vp,
	}
}

func (m *Model) Update(msg tea.Msg) (modal.ModalContent, tea.Cmd) {
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	content := m.buildContent()
	m.viewport.SetContent(content)
	return m.viewport.View()
}

func (m *Model) buildContent() string {
	var sections []string

	nameStyle := lipgloss.NewStyle().
		Foreground(theme.Colors.Primary).
		Bold(true)
	
	labelStyle := lipgloss.NewStyle().
		Foreground(theme.Colors.Text).
		Bold(true)
	
	valueStyle := lipgloss.NewStyle().
		Foreground(theme.Colors.Info)
	
	descStyle := lipgloss.NewStyle().
		Foreground(theme.Colors.Text)

	sections = append(sections, nameStyle.Render("Name"))
	sections = append(sections, valueStyle.Render(m.workItem.ShortName))
	sections = append(sections, "")

	sections = append(sections, nameStyle.Render("Description"))
	sections = append(sections, descStyle.Render(m.workItem.Description))
	sections = append(sections, "")

	isStarted := m.workItem.Status != "created" && m.workItem.Status != ""
	
	if isStarted {
		sections = append(sections, nameStyle.Render("Session Information"))
		
		safeName := util.ToSafeName(m.workItem.ShortName)
		
		sections = append(sections, labelStyle.Render("Tmux Window: ") + valueStyle.Render(safeName))
		
		sections = append(sections, labelStyle.Render("Git Branch: ") + valueStyle.Render(safeName))
		
		cwd, err := os.Getwd()
		var worktreePath string
		if err == nil {
			mainFolderName := filepath.Base(cwd)
			worktreePath = filepath.Join("..", fmt.Sprintf("%s-worktrees", mainFolderName), safeName)
			sections = append(sections, labelStyle.Render("Worktree Folder: ") + valueStyle.Render(worktreePath))
		}
		
		sections = append(sections, labelStyle.Render("Claude Session ID: ") + valueStyle.Render(m.workItem.Id))
		
		// Add git diff section
		if worktreePath != "" {
			sections = append(sections, "")
			sections = append(sections, nameStyle.
				Width(m.width).
				Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(theme.Colors.Muted).
				Render("Git Diff"))
			
			diff, err := util.GetColoredGitDiff(worktreePath)
			if err != nil {
				sections = append(sections, descStyle.Render(fmt.Sprintf("Error getting diff: %v", err)))
			} else {
				// The diff already contains ANSI color codes, so we append it directly
				sections = append(sections, diff)
			}
		}
	}

	return strings.Join(sections, "\n")
}

func (m *Model) ShouldCloseOnEscape() bool {
	return true
}

func (m *Model) WithWidth(width int) modal.ModalContent {
	m.width = width
	m.viewport.Width = width
	return m
}

func (m *Model) WithHeight(height int) modal.ModalContent {
	m.height = height
	m.viewport.Height = height - 4
	return m
}
