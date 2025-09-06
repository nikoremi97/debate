package storage

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/nikoremi97/debate/internal/models"

	"github.com/redis/go-redis/v9"
)

type Store interface {
	GetConversation(ctx context.Context, id string) (*models.Conversation, error)
	SaveConversation(ctx context.Context, c *models.Conversation) error
	Ping(ctx context.Context) error
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
