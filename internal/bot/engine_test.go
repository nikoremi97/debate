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
