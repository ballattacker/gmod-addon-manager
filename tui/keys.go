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
	Refresh KeyMapEntry
	Quit    KeyMapEntry
	Input   KeyMapEntry
	Detail  KeyMapEntry
	Enable  KeyMapEntry
	Disable KeyMapEntry
	Reload  KeyMapEntry
	Install KeyMapEntry
	Remove  KeyMapEntry
	Cancel  KeyMapEntry
}

// GlobalKeyMap is the single master keymap with all keybindings and actions
var GlobalKeyMap = KeyMap{
	Refresh: KeyMapEntry{
		Binding: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),
		Action: func(ctx *KeyContext) tea.Msg {
			return successMsg{msg: "refreshing..."}
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
	Input: KeyMapEntry{
		Binding: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "input"),
		),
		Action: func(ctx *KeyContext) tea.Msg {
			return requestInputViewMsg{}
		},
	},
	Detail: KeyMapEntry{
		Binding: key.NewBinding(
			key.WithKeys("enter", "v"),
			key.WithHelp("enter", "view detail info"),
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
	Reload: KeyMapEntry{
		Binding: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "reload"),
		),
		Action: func(ctx *KeyContext) tea.Msg {
			return reloadAddonMsg{addonID: ctx.AddonID}
		},
	},
	Install: KeyMapEntry{
		Binding: key.NewBinding(
			key.WithKeys("i"),
			key.WithHelp("i", "install"),
		),
		Action: func(ctx *KeyContext) tea.Msg {
			return installAddonMsg{addonID: ctx.AddonID}
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
	Cancel: KeyMapEntry{
		Binding: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
		Action: func(ctx *KeyContext) tea.Msg {
			return cancelMsg{}
		},
	},
}

// Update processes a key message against a subset of keys and executes the corresponding action
func (km KeyMap) Update(msg tea.KeyMsg, keyMaps []KeyMapEntry, ctx *KeyContext) tea.Msg {
	for _, entry := range keyMaps {
		if key.Matches(msg, entry.Binding) {
			if entry.Action != nil {
				return entry.Action(ctx)
			}
		}
	}
	return nil
}
