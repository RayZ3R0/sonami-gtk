package my_collection

import "sync"

// fetchAll fetches all items for the given ids concurrently (max 8 in-flight).
// Results preserve the original order of ids; items that fail to load are silently skipped.
func fetchAll[T any](ids []string, fetch func(string) (T, error)) []T {
	type indexed struct {
		i   int
		val T
	}

	results := make([]indexed, 0, len(ids))
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 8)

	for i, id := range ids {
		i, id := i, id
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			val, err := fetch(id)
			if err != nil {
				return
			}
			mu.Lock()
			results = append(results, indexed{i, val})
			mu.Unlock()
		}()
	}
	wg.Wait()

	// Reconstruct in original order (newest-first as returned by the DB).
	slots := make([]T, len(ids))
	filled := make([]bool, len(ids))
	for _, r := range results {
		slots[r.i] = r.val
		filled[r.i] = true
	}
	out := make([]T, 0, len(ids))
	for i, v := range slots {
		if filled[i] {
			out = append(out, v)
		}
	}
	return out
}
