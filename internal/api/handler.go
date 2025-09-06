package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/nikoremi97/debate/internal/bot"
	"github.com/nikoremi97/debate/internal/models"
	"github.com/nikoremi97/debate/internal/storage"
)

func RegisterRoutes(r *gin.Engine, store storage.Store, engine bot.Engine) {
	r.POST("/chat", handleChat(store, engine))
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

		conversation, err := getOrCreateConversation(ctx, store, req.ConversationID)
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
		}

		c.JSON(http.StatusOK, resp)
	}
}

func getOrCreateConversation(ctx context.Context, store storage.Store, conversationID *string) (*models.Conversation, error) {
	if conversationID == nil || *conversationID == "" {
		// start a new conversation
		conv := models.NewConversation(uuid.NewString())
		conv.Topic, conv.Stance = bot.PickTopicAndStance()

		return conv, nil
	}

	// try to get existing conversation
	conv, err := store.GetConversation(ctx, *conversationID)
	if err != nil {
		// conversation doesn't exist, create a new one with the provided ID
		conv = models.NewConversation(*conversationID)
		conv.Topic, conv.Stance = bot.PickTopicAndStance()

		return conv, nil
	}

	return conv, nil
}

func generateBotReply(ctx context.Context, engine bot.Engine, conv *models.Conversation, userMessage string) (string, error) {
	history := make([]bot.HistoryItem, len(conv.Messages))
	for i, msg := range conv.Messages {
		history[i] = bot.HistoryItem{Role: msg.Role, Message: msg.Message}
	}

	return engine.Generate(ctx, conv.Topic, conv.Stance, history, userMessage)
}
