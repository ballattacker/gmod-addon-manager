package tui

import (
	"fmt"

	"gmod-addon-manager/addon"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// DetailModel displays and manages the detail view for a selected addon
type DetailModel struct {
	addon        *addon.Addon
	delegateKeys *delegateKeyMap
	commonKeys   *commonKeyMap
	help         help.Model
	manager      *addon.Manager
}

func NewDetailModel(manager *addon.Manager) *DetailModel {
	return &DetailModel{
		delegateKeys: newDelegateKeyMap(),
		commonKeys:   newCommonKeyMap(),
		help:         help.New(),
		manager:      manager,
	}
}

func (m *DetailModel) Init() tea.Cmd {
	return nil
}

func (m *DetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.commonKeys.cancel):
			return m, func() tea.Msg {
				return requestViewMsg{view: "list"}
			}
		case key.Matches(msg, m.delegateKeys.enable):
			if m.addon != nil {
				return m, func() tea.Msg {
					err := m.manager.EnableAddon(m.addon.ID)
					if err != nil {
						return errorMsg{err}
					}
					return successMsg{fmt.Sprintf("Addon %s enabled", m.addon.ID)}
				}
			}
		case key.Matches(msg, m.delegateKeys.disable):
			if m.addon != nil {
				return m, func() tea.Msg {
					err := m.manager.DisableAddon(m.addon.ID)
					if err != nil {
						return errorMsg{err}
					}
					return successMsg{fmt.Sprintf("Addon %s disabled", m.addon.ID)}
				}
			}
		case key.Matches(msg, m.delegateKeys.refreshCache):
			if m.addon != nil {
				return m, func() tea.Msg {
					err := m.manager.RefreshCache(m.addon.ID)
					if err != nil {
						return errorMsg{err}
					}

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
		case key.Matches(msg, m.delegateKeys.remove):
			if m.addon != nil {
				return m, func() tea.Msg {
					err := m.manager.RemoveAddon(m.addon.ID)
					if err != nil {
						return errorMsg{err}
					}

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

	case tea.WindowSizeMsg:
		m.help.Width = msg.Width

	case selectAddonMsg:
		m.addon = msg.addon
	}

	return m, nil
}

func (m *DetailModel) View() string {
	if m.addon == nil {
		return "No addon selected\nPress [Esc] to return"
	}

	a := m.addon
	status := "❌ Disabled"
	if a.Enabled {
		status = "✅ Enabled"
	}

	return fmt.Sprintf(
		"Addon Details\n\n"+
			"Title: %s\n"+
			"ID: %s\n"+
			"Author: %s\n"+
			"Status: %s\n"+
			"Installed: %t\n"+
			"\n"+
			m.help.ShortHelpView([]key.Binding{
				m.delegateKeys.enable,
				m.delegateKeys.disable,
				m.delegateKeys.refreshCache,
				m.delegateKeys.remove,
				m.commonKeys.cancel,
			}),
		a.Title, a.ID, a.Author, status, a.Installed,
	)
}

func (m *DetailModel) SetAddon(addon *addon.Addon) {
	m.addon = addon
}
