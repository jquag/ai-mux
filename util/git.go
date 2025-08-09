package util

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// CreateWorktree creates a new git worktree with a new or existing branch
// Returns the worktree path
func CreateWorktree(branchName string) (string, error) {
	// Get the current directory name (main folder)
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}
	mainFolderName := filepath.Base(cwd)
	
	// Create worktree path in parent directory under worktrees folder
	worktreesDir := filepath.Join("..", fmt.Sprintf("%s-worktrees", mainFolderName))
	worktreePath := filepath.Join(worktreesDir, branchName)
	
	// Ensure the worktrees directory exists
	absWorktreesDir, err := filepath.Abs(worktreesDir)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}
	if err := os.MkdirAll(absWorktreesDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create worktrees directory: %w", err)
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
		return "", fmt.Errorf("failed to create worktree: %w - %s", err, string(output))
	}
	
	return worktreePath, nil
}

// RemoveWorktree removes a git worktree
func RemoveWorktree(worktreePath string) error {
	cmd := exec.Command("git", "worktree", "remove", worktreePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to remove worktree: %w - %s", err, string(output))
	}
	return nil
}

// IsWorktreeClean checks if a worktree has uncommitted changes
func IsWorktreeClean(worktreePath string) (bool, error) {
	cmd := exec.Command("git", "-C", worktreePath, "status", "--porcelain")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("failed to check worktree status: %w - %s", err, string(output))
	}
	
	// If output is empty, the worktree is clean
	return len(output) == 0, nil
}

// GetColoredGitDiff returns the git diff with ANSI color codes
func GetColoredGitDiff(worktreePath string) (string, error) {
	// Use git diff with color=always to force color output
	// Include both staged and unstaged changes
	cmd := exec.Command("git", "-C", worktreePath, "diff", "HEAD", "--color=always")
	output, err := cmd.CombinedOutput()
	if err != nil {
		// If there's an error but output exists, it might still be useful
		// (e.g., when there are no changes)
		if len(output) > 0 {
			return string(output), nil
		}
		return "", fmt.Errorf("failed to get git diff: %w", err)
	}
	
	// If there's no diff, return a message
	if len(output) == 0 {
		return "No changes in worktree", nil
	}
	
	return string(output), nil
}
