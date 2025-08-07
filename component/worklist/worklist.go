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
	"github.com/jquag/ai-mux/component/alert"
	"github.com/jquag/ai-mux/component/modal"
	"github.com/jquag/ai-mux/component/workform"
	"github.com/jquag/ai-mux/data"
	"github.com/jquag/ai-mux/service"
	"github.com/jquag/ai-mux/theme"
	"github.com/jquag/ai-mux/util"
	"slices"
)

type Model struct {
	width         int
	height        int
	viewport      viewport.Model
	workItems     []*data.WorkItem
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
		case "s":
			return m, m.startSelected()
		case "c":
			return m, m.closeSelected()
		}
	case data.NewWorkItemMsg:
		m.workItems = append(m.workItems, msg.WorkItem)
		return m, m.startStatusPoller(msg.WorkItem)
	case data.WorkItemRemovedMsg:
		m.removeWorkItem(msg.WorkItem.Id)
		return m, alert.Alert("Work item '"+msg.WorkItem.ShortName+"' closed successfully", alert.AlertTypeInfo)
	case loadItemsMsg:
		m.loading = false
		m.workItems = msg.items
		//TODO: handle error
		return m, m.startStatusPollers()
	case statusUpdateMsg:
		m.updateStatus(msg.item, msg.status)
		if msg.item.IsClosing && msg.status == "Stop" {
			//finished preping for close
			return m, tea.Batch(m.closeSelected(), calcStatus(msg.item, 3, false))
		}
		return m, calcStatus(msg.item, 3, false)
	}

	return m, nil
}

func (m *Model) View() string {
	titleColor := theme.Colors.Primary
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

func (m *Model) itemView(item *data.WorkItem, selected bool) string {
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
		Render(item.ShortName)
	descr := lipgloss.NewStyle().
		Height(2).MaxHeight(2).Width(centerWidth).
		Foreground(descriptionColor).
		Inherit(bg).
		Render(item.Description)
	status := m.statusView(item, selected)

	right := ""
	// Check if name was truncated
	if lipgloss.Width(item.ShortName) > centerWidth {
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
	return lipgloss.JoinHorizontal(lipgloss.Top, left, info, right)
}

func (m *Model) statusView(item *data.WorkItem, selected bool) string {
	bg := lipgloss.NewStyle()
	if selected {
		bg = bg.Background(theme.Colors.BgDark)
	}
	status := ""

	switch item.Status {
	case "PreToolUse", "PostToolUse", "UserPromptSubmit":
		status = "Working..."
	case "Starting":
		status = "Starting..."
	case "Notification":
		status = "Waiting for input"
	case "Stop":
		status = "Done"
	case "", "created":
		if item.PlanMode {
			status = "Plan Not Started"
		} else {
			status = "Not Started"
		}
	case "PrepForClosing":
		status = "Closing..."
	default:
		status = "Unknown"
	}

	if item.Status != "Notification" && item.IsClosing {
		status = "Closing..."
	}

	statusStyle := lipgloss.NewStyle().Foreground(m.colorForStatus(item)).Width(m.width - 3).Inherit(bg)
	return statusStyle.Render(fmt.Sprintf("[%s]", status))
}

func (m *Model) colorForStatus(item *data.WorkItem) lipgloss.TerminalColor {
	if m.Overlayed {
		return theme.Colors.Muted
	}

	if item.IsClosing && item.Status != "Notification" {
		return theme.Colors.Error
	}

	switch item.Status {
	case "PreToolUse", "PostToolUse", "UserPromptSubmit", "Starting":
		return theme.Colors.Success
	case "Notification":
		return theme.Colors.Primary
	case "Stop":
		return theme.Colors.Info
	case "", "created":
		return theme.Colors.Muted
	case "PrepForClosing":
		return theme.Colors.Error
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
		cmds[i] = m.startStatusPoller(item)
	}
	return tea.Batch(cmds...)
}

func (m *Model) startStatusPoller(item *data.WorkItem) tea.Cmd {
	// Start a poller for the specific item
	return calcStatus(item, 0, false)
}

func (m *Model) updateStatus(item *data.WorkItem, status string) {
	for _, existingItem := range m.workItems {
		if existingItem == item {
			existingItem.Status = status
			break
		}
	}
}

func (m *Model) startSelected() tea.Cmd {
	selected := m.getSelected()
	if selected == nil || (selected.Status != "created" && selected.Status != "") {
		return alert.Alert("This work item has alredy been started.", alert.AlertTypeWarning)
	}

	// Write PrepStarting status
	util.WriteStatusLog(selected.Id, "Starting", util.AiMuxDir)

	return tea.Batch(calcStatus(selected, 0, true), service.StartSession(selected))
}

func (m *Model) closeSelected() tea.Cmd {
	selected := m.getSelected()
	if selected == nil {
		return nil
	}

	// Write PrepStarting status
	util.WriteStatusLog(selected.Id, "PrepForClosing", util.AiMuxDir)

	return tea.Batch(calcStatus(selected, 0, true), service.CloseSession(selected))
}

func (m *Model) getSelected() *data.WorkItem {
	if m.selectedIndex >= 0 && m.selectedIndex < len(m.workItems) {
		return m.workItems[m.selectedIndex]
	}
	return nil
}

func (m *Model) removeWorkItem(id string) {
	for i, item := range m.workItems {
		if item.Id == id {
			m.workItems = slices.Delete(m.workItems, i, i+1)
			// Adjust selected index if necessary
			if m.selectedIndex >= len(m.workItems) && len(m.workItems) > 0 {
				m.selectedIndex = len(m.workItems) - 1
			} else if len(m.workItems) == 0 {
				m.selectedIndex = 0
			}
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
	item    *data.WorkItem
	status  string
	oneTime bool
}

func calcStatus(item *data.WorkItem, wait int, oneTime bool) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(time.Duration(wait) * time.Second)

		// Read the last line from status-log.txt
		status := readLastStatus(item.Id)

		return statusUpdateMsg{
			item:   item,
			status: status,
			oneTime: oneTime,
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
	var items []*data.WorkItem

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
		fileData, err := os.ReadFile(itemPath)
		if err != nil {
			// Skip items that can't be read
			continue
		}

		var item data.WorkItem
		if err := json.Unmarshal(fileData, &item); err != nil {
			// Skip items that can't be parsed
			continue
		}

		items = append(items, &item)
	}

	return loadItemsMsg{err: nil, items: items}
}

type loadItemsMsg struct {
	err   error
	items []*data.WorkItem
}
