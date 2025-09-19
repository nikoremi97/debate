package bot

import (
	"strings"
	"testing"
)

func TestPickTopicAndStance(t *testing.T) {
	topic, stance := PickTopicAndStance()

	if topic == "" {
		t.Fatal("topic should not be empty")
	}

	if stance != "PRO" {
		t.Fatalf("expected stance to be PRO, got %s", stance)
	}

	// Check that topic is one of the predefined topics
	validTopics := []string{
		"The Earth is flat",
		"Pineapple belongs on pizza",
		"Tabs are better than spaces",
		"Cats are better than dogs",
		"Remote work is superior to office work",
		"AI should be regulated heavily",
		"Soccer is more exciting than basketball",
	}

	found := false

	for _, validTopic := range validTopics {
		if topic == validTopic {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("topic %s is not in the predefined list", topic)
	}
}

func TestBuildMessages(t *testing.T) {
	topic := "Test topic"
	stance := "PRO"
	history := []HistoryItem{
		{Role: "user", Message: "Hello"},
		{Role: "bot", Message: "Hi there"},
	}
	userMessage := "Tell me more"

	messages := buildMessages(topic, stance, history, userMessage)

	// Should have system message + history + user message
	expectedLength := 1 + len(history) + 1
	if len(messages) != expectedLength {
		t.Fatalf("expected %d messages, got %d", expectedLength, len(messages))
	}

	// Check system message
	if messages[0]["role"] != "system" {
		t.Fatalf("first message should be system, got %s", messages[0]["role"])
	}

	// Check that system message contains topic and stance
	systemContent := messages[0]["content"]
	if !strings.Contains(systemContent, topic) {
		t.Fatalf("system message should contain topic %s", topic)
	}

	if !strings.Contains(systemContent, stance) {
		t.Fatalf("system message should contain stance %s", stance)
	}

	// Check history messages
	if messages[1]["role"] != "user" {
		t.Fatalf("second message should be user, got %s", messages[1]["role"])
	}

	if messages[2]["role"] != "assistant" {
		t.Fatalf("third message should be assistant (converted from bot), got %s", messages[2]["role"])
	}

	// Check user message
	lastMessage := messages[len(messages)-1]
	if lastMessage["role"] != "user" {
		t.Fatalf("last message should be user, got %s", lastMessage["role"])
	}

	if lastMessage["content"] != userMessage {
		t.Fatalf("last message content should be %s, got %s", userMessage, lastMessage["content"])
	}
}

func TestValidateTopic(t *testing.T) {
	tests := []struct {
		name     string
		topic    string
		expected bool
	}{
		{"Valid topic", "Climate change is real", true},
		{"Valid topic with spaces", "  Remote work is better  ", true},
		{"Empty topic", "", false},
		{"Whitespace only", "   ", false},
		{"Violence content", "Violence is good", false},
		{"Sexual content", "Sex education in schools", false},
		{"Hate content", "Racism is acceptable", false},
		{"Drug content", "Drugs should be legal", false},
		{"Illegal content", "Theft is justified", false},
		{"Too long", "This is a very long topic that exceeds the reasonable length limit and should be rejected because it goes beyond the maximum allowed characters for a topic which is set to 200 characters to prevent abuse and ensure reasonable debate topics", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateTopic(tt.topic)
			if result != tt.expected {
				t.Errorf("ValidateTopic(%q) = %v, expected %v", tt.topic, result, tt.expected)
			}
		})
	}
}

func TestDetermineStance(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{"Pro message", "I agree that climate change is real", "CON"},
		{"Con message", "I disagree with remote work", "PRO"},
		{"Pro keywords", "I support this idea and think it's good", "CON"},
		{"Con keywords", "This is bad and I hate it", "PRO"},
		{"Mixed message", "I think it's good but also bad", "CON"}, // proCount > conCount
		{"Unclear message", "Maybe it depends", "CON"},             // default to CON
		{"Empty message", "", "CON"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetermineStance(tt.message)
			if result != tt.expected {
				t.Errorf("DetermineStance(%q) = %s, expected %s", tt.message, result, tt.expected)
			}
		})
	}
}

func TestProcessUserTopic(t *testing.T) {
	tests := []struct {
		name           string
		topic          string
		message        string
		expectedTopic  string
		expectedStance string
		expectFallback bool
	}{
		{"Valid topic with pro message", "Climate change is real", "I agree with this", "Climate change is real", "CON", false},
		{"Valid topic with con message", "Remote work is better", "I disagree with this", "Remote work is better", "PRO", false},
		{"Invalid topic", "Violence is good", "I support this", "", "PRO", true}, // fallback topic
		{"Empty topic", "", "I agree", "", "PRO", true},                          // fallback topic
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			topic, stance := ProcessUserTopic(tt.topic, tt.message)

			// Check topic
			if tt.expectFallback {
				assertValidFallbackTopic(t, topic, tt.topic, tt.message)
			} else if topic != tt.expectedTopic {
				t.Errorf("ProcessUserTopic(%q, %q) topic = %q, expected %q", tt.topic, tt.message, topic, tt.expectedTopic)
			}

			// Check stance
			if stance != tt.expectedStance {
				t.Errorf("ProcessUserTopic(%q, %q) stance = %q, expected %q", tt.topic, tt.message, stance, tt.expectedStance)
			}
		})
	}
}

func assertValidFallbackTopic(t *testing.T, topic, inputTopic, inputMessage string) {
	t.Helper()

	validFallback := false

	for _, fallback := range fallbackTopics {
		if topic == fallback {
			validFallback = true
			break
		}
	}

	if !validFallback {
		t.Errorf("ProcessUserTopic(%q, %q) returned topic %q, expected a fallback topic", inputTopic, inputMessage, topic)
	}
}
