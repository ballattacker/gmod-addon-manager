package tui

import (
	"fmt"

	"gmod-addon-manager/addon"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// addonItem is a list item wrapper for addon.Addon
type addonItem struct {
	addon addon.Addon
}

func (i addonItem) Title() string {
	return fmt.Sprintf("%s - %s", i.addon.ID, i.addon.Title)
}

func (i addonItem) Description() string {
	status := "❌ Disabled"
	if i.addon.Enabled {
		status = "✅ Enabled"
	}
	return status
}

func (i addonItem) FilterValue() string { return i.addon.Title }

// newItemDelegate creates a list delegate with key bindings for addon items
func newItemDelegate(allowedKeys []KeyMapEntry, manager *addon.Manager) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		selected, ok := m.SelectedItem().(addonItem)
		if !ok {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			ctx := &KeyContext{
				AddonID: selected.addon.ID,
			}
			result := GlobalKeyMap.Update(msg, allowedKeys, ctx)
			if result != nil {
				return func() tea.Msg { return result }
			}
		}

		return nil
	}

	d.ShortHelpFunc = func() []key.Binding {
		return []key.Binding{
			GlobalKeyMap.View.Binding,
		}
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{{
			GlobalKeyMap.View.Binding,
			GlobalKeyMap.Enable.Binding,
			GlobalKeyMap.Disable.Binding,
			GlobalKeyMap.RefreshCache.Binding,
			GlobalKeyMap.Remove.Binding,
		}}
	}

	return d
}
