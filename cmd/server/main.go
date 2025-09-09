package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/nikoremi97/debate/internal/api"
	"github.com/nikoremi97/debate/internal/auth"
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

	// Initialize services
	authService := initializeAuthService()
	store := initializeStorage(redisAddr)
	llm := bot.NewOpenAIEngine(openAIKey, openAIModel)

	r := gin.Default()

	// Configure middleware
	r.Use(CORSMiddleware())
	setupAuthMiddleware(r, authService)

	// Register routes
	registerHealthRoutes(r, store)
	api.RegisterRoutes(r, store, llm)

	log.Printf("listening on :%s", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-API-Key, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func initializeAuthService() auth.ServiceInterface {
	// Check for local development mode
	if getenv("LOCAL_DEV", "false") == "true" {
		log.Println("LOCAL_DEV mode enabled - authentication disabled for local development")
		return nil
	}

	awsRegion := getenv("AWS_REGION", "us-east-2")
	apiKeySecretName := getenv("API_KEY_SECRET_NAME", "debate-chatbot-api-key")

	authService, err := auth.NewService(awsRegion, apiKeySecretName)
	if err != nil {
		log.Printf("WARNING: Failed to initialize auth service: %v. API will run without authentication.", err)
		return nil
	}

	return authService
}

func initializeStorage(redisAddr string) storage.Store {
	if redisAddr != "" {
		client, err := storage.NewRedisClient(redisAddr, getenv("REDIS_PASSWORD", ""))
		if err != nil {
			log.Printf("redis unavailable: %v — falling back to memory store", err)
			return storage.NewMemoryStore()
		}

		return storage.NewRedisStore(client)
	}

	return storage.NewMemoryStore()
}

func setupAuthMiddleware(r *gin.Engine, authService auth.ServiceInterface) {
	if authService != nil {
		r.Use(auth.Middleware(authService))
		log.Println("API key authentication enabled")
	} else {
		log.Println("WARNING: API key authentication disabled")
	}
}

func registerHealthRoutes(r *gin.Engine, store storage.Store) {
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
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}

	return def
}
