package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nikoremi97/debate/internal/bot"
	"github.com/nikoremi97/debate/internal/storage"
)

// mock engine avoids real OpenAI calls for tests
type mockEngine struct{}

func (m mockEngine) Generate(ctx context.Context, topic, stance string, history []bot.HistoryItem, userMessage string) (string, error) {
	return "Test reply on topic: " + topic + " (" + stance + ")", nil
}

func TestChatStart(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	store := storage.NewMemoryStore()
	RegisterRoutes(r, store, mockEngine{})

	req := httptest.NewRequest("POST", "/chat", strings.NewReader(`{"conversation_id":null,"message":"Hello"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	// quick shape check
	if !strings.Contains(w.Body.String(), "conversation_id") {
		t.Fatalf("missing conversation_id in response: %s", w.Body.String())
	}

	if !strings.Contains(w.Body.String(), "message") {
		t.Fatalf("missing message history in response: %s", w.Body.String())
	}
}

func TestChatContinue(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	store := storage.NewMemoryStore()
	RegisterRoutes(r, store, mockEngine{})

	// First, start a conversation
	startReq := httptest.NewRequest("POST", "/chat", strings.NewReader(`{"conversation_id":null,"message":"Hello"}`))
	startReq.Header.Set("Content-Type", "application/json")
	startW := httptest.NewRecorder()
	r.ServeHTTP(startW, startReq)

	if startW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", startW.Code, startW.Body.String())
	}

	// Extract conversation_id from response (simplified - in real test you'd parse JSON)
	responseBody := startW.Body.String()
	if !strings.Contains(responseBody, "conversation_id") {
		t.Fatalf("missing conversation_id in start response: %s", responseBody)
	}

	// Continue the conversation (using a mock conversation ID)
	continueReq := httptest.NewRequest("POST", "/chat", strings.NewReader(`{"conversation_id":"test-123","message":"Tell me more"}`))
	continueReq.Header.Set("Content-Type", "application/json")
	continueW := httptest.NewRecorder()
	r.ServeHTTP(continueW, continueReq)

	// Should get 200 and create a new conversation with the provided ID
	if continueW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", continueW.Code, continueW.Body.String())
	}
}
