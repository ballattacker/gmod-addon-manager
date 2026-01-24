package tui

import (
	"fmt"

	"gmod-addon-manager/addon"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/charmbracelet/lipgloss"
)

type model struct {
	list          list.Model
	selectedAddon *addon.Addon
	input         textinput.Model
	manager       *addon.Manager
	state         string
	error         error
	loading       bool
	keys          *listKeyMap
	delegateKeys  *delegateKeyMap
}

func NewModel(manager *addon.Manager) model {
	// Initialize the addon list
	items := []list.Item{}
	addons, err := manager.GetAddonsInfo()
	if err == nil {
		for _, a := range addons {
			items = append(items, addonItem{addon: a})
		}
	}

	// Create the list with custom delegate
	addonList := list.New(items, newItemDelegate(newDelegateKeyMap()), 0, 0)
	addonList.Title = "Garry's Mod Addons"

	// Initialize text input
	input := textinput.New()
	input.Placeholder = "Enter addon ID"
	input.Focus()

	return model{
		list:          addonList,
		manager:       manager,
		input:         input,
		state:         "list",
		loading:       false,
		keys:          newListKeyMap(),
		delegateKeys:  newDelegateKeyMap(),
	}
}

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

type listKeyMap struct {
	toggleSpinner    key.Binding
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
	installItem      key.Binding
	refreshList      key.Binding
	quit             key.Binding
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

func (k *listKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.installItem,
		k.refreshList,
		k.quit,
	}
}

func (k *listKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.installItem,
			k.refreshList,
			k.quit,
		},
	}
}

type delegateKeyMap struct {
	choose key.Binding
	enable key.Binding
	disable key.Binding
	refreshCache key.Binding
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
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
	}
}

func (d delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.choose,
		d.enable,
		d.disable,
		d.refreshCache,
	}
}

func (d delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.choose,
			d.enable,
			d.disable,
			d.refreshCache,
		},
	}
}

func newItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		if _, ok := m.SelectedItem().(addonItem); ok {
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.choose):
				if len(m.Items()) > 0 {
					selected := m.SelectedItem().(addonItem)
					return func() tea.Msg {
						return selectAddonMsg{addon: &selected.addon}
					}
				}
			case key.Matches(msg, keys.enable):
				if len(m.Items()) > 0 {
					selected := m.SelectedItem().(addonItem)
					return func() tea.Msg {
						return enableAddonMsg{id: selected.addon.ID}
					}
				}
			case key.Matches(msg, keys.disable):
				if len(m.Items()) > 0 {
					selected := m.SelectedItem().(addonItem)
					return func() tea.Msg {
						return disableAddonMsg{id: selected.addon.ID}
					}
				}
			case key.Matches(msg, keys.refreshCache):
				if len(m.Items()) > 0 {
					selected := m.SelectedItem().(addonItem)
					return func() tea.Msg {
						return refreshCacheMsg{id: selected.addon.ID}
					}
				}
			}
		}

		return nil
	}

	d.ShortHelpFunc = func() []key.Binding {
		return []key.Binding{
			keys.choose,
			keys.enable,
			keys.disable,
			keys.refreshCache,
		}
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{
			{
				keys.choose,
				keys.enable,
				keys.disable,
				keys.refreshCache,
			},
		}
	}

	return d
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.installItem):
			if m.state == "list" {
				m.state = "input"
				m.input.Reset()
			}

		case key.Matches(msg, m.keys.refreshList):
			if m.state == "list" {
				m.loading = true
				return m, func() tea.Msg {
					items := []list.Item{}
					addons, err := m.manager.GetAddonsInfo()
					if err == nil {
						for _, a := range addons {
							items = append(items, addonItem{addon: a})
						}
					}
					return refreshListMsg{items}
				}
			}
		}

		// Handle Esc key for detail view
		if msg.String() == "esc" {
			switch m.state {
			case "detail":
				m.state = "list"
				m.selectedAddon = nil
				return m, nil
			case "input":
				m.state = "list"
				m.error = nil
				return m, nil
			}
		}

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)

	case errorMsg:
		m.error = msg.err
		m.loading = false
		return m, nil

	case successMsg:
		m.error = nil
		m.loading = false
		// Refresh the list after successful operation
		items := []list.Item{}
		addons, err := m.manager.GetAddonsInfo()
		if err == nil {
			for _, a := range addons {
				items = append(items, addonItem{addon: a})
			}
		}
		m.list.SetItems(items)
		m.state = "list"
		return m, nil

	case refreshListMsg:
		m.list.SetItems(msg.items)
		m.loading = false
		return m, nil

	case selectAddonMsg:
		m.selectedAddon = msg.addon
		m.state = "detail"
		return m, nil

	case enableAddonMsg:
		m.loading = true
		return m, func() tea.Msg {
			err := m.manager.EnableAddon(msg.id)
			if err != nil {
				return errorMsg{err}
			}
			return successMsg{fmt.Sprintf("Addon %s enabled", msg.id)}
		}

	case disableAddonMsg:
		m.loading = true
		return m, func() tea.Msg {
			err := m.manager.DisableAddon(msg.id)
			if err != nil {
				return errorMsg{err}
			}
			return successMsg{fmt.Sprintf("Addon %s disabled", msg.id)}
		}

	case refreshCacheMsg:
		m.loading = true
		return m, func() tea.Msg {
			// Use the public RefreshCache method
			err := m.manager.RefreshCache(msg.id)
			if err != nil {
				return errorMsg{err}
			}

			// Refresh the list to show updated info
			items := []list.Item{}
			addons, err := m.manager.GetAddonsInfo()
			if err == nil {
				for _, a := range addons {
					items = append(items, addonItem{addon: a})
				}
			}
			return refreshListMsg{items}
		}
	}

	// Handle input updates
	if m.state == "input" {
		m.input, cmd = m.input.Update(msg)
		if msg, ok := msg.(tea.KeyMsg); ok {
			if msg.String() == "enter" {
				addonID := m.input.Value()
				if addonID != "" {
					m.loading = true
					return m, func() tea.Msg {
						err := m.manager.GetAddon(addonID)
						if err != nil {
							return errorMsg{err}
						}
						return successMsg{fmt.Sprintf("Addon %s installed successfully", addonID)}
					}
				}
			}
		}
	} else {
		m.list, cmd = m.list.Update(msg)
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.error != nil {
		return fmt.Sprintf("Error: %v\nPress any key to continue...", m.error)
	}

	if m.loading {
		return "Loading... Please wait.\n"
	}

	switch m.state {
	case "list":
		if len(m.list.Items()) == 0 {
			return "No addons installed.\n\nPress [i] to install a new addon or [q] to quit."
		}
		return m.list.View()

	case "input":
		return fmt.Sprintf(
			"Install new addon\n\n%s\n\n[Enter] Install  [Esc] Cancel\n",
			m.input.View(),
		)

	case "detail":
		if m.selectedAddon == nil {
			return "No addon selected\nPress [Esc] to return"
		}

		a := m.selectedAddon
		status := "❌ Disabled"
		if a.Enabled {
			status = "✅ Enabled"
		}

		return fmt.Sprintf(
			"Addon Details\n\n"+
				"Title: %s\n"+
				"ID: %s\n"+
				"Author: %s\n"+
				"Status: %s\n\n"+
				"Description:\n%s\n\n"+
				"[e] Enable  [d] Disable  [c] Refresh Cache  [Esc] Back\n",
			a.Title, a.ID, a.Author, status, a.Description,
		)

	default:
		return "Unknown state"
	}
}

type errorMsg struct{ err error }
type successMsg struct{ msg string }
type refreshListMsg struct{ items []list.Item }
type selectAddonMsg struct{ addon *addon.Addon }
type enableAddonMsg struct{ id string }
type disableAddonMsg struct{ id string }
type refreshCacheMsg struct{ id string }
