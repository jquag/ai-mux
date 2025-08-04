package util

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

// CreateWorktree creates a new git worktree with a new or existing branch
// Returns the worktree path and whether a new branch was created
func CreateWorktree(branchName string) (string, bool, error) {
	// Create worktree path in parent directory
	worktreePath := filepath.Join("..", branchName)
	
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