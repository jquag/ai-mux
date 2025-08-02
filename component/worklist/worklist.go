package worklist

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jquag/ai-mux/component/modal"
	"github.com/jquag/ai-mux/component/workform"
	workitem "github.com/jquag/ai-mux/data"
	"github.com/jquag/ai-mux/theme"
	"github.com/jquag/ai-mux/util"
)

type Model struct {
	width         int
	height        int
	viewport      viewport.Model
	workItems     []*workitem.WorkItem
	Overlayed     bool
	loading       bool
	selectedIndex int
}

func (m *Model) Init() tea.Cmd {
	m.loading = true
	return loadWorkItems
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			form := workform.New()
			initCmd := form.Init()
			return m, tea.Batch(initCmd, modal.ShowModal(form, "Add Work Item"))
		case "j", "down":
			if len(m.workItems) > m.selectedIndex+1 {
				m.selectedIndex++
			}
		case "k", "up":
			if m.selectedIndex > 0 {
				m.selectedIndex--
			}
		}
	case workitem.NewWorkItemMsg:
		m.workItems = append(m.workItems, msg.WorkItem)
	case loadItemsMsg:
		m.loading = false
		m.workItems = msg.items
		//TODO: handle error
		return m, m.startStatusPollers()
	case statusUpdateMsg:
		m.updateStatus(msg.item, msg.status)
		return m, calcStatus(msg.item, 3)
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

	if m.loading {
		body = "loading..."
	} else if len(m.workItems) == 0 {
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
		items[i] = m.itemView(item, i == m.selectedIndex)
	}
	return lipgloss.JoinVertical(lipgloss.Left, items...)
}

func (m *Model) itemView(item *workitem.WorkItem, selected bool) string {
	bg := lipgloss.NewStyle()
	if selected {
		bg = bg.Background(theme.Colors.BgDark)
	}
	lineStyle := lipgloss.NewStyle().Foreground(m.colorForStatus(item)).Inherit(bg)
	left := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Foreground(m.colorForStatus(item)).Inherit(bg).Render("● "),
		lineStyle.Render("│ "),
		lineStyle.Render("│ "),
		lineStyle.Render("╰─"),
	)

	nameColor := theme.Colors.Title
	descriptionColor := theme.Colors.Text

	if m.Overlayed {
		nameColor = theme.Colors.Muted
		descriptionColor = theme.Colors.Muted
	}

	centerWidth := m.width - lipgloss.Width(left) - 1
	name := lipgloss.NewStyle().
		Width(centerWidth).MaxWidth(centerWidth).MaxHeight(1).
		Foreground(nameColor).
		Inherit(bg).
		Render(item.BranchName)
	descr := lipgloss.NewStyle().
		Height(2).MaxHeight(2).Width(centerWidth).
		Foreground(descriptionColor).
		Inherit(bg).
		Render(item.Description)
	status := m.statusView(item, selected)

	right := ""
	// Check if name was truncated
	if lipgloss.Width(item.BranchName) > centerWidth {
		right = lipgloss.NewStyle().Foreground(theme.Colors.Muted).Inherit(bg).Render("…")
	} else {
		right = lipgloss.NewStyle().Foreground(theme.Colors.Muted).Inherit(bg).Render(" ")
	}
	// Check if description exceeds 2 lines when wrapped
	descrHeight := lipgloss.Height(lipgloss.NewStyle().Width(centerWidth).Render(item.Description))
	if descrHeight > 2 {
		right += lipgloss.NewStyle().Foreground(theme.Colors.Muted).Inherit(bg).Render("\n\n…")
	} else {
		right += lipgloss.NewStyle().Foreground(theme.Colors.Muted).Inherit(bg).Render("\n\n ")
	}

	right += lipgloss.NewStyle().Inherit(bg).Render("\n ")

	info := lipgloss.JoinVertical(lipgloss.Left, name, descr, status)
	return lipgloss.JoinHorizontal(lipgloss.Top, left, info, right) + "\n"
}

func (m *Model) statusView(item *workitem.WorkItem, selected bool) string {
	bg := lipgloss.NewStyle()
	if selected {
		bg = bg.Background(theme.Colors.BgDark)
	}
	status := ""

	switch item.Status {
	case "PreToolUse", "PostToolUse", "UserPromptSubmit":
		status = "Working"
	case "Notification":
		status = "Waiting for input"
	case "Stop":
		status = "Done"
	case "", "created":
		status = "Not Started"
	default:
		status = "Unknown"
	}

	statusStyle := lipgloss.NewStyle().Foreground(m.colorForStatus(item)).Width(m.width-3).Inherit(bg)
	return statusStyle.Render(fmt.Sprintf("[%s]", status))
}

func (m *Model) colorForStatus(item *workitem.WorkItem) lipgloss.TerminalColor {
	if m.Overlayed {
		return theme.Colors.Muted
	}

	switch item.Status {
	case "PreToolUse", "PostToolUse", "UserPromptSubmit":
		return theme.Colors.Success
	case "Notification":
		return theme.Colors.Primary
	case "Stop":
		return theme.Colors.Info
	case "", "created":
		return theme.Colors.Muted
	default:
		return theme.Colors.Error
	}
}

func (m *Model) SetWidth(width int) {
	m.viewport.Width = width
	m.width = width
}

func (m *Model) SetHeight(height int) {
	m.viewport.Height = height
	m.height = height
}

func (m *Model) startStatusPollers() tea.Cmd {
	cmds := make([]tea.Cmd, len(m.workItems))
	for i, item := range m.workItems {
		cmds[i] = calcStatus(item, 0)
	}
	return tea.Batch(cmds...)
}

func (m *Model) updateStatus(item *workitem.WorkItem, status string) {
	for _, existingItem := range m.workItems {
		if existingItem == item {
			existingItem.Status = status
			break
		}
	}
}

func New(width, height int) *Model {
	return &Model{
		width:    width,
		height:   height,
		viewport: viewport.New(width, height),
	}
}

type statusUpdateMsg struct {
	item   *workitem.WorkItem
	status string
}

func calcStatus(item *workitem.WorkItem, wait int) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(time.Duration(wait) * time.Second)

		// Read the last line from status-log.txt
		status := readLastStatus(item.Id)

		return statusUpdateMsg{
			item:   item,
			status: status,
		}
	}
}

func readLastStatus(itemId string) string {
	statusLogPath := filepath.Join(util.AiMuxDir, itemId, "state-log.txt")

	content, err := os.ReadFile(statusLogPath)
	if err != nil {
		return "unknown"
	}

	// Split content into lines and get the last non-empty line
	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	if len(lines) == 0 {
		return ""
	}

	lastLine := lines[len(lines)-1]
	status := lastLine
	return status
}

func loadWorkItems() tea.Msg {
	var items []*workitem.WorkItem

	// Check if .ai-mux directory exists
	if _, err := os.Stat(util.AiMuxDir); os.IsNotExist(err) {
		// No directory means no items to load
		return loadItemsMsg{err: err, items: items}
	}

	// Read all subdirectories in .ai-mux
	entries, err := os.ReadDir(util.AiMuxDir)
	if err != nil {
		return loadItemsMsg{err: fmt.Errorf("failed to read .ai-mux directory: %w", err), items: items}
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Read item.json from each subdirectory
		itemPath := filepath.Join(util.AiMuxDir, entry.Name(), "item.json")
		data, err := os.ReadFile(itemPath)
		if err != nil {
			// Skip items that can't be read
			continue
		}

		var item workitem.WorkItem
		if err := json.Unmarshal(data, &item); err != nil {
			// Skip items that can't be parsed
			continue
		}

		items = append(items, &item)
	}

	return loadItemsMsg{err: nil, items: items}
}

type loadItemsMsg struct {
	err   error
	items []*workitem.WorkItem
}
