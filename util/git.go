package util

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// CreateWorktree creates a new git worktree with a new or existing branch
// Returns the worktree path and whether a new branch was created
func CreateWorktree(branchName string) (string, bool, error) {
	// Get the current directory name (main folder)
	cwd, err := os.Getwd()
	if err != nil {
		return "", false, fmt.Errorf("failed to get current directory: %w", err)
	}
	mainFolderName := filepath.Base(cwd)
	
	// Create worktree path in parent directory under worktrees folder
	worktreesDir := filepath.Join("..", fmt.Sprintf("%s-worktrees", mainFolderName))
	worktreePath := filepath.Join(worktreesDir, branchName)
	
	// Ensure the worktrees directory exists
	absWorktreesDir, err := filepath.Abs(worktreesDir)
	if err != nil {
		return "", false, fmt.Errorf("failed to get absolute path: %w", err)
	}
	if err := os.MkdirAll(absWorktreesDir, 0755); err != nil {
		return "", false, fmt.Errorf("failed to create worktrees directory: %w", err)
	}
	
	// Check if branch already exists
	checkCmd := exec.Command("git", "show-ref", "--verify", "--quiet", fmt.Sprintf("refs/heads/%s", branchName))
	branchExists := checkCmd.Run() == nil
	
	var cmd *exec.Cmd
	if branchExists {
		// Use existing branch
		cmd = exec.Command("git", "worktree", "add", worktreePath, branchName)
	} else {
		// Create new branch
		cmd = exec.Command("git", "worktree", "add", worktreePath, "-b", branchName)
	}
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", false, fmt.Errorf("failed to create worktree: %w - %s", err, string(output))
	}
	
	return worktreePath, !branchExists, nil
}