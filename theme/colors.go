package theme

import "github.com/charmbracelet/lipgloss"

type ThemeColors struct {
	Border  lipgloss.Color
	Primary lipgloss.Color
	Title   lipgloss.Color
	Muted   lipgloss.Color
	Text    lipgloss.Color
	Success lipgloss.Color
	Error   lipgloss.Color
	Info    lipgloss.Color
	BgDark  lipgloss.Color
}

var Colors = ThemeColors{
	Border:  lipgloss.Color("#89b4fa"),
	Primary: lipgloss.Color("#f9b387"),
	Title:   lipgloss.Color("#cba6f7"),
	Muted:   lipgloss.Color("#9298b1"),
	Text:    lipgloss.Color("#c6cfec"),
	Success: lipgloss.Color("#a7e2a1"),
	Error:   lipgloss.Color("#eba0ac"),
	Info:    lipgloss.Color("#81d1e0"),
	// BgDark:  lipgloss.Color("#181825"),
	// BgDark:  lipgloss.Color("#3f3d3b"),
	BgDark:  lipgloss.Color("#243b40"),

}
