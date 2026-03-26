package cache

import (
	"testing"
)

func TestLRUBasic(t *testing.T) {
	cache := NewLRU[string, int](3)

	// Test set and get
	cache.Set("a", 1)
	cache.Set("b", 2)
	cache.Set("c", 3)

	if val, ok := cache.Get("a"); !ok || val != 1 {
		t.Errorf("expected 1, got %d (ok=%v)", val, ok)
	}
	if val, ok := cache.Get("b"); !ok || val != 2 {
		t.Errorf("expected 2, got %d (ok=%v)", val, ok)
	}
	if val, ok := cache.Get("c"); !ok || val != 3 {
		t.Errorf("expected 3, got %d (ok=%v)", val, ok)
	}

	// Test miss
	if _, ok := cache.Get("notfound"); ok {
		t.Error("expected cache miss for 'notfound'")
	}
}

func TestLRUEviction(t *testing.T) {
	cache := NewLRU[string, int](2)

	cache.Set("a", 1)
	cache.Set("b", 2)
	cache.Set("c", 3) // Should evict "a"

	if _, ok := cache.Get("a"); ok {
		t.Error("expected 'a' to be evicted")
	}
	if _, ok := cache.Get("b"); !ok {
		t.Error("expected 'b' to still be cached")
	}
	if _, ok := cache.Get("c"); !ok {
		t.Error("expected 'c' to still be cached")
	}
}

func TestLRUUpdate(t *testing.T) {
	cache := NewLRU[string, int](2)

	cache.Set("a", 1)
	cache.Set("b", 2)

	// Update a (moves to front)
	cache.Set("a", 10)

	cache.Set("c", 3) // Should evict "b" (oldest)

	if val, ok := cache.Get("a"); !ok || val != 10 {
		t.Errorf("expected a=10, got %d (ok=%v)", val, ok)
	}
	if _, ok := cache.Get("b"); ok {
		t.Error("expected 'b' to be evicted")
	}
	if _, ok := cache.Get("c"); !ok {
		t.Error("expected 'c' to still be cached")
	}
}

func TestLRUDelete(t *testing.T) {
	cache := NewLRU[string, int](3)

	cache.Set("a", 1)
	cache.Set("b", 2)

	cache.Delete("a")

	if _, ok := cache.Get("a"); ok {
		t.Error("expected 'a' to be deleted")
	}
	if _, ok := cache.Get("b"); !ok {
		t.Error("expected 'b' to still be cached")
	}
}

func TestLRUClear(t *testing.T) {
	cache := NewLRU[string, int](3)

	cache.Set("a", 1)
	cache.Set("b", 2)
	cache.Set("c", 3)

	cache.Clear()

	if cache.Len() != 0 {
		t.Errorf("expected len=0 after clear, got %d", cache.Len())
	}
	if _, ok := cache.Get("a"); ok {
		t.Error("expected 'a' to be cleared")
	}
}

func TestLRULen(t *testing.T) {
	cache := NewLRU[string, int](5)

	if cache.Len() != 0 {
		t.Errorf("expected empty cache, got len=%d", cache.Len())
	}

	cache.Set("a", 1)
	cache.Set("b", 2)

	if cache.Len() != 2 {
		t.Errorf("expected len=2, got %d", cache.Len())
	}

	cache.Delete("a")

	if cache.Len() != 1 {
		t.Errorf("expected len=1, got %d", cache.Len())
	}
}
