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

	"gmod-addon-manager/config"
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
	config *config.Config
}

func NewManager(cfg *config.Config) *Manager {
	return &Manager{
		config: cfg,
	}
}

func (m *Manager) DownloadAddon(id string) error {
	// Create temp directory for this addon
	addonTempDir := filepath.Join(m.config.TempDir, id)
	if err := os.MkdirAll(addonTempDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Run steamcmd to download the addon
	steamCmd := exec.Command(
		m.config.SteamCmdPath,
		"+login", "anonymous",
		"+workshop_download_item", "4000", id,
		"+quit",
	)

	if err := steamCmd.Run(); err != nil {
		return fmt.Errorf("failed to run steamcmd: %w", err)
	}

	// Find the downloaded file
	downloadDir := filepath.Join(m.config.DownloadDir, id)
	files, err := os.ReadDir(downloadDir)
	if err != nil {
		return fmt.Errorf("failed to read download directory: %w", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no files found in download directory")
	}

	// Get the first file (should be either .gma or _legacy.bin)
	downloadedFile := files[0]
	srcPath := filepath.Join(downloadDir, downloadedFile.Name())
	dstPath := filepath.Join(addonTempDir, id+".gma")

	// Handle .bin file (extract and rename to .gma)
	if strings.HasSuffix(downloadedFile.Name(), "_legacy.bin") {
		// For now, we'll just rename it to .gma
		// In a real implementation, we might need to extract it
		dstPath = filepath.Join(addonTempDir, id+".gma")
		if err := os.Rename(srcPath, dstPath); err != nil {
			return fmt.Errorf("failed to rename .bin file: %w", err)
		}
	} else if strings.HasSuffix(downloadedFile.Name(), ".gma") {
		// Move the .gma file to temp directory
		if err := os.Rename(srcPath, dstPath); err != nil {
			return fmt.Errorf("failed to move .gma file: %w", err)
		}
	} else {
		return fmt.Errorf("unknown file type: %s", downloadedFile.Name())
	}

	// Execute GMAD tool to extract the addon
	gmadCmd := exec.Command(
		m.config.GMADPath,
		"extract",
		"-file", dstPath,
		"-out", addonTempDir,
	)

	if err := gmadCmd.Run(); err != nil {
		return fmt.Errorf("failed to run gmad: %w", err)
	}

	// Move the extracted addon to the out directory
	extractedDir := filepath.Join(addonTempDir, id)
	outDir := filepath.Join(m.config.OutDir, id)

	if err := os.Rename(extractedDir, outDir); err != nil {
		return fmt.Errorf("failed to move extracted addon: %w", err)
	}

	// Create symlink to enable the addon
	addonDir := filepath.Join(m.config.AddonDir, id)
	if err := os.Symlink(outDir, addonDir); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	// Clean up temp directory
	if err := os.RemoveAll(addonTempDir); err != nil {
		return fmt.Errorf("failed to clean up temp directory: %w", err)
	}

	// Clean up download directory
	if err := os.RemoveAll(downloadDir); err != nil {
		return fmt.Errorf("failed to clean up download directory: %w", err)
	}

	return nil
}

func (m *Manager) ListAddons() ([]Addon, error) {
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
		// Check if addon is enabled by checking if symlink exists in addons directory
		addonSymlink := filepath.Join(m.config.AddonDir, addonID)
		_, err := os.Lstat(addonSymlink)
		isEnabled := !os.IsNotExist(err)

		addon := Addon{
			ID:        addonID,
			Installed: true,
			Enabled:   isEnabled,
		}

		// Try to get more info from Steam Workshop
		workshopAddon, err := m.getWorkshopAddonInfo(addonID)
		if err != nil {
			// If we can't get info from Steam, just use what we have
			addons = append(addons, addon)
			continue
		}

		// Merge the workshop info with our addon
		if workshopAddon != nil {
			addon.Title = workshopAddon.Title
			addon.Author = workshopAddon.Creator
			addon.Description = workshopAddon.Description
			addon.Tags = workshopAddon.GetTagsAsStrings()
		}

		addons = append(addons, addon)
	}

	return addons, nil
}

func (m *Manager) GetAddonInfo(id string) (*Addon, error) {
	// First check if the addon is installed
	addonDir := filepath.Join(m.config.OutDir, id)
	if _, err := os.Stat(addonDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("addon not installed")
	}

	// Get info from Steam Workshop
	workshopAddon, err := m.getWorkshopAddonInfo(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get workshop info: %w", err)
	}

	if workshopAddon == nil {
		return &Addon{
			ID:        id,
			Installed: true,
			Enabled:   true,
		}, nil
	}

	return &Addon{
		ID:          id,
		Title:       workshopAddon.Title,
		Author:      workshopAddon.Creator,
		Description: workshopAddon.Description,
		Tags:        workshopAddon.GetTagsAsStrings(),
		Installed:   true,
		Enabled:     true,
	}, nil
}

// Helper function to get addon info from Steam Workshop
func (m *Manager) getWorkshopAddonInfo(id string) (*WorkshopAddon, error) {
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

	return &result.Response.PublishedFileDetails[0], nil
}

// Steam Workshop API response structures
type WorkshopResponse struct {
	Response struct {
		PublishedFileDetails []WorkshopAddon `json:"publishedfiledetails"`
	} `json:"response"`
}

type WorkshopAddon struct {
	PublishedFileID string   `json:"publishedfileid"`
	Title           string   `json:"title"`
	Creator         string   `json:"creator"`
	TimeCreated     int64    `json:"time_created"`
	TimeUpdated     int64    `json:"time_updated"`
	Views           int      `json:"views"`
	Subscriptions   int      `json:"subscriptions"`
	Favorited       int      `json:"favorited"`
	Tags            []Tag    `json:"tags"`
	Description     string   `json:"description"`
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
