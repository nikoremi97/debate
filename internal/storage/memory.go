package storage

import (
	"context"
	"errors"
	"sync"

	"github.com/nikoremi97/debate/internal/models"
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
