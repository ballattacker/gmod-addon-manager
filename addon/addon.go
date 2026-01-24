package addon

import (
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
	// Implementation will go here
	return nil
}

func (m *Manager) ListAddons() ([]Addon, error) {
	// Implementation will go here
	return nil, nil
}

func (m *Manager) GetAddonInfo(id string) (*Addon, error) {
	// Implementation will go here
	return nil, nil
}
