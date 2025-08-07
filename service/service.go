package service

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jquag/ai-mux/component/alert"
	data "github.com/jquag/ai-mux/data"
	"github.com/jquag/ai-mux/util"
)

type StartSessionMsg struct {
	TmuxSessionMessage    string
	TmuxWindowMessage     string
	WorktreeFolderMessage string
	GitBranchMessage      string
	Error                 error
}

func StartSession(workitem *data.WorkItem) tea.Cmd {
	return func() tea.Msg {
		worktreePath, err := util.CreateWorktree(workitem.BranchName)
		if err != nil {
			return alert.Alert(fmt.Sprintf("Failed to create worktree: %v", err), alert.AlertTypeError)()
		}

		if err := setupTmuxWindow(workitem, worktreePath); err != nil {
			return alert.Alert(err.Error(), alert.AlertTypeError)()
		}

		// Start Claude Code in the tmux window
		if err := startClaudeInWindow(workitem); err != nil {
			return alert.Alert(fmt.Sprintf("Failed to start Claude: %v", err), alert.AlertTypeError)()
		}

		return nil
	}
}

func CloseSession(workitem *data.WorkItem) tea.Cmd {
	return func() tea.Msg {
		// Calculate session name once
		sessionName := ""
		if !util.InTmuxSession() {
			sessionName = "ai-mux"
		}

		// Check if work item has been started
		isStarted := workitem.Status != "created" && workitem.Status != ""

		if isStarted {
			// Get worktree path
			cwd, err := os.Getwd()
			if err != nil {
				return alert.Alert("Failed to get current directory: "+err.Error(), alert.AlertTypeError)()
			}
			mainFolderName := filepath.Base(cwd)
			worktreePath := filepath.Join("..", fmt.Sprintf("%s-worktrees", mainFolderName), workitem.BranchName)

			// Check if worktree is clean and tell claude to commit if needed
			if clean, err := util.IsWorktreeClean(worktreePath); err == nil && !clean {
				// Find Claude pane by custom variable
				claudePaneId, err := util.FindPaneByVariable(workitem.BranchName, sessionName, "role", "claude-ai")
				if err != nil {
					return alert.Alert("Could not find Claude pane: "+err.Error(), alert.AlertTypeError)()
				}
				
				err = util.RunCommandInTmuxPane(claudePaneId, "commit the changes")
				if err != nil {
					return alert.Alert("Failed to commit changes: "+err.Error(), alert.AlertTypeError)()
				}
				util.RunCommandInTmuxPane(claudePaneId, "C-m") // Have to send a carriage return by itself for claude
				workitem.IsClosing = true                      // Indicator so that when claude is done we will try the CloseSession again

				return nil // Need to wait for claude to finish commiting
			}

			// Remove tmux window
			util.KillTmuxWindow(workitem.BranchName, sessionName)

			// Remove git worktree
			util.RemoveWorktree(worktreePath)
		}

		// Always remove work item data at the end
		workItemDir := filepath.Join(util.AiMuxDir, workitem.Id)
		if err := os.RemoveAll(workItemDir); err != nil {
			return alert.Alert("Failed to remove work item data: "+err.Error(), alert.AlertTypeError)()
		}

		return data.WorkItemRemovedMsg{
			WorkItem: workitem,
		}
	}
}

func setupTmuxWindow(workitem *data.WorkItem, worktreePath string) error {
	sessionName := ""
	if !util.InTmuxSession() {
		sessionName = "ai-mux"
		// Ensure ai-mux session exists
		if _, err := util.EnsureTmuxSession("ai-mux"); err != nil {
			return fmt.Errorf("failed to create tmux session 'ai-mux': %w", err)
		}
	}
	
	// Create window
	if err := util.CreateTmuxWindow(workitem.BranchName, sessionName, worktreePath); err != nil {
		return fmt.Errorf("failed to create tmux window: %w", err)
	}
	
	// Start editor in the top pane (original pane)
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim" // Default fallback
	}
	if err := util.RunCommandInTmuxWindow(workitem.BranchName, sessionName, editor); err != nil {
		return fmt.Errorf("failed to start editor: %w", err)
	}
	
	// Create vertical split
	if err := util.SplitTmuxWindow(workitem.BranchName, sessionName, worktreePath); err != nil {
		return fmt.Errorf("failed to split tmux window: %w", err)
	}
	
	// Find the bottom pane (the newly created one) and set custom variables for identification
	// The new pane should be .1 (bottom pane)
	target := workitem.BranchName
	if sessionName != "" {
		target = sessionName + ":" + workitem.BranchName
	}
	bottomPane := target + ".{bottom}"
	
	util.SetPaneVariable(bottomPane, "role", "claude-ai")
	util.SetPaneVariable(bottomPane, "workitem-id", workitem.Id)
	
	return nil
}

func startClaudeInWindow(workitem *data.WorkItem) error {
	// Get the current directory name (main folder)
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	mainFolderName := filepath.Base(cwd)

	// Build the path to the ai-mux directory in the main tree
	// From worktree at ../mainFolder-worktrees/branchName to ../mainFolder/.ai-mux
	aiMuxDirPath := filepath.Join("..", "..", mainFolderName, ".ai-mux")
	settingsPath := filepath.Join(aiMuxDirPath, "claude-settings.json")

	// Determine permission mode
	permissionMode := "acceptEdits"
	if workitem.PlanMode {
		permissionMode = "plan"
	}

	// Build the claude command with initial prompt and AI_MUX_DIR environment variable
	claudeCmd := fmt.Sprintf("AI_MUX_DIR=%s claude --session-id %s --settings %s --permission-mode %s %s",
		aiMuxDirPath, workitem.Id, settingsPath, permissionMode, util.ShellQuote(workitem.Description))

	// Determine session name
	sessionName := ""
	if !util.InTmuxSession() {
		sessionName = "ai-mux"
	}
	
	// Find Claude pane by custom variable
	claudePaneId, err := util.FindPaneByVariable(workitem.BranchName, sessionName, "role", "claude-ai")
	if err != nil {
		return fmt.Errorf("could not find Claude pane: %w", err)
	}

	// Run the command in the Claude pane
	err = util.RunCommandInTmuxPane(claudePaneId, claudeCmd)

	// Wait for claude to bring up the trust prompt then automatically accept it
	// This is a hack but I didn't see any other way to do it and claude does not provide a notification when this prompt appears
	time.Sleep(2 * time.Second)
	util.RunCommandInTmuxPane(claudePaneId, "Enter")
	return err
}
