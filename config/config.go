package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	GModDir        string
	DownloadDir    string
	AddonDir       string
	TempDir        string
	OutDir         string
	SteamCmdPath   string
	GMADPath       string
	SteamAPIKey    string
}

func NewDefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()

	return &Config{
		GModDir:      "C:\\Local\\Garrys Mod",
		DownloadDir:  filepath.Join(homeDir, "AppData", "Local", "Microsoft", "WinGet", "Packages", "Valve.SteamCMD_Microsoft.Winget.Source_8wekyb3d8bbwe", "steamapps", "workshop", "content", "4000"),
		AddonDir:     filepath.Join("C:\\Local\\Garrys Mod", "garrysmod", "addons"),
		TempDir:      filepath.Join("C:\\Local\\Garrys Mod", "garrysmod", "addons", "0", "tmp"),
		OutDir:       filepath.Join("C:\\Local\\Garrys Mod", "garrysmod", "addons", "0", "out"),
		SteamCmdPath: "steamcmd.exe",
		GMADPath:     filepath.Join("C:\\Local\\Garrys Mod", "bin", "gmad.exe"),
		SteamAPIKey:  "",
	}
}
