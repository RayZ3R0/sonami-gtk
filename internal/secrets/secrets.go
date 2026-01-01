package secrets

import "github.com/zalando/go-keyring"

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
