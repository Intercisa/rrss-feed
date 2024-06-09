package table

import (
	"fmt"
	"os"
	"strconv"

	"bsipiczki.com/rss-feed/list"
	common "bsipiczki.com/rss-feed/model"
	"bsipiczki.com/rss-feed/scrap"
	"bsipiczki.com/rss-feed/util"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mmcdole/gofeed"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.DoubleBorder()).
	BorderForeground(lipgloss.Color("240"))

type comment struct {
	formatted []common.FormattedComments
}

type model struct {
	items         []*gofeed.Item
	comment       []comment
	table         table.Model
	showAltScreen bool
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
		case "c":
			list.Render(m.comment[m.table.Cursor()].formatted)
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
		{Title: "Reddit RSS Feed", Width: util.GetTermWidth() - 14},
		{Title: "Comments", Width: 10},
	}

	var rows []table.Row
	comment := getComments(inputs)

	for i, input := range inputs {
		rows = append(rows, table.Row{input.Title, strconv.Itoa(len(comment[i].formatted))})
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

	m := model{
		table:   t,
		items:   inputs,
		comment: comment,
	}

	util.ClearTerminal()

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func getComments(inputs []*gofeed.Item) []comment {
	var comments []comment
	for _, input := range inputs {
		fl := util.TransformToOldReddit(input.Link)
		cs := scrap.Scrap(fl)
		fc := getFormattedComments(cs)
		comments = append(comments, comment{fc})
	}
	return comments
}

func getFormattedComments(comments []scrap.Comment) []common.FormattedComments {
	var formatted []common.FormattedComments
	for _, comment := range comments {
		fc := common.FormattedComments{Author: comment.Author, Content: comment.Content}
		formatted = append(formatted, fc)
	}
	return formatted
}
