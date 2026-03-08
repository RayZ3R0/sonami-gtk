package secrets

// Account-free mode: these functions are stubs.
// The app behaves as if a user is always authenticated.

const (
	refreshTokenKey = "refresh_token"
)

// HasRefreshToken always returns true in account-free mode.
func HasRefreshToken() bool {
	return true
}

// GetRefreshToken returns an empty string in account-free mode.
func GetRefreshToken() string {
	return ""
}

// SetRefreshToken is a no-op in account-free mode.
func SetRefreshToken(token string) error {
	return nil
}

// DeleteRefreshToken is a no-op in account-free mode.
func DeleteRefreshToken() error {
	return nil
}

// UserID returns a placeholder user ID in account-free mode.
// This prevents nil/empty checks from blocking favourite operations.
func UserID() string {
	return "0"
}
