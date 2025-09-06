package storage

import (
	"context"
	"testing"

	"github.com/nikoremi97/debate/internal/models"
)

func TestRedisStoreIntegration(t *testing.T) {
	// Skip if Redis is not available
	client, err := NewRedisClient("localhost:6379", "")
	if err != nil {
		t.Skip("Redis not available, skipping integration test")
	}

	store := NewRedisStore(client)
	ctx := context.Background()

	// Test Ping
	if err := store.Ping(ctx); err != nil {
		t.Fatalf("ping should succeed: %v", err)
	}

	// Test SaveConversation
	conv := models.NewConversation("test-redis-123")
	conv.Topic = "Test topic"
	conv.Stance = "PRO"
	conv.Append(models.Message{Role: "user", Message: "Hello"})

	if err := store.SaveConversation(ctx, conv); err != nil {
		t.Fatalf("save conversation should succeed: %v", err)
	}

	// Test GetConversation
	retrieved, err := store.GetConversation(ctx, "test-redis-123")
	if err != nil {
		t.Fatalf("get conversation should succeed: %v", err)
	}

	if retrieved.ID != conv.ID {
		t.Fatalf("conversation ID mismatch: expected %s, got %s", conv.ID, retrieved.ID)
	}

	if retrieved.Topic != conv.Topic {
		t.Fatalf("conversation topic mismatch: expected %s, got %s", conv.Topic, retrieved.Topic)
	}

	if len(retrieved.Messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(retrieved.Messages))
	}

	// Test GetConversation with non-existent ID
	_, err = store.GetConversation(ctx, "non-existent")
	if err == nil {
		t.Fatal("get non-existent conversation should fail")
	}

	// Clean up
	client.Del(ctx, "convo:test-redis-123")
}
