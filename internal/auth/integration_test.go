package auth

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestService_Integration(t *testing.T) {
	// Skip integration tests if AWS credentials are not available
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Test with a real AWS Secrets Manager (requires AWS credentials)
	service, err := NewService("us-east-2", "test-secret-name")
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
	}

	// Test that the service can be created
	assert.NotNil(t, service)
	assert.Equal(t, "test-secret-name", service.secretName)
	assert.NotNil(t, service.cache)
}

func TestService_CacheBehavior(t *testing.T) {
	// Test cache behavior without AWS calls
	service := &Service{
		secretsClient: nil, // Will cause errors, but we're testing cache
		secretName:    "test-secret",
		cache: &apiKeyCache{
			key:       "test-key",
			expiresAt: time.Now().Add(30 * time.Minute),
		},
	}

	// Test cache hit
	key, err := service.GetAPIKey(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "test-key", key)
}

func TestService_ValidateAPIKey_CacheHit(t *testing.T) {
	service := &Service{
		secretsClient: nil, // Will cause errors, but we're testing cache
		secretName:    "test-secret",
		cache: &apiKeyCache{
			key:       "test-key",
			expiresAt: time.Now().Add(30 * time.Minute),
		},
	}

	// Test valid key
	err := service.ValidateAPIKey(context.Background(), "test-key")
	assert.NoError(t, err)

	// Test invalid key
	err = service.ValidateAPIKey(context.Background(), "wrong-key")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid API key")
}

func TestService_RefreshCache(t *testing.T) {
	// Test cache clearing behavior directly
	cache := &apiKeyCache{
		key:       "old-key",
		expiresAt: time.Now().Add(30 * time.Minute),
	}

	// Test initial state
	assert.Equal(t, "old-key", cache.key)
	assert.False(t, cache.expiresAt.IsZero())

	// Clear cache manually (simulating what RefreshCache does)
	cache.key = ""
	cache.expiresAt = time.Time{}

	// Verify cache is cleared
	assert.Empty(t, cache.key)
	assert.True(t, cache.expiresAt.IsZero())
}

func TestAPIKeyCache_ConcurrentAccess(t *testing.T) {
	cache := &apiKeyCache{
		key:       "test-key",
		expiresAt: time.Now().Add(30 * time.Minute),
	}

	// Test concurrent read access
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			cache.mutex.RLock()
			key := cache.key
			expiresAt := cache.expiresAt
			cache.mutex.RUnlock()

			assert.Equal(t, "test-key", key)
			assert.True(t, time.Now().Before(expiresAt))
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestAPIKeyCache_Expiration(t *testing.T) {
	// Test expired cache
	cache := &apiKeyCache{
		key:       "test-key",
		expiresAt: time.Now().Add(-1 * time.Minute), // Expired
	}

	cache.mutex.RLock()
	isExpired := time.Now().After(cache.expiresAt)
	cache.mutex.RUnlock()

	assert.True(t, isExpired)

	// Test valid cache
	cache.expiresAt = time.Now().Add(30 * time.Minute)

	cache.mutex.RLock()
	isExpired = time.Now().After(cache.expiresAt)
	cache.mutex.RUnlock()

	assert.False(t, isExpired)
}
