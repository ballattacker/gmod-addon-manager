package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// KeyContext holds context information for key action execution
type KeyContext struct {
	AddonID string
}

// KeyMapEntry is a keybinding with its associated action
type KeyMapEntry struct {
	Binding key.Binding
	Action  func(*KeyContext) tea.Msg
}

// KeyMap is the master registry of all keybindings and their actions
type KeyMap struct {
	Install      KeyMapEntry
	Refresh      KeyMapEntry
	Quit         KeyMapEntry
	View         KeyMapEntry
	Enable       KeyMapEntry
	Disable      KeyMapEntry
	RefreshCache KeyMapEntry
	Remove       KeyMapEntry
	Confirm      KeyMapEntry
	Info         KeyMapEntry
	Cancel       KeyMapEntry
}

// GlobalKeyMap is the single master keymap with all keybindings and actions
var GlobalKeyMap = KeyMap{
	Install: KeyMapEntry{
		Binding: key.NewBinding(
			key.WithKeys("i"),
			key.WithHelp("i", "install"),
		),
		Action: func(ctx *KeyContext) tea.Msg {
			return requestInputViewMsg{}
		},
	},
	Refresh: KeyMapEntry{
		Binding: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),
		Action: func(ctx *KeyContext) tea.Msg {
			return refreshListMsg{}
		},
	},
	Quit: KeyMapEntry{
		Binding: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Action: func(ctx *KeyContext) tea.Msg {
			return tea.Quit()
		},
	},
	View: KeyMapEntry{
		Binding: key.NewBinding(
			key.WithKeys("enter", "v"),
			key.WithHelp("enter", "view"),
		),
		Action: func(ctx *KeyContext) tea.Msg {
			return requestDetailViewMsg{addonID: ctx.AddonID}
		},
	},
	Enable: KeyMapEntry{
		Binding: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "enable"),
		),
		Action: func(ctx *KeyContext) tea.Msg {
			return enableAddonMsg{addonID: ctx.AddonID}
		},
	},
	Disable: KeyMapEntry{
		Binding: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "disable"),
		),
		Action: func(ctx *KeyContext) tea.Msg {
			return disableAddonMsg{addonID: ctx.AddonID}
		},
	},
	RefreshCache: KeyMapEntry{
		Binding: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "refresh cache"),
		),
		Action: func(ctx *KeyContext) tea.Msg {
			return refreshCacheMsg{addonID: ctx.AddonID}
		},
	},
	Remove: KeyMapEntry{
		Binding: key.NewBinding(
			key.WithKeys("x"),
			key.WithHelp("x", "remove"),
		),
		Action: func(ctx *KeyContext) tea.Msg {
			return removeAddonMsg{addonID: ctx.AddonID}
		},
	},
	Confirm: KeyMapEntry{
		Binding: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "confirm"),
		),
		Action: func(ctx *KeyContext) tea.Msg {
			return confirmMsg{}
		},
	},
	Info: KeyMapEntry{
		Binding: key.NewBinding(
			key.WithKeys("v"),
			key.WithHelp("v", "view info"),
		),
		Action: func(ctx *KeyContext) tea.Msg {
			return infoMsg{}
		},
	},
	Cancel: KeyMapEntry{
		Binding: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
		Action: func(ctx *KeyContext) tea.Msg {
			return requestListViewMsg{}
		},
	},
}

// Update processes a key message against a subset of keys and executes the corresponding action
func (km KeyMap) Update(msg tea.KeyMsg, allowedKeys []KeyMapEntry, ctx *KeyContext) tea.Msg {
	for _, entry := range allowedKeys {
		if key.Matches(msg, entry.Binding) {
			if entry.Action != nil {
				return entry.Action(ctx)
			}
		}
	}
	return nil
}

// ExtractBindings extracts key.Binding objects from a list of KeyMapEntry
func ExtractBindings(entries []KeyMapEntry) []key.Binding {
	bindings := make([]key.Binding, len(entries))
	for i, entry := range entries {
		bindings[i] = entry.Binding
	}
	return bindings
}
