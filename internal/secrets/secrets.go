package secrets

import (
	"errors"

	"codeberg.org/dergs/tonearm/internal/g"
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

var getService = g.Lazy(func() Service {
	return newService()
})

type Service interface {
	Delete(key string) error
	Get(key string) (Item, error)
	Has(key string) (bool, error)
	Set(key string, value Item) error
}

type Item struct {
	Label    string
	Password string
}
