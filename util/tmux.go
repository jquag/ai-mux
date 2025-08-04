package util

import (
	"os"
	"os/exec"
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

// CreateTmuxWindow creates a new tmux window in the specified session or current session
func CreateTmuxWindow(windowName string, sessionName string) error {
	var cmd *exec.Cmd
	
	if sessionName == "" {
		// Create window in current session
		cmd = exec.Command("tmux", "new-window", "-d", "-n", windowName)
	} else {
		// Create window in specified session
		cmd = exec.Command("tmux", "new-window", "-d", "-t", sessionName, "-n", windowName)
	}
	
	return cmd.Run()
}
