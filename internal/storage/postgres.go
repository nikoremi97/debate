package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nikoremi97/debate/internal/models"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/oklog/ulid/v2"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(connStr string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) GetConversation(ctx context.Context, id string) (*models.Conversation, error) {
	query := `
		SELECT c.id, c.topic_name, c.bot_stance,
		       COALESCE(json_agg(
		           json_build_object(
		               'role', m.role,
		               'message', m.content,
		               'ts', extract(epoch from m.created_at) * 1000
		           ) ORDER BY m.created_at
		       ) FILTER (WHERE m.id IS NOT NULL), '[]'::json) as messages
		FROM conversations c
		LEFT JOIN messages m ON c.id = m.conversation_id
		WHERE c.id = $1
		GROUP BY c.id, c.topic_name, c.bot_stance
	`

	var conv models.Conversation
	var messagesJSON string

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&conv.ID,
		&conv.Topic,
		&conv.Stance,
		&messagesJSON,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("conversation not found")
		}

		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	// Parse messages from JSON
	if err := json.Unmarshal([]byte(messagesJSON), &conv.Messages); err != nil {
		return nil, fmt.Errorf("failed to parse messages: %w", err)
	}

	return &conv, nil
}

func (s *PostgresStore) SaveConversation(ctx context.Context, c *models.Conversation) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			// Log rollback error but don't return it to avoid masking the original error
		}
	}()

	if err := s.updateConversationMetadata(ctx, tx, c); err != nil {
		return err
	}

	if err := s.clearMessages(ctx, tx, c.ID); err != nil {
		return err
	}

	if err := s.insertMessages(ctx, tx, c); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *PostgresStore) updateConversationMetadata(ctx context.Context, tx *sql.Tx, c *models.Conversation) error {
	updateConv := `
		UPDATE conversations
		SET topic_name = $2, bot_stance = $3, message_count = $4, updated_at = NOW()
		WHERE id = $1
	`

	_, err := tx.ExecContext(ctx, updateConv, c.ID, c.Topic, c.Stance, len(c.Messages))
	if err != nil {
		return fmt.Errorf("failed to update conversation: %w", err)
	}

	return nil
}

func (s *PostgresStore) clearMessages(ctx context.Context, tx *sql.Tx, conversationID string) error {
	_, err := tx.ExecContext(ctx, "DELETE FROM messages WHERE conversation_id = $1", conversationID)
	if err != nil {
		return fmt.Errorf("failed to clear messages: %w", err)
	}

	return nil
}

func (s *PostgresStore) insertMessages(ctx context.Context, tx *sql.Tx, c *models.Conversation) error {
	if len(c.Messages) == 0 {
		return nil
	}

	insertMsg := `
		INSERT INTO messages (conversation_id, role, content, created_at)
		VALUES ($1, $2, $3, to_timestamp($4 / 1000.0))
	`

	stmt, err := tx.PrepareContext(ctx, insertMsg)
	if err != nil {
		return fmt.Errorf("failed to prepare message insert: %w", err)
	}
	defer stmt.Close()

	for _, msg := range c.Messages {
		_, err = stmt.ExecContext(ctx, c.ID, msg.Role, msg.Message, msg.TS)
		if err != nil {
			return fmt.Errorf("failed to insert message: %w", err)
		}
	}

	return nil
}

func (s *PostgresStore) CreateConversation(ctx context.Context, topicName, botStance string) (*models.Conversation, error) {
	// Generate ULID for the conversation
	id := ulid.Make().String()

	query := `
		INSERT INTO conversations (id, topic_name, bot_stance, title)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at
	`

	title := fmt.Sprintf("Debate: %s (%s)", topicName, botStance)

	var createdAt time.Time

	err := s.db.QueryRowContext(ctx, query, id, topicName, botStance, title).Scan(&createdAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	return &models.Conversation{
		ID:       id,
		Topic:    topicName,
		Stance:   botStance,
		Messages: make([]models.Message, 0),
	}, nil
}

func (s *PostgresStore) ListConversations(ctx context.Context, limit, offset int) ([]ConversationSummary, error) {
	query := `
		SELECT id, topic_name, bot_stance, title, message_count, created_at, updated_at
		FROM conversations
		ORDER BY updated_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := s.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list conversations: %w", err)
	}
	defer rows.Close()

	var conversations []ConversationSummary

	for rows.Next() {
		var conv ConversationSummary

		err := rows.Scan(
			&conv.ID,
			&conv.TopicName,
			&conv.BotStance,
			&conv.Title,
			&conv.MessageCount,
			&conv.CreatedAt,
			&conv.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan conversation: %w", err)
		}

		conversations = append(conversations, conv)
	}

	return conversations, nil
}

func (s *PostgresStore) GetPopularTopics(ctx context.Context, limit int) ([]string, error) {
	query := `
		SELECT topic_name, COUNT(*) as count
		FROM conversations
		GROUP BY topic_name
		ORDER BY count DESC
		LIMIT $1
	`

	rows, err := s.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get popular topics: %w", err)
	}
	defer rows.Close()

	var topics []string

	for rows.Next() {
		var topic string
		var count int

		err := rows.Scan(&topic, &count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan topic: %w", err)
		}

		topics = append(topics, topic)
	}

	return topics, nil
}

func (s *PostgresStore) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s *PostgresStore) Close() error {
	return s.db.Close()
}
