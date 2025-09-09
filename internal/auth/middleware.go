package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	// APIKeyHeader is the header name for the API key
	APIKeyHeader = "X-API-Key" //nolint:gosec // This is a header name, not a credential
)

// Middleware creates a Gin middleware for API key authentication
func Middleware(authService ServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip authentication for health and ready endpoints
		if c.Request.URL.Path == "/health" || c.Request.URL.Path == "/ready" {
			c.Next()

			return
		}

		// Get API key from header
		apiKey := c.GetHeader(APIKeyHeader)
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "API key required",
				"code":  "MISSING_API_KEY",
			})
			c.Abort()

			return
		}

		// Validate API key
		if err := authService.ValidateAPIKey(c.Request.Context(), apiKey); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid API key",
				"code":  "INVALID_API_KEY",
			})
			c.Abort()

			return
		}

		// API key is valid, continue to next handler
		c.Next()
	}
}
