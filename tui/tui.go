package tui

import (
	"fmt"

	"gmod-addon-manager/addon"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	list          list.Model
	selectedAddon *addon.Addon
	input         textinput.Model
	manager       *addon.Manager
	state         string
	error         error
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

	// Create the list
	addonList := list.New(items, list.NewDefaultDelegate(), 0, 0)
	addonList.Title = "Garry's Mod Addons"
	addonList.SetFilteringEnabled(false)
	addonList.SetShowStatusBar(false)
	addonList.SetShowHelp(false)

	// Initialize text input
	input := textinput.New()
	input.Placeholder = "Enter addon ID"
	input.Focus()

	return model{
		list:    addonList,
		manager: manager,
		input:   input,
		state:   "list",
	}
}

type addonItem struct {
	addon addon.Addon
}

func (i addonItem) Title() string       { return i.addon.Title }
func (i addonItem) Description() string { return fmt.Sprintf("ID: %s | Author: %s", i.addon.ID, i.addon.Author) }
func (i addonItem) FilterValue() string { return i.addon.Title }

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "enter":
			if m.state == "list" {
				if len(m.list.Items()) > 0 {
					selected := m.list.SelectedItem().(addonItem)
					m.selectedAddon = &selected.addon
					m.state = "detail"
				}
			} else if m.state == "input" {
				addonID := m.input.Value()
				if addonID != "" {
					return m, func() tea.Msg {
						err := m.manager.GetAddon(addonID)
						if err != nil {
							return errorMsg{err}
						}
						return successMsg{fmt.Sprintf("Addon %s installed successfully", addonID)}
					}
				}
			}

		case "esc":
			if m.state != "list" {
				m.state = "list"
				m.error = nil
			}

		case "e":
			if m.state == "detail" && m.selectedAddon != nil {
				return m, func() tea.Msg {
					err := m.manager.EnableAddon(m.selectedAddon.ID)
					if err != nil {
						return errorMsg{err}
					}
					return successMsg{fmt.Sprintf("Addon %s enabled", m.selectedAddon.ID)}
				}
			}

		case "d":
			if m.state == "detail" && m.selectedAddon != nil {
				return m, func() tea.Msg {
					err := m.manager.DisableAddon(m.selectedAddon.ID)
					if err != nil {
						return errorMsg{err}
					}
					return successMsg{fmt.Sprintf("Addon %s disabled", m.selectedAddon.ID)}
				}
			}

		case "i":
			if m.state == "list" {
				m.state = "input"
				m.input.Reset()
			}
		}

	case errorMsg:
		m.error = msg.err
		return m, nil

	case successMsg:
		m.error = nil
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
	}

	// Handle input updates
	var cmd tea.Cmd
	if m.state == "input" {
		m.input, cmd = m.input.Update(msg)
	} else {
		m.list, cmd = m.list.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	if m.error != nil {
		return fmt.Sprintf("Error: %v\nPress any key to continue...", m.error)
	}

	switch m.state {
	case "list":
		return m.list.View() + "\n\n" +
			"[i] Install new addon  [q] Quit\n"

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
		status := "Disabled"
		if a.Enabled {
			status = "Enabled"
		}

		return fmt.Sprintf(
			"Addon Details\n\n"+
				"Title: %s\n"+
				"ID: %s\n"+
				"Author: %s\n"+
				"Status: %s\n\n"+
				"Description:\n%s\n\n"+
				"[e] Enable  [d] Disable  [Esc] Back\n",
			a.Title, a.ID, a.Author, status, a.Description,
		)

	default:
		return "Unknown state"
	}
}

type errorMsg struct{ err error }
type successMsg struct{ msg string }
