package secrets

import (
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/zalando/go-keyring"
)

const (
	service         = "dev.dergs.tidalwave"
	refreshTokenKey = "refresh_token"
)

func HasRefreshToken() bool {
	hasToken, timedOut := withTimeout(time.Second*2, func() bool {
		_, err := keyring.Get(service, refreshTokenKey)
		return err == nil
	})
	if timedOut {
		slog.Error("timed out checking refresh token in keyring")
		return false
	}
	return hasToken
}

func GetRefreshToken() string {
	token, timedOut := withTimeout(time.Second*2, func() string {
		token, err := keyring.Get(service, refreshTokenKey)
		if err != nil {
			return ""
		}
		return token
	})
	if timedOut {
		slog.Error("timed out reading refresh token from keyring")
		return ""
	}
	return token
}

func SetRefreshToken(token string) error {
	err, timedOut := withTimeout(time.Second*2, func() error {
		return keyring.Set(service, refreshTokenKey, token)
	})
	if timedOut {
		slog.Error("timed out setting refresh token in keyring")
		return nil
	}
	return err
}

func DeleteRefreshToken() error {
	err, timedOut := withTimeout(time.Second*2, func() error {
		return keyring.Delete(service, refreshTokenKey)
	})
	if timedOut {
		slog.Error("timed out deleting refresh token from keyring")
		return nil
	}
	return err
}

func UserID() string {
	if !HasRefreshToken() {
		return ""
	}

	token := GetRefreshToken()

	// Split JWT into parts (header.payload.signature)
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return ""
	}

	// Decode the payload (second part)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return ""
	}

	// Parse JSON payload
	var claims struct {
		UID int `json:"uid"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return ""
	}

	return strconv.Itoa(claims.UID)
}

// withTimeout runs a callback function with a timeout.
// It returns what the callback returned along with a boolean indicating if it timed out.
func withTimeout[T any](timeout time.Duration, callback func() T) (T, bool) {
	resultChan := make(chan T, 1)

	go func() {
		resultChan <- callback()
	}()

	select {
	case result := <-resultChan:
		return result, false
	case <-time.After(timeout):
		return *new(T), true
	}
}
