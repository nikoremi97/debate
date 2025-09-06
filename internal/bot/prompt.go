package bot

import (
	"crypto/rand"
	"math/big"
)

var topics = []string{
	"The Earth is flat",
	"Pineapple belongs on pizza",
	"Tabs are better than spaces",
	"Cats are better than dogs",
	"Remote work is superior to office work",
	"AI should be regulated heavily",
	"Soccer is more exciting than basketball",
}

// PickTopicAndStance selects a topic and assigns the bot the PRO stance by default.
// You can randomize stance as well if you prefer.
func PickTopicAndStance() (string, string) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(topics))))
	if err != nil {
		// Fallback to first topic if crypto/rand fails
		return topics[0], "PRO"
	}
	topic := topics[n.Int64()]
	stance := "PRO" // keep consistent; change to random if desired

	return topic, stance
}

func buildMessages(topic, stance string, history []HistoryItem, userMessage string) []map[string]string {
	// System prompt: fix topic and stance and the debate persona
	sys := map[string]string{"role": "system", "content": `You are a debate chatbot.\n\nTopic: ` + topic + `\nYour stance: ` + stance + ` (stand your ground, never switch sides).\n\nGoals:\n- Be persuasive, calm, and structured.\n- Stay on-topic for the defined Topic only.\n- Use short evidence and analogies.\n- Acknowledge counterpoints briefly, then reframe.\n- Keep responses concise (3-6 sentences).`}

	msgs := []map[string]string{sys}
	// Map history to OpenAI messages (convert "bot" -> "assistant")
	for _, h := range history {
		role := h.Role
		if role == "bot" {
			role = "assistant"
		}

		msgs = append(msgs, map[string]string{"role": role, "content": h.Message})
	}

	// Latest user message (already appended in handler, but include here to guide model)
	msgs = append(msgs, map[string]string{"role": "user", "content": userMessage})

	return msgs
}
