package util

import (
	"fmt"
	"os"
	"path/filepath"
)

type ClaudeHookPayload struct {
	SessionID      string `json:"session_id"`
	TranscriptPath string `json:"transcript_path"` // Path to conversation JSON
	CWD            string `json:"cwd"`             // The current working directory when the hook is invoked
	HookEventName  string `json:"hook_event_name"`
}

func HandleClaudeEvent(payload ClaudeHookPayload, aiMuxDir string) error {
	// Construct the status log path
	statusLogPath := filepath.Join(aiMuxDir, payload.SessionID, "state-log.txt")
	
	// Ensure the directory exists
	dir := filepath.Dir(statusLogPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	// Open the file in append mode
	file, err := os.OpenFile(statusLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open status log: %w", err)
	}
	defer file.Close()
	
	// Write a newline first, then the hook event name
	if _, err := file.WriteString("\n" + payload.HookEventName); err != nil {
		return fmt.Errorf("failed to write to status log: %w", err)
	}
	
	return nil
}
