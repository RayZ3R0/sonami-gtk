package tidalapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	cacheTTL          = 24 * time.Hour
	rateLimitSleepMs  = 2000
	maxStickyFailures = 3
	maxRetries        = 3
)

var endpointLogger = slog.With("module", "endpoint_manager")

// endpointCache is the JSON structure stored locally
type endpointCache struct {
	Endpoints []string  `json:"endpoints"`
	FetchedAt time.Time `json:"fetched_at"`
}

// endpointResponse is the JSON response from the instances endpoint.
// Format: { "api": { "<provider-name>": { "urls": [...], "cors": bool }, ... } }
type endpointResponse struct {
	API map[string]struct {
		URLs []string `json:"urls"`
	} `json:"api"`
}

// EndpointManager manages a rotating list of API endpoints with sticky preference,
// automatic failover, and rate-limit handling.
type EndpointManager struct {
	mu             sync.RWMutex
	endpoints      []string
	stickyIndex    int
	stickyFailures int
	endpointsURL   string
	cacheFilePath  string
	isInitialized  bool
}

// NewEndpointManager creates a new EndpointManager.
// The endpointsURL must be provided by the user (via Preferences → Streaming).
// Call Initialize() to prime it.
func NewEndpointManager(endpointsURL string) *EndpointManager {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = filepath.Join(os.Getenv("HOME"), ".config")
	}

	cacheDir := filepath.Join(configDir, "sonami")
	os.MkdirAll(cacheDir, 0755)

	return &EndpointManager{
		endpointsURL:  endpointsURL,
		cacheFilePath: filepath.Join(cacheDir, "tidal_cache.json"),
		stickyIndex:   -1,
	}
}

// Initialize primes the endpoint manager by loading cached or fetching fresh endpoints.
// This should be called at application startup to ensure smooth first playback.
func (em *EndpointManager) Initialize(ctx context.Context) error {
	em.mu.Lock()
	defer em.mu.Unlock()

	if em.endpointsURL == "" {
		return fmt.Errorf("no streaming instances URL configured — set one in Preferences → Streaming")
	}

	// Try loading from cache first
	if em.loadCache() {
		endpointLogger.Info("loaded endpoints from cache", "count", len(em.endpoints))
		em.isInitialized = true

		// Validate the first endpoint in the background
		go em.validateStickyEndpoint()
		return nil
	}

	// Fetch fresh
	if err := em.fetchEndpoints(ctx); err != nil {
		return fmt.Errorf("failed to initialize endpoints: %w", err)
	}

	endpointLogger.Info("fetched fresh endpoints", "count", len(em.endpoints))
	em.isInitialized = true
	return nil
}

// GetEndpoint returns the current best endpoint URL.
func (em *EndpointManager) GetEndpoint() (string, error) {
	em.mu.RLock()
	defer em.mu.RUnlock()

	if len(em.endpoints) == 0 {
		return "", fmt.Errorf("no endpoints available")
	}

	if em.stickyIndex >= 0 && em.stickyIndex < len(em.endpoints) {
		return em.endpoints[em.stickyIndex], nil
	}

	return em.endpoints[0], nil
}

// ReportSuccess marks the current endpoint as successful, making it sticky.
func (em *EndpointManager) ReportSuccess(endpoint string) {
	em.mu.Lock()
	defer em.mu.Unlock()

	for i, ep := range em.endpoints {
		if ep == endpoint {
			em.stickyIndex = i
			em.stickyFailures = 0
			return
		}
	}
}

// ReportFailure marks the current endpoint as failed. After maxStickyFailures,
// it rotates to the next endpoint.
func (em *EndpointManager) ReportFailure(endpoint string) {
	em.mu.Lock()
	defer em.mu.Unlock()

	for i, ep := range em.endpoints {
		if ep == endpoint && i == em.stickyIndex {
			em.stickyFailures++
			if em.stickyFailures >= maxStickyFailures {
				endpointLogger.Warn("sticky endpoint exceeded max failures, rotating",
					"endpoint", endpoint, "failures", em.stickyFailures)
				em.rotateEndpoint()
			}
			return
		}
	}
}

// DoRequest performs an HTTP GET request against the endpoints with automatic
// failover and rate-limit handling.
func (em *EndpointManager) DoRequest(ctx context.Context, path string) (*http.Response, error) {
	if !em.isInitialized {
		if err := em.Initialize(ctx); err != nil {
			return nil, err
		}
	}

	em.mu.RLock()
	endpoints := make([]string, len(em.endpoints))
	copy(endpoints, em.endpoints)
	startIdx := em.stickyIndex
	em.mu.RUnlock()

	if len(endpoints) == 0 {
		return nil, fmt.Errorf("no endpoints available")
	}

	// Reorder to try sticky endpoint first
	ordered := make([]string, 0, len(endpoints))
	if startIdx >= 0 && startIdx < len(endpoints) {
		ordered = append(ordered, endpoints[startIdx])
		for i, ep := range endpoints {
			if i != startIdx {
				ordered = append(ordered, ep)
			}
		}
	} else {
		ordered = endpoints
	}

	var lastErr error
	for _, endpoint := range ordered {
		for retry := 0; retry < maxRetries; retry++ {
			url := endpoint + path
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if err != nil {
				lastErr = err
				continue
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				lastErr = err
				em.ReportFailure(endpoint)
				break // Try next endpoint
			}

			if resp.StatusCode == http.StatusTooManyRequests {
				resp.Body.Close()
				endpointLogger.Debug("rate limited, sleeping", "endpoint", endpoint)
				time.Sleep(time.Duration(rateLimitSleepMs) * time.Millisecond)
				continue // Retry same endpoint
			}

			if resp.StatusCode != http.StatusOK {
				resp.Body.Close()
				lastErr = fmt.Errorf("unexpected status %d from %s", resp.StatusCode, endpoint)
				em.ReportFailure(endpoint)
				break // Try next endpoint
			}

			em.ReportSuccess(endpoint)
			return resp, nil
		}
	}

	return nil, fmt.Errorf("all endpoints failed: %w", lastErr)
}

func (em *EndpointManager) loadCache() bool {
	data, err := os.ReadFile(em.cacheFilePath)
	if err != nil {
		return false
	}

	var cache endpointCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return false
	}

	if time.Since(cache.FetchedAt) > cacheTTL {
		endpointLogger.Debug("cache expired", "fetched_at", cache.FetchedAt)
		return false
	}

	if len(cache.Endpoints) == 0 {
		return false
	}

	em.endpoints = cache.Endpoints
	// Start with a random endpoint for load distribution
	em.stickyIndex = rand.IntN(len(em.endpoints))
	return true
}

func (em *EndpointManager) saveCache() {
	cache := endpointCache{
		Endpoints: em.endpoints,
		FetchedAt: time.Now(),
	}

	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		endpointLogger.Error("failed to marshal endpoint cache", "error", err)
		return
	}

	if err := os.WriteFile(em.cacheFilePath, data, 0644); err != nil {
		endpointLogger.Error("failed to write endpoint cache", "error", err)
	}
}

func (em *EndpointManager) fetchEndpoints(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, em.endpointsURL, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch endpoints: HTTP %d", resp.StatusCode)
	}

	var endpointResp endpointResponse
	if err := json.NewDecoder(resp.Body).Decode(&endpointResp); err != nil {
		return err
	}

	// Collect URLs from ALL provider groups into one flat pool
	var allURLs []string
	for providerName, provider := range endpointResp.API {
		endpointLogger.Debug("loaded provider", "name", providerName, "count", len(provider.URLs))
		allURLs = append(allURLs, provider.URLs...)
	}

	if len(allURLs) == 0 {
		return fmt.Errorf("no endpoints found in response")
	}

	em.endpoints = allURLs
	em.stickyIndex = rand.IntN(len(em.endpoints))
	em.stickyFailures = 0

	em.saveCache()
	return nil
}

func (em *EndpointManager) rotateEndpoint() {
	if len(em.endpoints) <= 1 {
		return
	}

	em.stickyIndex = (em.stickyIndex + 1) % len(em.endpoints)
	em.stickyFailures = 0
	endpointLogger.Info("rotated to new endpoint", "index", em.stickyIndex, "endpoint", em.endpoints[em.stickyIndex])
}

func (em *EndpointManager) validateStickyEndpoint() {
	em.mu.RLock()
	if len(em.endpoints) == 0 || em.stickyIndex < 0 {
		em.mu.RUnlock()
		return
	}
	endpoints := make([]string, len(em.endpoints))
	copy(endpoints, em.endpoints)
	startIdx := em.stickyIndex
	em.mu.RUnlock()

	// Try validating endpoints until we find a reachable one
	for attempt := 0; attempt < len(endpoints); attempt++ {
		idx := (startIdx + attempt) % len(endpoints)
		endpoint := endpoints[idx]

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint+"/search/?s=test", nil)
		if err != nil {
			cancel()
			continue
		}

		resp, err := http.DefaultClient.Do(req)
		cancel()
		if err != nil {
			endpointLogger.Debug("startup validation: endpoint unreachable, trying next", "endpoint", endpoint)
			em.ReportFailure(endpoint)
			continue
		}
		resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			endpointLogger.Info("startup endpoint validation successful", "endpoint", endpoint)
			em.ReportSuccess(endpoint)
			return
		}
		em.ReportFailure(endpoint)
	}

	endpointLogger.Warn("startup validation: no reachable endpoint found after trying multiple, will retry on first request")
}
