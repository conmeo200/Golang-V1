package repository

// AuthRepository defines the interface for authentication-related database operations.
// It is currently empty but can be extended with methods like FindByAPIKey, etc.
type AuthRepository interface {
	// Add methods here in the future, e.g.:
	// ValidateAPIKey(ctx context.Context, apiKey string) (bool, error)
}
