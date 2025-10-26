package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

func main() {
	store := new(Store)
	if err := store.Init(); err != nil {
		log.Error("unable to init store", "err", err)
		os.Exit(1)
	}

	m, err := NewModel(store)
	if err != nil {
		log.Error("unable to init model", "err", err)
		os.Exit(1)
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Error("unable to run tui", "err", err)
		os.Exit(1)
	}
}
