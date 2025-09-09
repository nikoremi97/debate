package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockAuthService is a mock implementation of the auth service
type MockAuthService struct {
	ValidateAPIKeyFunc func(ctx context.Context, providedKey string) error
}

func (m *MockAuthService) GetAPIKey(ctx context.Context) (string, error) {
	return "", nil
}

func (m *MockAuthService) ValidateAPIKey(ctx context.Context, providedKey string) error {
	if m.ValidateAPIKeyFunc != nil {
		return m.ValidateAPIKeyFunc(ctx, providedKey)
	}

	return nil
}

func (m *MockAuthService) RefreshCache(ctx context.Context) error {
	return nil
}

// Ensure MockAuthService implements ServiceInterface
var _ ServiceInterface = (*MockAuthService)(nil)

func TestAuthMiddleware_HealthEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAuthService := &MockAuthService{}

	router := gin.New()
	router.Use(Middleware(mockAuthService))
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
}

func TestAuthMiddleware_ReadyEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAuthService := &MockAuthService{}

	router := gin.New()
	router.Use(Middleware(mockAuthService))
	router.GET("/ready", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/ready", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
}

func TestAuthMiddleware_MissingAPIKey(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAuthService := &MockAuthService{}

	router := gin.New()
	router.Use(Middleware(mockAuthService))
	router.POST("/chat", func(c *gin.Context) {
		t.Error("Handler should not be called when authentication fails")
	})

	req := httptest.NewRequest("POST", "/chat", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "API key required", response["error"])
	assert.Equal(t, "MISSING_API_KEY", response["code"])
}

func TestAuthMiddleware_ValidAPIKey(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAuthService := &MockAuthService{}
	mockAuthService.ValidateAPIKeyFunc = func(ctx context.Context, providedKey string) error {
		if providedKey == "valid-key" {
			return nil
		}

		return assert.AnError
	}

	router := gin.New()
	router.Use(Middleware(mockAuthService))
	router.POST("/chat", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("POST", "/chat", nil)
	req.Header.Set("X-API-Key", "valid-key")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
}

func TestAuthMiddleware_InvalidAPIKey(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAuthService := &MockAuthService{}
	mockAuthService.ValidateAPIKeyFunc = func(ctx context.Context, providedKey string) error {
		if providedKey == "valid-key" {
			return nil
		}

		return assert.AnError
	}

	router := gin.New()
	router.Use(Middleware(mockAuthService))
	router.POST("/chat", func(c *gin.Context) {
		t.Error("Handler should not be called when authentication fails")
	})

	req := httptest.NewRequest("POST", "/chat", nil)
	req.Header.Set("X-API-Key", "invalid-key")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid API key", response["error"])
	assert.Equal(t, "INVALID_API_KEY", response["code"])
}

func TestAuthMiddleware_CaseInsensitiveHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAuthService := &MockAuthService{}
	mockAuthService.ValidateAPIKeyFunc = func(ctx context.Context, providedKey string) error {
		if providedKey == "valid-key" {
			return nil
		}

		return assert.AnError
	}

	router := gin.New()
	router.Use(Middleware(mockAuthService))
	router.POST("/chat", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("POST", "/chat", nil)
	req.Header.Set("x-api-key", "valid-key") // lowercase
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
}

func TestAuthMiddleware_ContextPassing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAuthService := &MockAuthService{}
	mockAuthService.ValidateAPIKeyFunc = func(ctx context.Context, providedKey string) error {
		// Verify that the context is passed correctly
		assert.NotNil(t, ctx)
		return nil
	}

	router := gin.New()
	router.Use(Middleware(mockAuthService))
	router.POST("/chat", func(c *gin.Context) {
		// Verify that the context is passed correctly
		ctx := c.Request.Context()
		assert.NotNil(t, ctx)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("POST", "/chat", nil)
	req.Header.Set("X-API-Key", "valid-key")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddleware_HeaderName(t *testing.T) {
	assert.Equal(t, "X-API-Key", APIKeyHeader)
}

func TestAuthMiddleware_AbortBehavior(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAuthService := &MockAuthService{}

	router := gin.New()
	router.Use(Middleware(mockAuthService))
	router.POST("/chat", func(c *gin.Context) {
		t.Error("Handler should not be called when authentication fails")
	})

	req := httptest.NewRequest("POST", "/chat", nil)
	// No API key header
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// Benchmark tests
func BenchmarkAuthMiddleware_ValidKey(b *testing.B) {
	gin.SetMode(gin.TestMode)

	mockAuthService := &MockAuthService{}
	mockAuthService.ValidateAPIKeyFunc = func(ctx context.Context, providedKey string) error {
		return nil
	}

	router := gin.New()
	router.Use(Middleware(mockAuthService))
	router.POST("/chat", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("POST", "/chat", nil)
	req.Header.Set("X-API-Key", "valid-key")
	w := httptest.NewRecorder()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}

func BenchmarkAuthMiddleware_HealthEndpoint(b *testing.B) {
	gin.SetMode(gin.TestMode)

	mockAuthService := &MockAuthService{}

	router := gin.New()
	router.Use(Middleware(mockAuthService))
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}
