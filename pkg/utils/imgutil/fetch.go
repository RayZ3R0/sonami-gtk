package imgutil

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
)

func (i *ImgUtil) fetch(url string) ([]byte, error) {
	hash := sha256.Sum256([]byte(url))
	key := base64.StdEncoding.EncodeToString(hash[:])
	if i.cache.Has(key) {
		return i.cache.Retrieve(key)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch image: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	i.cache.Store(key, data)
	return data, nil
}
