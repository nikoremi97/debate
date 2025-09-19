package bot

import (
	"crypto/rand"
	"math/big"
	"regexp"
	"strings"
)

// Safety validation patterns for inappropriate content
var (
	// Patterns for violence, explicit content, and inappropriate behavior
	inappropriatePatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)\b(violence|violent|kill|murder|assault|attack|harm|hurt|abuse|torture|suicide|self-harm)\b`),
		regexp.MustCompile(`(?i)\b(sex|sexual|porn|pornography|nude|naked|prostitution|rape|molest)\b`),
		regexp.MustCompile(`(?i)\b(hate|racist|racism|discrimination|bigotry|slur|offensive|inappropriate)\b`),
		regexp.MustCompile(`(?i)\b(drug|drugs|cocaine|heroin|marijuana|weed|alcohol abuse|addiction)\b`),
		regexp.MustCompile(`(?i)\b(illegal|crime|criminal|theft|fraud|scam|exploit)\b`),
	}

	// Fallback topics if user topic is inappropriate
	fallbackTopics = []string{
		"The Earth is flat",
		"Pineapple belongs on pizza",
		"Tabs are better than spaces",
		"Cats are better than dogs",
		"Remote work is superior to office work",
		"AI should be regulated heavily",
		"Soccer is more exciting than basketball",
	}
)

// ValidateTopic checks if a topic is appropriate for debate
func ValidateTopic(topic string) bool {
	if strings.TrimSpace(topic) == "" {
		return false
	}

	// Check against inappropriate patterns
	for _, pattern := range inappropriatePatterns {
		if pattern.MatchString(topic) {
			return false
		}
	}

	// Additional checks for length and content
	if len(topic) > 200 { // reasonable length limit
		return false
	}

	return true
}

// GetFallbackTopic returns a random safe topic
func GetFallbackTopic() string {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(fallbackTopics))))
	if err != nil {
		return fallbackTopics[0]
	}

	return fallbackTopics[n.Int64()]
}

// DetermineStance analyzes user message to determine if they're PRO or CON
func DetermineStance(userMessage string) string {
	// Simple keyword-based stance detection
	// This could be enhanced with more sophisticated NLP
	proKeywords := []string{"agree", "support", "pro", "for", "yes", "good", "benefit", "advantage", "favor", "like", "love", "prefer", "better", "best", "should", "must", "need", "important", "valuable", "worth", "right", "correct", "true"}
	conKeywords := []string{"disagree", "against", "con", "no", "bad", "harm", "disadvantage", "oppose", "hate", "dislike", "worse", "worst", "shouldn't", "mustn't", "don't need", "unimportant", "worthless", "wrong", "incorrect", "false"}

	userLower := strings.ToLower(userMessage)

	proCount := 0
	conCount := 0

	// Count pro keywords using word boundaries
	for _, keyword := range proKeywords {
		pattern := `\b` + regexp.QuoteMeta(keyword) + `\b`
		matched, _ := regexp.MatchString(pattern, userLower)

		if matched {
			proCount++
		}
	}

	// Count con keywords using word boundaries
	for _, keyword := range conKeywords {
		pattern := `\b` + regexp.QuoteMeta(keyword) + `\b`
		matched, _ := regexp.MatchString(pattern, userLower)

		if matched {
			conCount++
		}
	}

	// If user seems to be PRO, bot should be CON, and vice versa
	if proCount > conCount {
		return "CON"
	} else if conCount > proCount {
		return "PRO"
	}

	// Default to CON if unclear (bot takes opposite stance)
	return "CON"
}

// PickTopicAndStance selects a topic and assigns the bot the PRO stance by default.
// You can randomize stance as well if you prefer.
func PickTopicAndStance() (string, string) {
	topic := GetFallbackTopic()
	stance := "PRO" // keep consistent; change to random if desired

	return topic, stance
}

// ProcessUserTopic validates and processes a user-provided topic
func ProcessUserTopic(userTopic string, userMessage string) (string, string) {
	// Validate the topic
	if !ValidateTopic(userTopic) {
		// Use fallback topic if user topic is inappropriate
		return GetFallbackTopic(), "PRO"
	}

	// Determine bot stance based on user's position
	stance := DetermineStance(userMessage)

	return userTopic, stance
}

func buildMessages(topic, stance string, history []HistoryItem, userMessage string) []map[string]string {
	// System prompt: fix topic and stance and the debate persona
	sys := map[string]string{"role": "system", "content": `You are a debate chatbot.\n\nTopic: ` + topic + `\nYour stance: ` + stance + ` (stand your ground, never switch sides).\n\nCRITICAL RULES:\n- You MUST ONLY debate about the specified Topic: ` + topic + `\n- NEVER respond to or engage with different topics mentioned by the user\n- If the user mentions a different topic, politely redirect them back to the original debate topic\n- Stay focused on the original debate topic throughout the entire conversation\n\nGoals:\n- Be persuasive, calm, and structured.\n- Use short evidence and analogies.\n- Acknowledge counterpoints briefly, then reframe.\n- Keep responses concise (3-6 sentences).\n- Always bring the conversation back to the original topic if the user tries to change subjects.`}

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
