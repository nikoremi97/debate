package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/nikoremi97/debate/internal/api"
	"github.com/nikoremi97/debate/internal/bot"
	"github.com/nikoremi97/debate/internal/storage"
)

func main() {
	port := getenv("PORT", "8080")
	openAIKey := os.Getenv("OPENAI_API_KEY")

	if openAIKey == "" {
		log.Println("WARNING: OPENAI_API_KEY is not set. The /chat endpoint will fail without it.")
	}

	openAIModel := getenv("OPENAI_MODEL", "gpt-4o-mini")
	redisAddr := os.Getenv("REDIS_ADDR") // e.g. "redis:6379"

	// Storage: try Redis, fallback to in-memory
	var store storage.Store

	if redisAddr != "" {
		client, err := storage.NewRedisClient(redisAddr, getenv("REDIS_PASSWORD", ""))
		if err != nil {
			log.Printf("redis unavailable: %v — falling back to memory store", err)
			store = storage.NewMemoryStore()
		} else {
			store = storage.NewRedisStore(client)
		}
	} else {
		store = storage.NewMemoryStore()
	}

	// Bot engine (OpenAI-backend)
	llm := bot.NewOpenAIEngine(openAIKey, openAIModel)

	r := gin.Default()

	// Configure CORS - must be before all routes
	r.Use(CORSMiddleware())

	// health endpoints
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })
	r.GET("/ready", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()

		if err := store.Ping(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"ok": false, "error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	api.RegisterRoutes(r, store, llm)

	log.Printf("listening on :%s", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}

	return def
}
