package auth

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// Service handles API key authentication with caching
type Service struct {
	secretsClient *secretsmanager.Client
	secretName    string
	cache         *apiKeyCache
}

// apiKeyCache holds the cached API key and metadata
type apiKeyCache struct {
	key       string
	expiresAt time.Time
	mutex     sync.RWMutex
}

// NewService creates a new authentication service
func NewService(region, secretName string) (*Service, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := secretsmanager.NewFromConfig(cfg)

	return &Service{
		secretsClient: client,
		secretName:    secretName,
		cache: &apiKeyCache{
			expiresAt: time.Time{}, // Will be set on first fetch
		},
	}, nil
}

// GetAPIKey retrieves the API key from cache or Secrets Manager
func (s *Service) GetAPIKey(ctx context.Context) (string, error) {
	// Check cache first
	s.cache.mutex.RLock()
	if s.cache.key != "" && time.Now().Before(s.cache.expiresAt) {
		key := s.cache.key
		s.cache.mutex.RUnlock()

		return key, nil
	}
	s.cache.mutex.RUnlock()

	// Cache miss or expired, fetch from Secrets Manager
	return s.fetchAndCacheAPIKey(ctx)
}

// fetchAndCacheAPIKey retrieves the API key from Secrets Manager and caches it
func (s *Service) fetchAndCacheAPIKey(ctx context.Context) (string, error) {
	s.cache.mutex.Lock()
	defer s.cache.mutex.Unlock()

	// Double-check in case another goroutine already fetched it
	if s.cache.key != "" && time.Now().Before(s.cache.expiresAt) {
		return s.cache.key, nil
	}

	// Fetch from Secrets Manager
	input := &secretsmanager.GetSecretValueInput{
		SecretId: &s.secretName,
	}

	result, err := s.secretsClient.GetSecretValue(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to get secret value: %w", err)
	}

	if result.SecretString == nil {
		return "", fmt.Errorf("secret value is nil")
	}

	// Cache the API key for 30 minutes
	s.cache.key = *result.SecretString
	s.cache.expiresAt = time.Now().Add(30 * time.Minute)

	log.Printf("API key refreshed from Secrets Manager, cached until %s", s.cache.expiresAt.Format(time.RFC3339))

	return s.cache.key, nil
}

// ValidateAPIKey validates the provided API key against the stored one
func (s *Service) ValidateAPIKey(ctx context.Context, providedKey string) error {
	expectedKey, err := s.GetAPIKey(ctx)
	if err != nil {
		return fmt.Errorf("failed to get API key: %w", err)
	}

	if providedKey != expectedKey {
		return fmt.Errorf("invalid API key")
	}

	return nil
}

// RefreshCache forces a refresh of the cached API key
func (s *Service) RefreshCache(ctx context.Context) error {
	s.cache.mutex.Lock()
	defer s.cache.mutex.Unlock()

	// Clear current cache
	s.cache.key = ""
	s.cache.expiresAt = time.Time{}

	// Fetch new key
	_, err := s.fetchAndCacheAPIKey(ctx)

	return err
}
