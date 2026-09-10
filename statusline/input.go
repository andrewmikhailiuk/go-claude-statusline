// Package statusline implements the rendering pipeline for the Claude Code
// statusline binary. It exposes a Segment registry so new rows can be added
// in their own files without touching the core, plus chip helpers that work
// against any registered theme.Palette.
package statusline

// Input mirrors the JSON contract Claude Code writes to the binary's stdin.
// Only fields currently used by built-in segments are declared; extra fields
// in the payload are silently ignored by the JSON decoder.
type Input struct {
	Workspace Workspace `json:"workspace"`
	Model     Model     `json:"model"`

	// ContextWindow is absent from the payload on the first render of a
	// session, hence the pointer: nil means "no percentage yet".
	ContextWindow *ContextWindow `json:"context_window"`
}

// Workspace describes where the session is running.
type Workspace struct {
	CurrentDir string `json:"current_dir"`
}

// Model identifies the model backing the session.
type Model struct {
	DisplayName string `json:"display_name"`
	ID          string `json:"id"`
}

// ContextWindow carries context-window usage for the session.
type ContextWindow struct {
	UsedPercentage float64 `json:"used_percentage"`
}
