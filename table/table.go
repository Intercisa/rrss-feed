package table

import (
	"fmt"
	"os"

	"bsipiczki.com/rss-feed/util"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mmcdole/gofeed"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.DoubleBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	items []*gofeed.Item
	table table.Model
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Sequence(clearTerminalCmd(), tea.Quit)
		case "enter":
			idx := m.table.Cursor()
			return m, tea.Sequence(
				openLinkCmd(m.items[idx].Link),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func openLinkCmd(path string) tea.Cmd {
	return func() tea.Msg {
		util.OpenFile(path)
		return nil
	}
}

func clearTerminalCmd() tea.Cmd {
	return func() tea.Msg {
		util.ClearTerminal()
		return nil
	}
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func getLongestNameLen(inputs []*gofeed.Item) int {
	max := 0

	for _, input := range inputs {
		if len(input.Title) > max {
			max = len(input.Title)
		}
	}
	termWidth := util.GetTermWidth()
	if max > termWidth {
		return termWidth
	}
	return max
}

func Render(inputs []*gofeed.Item) {
	columns := []table.Column{
		{Title: "Reddit RSS Feed", Width: util.GetTermWidth() - 4},
	}

	var rows []table.Row

	for _, input := range inputs {
		rows = append(rows, table.Row{input.Title})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(inputs)),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false).
		BorderLeft(false).
		BorderRight(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := model{table: t, items: inputs}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
