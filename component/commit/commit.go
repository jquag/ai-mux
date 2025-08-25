package commit

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/jquag/ai-mux/component/modal"
	"github.com/jquag/ai-mux/data"
	"github.com/jquag/ai-mux/theme"
	"github.com/jquag/ai-mux/util"
)

type Model struct {
	diffViewport viewport.Model
	width        int
	height       int
	form         *huh.Form
	submitted    bool
	item         *data.WorkItem
	diff         string
}

func New(item *data.WorkItem) *Model {
	vp := viewport.New(0, 0)

	confirmValue := true // Default to Submit

	keyMap := huh.NewDefaultKeyMap()
	keyMap.Text.NewLine = key.NewBinding(key.WithKeys("alt+enter"), key.WithHelp("alt+enter", "new line"))
	
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Key("message").
				Title("Commit Message").
				Lines(3),
			huh.NewConfirm().
				Key("done").
				Value(&confirmValue).
				Title("Commit changes?"),
		),
	).WithWidth(0).WithHeight(0).WithKeyMap(keyMap)

	form.Init()
	return &Model{
		diffViewport: vp,
		item:         item,
		form:         form,
	}
}

func (m Model) Init() tea.Cmd {
	diffCmd := func() tea.Msg {
		var err error
		var diff string

		folder, err := util.GetWortreeFolder(m.item)
		if err == nil {
			diff, err = util.GetColoredGitDiff(folder)
		}
		return DiffForCommitMsg{
			Diff: diff,
			Err:  err,
		}
	}
	return diffCmd
}

func (m *Model) Update(msg tea.Msg) (modal.ModalContent, tea.Cmd) {
	if m.submitted {
		return m, nil
	}

	if diffMsg, ok := msg.(DiffForCommitMsg); ok {
		if diffMsg.Err == nil {
			m.diff = diffMsg.Diff
		} else {
			m.diff = diffMsg.Err.Error()
		}
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "ctrl+j":
			m.diffViewport.ScrollDown(1)
			return m, nil
		case "ctrl+k":
			m.diffViewport.ScrollUp(1)
			return m, nil
		case "ctrl+d":
			m.diffViewport.ScrollDown(m.diffViewport.Height-1)
			return m, nil
		case "ctrl+u":
			m.diffViewport.ScrollUp(m.diffViewport.Height-1)
			return m, nil
		}
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f

		if m.form.State == huh.StateCompleted {
			m.submitted = true
		}
	}

	return m, cmd
}

func (m *Model) View() string {
	msg := lipgloss.
		NewStyle().
		Foreground(theme.Colors.Text).
		Render(m.item.ShortName + " - contains uncommitted changes!\n")

	formView := m.form.View()

	divider := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(theme.Colors.Border).
		Width(m.width).
		Render("")

	helpLabelStyle := lipgloss.NewStyle().Foreground(theme.Colors.Muted)
	helpKeyStyle := lipgloss.NewStyle().Foreground(theme.Colors.Primary)
	helpMsg := lipgloss.JoinHorizontal(
		lipgloss.Top, helpLabelStyle.Render("Scroll down: "),
		helpKeyStyle.Render("ctrl-j"),
		helpLabelStyle.Render(" | Scroll up: "),
		helpKeyStyle.Render("ctrl-k"),
		helpLabelStyle.Render(" | Page down: "),
		helpKeyStyle.Render("ctrl-d"),
		helpLabelStyle.Render(" | Page up: "),
		helpKeyStyle.Render("ctrl-u"),
	)

	topContent := lipgloss.JoinVertical(
		lipgloss.Left,
		msg,
		formView,
		divider,
		helpMsg,
		"",
	)

	m.diffViewport.Height = m.height - lipgloss.Height(topContent) - 4
	m.diffViewport.SetContent(m.diff)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		topContent,
		m.diffViewport.View(),
	)
}

func (m *Model) ShouldCloseOnEscape() bool {
	return true
}

func (m *Model) WithWidth(width int) modal.ModalContent {
	m.width = width
	m.form = m.form.WithWidth(m.width)
	return m
}

func (m *Model) WithHeight(height int) modal.ModalContent {
	m.height = height
	return m
}

type DiffForCommitMsg struct {
	Diff string
	Err  error
}
