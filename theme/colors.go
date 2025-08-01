package theme

import "github.com/charmbracelet/lipgloss"

type ThemeColors struct {
	Border  lipgloss.Color
	Primary lipgloss.Color
	Title   lipgloss.Color
	Muted   lipgloss.Color
	Text    lipgloss.Color
}

var Colors = ThemeColors{
	Border:  lipgloss.Color("#89b4fa"),
	Primary: lipgloss.Color("#f9b387"),
	Title:   lipgloss.Color("#cba6f7"),
	Muted:   lipgloss.Color("#9298b1"),
	Text:    lipgloss.Color("#c6cfec"),
}
