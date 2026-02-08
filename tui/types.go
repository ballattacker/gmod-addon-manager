package tui

// Message types for the TUI application

type errorMsg struct{ err error }
type successMsg struct{ msg string }

// View transition messages
type requestListViewMsg struct{}
type requestInputViewMsg struct{}
type requestDetailViewMsg struct{ addonID string }
type refreshListMsg struct{}

// Action messages
type enableAddonMsg struct{ addonID string }
type disableAddonMsg struct{ addonID string }
type refreshCacheMsg struct{ addonID string }
type removeAddonMsg struct{ addonID string }
type confirmInstallMsg struct{ addonID string }
type getAddonInfoMsg struct{ addonID string }
type confirmMsg struct{}
type infoMsg struct{}
