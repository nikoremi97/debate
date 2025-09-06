package models

import "time"

// ChatRequest is the incoming API payload.
type ChatRequest struct {
	ConversationID *string `json:"conversation_id"`
	Message        string  `json:"message"` // text
}

// ChatResponse is the outgoing API payload.
type ChatResponse struct {
	ConversationID string    `json:"conversation_id"`
	Messages       []Message `json:"message"`
}

// Message is a single turn.
type Message struct {
	Role    string `json:"role"` // "user" | "bot"
	Message string `json:"message"`
	TS      int64  `json:"ts"` // unix ms (for ordering if needed)
}

// Conversation state stored in the DB.
type Conversation struct {
	ID       string    `json:"id"`
	Topic    string    `json:"topic"`
	Stance   string    `json:"stance"` // e.g., PRO/CON
	Messages []Message `json:"messages"`
}

func NewConversation(id string) *Conversation {
	return &Conversation{ID: id, Messages: make([]Message, 0, 16)}
}

func (c *Conversation) Append(m Message) {
	m.TS = time.Now().UnixMilli()
	c.Messages = append(c.Messages, m)

	if len(c.Messages) > 200 { // cap growth defensively
		c.Messages = c.Messages[len(c.Messages)-200:]
	}
}

func (c *Conversation) History() []Message { return c.Messages }

// LastN returns at most n most-recent messages (user/bot mixed), newest last.
func (c *Conversation) LastN(n int) []Message {
	if n <= 0 || len(c.Messages) <= n {
		return append([]Message(nil), c.Messages...)
	}

	return append([]Message(nil), c.Messages[len(c.Messages)-n:]...)
}
