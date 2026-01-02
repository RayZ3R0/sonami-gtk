package secrets

import (
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/zalando/go-keyring"
)

const (
	service         = "org.codeberg.dergs.tidalwave"
	refreshTokenKey = "refresh_token"
)

func HasRefreshToken() bool {
	_, err := keyring.Get(service, refreshTokenKey)
	return err == nil
}

func GetRefreshToken() string {
	token, err := keyring.Get(service, refreshTokenKey)
	if err != nil {
		return ""
	}
	return token
}

func SetRefreshToken(token string) error {
	return keyring.Set(service, refreshTokenKey, token)
}

func DeleteRefreshToken() error {
	return keyring.Delete(service, refreshTokenKey)
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
