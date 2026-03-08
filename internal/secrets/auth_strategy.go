package secrets

import (
	"net/http"
)

type AuthStrategy struct{}

func NewTokenAuthStrategy() *AuthStrategy {
	return &AuthStrategy{}
}

// GetToken is a no-op in account-free mode. No OAuth tokens are needed.
func (s *AuthStrategy) GetToken(clientID, clientSecret string) (string, error) {
	return "", nil
}

// Authenticate is a no-op in account-free mode. No Authorization header is set.
func (s *AuthStrategy) Authenticate(req *http.Request, clientID, clientSecret string) error {
	return nil
}
