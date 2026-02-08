package tui

import (
	"fmt"

	"gmod-addon-manager/addon"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// DetailModel displays and manages the detail view for a selected addon
type DetailModel struct {
	addon   *addon.Addon
	keyMaps []KeyMapEntry
	help    help.Model
	manager *addon.Manager
}

func NewDetailModel(manager *addon.Manager) *DetailModel {
	keyMaps := []KeyMapEntry{
		GlobalKeyMap.Enable,
		GlobalKeyMap.Disable,
		GlobalKeyMap.Reload,
		GlobalKeyMap.Remove,
		GlobalKeyMap.Cancel,
	}

	return &DetailModel{
		keyMaps: keyMaps,
		help:    help.New(),
		manager: manager,
	}
}

func (m *DetailModel) Init() tea.Cmd {
	return nil
}

func (m *DetailModel) updateAddonInfo(addonID string) {
	if m.manager != nil {
		addonInfo, err := m.manager.GetAddonInfo(addonID)
		if err == nil {
			m.addon = addonInfo
		}
	}
}

func (m *DetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.addon != nil {
			ctx := &KeyContext{
				AddonID: m.addon.ID,
			}
			result := GlobalKeyMap.Update(msg, m.keyMaps, ctx)
			if result != nil {
				return m, func() tea.Msg { return result }
			}
		}

	case tea.WindowSizeMsg:
		m.help.Width = msg.Width

	case successMsg:
		m.updateAddonInfo(m.addon.ID)

	case requestDetailViewMsg:
		m.updateAddonInfo(msg.addonID)
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
				GlobalKeyMap.Enable.Binding,
				GlobalKeyMap.Disable.Binding,
				GlobalKeyMap.Reload.Binding,
				GlobalKeyMap.Remove.Binding,
				GlobalKeyMap.Cancel.Binding,
			}),
		a.Title, a.ID, a.Author, status, a.Installed,
	)
}
