package tui

import (
	"gmod-addon-manager/addon"
)

// Message types for the TUI application

type errorMsg struct{ err error }
type successMsg struct{ msg string }

// View transition messages
type requestListViewMsg struct{}
type requestInputViewMsg struct{}
type requestDetailViewMsg struct{ addon *addon.Addon }
