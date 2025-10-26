package main

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

const (
	stateListView uint = iota
	stateTitleView
	stateBodyView
)

type model struct {
	state       uint
	store       *Store
	notes       []Note
	currentNote Note
	listIndex   int
	textinput   textinput.Model
	textarea    textarea.Model
	exitCode    int
}

func NewModel(store *Store) (model, error) {
	notes, err := store.GetAllNotes()
	if err != nil {
		return model{}, err
	}

	return model{
		state:     stateListView,
		store:     store,
		notes:     notes,
		textinput: textinput.New(),
		textarea:  textarea.New(),
	}, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	m.textinput, cmd = m.textinput.Update(msg)
	cmds = append(cmds, cmd)

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch m.state {
		case stateListView:
			switch key {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "n":
				m.textinput.SetValue("")
				m.textinput.Focus()
				m.currentNote = Note{}
				m.state = stateTitleView
			case "up", "k":
				if m.listIndex > 0 {
					m.listIndex--
				}
			case "down", "j":
				if m.listIndex < len(m.notes)-1 {
					m.listIndex++
				}
			case "enter":
				m.textarea.SetValue(m.notes[m.listIndex].Body)
				m.textarea.Focus()
				m.textarea.CursorEnd()
				m.currentNote = m.notes[m.listIndex]
				m.state = stateBodyView
			}
		case stateTitleView:
			key := msg.String()
			switch key {
			case "enter":
				title := m.textinput.Value()
				if len(title) > 0 {
					m.currentNote.Title = title

					m.textarea.SetValue("")
					m.textarea.Focus()
					m.textarea.CursorEnd()

					m.state = stateBodyView
				}
			case "esc":
				m.state = stateListView

			}
		case stateBodyView:
			key := msg.String()
			switch key {
			case "ctrl+s":
				body := m.textarea.Value()
				m.currentNote.Body = body
				var err error
				if err = m.store.SaveNote(m.currentNote); err != nil {
					log.Error("error saving note", "err", err)
					m.exitCode = 1
					return m, tea.Quit
				}
				m.notes, err = m.store.GetAllNotes()
				if err != nil {
					log.Error("error fetching notes", "err", err)
					m.exitCode = 1
					return m, tea.Quit
				}
				m.listIndex = len(m.notes) - 1
				m.currentNote = Note{}
				m.state = stateListView
			case "esc":
				m.state = stateBodyView
			}
		}
	}
	return m, tea.Batch(cmds...)
}
