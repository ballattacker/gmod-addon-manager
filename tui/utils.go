package tui

import (
	"fmt"

	"gmod-addon-manager/addon"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// Key bindings for different views
type listKeyMap struct {
	installItem key.Binding
	refreshList key.Binding
	quit        key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		installItem: key.NewBinding(
			key.WithKeys("i"),
			key.WithHelp("i", "install"),
		),
		refreshList: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),
		quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

type delegateKeyMap struct {
	choose       key.Binding
	enable       key.Binding
	disable      key.Binding
	refreshCache key.Binding
	remove       key.Binding
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter", "v"),
			key.WithHelp("enter", "view"),
		),
		enable: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "enable"),
		),
		disable: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "disable"),
		),
		refreshCache: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "refresh cache"),
		),
		remove: key.NewBinding(
			key.WithKeys("x"),
			key.WithHelp("x", "remove"),
		),
	}
}

type commonKeyMap struct {
	cancel key.Binding
}

func newCommonKeyMap() *commonKeyMap {
	return &commonKeyMap{
		cancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
	}
}

type inputKeyMap struct {
	install key.Binding
	info    key.Binding
	cancel  key.Binding
}

func newInputKeyMap() *inputKeyMap {
	return &inputKeyMap{
		install: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "install"),
		),
		info: key.NewBinding(
			key.WithKeys("v"),
			key.WithHelp("v", "view info"),
		),
		cancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
	}
}

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
func newItemDelegate(delegateKeys *delegateKeyMap, manager *addon.Manager) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		selected, ok := m.SelectedItem().(addonItem)
		if !ok {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, delegateKeys.choose):
				return func() tea.Msg {
					return requestDetailViewMsg{addon: &selected.addon}
				}
			case key.Matches(msg, delegateKeys.enable):
				return func() tea.Msg {
					err := manager.EnableAddon(selected.addon.ID)
					if err != nil {
						return errorMsg{err}
					}
					return successMsg{msg: "Addon enabled", refreshList: true}
				}
			case key.Matches(msg, delegateKeys.disable):
				return func() tea.Msg {
					err := manager.DisableAddon(selected.addon.ID)
					if err != nil {
						return errorMsg{err}
					}
					return successMsg{msg: "Addon disabled", refreshList: true}
				}
			case key.Matches(msg, delegateKeys.refreshCache):
				return func() tea.Msg {
					err := manager.RefreshCache(selected.addon.ID)
					if err != nil {
						return errorMsg{err}
					}
					return successMsg{msg: "Cache refreshed", refreshList: true}
				}
			case key.Matches(msg, delegateKeys.remove):
				return func() tea.Msg {
					err := manager.RemoveAddon(selected.addon.ID)
					if err != nil {
						return errorMsg{err}
					}
					return successMsg{msg: "Addon removed", refreshList: true}
				}
			}
		}

		return nil
	}

	d.ShortHelpFunc = func() []key.Binding {
		return []key.Binding{
			delegateKeys.choose,
		}
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{{
			delegateKeys.choose,
			delegateKeys.enable,
			delegateKeys.disable,
			delegateKeys.refreshCache,
			delegateKeys.remove,
		}}
	}

	return d
}
