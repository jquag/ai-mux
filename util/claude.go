package util

type ClaudeHookPayload struct {
	SessionID      string `json:"session_id"`
	TranscriptPath string `json:"transcript_path"` // Path to conversation JSON
	CWD            string `json:"cwd"`             // The current working directory when the hook is invoked
	HookEventName  string `json:"hook_event_name"`
}

func HandleClaudeEvent(payload ClaudeHookPayload, aiMuxDir string) error {
	// Use the utility function to write status
	return WriteStatusLog(payload.SessionID, payload.HookEventName, aiMuxDir)
}
