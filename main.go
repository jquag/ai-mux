package main

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jquag/ai-mux/component/app"
	"github.com/jquag/ai-mux/util"
)

//go:embed claudeSettings.json
var claudeSettings string


func checkCommand(name string) error {
	_, err := exec.LookPath(name)
	if err != nil {
		return fmt.Errorf("%s is not installed or not in PATH", name)
	}
	return nil
}

func checkSystemRequirements() error {
	requiredCommands := []string{"git", "claude", "tmux"}
	missingCommands := []string{}
	
	for _, cmd := range requiredCommands {
		if err := checkCommand(cmd); err != nil {
			missingCommands = append(missingCommands, cmd)
		}
	}
	
	if len(missingCommands) > 0 {
		return fmt.Errorf("missing required commands: %v", missingCommands)
	}
	
	return nil
}

func main() {
	// Check for --notification flag
	if len(os.Args) > 1 && os.Args[1] == "--notification" {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(data))
		os.Exit(0)
	}

	// Check and create .ai-mux directory and claude-settings.json if needed
	if err := util.EnsureAiMuxDir(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	
	settingsPath := filepath.Join(util.AiMuxDir, "claude-settings.json")
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		if err := os.WriteFile(settingsPath, []byte(claudeSettings), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing claude-settings.json: %v\n", err)
			os.Exit(1)
		}
	}

	if err := checkSystemRequirements(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	model := app.New()

	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("There's been an error: %v", err)
		os.Exit(1)
	}
}
