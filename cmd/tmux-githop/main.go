package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

func runCmd(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.Output()
	return strings.TrimSpace(string(out)), err
}

func tmuxRunning() bool {
	return exec.Command("pgrep", "tmux").Run() == nil
}

func inTmux() bool {
	return os.Getenv("TMUX") != ""
}

func sessionExists(name string) bool {
	return exec.Command("tmux", "has-session", "-t="+name).Run() == nil
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
	ti.Width = 0
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

func makeItems(data []string, descFunc func(string) string) []list.Item {
	items := make([]list.Item, len(data))
	for i, d := range data {
		items[i] = item{title: d, desc: descFunc(d)}
	}
	return items
}

func splitOutput(out string) []string {
	if out == "" {
		return nil
	}
	return strings.Split(out, "\n")
}

func pickSession() error {
	out, err := runCmd("ghq", "list", "-p")
	if err != nil {
		return err
	}

	paths := splitOutput(out)
	if len(paths) == 0 {
		return nil
	}

	items := makeItems(paths, filepath.Base)
	selected, err := selectItem(items, "Select Project")
	if err != nil || selected == "" {
		return err
	}

	selectedName := strings.ReplaceAll(filepath.Base(selected), ".", "_")

	if !inTmux() && !tmuxRunning() {
		return exec.Command("tmux", "new-session", "-s", selectedName, "-c", selected).Run()
	}

	if !sessionExists(selectedName) {
		if err := exec.Command("tmux", "new-session", "-ds", selectedName, "-c", selected).Run(); err != nil {
			return err
		}
	}

	return exec.Command("tmux", "switch-client", "-t", selectedName).Run()
}

func switchSession() error {
	out, err := runCmd("tmux", "list-sessions", "-F", "#{session_name}")
	if err != nil {
		return err
	}

	sessions := splitOutput(out)
	if len(sessions) == 0 {
		return nil
	}

	items := makeItems(sessions, func(string) string { return "tmux session" })
	selected, err := selectItem(items, "Switch Session")
	if err != nil || selected == "" {
		return err
	}

	return exec.Command("tmux", "switch-client", "-t", selected).Run()
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func main() {
	var pick = flag.Bool("p", false, "pick session from ghq repos")
	var switch_ = flag.Bool("s", false, "switch between existing sessions")
	flag.Parse()

	if !commandExists("ghq") {
		fmt.Fprintf(os.Stderr, "Error: ghq is not installed\n")
		os.Exit(1)
	}

	if !commandExists("tmux") {
		fmt.Fprintf(os.Stderr, "Error: tmux is not installed\n")
		os.Exit(1)
	}

	selections := 0
	if *pick {
		selections++
	}
	if *switch_ {
		selections++
	}

	if selections != 1 {
		fmt.Fprintln(os.Stderr, "Please make exactly one selection (-p or -s)")
		os.Exit(1)
	}

	var err error
	if *pick {
		err = pickSession()
	} else {
		err = switchSession()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
