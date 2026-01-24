package addon

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type CacheEntry struct {
	WorkshopAddon *WorkshopAddon
	Timestamp     time.Time
}

type PersistentCache struct {
	cacheDir string
	ttl      time.Duration
}

func NewPersistentCache(ttl time.Duration) (*PersistentCache, error) {
	// Get user's cache directory
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get cache directory: %w", err)
	}

	// Create app-specific cache directory
	appCacheDir := filepath.Join(cacheDir, "gmod-addon-manager")
	if err := os.MkdirAll(appCacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	return &PersistentCache{
		cacheDir: appCacheDir,
		ttl:      ttl,
	}, nil
}

func (c *PersistentCache) cacheFilePath(id string) string {
	return filepath.Join(c.cacheDir, fmt.Sprintf("%s.json", id))
}

func (c *PersistentCache) Get(id string) (*WorkshopAddon, bool, error) {
	cacheFile := c.cacheFilePath(id)

	// Check if cache file exists
	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		return nil, false, nil
	}

	// Read cache file
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil, false, fmt.Errorf("failed to read cache file: %w", err)
	}

	// Parse cache entry
	var entry CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, false, fmt.Errorf("failed to parse cache entry: %w", err)
	}

	// Check if entry has expired
	if time.Since(entry.Timestamp) > c.ttl {
		// Remove expired cache file
		if err := os.Remove(cacheFile); err != nil {
			return nil, false, fmt.Errorf("failed to remove expired cache file: %w", err)
		}
		return nil, false, nil
	}

	return entry.WorkshopAddon, true, nil
}

func (c *PersistentCache) Set(id string, workshopAddon *WorkshopAddon) error {
	entry := CacheEntry{
		WorkshopAddon: workshopAddon,
		Timestamp:     time.Now(),
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal cache entry: %w", err)
	}

	cacheFile := c.cacheFilePath(id)
	if err := os.WriteFile(cacheFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	return nil
}
