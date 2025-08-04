package service

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
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
		// Create worktree folder and branch
		worktreePath, newBranch, err := util.CreateWorktree(workitem.BranchName)
		if err != nil {
			return startinfo.Alert(startinfo.Model{
				Error: err,
			})
		}
		
		info := setupTmuxWindow(workitem)
		info.WorktreeFolderMessage = worktreePath
		if newBranch {
			info.GitBranchMessage = fmt.Sprintf("%s (new)", workitem.BranchName)
		} else {
			info.GitBranchMessage = workitem.BranchName
		}
		
		//TODO: kick off claude and editor
		return startinfo.Alert(info)
	}
}

func setupTmuxWindow(workitem *data.WorkItem) startinfo.Model {
		if util.InTmuxSession() {
			// Create window in current session
			if err := util.CreateTmuxWindow(workitem.BranchName, ""); err != nil {
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
				if err := util.CreateTmuxWindow(workitem.BranchName, "ai-mux"); err != nil {
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

