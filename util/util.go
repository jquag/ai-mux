package util

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jquag/ai-mux/data"
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

// ToSafeName converts a short name to be safe for tmux window names and git branch names
// by replacing spaces with dashes
func ToSafeName(shortName string) string {
	return strings.ReplaceAll(shortName, " ", "-")
}

// WriteStatusLog writes a status to the work item's status log
func WriteStatusLog(workItemId string, status string, aiMuxDir string) error {
	statusLogPath := filepath.Join(aiMuxDir, workItemId, "state-log.txt")
	
	// Ensure the directory exists
	dir := filepath.Dir(statusLogPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	// Open the file in append mode
	file, err := os.OpenFile(statusLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open status log: %w", err)
	}
	defer file.Close()
	
	// Write a newline first, then the status
	if _, err := file.WriteString("\n" + status); err != nil {
		return fmt.Errorf("failed to write to status log: %w", err)
	}
	
	return nil
}

// UpdateWorkItem updates an existing work item's JSON file
func UpdateWorkItem(item *data.WorkItem) error {
	// Update the item.json file
	itemPath := filepath.Join(AiMuxDir, item.Id, "item.json")
	itemData, err := json.MarshalIndent(item, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}
	if err := os.WriteFile(itemPath, itemData, 0644); err != nil {
		return fmt.Errorf("failed to write item.json: %w", err)
	}
	
	return nil
}

// SaveWorkItem creates a new work item with its directory and files
func SaveWorkItem(item *data.WorkItem) error {
	// Ensure .ai-mux directory exists
	if err := EnsureAiMuxDir(); err != nil {
		return err
	}

	// Create directory for this item using its UUID
	itemDir := filepath.Join(AiMuxDir, item.Id)
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
