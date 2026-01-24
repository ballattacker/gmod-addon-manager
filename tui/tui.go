package tui

import (
	"fmt"
	"io"

	"gmod-addon-manager/addon"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	tableHeaderStyle  = lipgloss.NewStyle().Bold(true).Padding(0, 1)
	tableRowStyle     = lipgloss.NewStyle().Padding(0, 1)
	selectedRowStyle  = lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Color("170"))
)

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(addonItem)
	if !ok {
		return
	}

	status := "❌ Disabled"
	if i.addon.Enabled {
		status = "✅ Enabled"
	}

	str := fmt.Sprintf("%-10s | %-30s | %-20s | %s",
		i.addon.ID, i.addon.Title, i.addon.Author, status)

	fn := tableRowStyle.Render
	if index == m.Index() {
		fn = selectedRowStyle.Render
	}

	fmt.Fprintf(w, fn(str))
}

type model struct {
	list          list.Model
	selectedAddon *addon.Addon
	input         textinput.Model
	manager       *addon.Manager
	state         string
	error         error
	loading       bool
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
	addonList := list.New(items, itemDelegate{}, 0, 0)
	addonList.Title = "Garry's Mod Addons"
	addonList.SetShowStatusBar(false)
	addonList.SetShowHelp(false)
	addonList.SetFilteringEnabled(false)
	addonList.Styles.Title = titleStyle
	addonList.Styles.PaginationStyle = paginationStyle

	// Initialize text input
	input := textinput.New()
	input.Placeholder = "Enter addon ID"
	input.Focus()

	return model{
		list:    addonList,
		manager: manager,
		input:   input,
		state:   "list",
		loading: false,
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

		case "esc":
			if m.state != "list" {
				m.state = "list"
				m.error = nil
			}

		case "e":
			if m.state == "detail" && m.selectedAddon != nil {
				m.loading = true
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
				m.loading = true
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

		case "r":
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

	if m.loading {
		return "Loading... Please wait.\n"
	}

	switch m.state {
	case "list":
		if len(m.list.Items()) == 0 {
			return "No addons installed.\n\nPress [i] to install a new addon or [q] to quit."
		}

		// Add table header
		header := tableHeaderStyle.Render(fmt.Sprintf("%-10s | %-30s | %-20s | %s",
			"ID", "Title", "Author", "Status")) + "\n"

		return m.list.Title + "\n\n" + header + m.list.View() + "\n\n" +
			"[i] Install new addon  [r] Refresh  [q] Quit\n"

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
				"[e] Enable  [d] Disable  [Esc] Back\n",
			a.Title, a.ID, a.Author, status, a.Description,
		)

	default:
		return "Unknown state"
	}
}

type errorMsg struct{ err error }
type successMsg struct{ msg string }
type refreshListMsg struct{ items []list.Item }
