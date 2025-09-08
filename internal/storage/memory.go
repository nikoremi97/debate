package storage

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/nikoremi97/debate/internal/models"

	"github.com/oklog/ulid/v2"
)

type memoryStore struct {
	mu   sync.RWMutex
	data map[string]*models.Conversation
}

func NewMemoryStore() Store { return &memoryStore{data: map[string]*models.Conversation{}} }

func (m *memoryStore) GetConversation(_ context.Context, id string) (*models.Conversation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.data[id]

	if !ok {
		return nil, errors.New("not found")
	}

	// return a shallow copy to avoid external mutation of map pointer
	copy := *c

	return &copy, nil
}

func (m *memoryStore) SaveConversation(_ context.Context, conv *models.Conversation) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// store a copy
	copy := *conv
	m.data[conv.ID] = &copy

	return nil
}

func (m *memoryStore) Ping(_ context.Context) error { return nil }

// CreateConversation creates a new conversation (memory implementation)
func (m *memoryStore) CreateConversation(_ context.Context, topicName, botStance string) (*models.Conversation, error) {
	// Generate ULID for the conversation
	id := ulid.Make().String()

	conv := &models.Conversation{
		ID:       id,
		Topic:    topicName,
		Stance:   botStance,
		Messages: make([]models.Message, 0),
	}

	// Save the conversation
	if err := m.SaveConversation(context.Background(), conv); err != nil {
		return nil, err
	}

	return conv, nil
}

// ListConversations lists conversations (memory implementation)
func (m *memoryStore) ListConversations(_ context.Context, limit, offset int) ([]ConversationSummary, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var conversations []ConversationSummary
	count := 0

	for _, conv := range m.data {
		if count < offset {
			count++
			continue
		}

		if len(conversations) >= limit {
			break
		}

		conversations = append(conversations, ConversationSummary{
			ID:           conv.ID,
			TopicName:    conv.Topic,
			BotStance:    conv.Stance,
			Title:        "Debate: " + conv.Topic + " (" + conv.Stance + ")",
			MessageCount: len(conv.Messages),
			CreatedAt:    time.Now(), // Memory store doesn't track creation time
			UpdatedAt:    time.Now(),
		})
		count++
	}

	return conversations, nil
}

// GetPopularTopics gets popular topics (memory implementation)
func (m *memoryStore) GetPopularTopics(_ context.Context, limit int) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	topicCount := make(map[string]int)
	for _, conv := range m.data {
		topicCount[conv.Topic]++
	}

	// Simple implementation - just return the first N topics
	var topics []string
	for topic := range topicCount {
		topics = append(topics, topic)
		if len(topics) >= limit {
			break
		}
	}

	return topics, nil
}
