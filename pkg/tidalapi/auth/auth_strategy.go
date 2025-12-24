package auth

import (
	"net/http"
)

type AuthStrategy interface {
	Authenticate(req *http.Request) error
}
