package bot

import "context"

// Engine is the debate text generator.
type Engine interface {
	// Generate returns the bot's reply given the topic, stance, history and latest user input.
	Generate(ctx context.Context, topic, stance string, history []HistoryItem, userMessage string) (string, error)
}

// HistoryItem is a compact view for prompts.
type HistoryItem struct {
	Role    string // "user" or "bot"
	Message string
}
