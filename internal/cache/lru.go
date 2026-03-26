package cache

import (
	"container/list"
	"sync"
)

// LRUCache is a thread-safe Least Recently Used cache
// with fixed capacity and O(1) get/set operations
type LRUCache[K comparable, V any] struct {
	capacity  int
	items     map[K]*list.Element
	evictList *list.List
	mu        sync.RWMutex
}

type entry[K comparable, V any] struct {
	key   K
	value V
}

// NewLRU creates a new LRU cache with the given capacity
func NewLRU[K comparable, V any](capacity int) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		capacity:  capacity,
		items:     make(map[K]*list.Element, capacity),
		evictList: list.New(),
	}
}

// Get retrieves a value from the cache
// Returns (value, true) if found, (zero, false) if not found
func (c *LRUCache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	elem, ok := c.items[key]
	c.mu.RUnlock()

	if !ok {
		var zero V
		return zero, false
	}

	// Move to front (most recently used)
	c.mu.Lock()
	c.evictList.MoveToFront(elem)
	c.mu.Unlock()

	return elem.Value.(*entry[K, V]).value, true
}

// Set adds or updates a value in the cache
func (c *LRUCache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Update existing entry
	if elem, ok := c.items[key]; ok {
		c.evictList.MoveToFront(elem)
		elem.Value.(*entry[K, V]).value = value
		return
	}

	// Add new entry
	elem := c.evictList.PushFront(&entry[K, V]{key: key, value: value})
	c.items[key] = elem

	// Evict oldest if over capacity
	if c.evictList.Len() > c.capacity {
		c.evictOldest()
	}
}

// Delete removes a key from the cache
func (c *LRUCache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		c.removeElement(elem)
	}
}

// Clear removes all entries from the cache
func (c *LRUCache[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[K]*list.Element, c.capacity)
	c.evictList.Init()
}

// Len returns the current number of cached items
func (c *LRUCache[K, V]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.evictList.Len()
}

// evictOldest removes the least recently used item
// Must be called with lock held
func (c *LRUCache[K, V]) evictOldest() {
	elem := c.evictList.Back()
	if elem != nil {
		c.removeElement(elem)
	}
}

// removeElement removes an element from the cache
// Must be called with lock held
func (c *LRUCache[K, V]) removeElement(elem *list.Element) {
	c.evictList.Remove(elem)
	kv := elem.Value.(*entry[K, V])
	delete(c.items, kv.key)
}
