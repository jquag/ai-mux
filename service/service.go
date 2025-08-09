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

func StartSession(workitem *data.WorkItem, mode string) tea.Cmd {
	return func() tea.Msg {
		safeName := util.ToSafeName(workitem.ShortName)
		worktreePath, err := util.CreateWorktree(safeName)
		if err != nil {
			return alert.Alert(fmt.Sprintf("Failed to create worktree: %v", err), alert.AlertTypeError)()
		}

		if err := setupTmuxWindow(workitem, worktreePath); err != nil {
			return alert.Alert(err.Error(), alert.AlertTypeError)()
		}

		// Start Claude Code in the tmux window
		if err := startClaudeInWindow(workitem, mode); err != nil {
			return alert.Alert(fmt.Sprintf("Failed to start Claude: %v", err), alert.AlertTypeError)()
		}

		return nil
	}
}

func ResumeSession(workitem *data.WorkItem) tea.Cmd {
	return func() tea.Msg {
		// Get worktree path
		cwd, err := os.Getwd()
		if err != nil {
			return alert.Alert(fmt.Sprintf("Failed to get current directory: %v", err), alert.AlertTypeError)()
		}
		mainFolderName := filepath.Base(cwd)
		safeName := util.ToSafeName(workitem.ShortName)
		worktreePath := filepath.Join("..", fmt.Sprintf("%s-worktrees", mainFolderName), safeName)
		
		// Ensure tmux window and panes are set up (will reuse existing if present)
		if err := setupTmuxWindow(workitem, worktreePath); err != nil {
			return alert.Alert(fmt.Sprintf("Failed to setup tmux window: %v", err), alert.AlertTypeError)()
		}
		
		// Resume Claude Code in the tmux window
		if err := startClaudeInWindow(workitem, "resume"); err != nil {
			return alert.Alert(fmt.Sprintf("Failed to resume Claude: %v", err), alert.AlertTypeError)()
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
			safeName := util.ToSafeName(workitem.ShortName)
			worktreePath := filepath.Join("..", fmt.Sprintf("%s-worktrees", mainFolderName), safeName)

			// Check if worktree is clean and tell claude to commit if needed
			if clean, err := util.IsWorktreeClean(worktreePath); err == nil && !clean {
				// Find Claude pane by custom variable
				claudePaneId, err := util.FindPaneByVariable(safeName, sessionName, "role", "claude-ai")
				if err != nil {
					return alert.Alert("Could not find Claude pane: "+err.Error(), alert.AlertTypeError)()
				}
				
				err = util.RunCommandInTmuxPane(claudePaneId, "commit the changes")
				if err != nil {
					return alert.Alert("Failed to commit changes: "+err.Error(), alert.AlertTypeError)()
				}
				workitem.IsClosing = true                      // Indicator so that when claude is done we will try the CloseSession again

				return nil // Need to wait for claude to finish commiting
			}

			// Remove tmux window
			util.KillTmuxWindow(safeName, sessionName)

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
	
	safeName := util.ToSafeName(workitem.ShortName)
	
	// Check if window already exists
	windowExists := util.WindowExists(safeName, sessionName)
	
	if !windowExists {
		// Create window if it doesn't exist
		if err := util.CreateTmuxWindow(safeName, sessionName, worktreePath); err != nil {
			return fmt.Errorf("failed to create tmux window: %w", err)
		}
		
		// Start editor in the top pane (original pane)
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vim" // Default fallback
		}
		if err := util.RunCommandInTmuxWindow(safeName, sessionName, editor); err != nil {
			return fmt.Errorf("failed to start editor: %w", err)
		}
	}
	
	// Check if Claude pane exists
	_, claudePaneErr := util.FindPaneByVariable(safeName, sessionName, "role", "claude-ai")
	
	if claudePaneErr != nil {
		// Claude pane doesn't exist, create the split
		if err := util.SplitTmuxWindow(safeName, sessionName, worktreePath); err != nil {
			return fmt.Errorf("failed to split tmux window: %w", err)
		}
		
		// Set custom variables for the bottom pane
		target := safeName
		if sessionName != "" {
			target = sessionName + ":" + safeName
		}
		bottomPane := target + ".{bottom}"
		
		util.SetPaneVariable(bottomPane, "role", "claude-ai")
		util.SetPaneVariable(bottomPane, "workitem-id", workitem.Id)
	}
	
	return nil
}

func startClaudeInWindow(workitem *data.WorkItem, mode string) error {
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

	// Build the claude command based on mode
	var claudeCmd string
	if mode == "resume" {
		// For resume, use --resume instead of --session-id and don't pass the prompt
		claudeCmd = fmt.Sprintf("AI_MUX_DIR=%s claude --resume %s --settings %s",
			aiMuxDirPath, workitem.Id, settingsPath)
	} else {
		// For start modes (default, plan, acceptEdits), use --session-id and pass the prompt
		claudeCmd = fmt.Sprintf("AI_MUX_DIR=%s claude --session-id %s --settings %s --permission-mode %s %s",
			aiMuxDirPath, workitem.Id, settingsPath, mode, util.ShellQuote(workitem.Description))
	}

	// Determine session name
	sessionName := ""
	if !util.InTmuxSession() {
		sessionName = "ai-mux"
	}
	
	// Find Claude pane by custom variable
	safeName := util.ToSafeName(workitem.ShortName)
	claudePaneId, err := util.FindPaneByVariable(safeName, sessionName, "role", "claude-ai")
	if err != nil {
		return fmt.Errorf("could not find Claude pane: %w", err)
	}

	// Run the command in the Claude pane
	err = util.RunCommandInTmuxPane(claudePaneId, claudeCmd)

	// For non-resume modes, wait for claude to bring up the trust prompt then automatically accept it
	// This is a hack but I didn't see any other way to do it and claude does not provide a notification when this prompt appears
	if mode != "resume" {
		time.Sleep(2 * time.Second)
		util.RunCommandInTmuxPane(claudePaneId, "Enter")
	}
	return err
}
