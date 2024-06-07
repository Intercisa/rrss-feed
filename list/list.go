package list

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	common "bsipiczki.com/rss-feed/model"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return wrapText(i.desc, 150) }
func (i item) FilterValue() string { return i.title }

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func Render(comments []common.FormattedComments) {
	var items []list.Item
	for _, comment := range comments {
		c := item{title: comment.Author, desc: comment.Content}
		items = append(items, c)
	}

	delegate := list.NewDefaultDelegate()
	delegate.SetHeight(7)

	m := model{list: list.New(items, delegate, 0, 0)}
	m.list.Title = "Comments"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func wrapText(text string, width int) string {
	var wrappedLines []string

	// Split text into words
	words := strings.Fields(text)

	// Initialize current line
	var currentLine string

	// Iterate over words
	for _, word := range words {
		// If adding the word would exceed the width, start a new line
		if utf8.RuneCountInString(currentLine+word) > width {
			wrappedLines = append(wrappedLines, currentLine)
			currentLine = ""
		}

		// Add word to current line
		if currentLine != "" {
			currentLine += " "
		}
		currentLine += word
	}

	// Add the last line
	wrappedLines = append(wrappedLines, currentLine)

	// Join lines into a single string
	return strings.Join(wrappedLines, "\n")
}
