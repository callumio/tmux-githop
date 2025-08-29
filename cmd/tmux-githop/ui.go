package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list     list.Model
	input    textinput.Model
	allItems []list.Item
	choice   string
	quitting bool
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 2)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			if i, ok := m.list.SelectedItem().(item); ok {
				m.choice = i.title
			}
			return m, tea.Quit
		case "up", "down":
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			return m, cmd
		default:
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			m.filterItems()
			return m, cmd
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *model) filterItems() {
	filter := strings.ToLower(m.input.Value())
	if filter == "" {
		m.list.SetItems(m.allItems)
		return
	}

	var filtered []list.Item
	for _, listItem := range m.allItems {
		if strings.Contains(strings.ToLower(listItem.(item).title), filter) {
			filtered = append(filtered, listItem)
		}
	}
	m.list.SetItems(filtered)
}

func (m model) View() string {
	if m.choice != "" {
		return ""
	}
	if m.quitting {
		return "Cancelled.\n"
	}
	return fmt.Sprintf("%s\n\n%s", m.input.View(), m.list.View())
}

func createList(title string) list.Model {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = title
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = lipgloss.NewStyle().MarginLeft(2)
	return l
}

func createInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "Type to filter..."
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 100
	return ti
}

func selectItem(items []list.Item, title string) (string, error) {
	if len(items) == 0 {
		return "", nil
	}

	l := createList(title)
	l.SetItems(items)

	m := model{
		list:     l,
		input:    createInput(),
		allItems: items,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	final, err := p.Run()
	if err != nil {
		return "", err
	}

	result := final.(model)
	if result.quitting || result.choice == "" {
		return "", nil
	}
	return result.choice, nil
}
