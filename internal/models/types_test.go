package models

import (
	"fmt"
	"testing"
	"time"
)

func TestNewConversation(t *testing.T) {
	id := "test-123"
	conv := NewConversation(id)

	if conv.ID != id {
		t.Fatalf("expected ID %s, got %s", id, conv.ID)
	}

	if conv.Messages == nil {
		t.Fatal("Messages should be initialized")
	}

	if len(conv.Messages) != 0 {
		t.Fatalf("expected empty messages, got %d", len(conv.Messages))
	}
}

func TestAppend(t *testing.T) {
	conv := NewConversation("test-123")

	msg := Message{Role: "user", Message: "Hello"}

	conv.Append(msg)

	if len(conv.Messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(conv.Messages))
	}

	if conv.Messages[0].Role != msg.Role {
		t.Fatalf("expected role %s, got %s", msg.Role, conv.Messages[0].Role)
	}

	if conv.Messages[0].Message != msg.Message {
		t.Fatalf("expected message %s, got %s", msg.Message, conv.Messages[0].Message)
	}

	if conv.Messages[0].TS == 0 {
		t.Fatal("timestamp should be set")
	}
}

func TestAppendTimestamp(t *testing.T) {
	conv := NewConversation("test-123")

	before := time.Now().UnixMilli()

	conv.Append(Message{Role: "user", Message: "Hello"})
	after := time.Now().UnixMilli()

	ts := conv.Messages[0].TS
	if ts < before || ts > after {
		t.Fatalf("timestamp %d should be between %d and %d", ts, before, after)
	}
}

func TestAppendGrowthCap(t *testing.T) {
	conv := NewConversation("test-123")

	// Add more than 200 messages
	for i := 0; i < 250; i++ {
		conv.Append(Message{Role: "user", Message: "Message"})
	}

	if len(conv.Messages) != 200 {
		t.Fatalf("expected 200 messages after cap, got %d", len(conv.Messages))
	}

	// Check that the last 200 messages are kept
	expectedMessage := "Message"
	if conv.Messages[0].Message != expectedMessage {
		t.Fatalf("expected first message to be %s, got %s", expectedMessage, conv.Messages[0].Message)
	}
}

func TestHistory(t *testing.T) {
	conv := NewConversation("test-123")

	conv.Append(Message{Role: "user", Message: "Hello"})
	conv.Append(Message{Role: "bot", Message: "Hi"})

	history := conv.History()

	if len(history) != 2 {
		t.Fatalf("expected 2 messages in history, got %d", len(history))
	}

	if history[0].Message != "Hello" {
		t.Fatalf("expected first message to be Hello, got %s", history[0].Message)
	}

	if history[1].Message != "Hi" {
		t.Fatalf("expected second message to be Hi, got %s", history[1].Message)
	}
}

func TestLastN(t *testing.T) {
	conv := NewConversation("test-123")

	// Add 5 messages
	for i := 0; i < 5; i++ {
		conv.Append(Message{Role: "user", Message: fmt.Sprintf("Message %d", i)})
	}

	// Test LastN with n = 3
	last3 := conv.LastN(3)
	if len(last3) != 3 {
		t.Fatalf("expected 3 messages, got %d", len(last3))
	}

	if last3[0].Message != "Message 2" {
		t.Fatalf("expected first message to be Message 2, got %s", last3[0].Message)
	}

	if last3[2].Message != "Message 4" {
		t.Fatalf("expected last message to be Message 4, got %s", last3[2].Message)
	}

	// Test LastN with n = 0 (should return all messages when n <= 0)
	last0 := conv.LastN(0)
	if len(last0) != 5 {
		t.Fatalf("expected 5 messages when n=0, got %d", len(last0))
	}

	// Test LastN with n > len(messages)
	last10 := conv.LastN(10)
	if len(last10) != 5 {
		t.Fatalf("expected 5 messages, got %d", len(last10))
	}
}
