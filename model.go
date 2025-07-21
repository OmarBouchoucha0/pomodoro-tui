package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	msg string
}

func InitialModel() model {
	return model{msg: "Hello"}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	s := m.msg
	// The footer
	s += "\nPress q to quit.\n"
	// Send the UI for rendering
	return s
}
