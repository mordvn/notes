package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	styleAppName    = lipgloss.NewStyle().Bold(true).Padding(0, 1).Background(lipgloss.Color("#F8CE26")).Foreground(lipgloss.Color("#515151"))
	styleFaint      = lipgloss.NewStyle().Faint(true)
	styleEnumerator = lipgloss.NewStyle().Foreground(lipgloss.Color("#F8CE26"))
)

const shortBodyLength = 30

func (m model) View() string {
	s := styleAppName.Render("NOTES APP") + "\n\n"

	switch m.state {
	case stateListView:
		for i, note := range m.notes {
			prefix := " "
			if i == m.listIndex {
				prefix = ">"
			}

			shortenBody := strings.ReplaceAll(note.Body, "\n", " ")
			if len(shortenBody) > shortBodyLength {
				shortenBody = shortenBody[:shortBodyLength] + "..."
			}

			s += styleEnumerator.Render(prefix+" ") + note.Title + " | " + styleFaint.Render(shortenBody) + "\n\n"
		}

		s += styleFaint.Render("n — new node, q — quit")

	case stateTitleView:
		s += "Note title:\n\n"
		s += m.textinput.View() + "\n\n"

		s += styleFaint.Render("enter — save, esc — discard")
	case stateBodyView:
		s += "Note:\n\n"
		s += m.textarea.View() + "\n\n"
		s += styleFaint.Render("ctrl+s — save, esc — discard")
	}

	return s
}
