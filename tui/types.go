package tui

import (
	"gmod-addon-manager/addon"

	"github.com/charmbracelet/bubbles/list"
)

// Message types for the TUI application
type errorMsg struct{ err error }
type successMsg struct{ msg string }
type refreshListMsg struct{ items []list.Item }
type selectAddonMsg struct{ addon *addon.Addon }
type enableAddonMsg struct{ id string }
type disableAddonMsg struct{ id string }
type refreshCacheMsg struct{ id string }
type removeAddonMsg struct{ id string }

// requestViewMsg tells main model to switch to a different view
type requestViewMsg struct {
	view string // "list", "input", "detail"
}

