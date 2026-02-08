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
	input   textinput.Model
	keyMaps []KeyMapEntry
	help    help.Model
	manager *addon.Manager
}

func NewInputModel(manager *addon.Manager) *InputModel {
	input := textinput.New()
	input.Placeholder = "Enter addon ID"
	input.Focus()

	keyMaps := []KeyMapEntry{
		GlobalKeyMap.Install,
		GlobalKeyMap.Detail,
		GlobalKeyMap.Cancel,
	}

	return &InputModel{
		input:   input,
		keyMaps: keyMaps,
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
		ctx := &KeyContext{
			AddonID: m.input.Value(),
		}
		result := GlobalKeyMap.Update(msg, m.keyMaps, ctx)
		if result != nil {
			return m, func() tea.Msg { return result }
		}

	case tea.WindowSizeMsg:
		m.input.Width = msg.Width
		m.help.Width = msg.Width

	case successMsg, requestListViewMsg:
		m.input.Reset()
		m.input.Focus()
	}

	// only allow numeric input for addon ID
	if msg, ok := msg.(tea.KeyMsg); ok {
		if msg.Type == tea.KeyRunes {
			for _, r := range msg.Runes {
				if r < '0' || r > '9' {
					return m, nil
				}
			}
		}
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *InputModel) View() string {
	return strings.Join([]string{
		"Install new addon",
		m.input.View(),
		m.help.ShortHelpView([]key.Binding{
			GlobalKeyMap.Install.Binding,
			GlobalKeyMap.Detail.Binding,
			GlobalKeyMap.Cancel.Binding,
		}),
	}, "\n\n")
}
