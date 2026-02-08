package tui

// Message types for the TUI application

type errorMsg struct{ err error }
type successMsg struct{ msg string }
type cancelMsg struct{}

// View transition messages
type requestListViewMsg struct{}
type requestInputViewMsg struct{}
type requestDetailViewMsg struct{ addonID string }

// Action messages
type enableAddonMsg struct{ addonID string }
type disableAddonMsg struct{ addonID string }
type reloadAddonMsg struct{ addonID string }
type installAddonMsg struct{ addonID string }
type removeAddonMsg struct{ addonID string }
