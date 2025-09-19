package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/nikoremi97/debate/internal/bot"
	"github.com/nikoremi97/debate/internal/models"
	"github.com/nikoremi97/debate/internal/storage"
	"github.com/oklog/ulid/v2"
)

func RegisterRoutes(r *gin.Engine, store storage.Store, engine bot.Engine) {
	r.POST("/chat", handleChat(store, engine))
	RegisterConversationRoutes(r, store)
}

func handleChat(store storage.Store, engine bot.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.ChatRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 25*time.Second) // keep under 30s
		defer cancel()

		conversation, err := getOrCreateConversation(ctx, store, req.ConversationID, req.Topic, req.Message)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "conversation not found"})
			return
		}

		// append user message
		conversation.Append(models.Message{Role: "user", Message: req.Message})

		// generate bot reply
		reply, err := generateBotReply(ctx, engine, conversation, req.Message)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "llm error: " + err.Error()})
			return
		}

		conversation.Append(models.Message{Role: "bot", Message: reply})

		// persist (best effort)
		_ = store.SaveConversation(ctx, conversation)

		// build response with last 5 messages on both sides (max 10 total)
		resp := models.ChatResponse{
			ConversationID: conversation.ID,
			Messages:       conversation.LastN(10),
			Topic:          conversation.Topic,
			Stance:         conversation.Stance,
		}

		c.JSON(http.StatusOK, resp)
	}
}

func getOrCreateConversation(ctx context.Context, store storage.Store, conversationID *string, userTopic *string, userMessage string) (*models.Conversation, error) {
	// Determine conversation ID
	var convID string
	if conversationID == nil || *conversationID == "" {
		convID = ulid.Make().String()
	} else {
		convID = *conversationID
	}

	// Try to get existing conversation if we have a specific ID
	if conversationID != nil && *conversationID != "" {
		conv, err := store.GetConversation(ctx, *conversationID)
		if err == nil {
			return conv, nil
		}
	}

	// Create new conversation
	conv := models.NewConversation(convID)
	setConversationTopicAndStance(conv, userTopic, userMessage)

	return conv, nil
}

func setConversationTopicAndStance(conv *models.Conversation, userTopic *string, userMessage string) {
	if userTopic != nil && *userTopic != "" {
		conv.Topic, conv.Stance = bot.ProcessUserTopic(*userTopic, userMessage)
	} else {
		conv.Topic, conv.Stance = bot.PickTopicAndStance()
	}
}

func generateBotReply(ctx context.Context, engine bot.Engine, conv *models.Conversation, userMessage string) (string, error) {
	history := make([]bot.HistoryItem, len(conv.Messages))
	for i, msg := range conv.Messages {
		history[i] = bot.HistoryItem{Role: msg.Role, Message: msg.Message}
	}

	return engine.Generate(ctx, conv.Topic, conv.Stance, history, userMessage)
}
