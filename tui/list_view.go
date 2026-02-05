package tui

import (
	"gmod-addon-manager/addon"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// ListModel displays and manages the addon list view
type ListModel struct {
	list         list.Model
	manager      *addon.Manager
	keys         *listKeyMap
	delegateKeys *delegateKeyMap
	commonKeys   *commonKeyMap
	help         help.Model
}

func NewListModel(manager *addon.Manager) *ListModel {
	// Create the list with custom delegate
	keys := newListKeyMap()
	delegateKeys := newDelegateKeyMap()
	commonKeys := newCommonKeyMap()
	addonList := list.New(buildAddonItems(manager), newItemDelegate(delegateKeys, manager), 0, 0)
	addonList.Title = "Garry's Mod Addons"
	addonList.KeyMap.PrevPage = key.NewBinding(
		key.WithKeys("left", "h", "pgup"),
		key.WithHelp("←/h/pgup", "prev page"),
	)
	addonList.KeyMap.NextPage = key.NewBinding(
		key.WithKeys("right", "l", "pgdown"),
		key.WithHelp("→/l/pgdn", "next page"),
	)
	addonList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			keys.installItem,
		}
	}
	addonList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			keys.installItem,
			keys.refreshList,
			keys.quit,
		}
	}

	return &ListModel{
		list:         addonList,
		manager:      manager,
		keys:         keys,
		delegateKeys: delegateKeys,
		commonKeys:   commonKeys,
		help:         help.New(),
	}
}

func (m *ListModel) Init() tea.Cmd {
	return nil
}

func (m *ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.installItem):
			return m, func() tea.Msg {
				return requestInputViewMsg{}
			}
		case key.Matches(msg, m.keys.refreshList):
			m.RefreshItems()
		}

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
		m.help.Width = msg.Width

	case successMsg:
		m.RefreshItems()
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *ListModel) RefreshItems() {
	m.list.SetItems(buildAddonItems(m.manager))
}

// buildAddonItems creates list items from addon manager data
func buildAddonItems(manager *addon.Manager) []list.Item {
	items := []list.Item{}
	addons, err := manager.GetAddonsInfo()
	if err == nil {
		for _, a := range addons {
			items = append(items, addonItem{addon: a})
		}
	}
	return items
}

func (m *ListModel) View() string {
	if len(m.list.Items()) == 0 {
		return "No addons installed.\n\nPress [i] to install a new addon or [q] to quit."
	}
	return m.list.View()
}
