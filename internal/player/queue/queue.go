package queue

import (
	"codeberg.org/dergs/tonearm/internal/signals"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/models/openapi"
)

// Queue represents a queue that holds tracks temporarily.
// Once the tracks are played, they are removed from the queue.
type Queue interface {
	// Add a track to the end of the queue
	Append(track *openapi.Track)
	// Remove all tracks from the queue
	Clear()
	// Check if a track with the given ID is in the queue
	Contains(trackID string) bool
	// Entries returns a stateful signal that emits the current queue entries.
	Entries() *signals.StatefulSignal[[]*openapi.Track]
	// Get track at the specified index
	Get(index int) *openapi.Track
	// Insert a track at a specific index in the queue
	Insert(track *openapi.Track, index int) error
	// Peek at the next track in the queue without removing it
	Peek() *openapi.Track
	// Pop the next track from the queue
	Pop() *openapi.Track
	// Prepend a track to the beginning of the queue
	Prepend(track *openapi.Track)
	// Remove a track from the queue
	RemoveAt(index int) error
	// Replace the current queue with a new set of tracks
	Set(tracks []*openapi.Track)
	// Skip a specificed number of tracks in the queue
	Skip(n int) ([]*openapi.Track, error)
}

// DurableQueue represents a queue that holds tracks persistently.
// Once the tracks are played, they are removed from the queue,
// but the original queue is preserved for future use like looping or shuffling.
type DurableQueue interface {
	Queue

	// Restore all original tracks to the queue.
	// All tracks before and including the current track are not re-added.
	Restore(currentTrackID string)

	// Re-insert all track from the original queue and randomize their order.
	// If the current track is found in the original queue, it not be re-added
	Shuffle(currentTrackID string)
}
