package cacheutil

import (
	"log/slog"
	"os"
	"path/filepath"
)

var cacheDir string

func GetCacheDir(appId string) string {
	if cacheDir == "" {
		d, err := os.UserCacheDir()
		if err != nil {
			cacheDir = os.TempDir()
		}
		cacheDir = d

		if err := os.MkdirAll(cacheDir, 0755); err != nil {
			slog.Error("error making cache directory", "cache_dir", cacheDir, "error", err)
		}
	}
	return filepath.Join(cacheDir, appId)
}

type Cache struct {
	path string
}

func (c *Cache) Path() string {
	return c.path
}

func (c *Cache) Has(key string) bool {
	stat, err := os.Stat(filepath.Join(c.path, key))
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	if stat.IsDir() {
		return false
	}
	return true
}

func (c *Cache) Retrieve(key string) ([]byte, error) {
	return os.ReadFile(filepath.Join(c.path, key))
}

func (c *Cache) Store(key string, data []byte) error {
	return os.WriteFile(filepath.Join(c.path, key), data, 0644)
}

func NewCache(appId string, subdir string) *Cache {
	cacheDir := GetCacheDir(appId)
	return &Cache{
		path: filepath.Join(cacheDir, subdir),
	}
}
