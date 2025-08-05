package service

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jquag/ai-mux/component/alert"
	"github.com/jquag/ai-mux/component/startinfo"
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
		worktreePath, newBranch, err := util.CreateWorktree(workitem.BranchName)
		if err != nil {
			return startinfo.Alert(startinfo.Model{
				Error: err,
			})
		}
		
		info := setupTmuxWindow(workitem, worktreePath)
		if info.Error != nil {
			return startinfo.Alert(info)
		}
		
		info.WorktreeFolderMessage = worktreePath
		if newBranch {
			info.GitBranchMessage = fmt.Sprintf("%s (new)", workitem.BranchName)
		} else {
			info.GitBranchMessage = workitem.BranchName
		}
		
		// Start Claude Code in the tmux window
		if err := startClaudeInWindow(workitem, info); err != nil {
			info.Error = err
			return startinfo.Alert(info)
		}
		
		return startinfo.Alert(info)
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
				util.RunCommandInTmuxWindow(workitem.BranchName, sessionName, "commit the changes")
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
		
		return alert.Alert("Work item '"+workitem.BranchName+"' closed successfully", alert.AlertTypeInfo)()
	}
}

func setupTmuxWindow(workitem *data.WorkItem, worktreePath string) startinfo.Model {
		if util.InTmuxSession() {
			// Create window in current session
			if err := util.CreateTmuxWindow(workitem.BranchName, "", worktreePath); err != nil {
				return startinfo.Model{
					Error: fmt.Errorf("failed to create tmux window: %w", err),
				}
			}
			return startinfo.Model{
				TmuxSessionMessage:    "(current)",
				TmuxWindowMessage:     fmt.Sprintf("%s (new)", workitem.BranchName),
				WorktreeFolderMessage: "",
				GitBranchMessage:      "",
			}
		} else {
			// Ensure ai-mux session exists
			if created, err := util.EnsureTmuxSession("ai-mux"); err != nil {
				return startinfo.Model{
					Error: fmt.Errorf("failed to create tmux session 'ai-mux': %w", err),
				}
			} else {
				// Create window in ai-mux session
				if err := util.CreateTmuxWindow(workitem.BranchName, "ai-mux", worktreePath); err != nil {
					return startinfo.Model{
						Error: fmt.Errorf("failed to create tmux window in ai-mux session: %w", err),
					}
				}
				sessionMsg := "ai-mux"
				if created {
					sessionMsg += " (new)"
				}
				return startinfo.Model{
					TmuxSessionMessage:    sessionMsg,
					TmuxWindowMessage:     fmt.Sprintf("%s (new)", workitem.BranchName),
					WorktreeFolderMessage: "",
					GitBranchMessage:      "",
				}
			}
		}
}

func startClaudeInWindow(workitem *data.WorkItem, info startinfo.Model) error {
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
	
	// Determine the session name based on tmux context
	sessionName := ""
	if !util.InTmuxSession() {
		sessionName = "ai-mux"
	}
	
	// Run the command in the tmux window
	return util.RunCommandInTmuxWindow(workitem.BranchName, sessionName, claudeCmd)
}

