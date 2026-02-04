package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	GModDir      string `json:"gmod_dir"`
	DownloadDir  string `json:"download_dir"`
	AddonDir     string `json:"addon_dir"`
	OutDir       string `json:"out_dir"`
	TmpDir       string `json:"tmp_dir"`
	SteamCmdPath string `json:"steamcmd_path"`
	GMADPath     string `json:"gmad_path"`
	SteamAPIKey  string `json:"steam_api_key"`
}

const ConfigFileName = "gmod-addon-manager.json"

func NewDefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()

	return &Config{
		GModDir:      "C:\\Games\\GarrysMod",
		DownloadDir:  filepath.Join(homeDir, "AppData", "Local", "Microsoft", "WinGet", "Packages", "Valve.SteamCMD_Microsoft.Winget.Source_8wekyb3d8bbwe", "steamapps", "workshop", "content", "4000"),
		AddonDir:     "",
		OutDir:       "",
		TmpDir:       "",
		SteamCmdPath: "steamcmd.exe",
		GMADPath:     "",
		SteamAPIKey:  "",
	}
}

func LoadConfig() (*Config, error) {
	// Get config file path
	configPath, err := getConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get config path: %w", err)
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config and save it
		config := NewDefaultConfig()
		if err := SaveConfig(config); err != nil {
			return nil, fmt.Errorf("failed to save default config: %w", err)
		}
		return fillInDefaultPaths(config), nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse config
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Fill in any missing paths based on GModDir
	return fillInDefaultPaths(&config), nil
}

func fillInDefaultPaths(config *Config) *Config {
	// If GModDir is empty, use default
	if config.GModDir == "" {
		config.GModDir = "C:\\Games\\GarrysMod"
	}

	// Fill in AddonDir if empty
	if config.AddonDir == "" {
		config.AddonDir = filepath.Join(config.GModDir, "garrysmod", "addons")
	}

	// Fill in OutDir if empty
	if config.OutDir == "" {
		config.OutDir = filepath.Join(config.AddonDir, "0", "out")
	}

	// Fill in TmpDir if empty
	if config.TmpDir == "" {
		config.TmpDir = filepath.Join(config.AddonDir, "0", "tmp")
	}

	// Fill in GMADPath if empty
	if config.GMADPath == "" {
		config.GMADPath = filepath.Join(config.GModDir, "bin", "gmad.exe")
	}

	// Ensure DownloadDir is set (this one can't be deduced from GModDir)
	if config.DownloadDir == "" {
		homeDir, _ := os.UserHomeDir()
		config.DownloadDir = filepath.Join(homeDir, "AppData", "Local", "Microsoft", "WinGet", "Packages", "Valve.SteamCMD_Microsoft.Winget.Source_8wekyb3d8bbwe", "steamapps", "workshop", "content", "4000")
	}

	return config
}

func SaveConfig(config *Config) error {
	// Get config file path
	configPath, err := getConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal config to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write config file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func getConfigPath() (string, error) {
	// Get user's config directory
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get config directory: %w", err)
	}

	// Create app-specific config directory
	appConfigDir := filepath.Join(configDir, "gmod-addon-manager")
	return filepath.Join(appConfigDir, ConfigFileName), nil
}

// GetConfigPath returns the path to the config file
// This is a public version of getConfigPath for external use
func GetConfigPath() (string, error) {
	return getConfigPath()
}
