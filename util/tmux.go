package util

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func InTmuxSession() bool {
	_, exists := os.LookupEnv("TMUX")
	return exists
}

// EnsureTmuxSession creates a tmux session if it doesn't exist
func EnsureTmuxSession(sessionName string) (bool, error) {
	// Check if session exists
	checkCmd := exec.Command("tmux", "has-session", "-t", sessionName)
	err := checkCmd.Run()
	
	// If session doesn't exist, create it
	if err != nil {
		createCmd := exec.Command("tmux", "new-session", "-d", "-s", sessionName)
		if err := createCmd.Run(); err != nil {
			return false, err
		}
		return true, nil
	}
	
	return false, nil
}

// WindowExists checks if a tmux window exists
func WindowExists(windowName string, sessionName string) bool {
	target := windowName
	if sessionName != "" {
		target = sessionName + ":" + windowName
	}
	
	cmd := exec.Command("tmux", "list-windows", "-t", target, "-F", "#{window_name}")
	err := cmd.Run()
	return err == nil
}

// CreateTmuxWindow creates a new tmux window in the specified session or current session
// If workingDir is provided, the window will start in that directory
func CreateTmuxWindow(windowName string, sessionName string, workingDir string) error {
	args := []string{"new-window", "-d", "-n", windowName}
	
	if sessionName != "" {
		args = append(args, "-t", sessionName)
	}
	
	if workingDir != "" {
		args = append(args, "-c", workingDir)
	}
	
	cmd := exec.Command("tmux", args...)
	return cmd.Run()
}

// RunCommandInTmuxWindow runs a command in a specific tmux window
func RunCommandInTmuxWindow(windowName string, sessionName string, command string) error {
	target := windowName
	if sessionName != "" {
		target = sessionName + ":" + windowName
	}
	
	// Send the command to the tmux window
	cmd := exec.Command("tmux", "send-keys", "-t", target, command)
	cmd.Run()
	cmd = exec.Command("tmux", "send-keys", "-t", target, "Enter")
	cmd.Run()
	return nil
}

// RunCommandInTmuxPane runs a command in a specific tmux pane by pane ID
func RunCommandInTmuxPane(paneId string, command string) error {
	// Send the command to the tmux pane
	cmd := exec.Command("tmux", "send-keys", "-t", paneId, command)
	cmd.Run()
	cmd = exec.Command("tmux", "send-keys", "-t", paneId, "Enter")
	cmd.Run()
	return nil
}

// SplitTmuxWindow creates a vertical split in the specified window
func SplitTmuxWindow(windowName string, sessionName string, folder string) error {
	target := windowName
	if sessionName != "" {
		target = sessionName + ":" + windowName
	}
	
	// Create vertical split
	cmd := exec.Command("tmux", "split-window", "-v", "-t", target, "-c", folder)
	return cmd.Run()
}

// SetPaneVariable sets a custom variable on a specific pane
func SetPaneVariable(paneId string, variable string, value string) error {
	cmd := exec.Command("tmux", "set", "-p", "-t", paneId, "@"+variable, value)
	return cmd.Run()
}

// FindPaneByVariable finds a pane by searching for a custom variable value in a specific window
func FindPaneByVariable(windowName string, sessionName string, variable string, value string) (string, error) {
	target := windowName
	if sessionName != "" {
		target = sessionName + ":" + windowName
	}
	
	cmd := exec.Command("tmux", "list-panes", "-t", target, "-F", "#{pane_id} #{@"+variable+"}")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) == 2 && parts[1] == value {
			return parts[0], nil
		}
	}
	
	return "", fmt.Errorf("no pane found with %s=%s in window %s", variable, value, target)
}

// SwitchToTmuxWindow switches focus to a specific tmux window
func SwitchToTmuxWindow(windowName string, sessionName string) error {
	target := windowName
	if sessionName != "" {
		target = sessionName + ":" + windowName
	}
	
	// Switch to the tmux window
	cmd := exec.Command("tmux", "select-window", "-t", target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to switch to window '%s': %w - output: %s", target, err, string(output))
	}
	return nil
}

// KillTmuxWindow kills a specific tmux window
func KillTmuxWindow(windowName string, sessionName string) error {
	target := windowName
	if sessionName != "" {
		target = sessionName + ":" + windowName
	}
	
	// Kill the tmux window
	cmd := exec.Command("tmux", "kill-window", "-t", target)
	return cmd.Run()
}
