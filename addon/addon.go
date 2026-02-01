package addon

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"gmod-addon-manager/config"
	"gmod-addon-manager/file"
)

type Addon struct {
	ID          string
	Title       string
	Author      string
	Description string
	Tags        []string
	Installed   bool
	Enabled     bool
}

type Manager struct {
	config  *config.Config
	cache   *PersistentCache
	verbose bool
}

func NewManager(cfg *config.Config) (*Manager, error) {
	// Initialize persistent cache with 24-hour TTL
	cache, err := NewPersistentCache(24 * time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cache: %w", err)
	}

	return &Manager{
		config:  cfg,
		cache:   cache,
		verbose: true, // Default to verbose for CLI mode
	}, nil
}

func (m *Manager) SetVerbose(verbose bool) {
	m.verbose = verbose
}

func (m *Manager) log(message string) {
	if m.verbose {
		fmt.Println(message)
	}
}

func (m *Manager) GetAddon(id string) error {
	// Run steamcmd to get the addon with output
	steamCmd := exec.Command(
		m.config.SteamCmdPath,
		"+login", "anonymous",
		"+workshop_download_item", "4000", id,
		"+quit",
	)

	// Set up output pipes to capture and display SteamCMD output
	steamCmd.Stdout = os.Stdout
	steamCmd.Stderr = os.Stderr

	m.log(fmt.Sprintf("Downloading addon %s...", id))
	if err := steamCmd.Run(); err != nil {
		return fmt.Errorf("failed to run steamcmd: %w", err)
	}
	m.log("Download completed.")

	// Find the downloaded file
	downloadDir := filepath.Join(m.config.DownloadDir, id)

	// Get the first file (should be either .gma or _legacy.bin)
	downloadedFileName, err := file.First(downloadDir)
	if err != nil {
		return fmt.Errorf("failed to get the first file: %w", err)
	}
	downloadedFilePath := filepath.Join(downloadDir, downloadedFileName)

	// Create output directory
	outDir := filepath.Join(m.config.OutDir, id)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create tmp directory
	tmpDir := filepath.Join(m.config.TmpDir, id)
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return fmt.Errorf("failed to create tmp directory: %w", err)
	}

	gmaPath := filepath.Join(tmpDir, id+".gma")

	// Handle .bin file (extract and rename to .gma)
	if strings.HasSuffix(downloadedFileName, "_legacy.bin") {
		if err := file.ExtractLZMA(downloadedFilePath, gmaPath); err != nil {
			return fmt.Errorf("failed to extract bin file: %w", err)
		}
	} else if strings.HasSuffix(downloadedFileName, ".gma") {
		// gmaPath = downloadedFilePath
		if err := file.Copy(downloadedFilePath, gmaPath); err != nil {
			return fmt.Errorf("failed to copy gma file: %w", err)
		}
	} else {
		return fmt.Errorf("unknown file type: %s", downloadedFileName)
	}

	// Execute GMAD tool to extract directly to output directory
	gmadCmd := exec.Command(
		m.config.GMADPath,
		"extract",
		"-file", gmaPath,
		"-out", outDir,
	)

	m.log(fmt.Sprintf("Extracting addon %s...", id))
	if err := gmadCmd.Run(); err != nil {
		return fmt.Errorf("failed to run gmad: %w", err)
	}
	m.log("Extraction completed.")

	// Clean up tmp directory
	if err := os.RemoveAll(tmpDir); err != nil {
		return fmt.Errorf("failed to clean up tmp directory: %w", err)
	}

	m.EnableAddon(id)

	m.log(fmt.Sprintf("Addon %s installed and enabled successfully.", id))
	return nil
}

func (m *Manager) EnableAddon(id string) error {
	// Check if addon is installed
	addonDir := filepath.Join(m.config.OutDir, id)
	if _, err := os.Stat(addonDir); os.IsNotExist(err) {
		return fmt.Errorf("addon %s is not installed", id)
	}

	// Check if already enabled
	addonSymlink := filepath.Join(m.config.AddonDir, id)
	if _, err := os.Lstat(addonSymlink); err == nil {
		return fmt.Errorf("addon %s is already enabled", id)
	}

	// Create symlink to enable the addon
	if err := os.Symlink(addonDir, addonSymlink); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	m.log(fmt.Sprintf("Addon %s enabled successfully.", id))
	return nil
}

func (m *Manager) DisableAddon(id string) error {
	// Check if addon is installed
	addonDir := filepath.Join(m.config.OutDir, id)
	if _, err := os.Stat(addonDir); os.IsNotExist(err) {
		return fmt.Errorf("addon %s is not installed", id)
	}

	// Check if already disabled
	addonSymlink := filepath.Join(m.config.AddonDir, id)
	if _, err := os.Lstat(addonSymlink); os.IsNotExist(err) {
		return fmt.Errorf("addon %s is already disabled", id)
	}

	// Remove symlink to disable the addon
	if err := os.Remove(addonSymlink); err != nil {
		return fmt.Errorf("failed to remove symlink: %w", err)
	}

	m.log(fmt.Sprintf("Addon %s disabled successfully.", id))
	return nil
}

func (m *Manager) RemoveAddon(id string) error {
	// Check if addon is installed
	addonDir := filepath.Join(m.config.OutDir, id)
	if _, err := os.Stat(addonDir); os.IsNotExist(err) {
		return fmt.Errorf("addon %s is not installed", id)
	}

	// First disable the addon if it's enabled
	addonSymlink := filepath.Join(m.config.AddonDir, id)
	if _, err := os.Lstat(addonSymlink); err == nil {
		if err := os.Remove(addonSymlink); err != nil {
			return fmt.Errorf("failed to remove symlink: %w", err)
		}
	}

	// Remove the addon directory
	if err := os.RemoveAll(addonDir); err != nil {
		return fmt.Errorf("failed to remove addon directory: %w", err)
	}

	// Clear the cache for this addon
	if err := m.RefreshCache(id); err != nil {
		return fmt.Errorf("failed to clear cache: %w", err)
	}

	m.log(fmt.Sprintf("Addon %s removed successfully.", id))
	return nil
}

func (m *Manager) RefreshCache(id string) error {
	// Clear the cache for this addon
	cacheFile := m.cache.cacheFilePath(id)
	if err := os.Remove(cacheFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to clear cache: %w", err)
	}

	// Refresh the addon info by fetching it again
	_, err := m.GetAddonInfo(id)
	if err != nil {
		return fmt.Errorf("failed to refresh addon info: %w", err)
	}

	m.log(fmt.Sprintf("Cache refreshed for addon %s.", id))
	return nil
}

func (m *Manager) GetAddonsInfo() ([]Addon, error) {
	var addons []Addon

	// Read the out directory to find installed addons
	entries, err := os.ReadDir(m.config.OutDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read out directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		addonID := entry.Name()
		// Get addon info using the existing GetAddonInfo method
		addonInfo, err := m.GetAddonInfo(addonID)
		if err != nil {
			// Create addon with empty/default fields when we can't get info
			addon := Addon{
				ID:          addonID,
				Installed:   true,
				Enabled:     false,
				Title:       "",
				Author:      "",
				Description: "",
				Tags:        []string{},
			}
			addons = append(addons, addon)
			continue
		}

		// Only include installed addons in the list
		if addonInfo.Installed {
			addons = append(addons, *addonInfo)
		}
	}

	return addons, nil
}

func (m *Manager) GetAddonInfo(id string) (*Addon, error) {
	// Check if addon is installed
	addonDir := filepath.Join(m.config.OutDir, id)
	isInstalled := true
	if _, err := os.Stat(addonDir); os.IsNotExist(err) {
		isInstalled = false
	}

	// Check if addon is enabled (only if installed)
	isEnabled := false
	if isInstalled {
		addonSymlink := filepath.Join(m.config.AddonDir, id)
		_, err := os.Lstat(addonSymlink)
		isEnabled = !os.IsNotExist(err)
	}

	// Create base addon with local info
	addon := &Addon{
		ID:        id,
		Installed: isInstalled,
		Enabled:   isEnabled,
	}

	// Try to get more info from Steam Workshop
	workshopAddon, err := m.getWorkshopAddonInfo(id)
	if err != nil {
		return addon, nil
	}

	// Merge the workshop info with our addon
	if workshopAddon != nil {
		addon.Title = workshopAddon.Title
		addon.Author = workshopAddon.Creator
		addon.Description = workshopAddon.Description
		addon.Tags = workshopAddon.GetTagsAsStrings()
	}

	return addon, nil
}

// Helper function to get addon info from Steam Workshop with caching
func (m *Manager) getWorkshopAddonInfo(id string) (*WorkshopAddon, error) {
	// Check cache first
	if cachedAddon, found, err := m.cache.Get(id); err != nil {
		return nil, fmt.Errorf("cache error: %w", err)
	} else if found {
		return cachedAddon, nil
	}

	apiURL := "https://api.steampowered.com/ISteamRemoteStorage/GetPublishedFileDetails/v1/"

	var requestBody string
	if m.config.SteamAPIKey != "" {
		requestBody = fmt.Sprintf("itemcount=1&publishedfileids[0]=%s&key=%s", id, m.config.SteamAPIKey)
	} else {
		requestBody = fmt.Sprintf("itemcount=1&publishedfileids[0]=%s", id)
	}

	resp, err := http.Post(apiURL, "application/x-www-form-urlencoded", strings.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to make API request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result WorkshopResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(result.Response.PublishedFileDetails) == 0 {
		return nil, nil
	}

	workshopAddon := &result.Response.PublishedFileDetails[0]

	// Cache the result
	if err := m.cache.Set(id, workshopAddon); err != nil {
		return nil, fmt.Errorf("failed to cache workshop addon: %w", err)
	}

	return workshopAddon, nil
}

// Steam Workshop API response structures
type WorkshopResponse struct {
	Response struct {
		PublishedFileDetails []WorkshopAddon `json:"publishedfiledetails"`
	} `json:"response"`
}

type WorkshopAddon struct {
	PublishedFileID string `json:"publishedfileid"`
	Title           string `json:"title"`
	Creator         string `json:"creator"`
	TimeCreated     int64  `json:"time_created"`
	TimeUpdated     int64  `json:"time_updated"`
	Views           int    `json:"views"`
	Subscriptions   int    `json:"subscriptions"`
	Favorited       int    `json:"favorited"`
	Tags            []Tag  `json:"tags"`
	Description     string `json:"description"`
}

type Tag struct {
	Tag string `json:"tag"`
}

// Method to convert []Tag to []string
func (w *WorkshopAddon) GetTagsAsStrings() []string {
	tags := make([]string, len(w.Tags))
	for i, tag := range w.Tags {
		tags[i] = tag.Tag
	}
	return tags
}
