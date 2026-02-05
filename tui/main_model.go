package tui

import (
	"fmt"

	"gmod-addon-manager/addon"

	tea "github.com/charmbracelet/bubbletea"
)

// Model is the root TUI model that orchestrates all views
type Model struct {
	manager     *addon.Manager
	state       string // "list", "input", "detail"
	error       error
	loading     bool
	listModel   *ListModel
	inputModel  *InputModel
	detailModel *DetailModel
}

func NewModel(manager *addon.Manager) Model {
	return Model{
		manager:     manager,
		state:       "list",
		loading:     false,
		listModel:   NewListModel(manager),
		inputModel:  NewInputModel(manager),
		detailModel: NewDetailModel(manager),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// If there's an error, any key press dismisses it
	if m.error != nil {
		if _, ok := msg.(tea.KeyMsg); ok {
			m.error = nil
			return m, nil
		}
	}

	switch msg := msg.(type) {
	case errorMsg:
		m.error = msg.err
		m.loading = false
		return m, nil

	case successMsg:
		m.error = nil
		m.loading = false
		return m, nil

	case requestListViewMsg:
		m.state = "list"
		return m, nil

	case requestInputViewMsg:
		m.state = "input"
		m.inputModel.Reset()
		return m, nil

	case requestDetailViewMsg:
		m.state = "detail"
		m.detailModel.Update(msg)
		return m, nil
	}

	// Delegate to the active component
	switch m.state {
	case "list":
		_, cmd = m.listModel.Update(msg)
	case "input":
		_, cmd = m.inputModel.Update(msg)
	case "detail":
		_, cmd = m.detailModel.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	if m.error != nil {
		return fmt.Sprintf("Error: %v\nPress any key to continue...", m.error)
	}

	if m.loading {
		return "Loading... Please wait.\n"
	}

	switch m.state {
	case "list":
		return m.listModel.View()
	case "input":
		return m.inputModel.View()
	case "detail":
		return m.detailModel.View()
	default:
		return "Unknown state"
	}
}
