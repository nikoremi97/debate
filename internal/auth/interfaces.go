package auth

import "context"

// ServiceInterface defines the interface for authentication services
type ServiceInterface interface {
	GetAPIKey(ctx context.Context) (string, error)
	ValidateAPIKey(ctx context.Context, providedKey string) error
	RefreshCache(ctx context.Context) error
}
