package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nikoremi97/debate/internal/storage"
)

// ListConversationsResponse represents the response for listing conversations
type ListConversationsResponse struct {
	Conversations []storage.ConversationSummary `json:"conversations"`
	Total         int                           `json:"total"`
	Page          int                           `json:"page"`
	Limit         int                           `json:"limit"`
}

// PopularTopicsResponse represents the response for popular topics
type PopularTopicsResponse struct {
	Topics []string `json:"topics"`
}

// RegisterConversationRoutes registers conversation-related routes
func RegisterConversationRoutes(r *gin.Engine, store storage.Store) {
	conversations := r.Group("/conversations")
	{
		conversations.GET("", listConversations(store))
		conversations.GET("/topics", getPopularTopics(store))
		conversations.GET("/:id", getConversation(store))
	}
}

// listConversations handles GET /conversations
func listConversations(store storage.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse query parameters
		limitStr := c.DefaultQuery("limit", "20")
		offsetStr := c.DefaultQuery("offset", "0")

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 || limit > 100 {
			limit = 20
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			offset = 0
		}

		// Get conversations from store
		conversations, err := store.ListConversations(c.Request.Context(), limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to list conversations: " + err.Error(),
			})

			return
		}

		// For now, we don't have a total count method, so we'll use the length
		// In a real implementation, you'd want a separate method to get total count
		total := len(conversations)
		if len(conversations) == limit {
			// If we got exactly the limit, there might be more
			total = offset + len(conversations) + 1
		} else {
			total = offset + len(conversations)
		}

		response := ListConversationsResponse{
			Conversations: conversations,
			Total:         total,
			Page:          (offset / limit) + 1,
			Limit:         limit,
		}

		c.JSON(http.StatusOK, response)
	}
}

// getPopularTopics handles GET /conversations/topics
func getPopularTopics(store storage.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse query parameters
		limitStr := c.DefaultQuery("limit", "10")

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 || limit > 50 {
			limit = 10
		}

		// Get popular topics from store
		topics, err := store.GetPopularTopics(c.Request.Context(), limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get popular topics: " + err.Error(),
			})
			return
		}

		response := PopularTopicsResponse{
			Topics: topics,
		}

		c.JSON(http.StatusOK, response)
	}
}

// getConversation handles GET /conversations/:id
func getConversation(store storage.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		conversationID := c.Param("id")
		if conversationID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "conversation ID is required"})
			return
		}

		conversation, err := store.GetConversation(c.Request.Context(), conversationID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "conversation not found"})
			return
		}

		c.JSON(http.StatusOK, conversation)
	}
}
