package cache

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/RayZ3R0/sonami-gtk/internal/settings"
	"github.com/RayZ3R0/sonami-gtk/pkg/utils/imgutil"
)

// ImagePrefetcher prefetches cover images in the background with rate limiting
type ImagePrefetcher struct {
	imgUtil     *imgutil.ImgUtil
	rateLimiter *time.Ticker
	queue       chan string
	inFlight    sync.Map // map[string]bool - tracks in-progress fetches
	logger      *slog.Logger
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
}

// NewImagePrefetcher creates a new image prefetcher
func NewImagePrefetcher(imgUtil *imgutil.ImgUtil) *ImagePrefetcher {
	ctx, cancel := context.WithCancel(context.Background())

	return &ImagePrefetcher{
		imgUtil:     imgUtil,
		rateLimiter: time.NewTicker(100 * time.Millisecond), // 10 images per second
		queue:       make(chan string, 1000),                // Buffer up to 1000 URLs
		logger:      slog.With("module", "image-prefetcher"),
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Start begins the background prefetch worker
func (p *ImagePrefetcher) Start() {
	p.wg.Add(1)
	go p.worker()
	p.logger.Info("image prefetcher started")
}

// Stop gracefully shuts down the prefetcher
func (p *ImagePrefetcher) Stop() {
	p.cancel()
	p.wg.Wait()
	p.rateLimiter.Stop()
	p.logger.Info("image prefetcher stopped")
}

// Prefetch adds URLs to the prefetch queue (non-blocking)
// Skips URLs that are already in-flight or empty
func (p *ImagePrefetcher) Prefetch(urls []string) {
	if !settings.Performance().ShouldCacheImages() {
		p.logger.Debug("image caching disabled, skipping prefetch")
		return
	}

	added := 0
	for _, url := range urls {
		if url == "" {
			continue
		}

		// Skip if already in-flight
		if _, exists := p.inFlight.LoadOrStore(url, true); exists {
			continue
		}

		// Try to add to queue (non-blocking)
		select {
		case p.queue <- url:
			added++
		default:
			// Queue is full, skip this URL
			p.inFlight.Delete(url)
			p.logger.Warn("prefetch queue full, dropping URL", "url", url)
		}
	}

	if added > 0 {
		p.logger.Debug("queued images for prefetch", "count", added)
	}
}

// worker runs in the background and processes the prefetch queue
func (p *ImagePrefetcher) worker() {
	defer p.wg.Done()

	for {
		select {
		case <-p.ctx.Done():
			p.logger.Info("worker stopping", "queued", len(p.queue))
			return

		case <-p.rateLimiter.C:
			// Rate limit tick - try to fetch one image
			select {
			case url := <-p.queue:
				p.prefetchOne(url)
			default:
				// Queue is empty, nothing to do
			}
		}
	}
}

// prefetchOne fetches a single image (called by worker with rate limiting)
func (p *ImagePrefetcher) prefetchOne(url string) {
	defer p.inFlight.Delete(url)

	// Check if caching is still enabled
	if !settings.Performance().ShouldCacheImages() {
		return
	}

	// Use ImgUtil.Load which handles caching internally
	_, err := p.imgUtil.Load(url)
	if err != nil {
		p.logger.Debug("failed to prefetch image", "url", url, "error", err)
		return
	}

	p.logger.Debug("prefetched image", "url", url)
}

// Stats returns current prefetcher statistics
func (p *ImagePrefetcher) Stats() map[string]int {
	inFlightCount := 0
	p.inFlight.Range(func(key, value interface{}) bool {
		inFlightCount++
		return true
	})

	return map[string]int{
		"queued":    len(p.queue),
		"in_flight": inFlightCount,
	}
}
