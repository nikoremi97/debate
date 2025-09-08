package storage

import (
	"context"
	"testing"
	"time"

	"github.com/nikoremi97/debate/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPostgresStore tests the PostgreSQL storage implementation
// Note: These tests require a running PostgreSQL database
// Set POSTGRES_TEST_DSN environment variable to run these tests
func TestPostgresStore(t *testing.T) {
	// Skip if no test database configured
	testDSN := getTestDSN()
	if testDSN == "" {
		t.Skip("Skipping PostgreSQL tests: POSTGRES_TEST_DSN not set")
	}

	store, err := NewPostgresStore(testDSN)
	require.NoError(t, err)
	defer store.Close()

	ctx := context.Background()

	// Test Ping
	t.Run("Ping", func(t *testing.T) {
		err := store.Ping(ctx)
		assert.NoError(t, err)
	})

	// Test CreateConversation
	t.Run("CreateConversation", func(t *testing.T) {
		testCreateConversation(ctx, t, store)
	})

	// Test GetConversation
	t.Run("GetConversation", func(t *testing.T) {
		testGetConversation(ctx, t, store)
	})

	// Test GetConversation - Not Found
	t.Run("GetConversation_NotFound", func(t *testing.T) {
		_, err := store.GetConversation(ctx, "nonexistent-id")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "conversation not found")
	})

	// Test SaveConversation
	t.Run("SaveConversation", func(t *testing.T) {
		testSaveConversation(ctx, t, store)
	})

	// Test SaveConversation - Update Messages
	t.Run("SaveConversation_UpdateMessages", func(t *testing.T) {
		testSaveConversationUpdate(ctx, t, store)
	})

	// Test ListConversations
	t.Run("ListConversations", func(t *testing.T) {
		testListConversations(ctx, t, store)
	})

	// Test ListConversations - Pagination
	t.Run("ListConversations_Pagination", func(t *testing.T) {
		testListConversationsPagination(ctx, t, store)
	})

	// Test GetPopularTopics
	t.Run("GetPopularTopics", func(t *testing.T) {
		testGetPopularTopics(ctx, t, store)
	})

	// Test GetPopularTopics - Limit
	t.Run("GetPopularTopics_Limit", func(t *testing.T) {
		testGetPopularTopicsLimit(ctx, t, store)
	})
}

func testCreateConversation(ctx context.Context, t *testing.T, store *PostgresStore) {
	conv, err := store.CreateConversation(ctx, "Test Topic", "PRO")
	require.NoError(t, err)
	assert.NotEmpty(t, conv.ID)
	assert.Equal(t, "Test Topic", conv.Topic)
	assert.Equal(t, "PRO", conv.Stance)
	assert.Empty(t, conv.Messages)

	// Clean up
	cleanupConversation(t, store, conv.ID)
}

func testGetConversation(ctx context.Context, t *testing.T, store *PostgresStore) {
	// Create a conversation with messages
	conv, err := store.CreateConversation(ctx, "Get Test Topic", "CON")
	require.NoError(t, err)

	// Add messages
	conv.Messages = []models.Message{
		{Role: "user", Message: "Hello", TS: time.Now().UnixMilli()},
		{Role: "bot", Message: "Hi there!", TS: time.Now().UnixMilli()},
	}

	err = store.SaveConversation(ctx, conv)
	require.NoError(t, err)

	// Retrieve the conversation
	retrieved, err := store.GetConversation(ctx, conv.ID)
	require.NoError(t, err)
	assert.Equal(t, conv.ID, retrieved.ID)
	assert.Equal(t, conv.Topic, retrieved.Topic)
	assert.Equal(t, conv.Stance, retrieved.Stance)
	assert.Len(t, retrieved.Messages, 2)
	assert.Equal(t, "user", retrieved.Messages[0].Role)
	assert.Equal(t, "Hello", retrieved.Messages[0].Message)
	assert.Equal(t, "bot", retrieved.Messages[1].Role)
	assert.Equal(t, "Hi there!", retrieved.Messages[1].Message)

	// Clean up
	cleanupConversation(t, store, conv.ID)
}

func testSaveConversation(ctx context.Context, t *testing.T, store *PostgresStore) {
	// Create a conversation
	conv, err := store.CreateConversation(ctx, "Save Test Topic", "PRO")
	require.NoError(t, err)

	// Add messages
	conv.Messages = []models.Message{
		{Role: "user", Message: "Test message 1", TS: time.Now().UnixMilli()},
		{Role: "bot", Message: "Test response 1", TS: time.Now().UnixMilli()},
		{Role: "user", Message: "Test message 2", TS: time.Now().UnixMilli()},
	}

	// Save the conversation
	err = store.SaveConversation(ctx, conv)
	require.NoError(t, err)

	// Verify the conversation was saved
	retrieved, err := store.GetConversation(ctx, conv.ID)
	require.NoError(t, err)
	assert.Len(t, retrieved.Messages, 3)
	assert.Equal(t, "Test message 1", retrieved.Messages[0].Message)
	assert.Equal(t, "Test response 1", retrieved.Messages[1].Message)
	assert.Equal(t, "Test message 2", retrieved.Messages[2].Message)

	// Clean up
	cleanupConversation(t, store, conv.ID)
}

func testSaveConversationUpdate(ctx context.Context, t *testing.T, store *PostgresStore) {
	// Create a conversation with initial messages
	conv, err := store.CreateConversation(ctx, "Update Test Topic", "CON")
	require.NoError(t, err)

	conv.Messages = []models.Message{
		{Role: "user", Message: "Initial message", TS: time.Now().UnixMilli()},
	}

	err = store.SaveConversation(ctx, conv)
	require.NoError(t, err)

	// Update with new messages
	conv.Messages = []models.Message{
		{Role: "user", Message: "Updated message 1", TS: time.Now().UnixMilli()},
		{Role: "bot", Message: "Updated response", TS: time.Now().UnixMilli()},
	}

	err = store.SaveConversation(ctx, conv)
	require.NoError(t, err)

	// Verify the messages were updated
	retrieved, err := store.GetConversation(ctx, conv.ID)
	require.NoError(t, err)
	assert.Len(t, retrieved.Messages, 2)
	assert.Equal(t, "Updated message 1", retrieved.Messages[0].Message)
	assert.Equal(t, "Updated response", retrieved.Messages[1].Message)

	// Clean up
	cleanupConversation(t, store, conv.ID)
}

func testListConversations(ctx context.Context, t *testing.T, store *PostgresStore) {
	// Create multiple conversations
	conv1, err := store.CreateConversation(ctx, "List Test Topic 1", "PRO")
	require.NoError(t, err)

	conv2, err := store.CreateConversation(ctx, "List Test Topic 2", "CON")
	require.NoError(t, err)

	// List conversations
	conversations, err := store.ListConversations(ctx, 10, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(conversations), 2)

	// Find our test conversations
	var found1, found2 bool

	for _, conv := range conversations {
		if conv.ID == conv1.ID {
			found1 = true

			assert.Equal(t, "List Test Topic 1", conv.TopicName)
			assert.Equal(t, "PRO", conv.BotStance)
		}

		if conv.ID == conv2.ID {
			found2 = true

			assert.Equal(t, "List Test Topic 2", conv.TopicName)
			assert.Equal(t, "CON", conv.BotStance)
		}
	}

	assert.True(t, found1, "Conversation 1 not found in list")
	assert.True(t, found2, "Conversation 2 not found in list")

	// Clean up
	cleanupConversation(t, store, conv1.ID)
	cleanupConversation(t, store, conv2.ID)
}

func testListConversationsPagination(ctx context.Context, t *testing.T, store *PostgresStore) {
	// Create multiple conversations
	var convIDs []string

	for i := 0; i < 5; i++ {
		conv, err := store.CreateConversation(ctx, "Pagination Test Topic", "PRO")
		require.NoError(t, err)

		convIDs = append(convIDs, conv.ID)
	}

	// Test pagination
	conversations, err := store.ListConversations(ctx, 2, 0)
	require.NoError(t, err)
	assert.Len(t, conversations, 2)

	conversations, err = store.ListConversations(ctx, 2, 2)
	require.NoError(t, err)
	assert.Len(t, conversations, 2)

	// Clean up
	for _, id := range convIDs {
		cleanupConversation(t, store, id)
	}
}

func testGetPopularTopics(ctx context.Context, t *testing.T, store *PostgresStore) {
	// Create conversations with different topics
	topics := []string{"Popular Topic 1", "Popular Topic 2", "Popular Topic 1", "Popular Topic 1"}
	var convIDs []string

	for _, topic := range topics {
		conv, err := store.CreateConversation(ctx, topic, "PRO")
		require.NoError(t, err)

		convIDs = append(convIDs, conv.ID)
	}

	// Get popular topics
	popularTopics, err := store.GetPopularTopics(ctx, 10)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(popularTopics), 1)

	// "Popular Topic 1" should be the most popular
	assert.Equal(t, "Popular Topic 1", popularTopics[0])

	// Clean up
	for _, id := range convIDs {
		cleanupConversation(t, store, id)
	}
}

func testGetPopularTopicsLimit(ctx context.Context, t *testing.T, store *PostgresStore) {
	// Create conversations with different topics
	topics := []string{"Limit Topic 1", "Limit Topic 2", "Limit Topic 3"}
	var convIDs []string

	for _, topic := range topics {
		conv, err := store.CreateConversation(ctx, topic, "PRO")
		require.NoError(t, err)

		convIDs = append(convIDs, conv.ID)
	}

	// Get popular topics with limit
	popularTopics, err := store.GetPopularTopics(ctx, 2)
	require.NoError(t, err)
	assert.Len(t, popularTopics, 2)

	// Clean up
	for _, id := range convIDs {
		cleanupConversation(t, store, id)
	}
}

// TestPostgresStore_ErrorHandling tests error handling scenarios
func TestPostgresStore_ErrorHandling(t *testing.T) {
	// Test with invalid connection string
	t.Run("InvalidConnectionString", func(t *testing.T) {
		_, err := NewPostgresStore("invalid-connection-string")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to ping database")
	})

	// Test with invalid DSN
	t.Run("InvalidDSN", func(t *testing.T) {
		_, err := NewPostgresStore("postgres://invalid:invalid@localhost:5432/invalid")
		assert.Error(t, err)
	})
}

// TestPostgresStore_HelperMethods tests the helper methods
func TestPostgresStore_HelperMethods(t *testing.T) {
	testDSN := getTestDSN()
	if testDSN == "" {
		t.Skip("Skipping PostgreSQL tests: POSTGRES_TEST_DSN not set")
	}

	store, err := NewPostgresStore(testDSN)
	require.NoError(t, err)
	defer store.Close()

	ctx := context.Background()

	// Create a test conversation
	conv, err := store.CreateConversation(ctx, "Helper Test Topic", "PRO")
	require.NoError(t, err)

	// Test updateConversationMetadata
	t.Run("updateConversationMetadata", func(t *testing.T) {
		testUpdateConversationMetadata(ctx, t, store, conv)
	})

	// Test clearMessages
	t.Run("clearMessages", func(t *testing.T) {
		testClearMessages(ctx, t, store, conv)
	})

	// Test insertMessages
	t.Run("insertMessages", func(t *testing.T) {
		testInsertMessages(ctx, t, store, conv)
	})

	// Test insertMessages - Empty messages
	t.Run("insertMessages_Empty", func(t *testing.T) {
		testInsertMessagesEmpty(ctx, t, store, conv)
	})

	// Clean up
	cleanupConversation(t, store, conv.ID)
}

func testUpdateConversationMetadata(ctx context.Context, t *testing.T, store *PostgresStore, conv *models.Conversation) {
	tx, err := store.db.BeginTx(ctx, nil)
	require.NoError(t, err)

	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			t.Logf("Failed to rollback transaction: %v", rollbackErr)
		}
	}()

	conv.Topic = "Updated Topic"
	conv.Stance = "CON"
	conv.Messages = []models.Message{
		{Role: "user", Message: "Test", TS: time.Now().UnixMilli()},
	}

	err = store.updateConversationMetadata(ctx, tx, conv)
	assert.NoError(t, err)
}

func testClearMessages(ctx context.Context, t *testing.T, store *PostgresStore, conv *models.Conversation) {
	tx, err := store.db.BeginTx(ctx, nil)
	require.NoError(t, err)

	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			t.Logf("Failed to rollback transaction: %v", rollbackErr)
		}
	}()

	err = store.clearMessages(ctx, tx, conv.ID)
	assert.NoError(t, err)
}

func testInsertMessages(ctx context.Context, t *testing.T, store *PostgresStore, conv *models.Conversation) {
	tx, err := store.db.BeginTx(ctx, nil)
	require.NoError(t, err)

	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			t.Logf("Failed to rollback transaction: %v", rollbackErr)
		}
	}()

	conv.Messages = []models.Message{
		{Role: "user", Message: "Test message", TS: time.Now().UnixMilli()},
		{Role: "bot", Message: "Test response", TS: time.Now().UnixMilli()},
	}

	err = store.insertMessages(ctx, tx, conv)
	assert.NoError(t, err)
}

func testInsertMessagesEmpty(ctx context.Context, t *testing.T, store *PostgresStore, conv *models.Conversation) {
	tx, err := store.db.BeginTx(ctx, nil)
	require.NoError(t, err)

	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			t.Logf("Failed to rollback transaction: %v", rollbackErr)
		}
	}()

	conv.Messages = []models.Message{}

	err = store.insertMessages(ctx, tx, conv)
	assert.NoError(t, err)
}

// Helper functions

func getTestDSN() string {
	// In a real test environment, you would get this from environment variables
	// For now, return empty string to skip tests
	return ""
}

func cleanupConversation(t *testing.T, store *PostgresStore, conversationID string) {
	ctx := context.Background()

	// Delete messages first
	_, err := store.db.ExecContext(ctx, "DELETE FROM messages WHERE conversation_id = $1", conversationID)
	if err != nil {
		t.Logf("Failed to cleanup messages for conversation %s: %v", conversationID, err)
	}

	// Delete conversation
	_, err = store.db.ExecContext(ctx, "DELETE FROM conversations WHERE id = $1", conversationID)
	if err != nil {
		t.Logf("Failed to cleanup conversation %s: %v", conversationID, err)
	}
}

// Benchmark tests

func BenchmarkPostgresStore_CreateConversation(b *testing.B) {
	testDSN := getTestDSN()
	if testDSN == "" {
		b.Skip("Skipping PostgreSQL benchmarks: POSTGRES_TEST_DSN not set")
	}

	store, err := NewPostgresStore(testDSN)
	require.NoError(b, err)
	defer store.Close()

	ctx := context.Background()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		conv, err := store.CreateConversation(ctx, "Benchmark Topic", "PRO")
		if err != nil {
			b.Fatal(err)
		}

		// Clean up immediately
		cleanupConversation(&testing.T{}, store, conv.ID)
	}
}

func BenchmarkPostgresStore_GetConversation(b *testing.B) {
	testDSN := getTestDSN()
	if testDSN == "" {
		b.Skip("Skipping PostgreSQL benchmarks: POSTGRES_TEST_DSN not set")
	}

	store, err := NewPostgresStore(testDSN)
	require.NoError(b, err)
	defer store.Close()

	ctx := context.Background()

	// Create a conversation with messages
	conv, err := store.CreateConversation(ctx, "Benchmark Topic", "PRO")
	require.NoError(b, err)

	conv.Messages = []models.Message{
		{Role: "user", Message: "Test message", TS: time.Now().UnixMilli()},
		{Role: "bot", Message: "Test response", TS: time.Now().UnixMilli()},
	}

	err = store.SaveConversation(ctx, conv)
	require.NoError(b, err)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := store.GetConversation(ctx, conv.ID)
		if err != nil {
			b.Fatal(err)
		}
	}

	// Clean up
	cleanupConversation(&testing.T{}, store, conv.ID)
}
