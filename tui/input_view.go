package tui

import (
	"fmt"
	"strings"

	"gmod-addon-manager/addon"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// InputModel displays and manages the input view for installing addons
type InputModel struct {
	input   textinput.Model
	keys    *inputKeyMap
	help    help.Model
	manager *addon.Manager
}

func NewInputModel(manager *addon.Manager) *InputModel {
	input := textinput.New()
	input.Placeholder = "Enter addon ID"
	input.Focus()
	inputKeys := newInputKeyMap()

	return &InputModel{
		input:   input,
		keys:    inputKeys,
		help:    help.New(),
		manager: manager,
	}
}

func (m *InputModel) Init() tea.Cmd {
	return nil
}

func (m *InputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.cancel):
			return m, func() tea.Msg {
				return requestViewMsg{view: "list"}
			}
		case key.Matches(msg, m.keys.install):
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
		case key.Matches(msg, m.keys.info):
			addonID := m.input.Value()
			if addonID != "" {
				return m, func() tea.Msg {
					addonInfo, err := m.manager.GetAddonInfo(addonID)
					if err != nil {
						return errorMsg{err}
					}
					return selectAddonMsg{addon: addonInfo}
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
		m.help.ShortHelpView([]key.Binding{
			m.keys.install,
			m.keys.info,
			m.keys.cancel,
		}),
	}, "\n\n")
}

func (m *InputModel) Reset() {
	m.input.Reset()
	m.input.Focus()
}

func (m *InputModel) Value() string {
	return m.input.Value()
}
