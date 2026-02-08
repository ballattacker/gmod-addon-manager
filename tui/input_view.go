package tui

import (
	"strings"

	"gmod-addon-manager/addon"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// InputModel displays and manages the input view for installing addons
type InputModel struct {
	input       textinput.Model
	allowedKeys []KeyMapEntry
	help        help.Model
	manager     *addon.Manager
}

func NewInputModel(manager *addon.Manager) *InputModel {
	input := textinput.New()
	input.Placeholder = "Enter addon ID"
	input.Focus()

	allowedKeys := []KeyMapEntry{
		GlobalKeyMap.Confirm,
		GlobalKeyMap.Info,
		GlobalKeyMap.Cancel,
	}

	return &InputModel{
		input:       input,
		allowedKeys: allowedKeys,
		help:        help.New(),
		manager:     manager,
	}
}

func (m *InputModel) Init() tea.Cmd {
	return nil
}

func (m *InputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		addonID := m.input.Value()
		switch {
		case key.Matches(msg, GlobalKeyMap.Cancel.Binding):
			return m, func() tea.Msg {
				return requestListViewMsg{}
			}
		case key.Matches(msg, GlobalKeyMap.Confirm.Binding):
			if addonID != "" {
				return m, func() tea.Msg {
					return confirmInstallMsg{addonID: addonID}
				}
			}
		case key.Matches(msg, GlobalKeyMap.Info.Binding):
			if addonID != "" {
				return m, func() tea.Msg {
					return getAddonInfoMsg{addonID: addonID}
				}
			}
		}

	case tea.WindowSizeMsg:
		m.input.Width = msg.Width
		m.help.Width = msg.Width
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *InputModel) View() string {
	return strings.Join([]string{
		"Install new addon",
		m.input.View(),
		m.help.ShortHelpView(ExtractBindings(m.allowedKeys)),
	}, "\n\n")
}

func (m *InputModel) Reset() {
	m.input.Reset()
	m.input.Focus()
}

func (m *InputModel) Value() string {
	return m.input.Value()
}
