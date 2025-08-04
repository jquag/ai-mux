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
