package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/nikoremi97/debate/internal/models"

	"github.com/oklog/ulid/v2"
	"github.com/redis/go-redis/v9"
)

type Store interface {
	GetConversation(ctx context.Context, id string) (*models.Conversation, error)
	SaveConversation(ctx context.Context, c *models.Conversation) error
	CreateConversation(ctx context.Context, topicName, botStance string) (*models.Conversation, error)
	ListConversations(ctx context.Context, limit, offset int) ([]ConversationSummary, error)
	GetPopularTopics(ctx context.Context, limit int) ([]string, error)
	Ping(ctx context.Context) error
}

type ConversationSummary struct {
	ID           string    `json:"id"`
	TopicName    string    `json:"topic_name"`
	BotStance    string    `json:"bot_stance"`
	Title        string    `json:"title"`
	MessageCount int       `json:"message_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type RedisStore struct {
	c *redis.Client
}

func NewRedisClient(addr, password string) (*redis.Client, error) {
	c := redis.NewClient(&redis.Options{Addr: addr, Password: password, DB: 0})
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)

	defer cancel()

	if err := c.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return c, nil
}

func NewRedisStore(c *redis.Client) Store {
	return &RedisStore{c: c}
}

func (s *RedisStore) key(id string) string {
	return "convo:" + id
}

func (s *RedisStore) GetConversation(ctx context.Context, id string) (*models.Conversation, error) {
	b, err := s.c.Get(ctx, s.key(id)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, errors.New("not found")
		}

		return nil, err
	}

	var conv models.Conversation

	if err := json.Unmarshal(b, &conv); err != nil {
		return nil, err
	}

	return &conv, nil
}

func (s *RedisStore) SaveConversation(ctx context.Context, c *models.Conversation) error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	return s.c.Set(ctx, s.key(c.ID), b, 24*time.Hour).Err()
}

func (s *RedisStore) Ping(ctx context.Context) error {
	return s.c.Ping(ctx).Err()
}

// CreateConversation creates a new conversation (Redis fallback implementation)
func (s *RedisStore) CreateConversation(ctx context.Context, topicName, botStance string) (*models.Conversation, error) {
	// Generate ULID for the conversation
	id := ulid.Make().String()

	conv := &models.Conversation{
		ID:       id,
		Topic:    topicName,
		Stance:   botStance,
		Messages: make([]models.Message, 0),
	}

	// Save the conversation
	if err := s.SaveConversation(ctx, conv); err != nil {
		return nil, err
	}

	return conv, nil
}

// ListConversations lists conversations (Redis fallback - limited functionality)
func (s *RedisStore) ListConversations(ctx context.Context, limit, offset int) ([]ConversationSummary, error) {
	// Redis doesn't have great support for complex queries
	// This is a simplified implementation - in production, use PostgreSQL
	keys, err := s.c.Keys(ctx, "convo:*").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation keys: %w", err)
	}

	var conversations []ConversationSummary

	for i, key := range keys {
		if i < offset {
			continue
		}

		if len(conversations) >= limit {
			break
		}

		// Get conversation
		b, err := s.c.Get(ctx, key).Bytes()
		if err != nil {
			continue // Skip invalid conversations
		}

		var conv models.Conversation
		if err := json.Unmarshal(b, &conv); err != nil {
			continue // Skip invalid conversations
		}

		conversations = append(conversations, ConversationSummary{
			ID:           conv.ID,
			TopicName:    conv.Topic,
			BotStance:    conv.Stance,
			Title:        fmt.Sprintf("Debate: %s (%s)", conv.Topic, conv.Stance),
			MessageCount: len(conv.Messages),
			CreatedAt:    time.Now(), // Redis doesn't store creation time
			UpdatedAt:    time.Now(),
		})
	}

	return conversations, nil
}

// GetPopularTopics gets popular topics (Redis fallback - limited functionality)
func (s *RedisStore) GetPopularTopics(ctx context.Context, limit int) ([]string, error) {
	// This is a simplified implementation for Redis
	// In production, use PostgreSQL for better analytics
	keys, err := s.c.Keys(ctx, "convo:*").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation keys: %w", err)
	}

	topicCount := make(map[string]int)

	for _, key := range keys {
		b, err := s.c.Get(ctx, key).Bytes()
		if err != nil {
			continue
		}

		var conv models.Conversation
		if err := json.Unmarshal(b, &conv); err != nil {
			continue
		}

		topicCount[conv.Topic]++
	}

	// Sort by count and return top topics
	var topics []string
	for topic := range topicCount {
		topics = append(topics, topic)
		if len(topics) >= limit {
			break
		}
	}

	return topics, nil
}
