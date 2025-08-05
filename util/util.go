package util

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func TitledBorderStyle(color lipgloss.Color, title string, width int) lipgloss.Style {
	style := lipgloss.NewStyle()
	if width <= 0 {
		return style
	}

	topLength := width + 2
	title = title[:min(len(title), topLength)]

	top := fmt.Sprintf("─%s%s", title, strings.Repeat("─", max(0, topLength-len(title)-1)))

	var borderWithTitle = lipgloss.Border{
		Top:          top,
		Bottom:       "─",
		Left:         "│",
		Right:        "│",
		TopLeft:      "╭",
		TopRight:     "╮",
		BottomLeft:   "╰",
		BottomRight:  "╯",
		MiddleLeft:   "├",
		MiddleRight:  "┤",
		Middle:       "┼",
		MiddleTop:    "┬",
		MiddleBottom: "┴",
	}
	return style.BorderStyle(borderWithTitle).BorderForeground(color)
}

const AiMuxDir = ".ai-mux"

func EnsureAiMuxDir() error {
	if _, err := os.Stat(AiMuxDir); os.IsNotExist(err) {
		if err := os.MkdirAll(AiMuxDir, 0755); err != nil {
			return fmt.Errorf("failed to create .ai-mux directory: %w", err)
		}
	}
	return nil
}

// ShellQuote escapes a string for safe use in shell commands
func ShellQuote(s string) string {
	// Replace single quotes with '\'' and wrap in single quotes
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}
