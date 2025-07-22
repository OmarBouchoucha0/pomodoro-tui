package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	workTime     int
	isWorking    bool
	workTimeLeft int
}
type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func InitialModel() model {
	timeTotal := 1 * 60
	return model{workTime: timeTotal, isWorking: false, workTimeLeft: timeTotal}
}

func (m model) Init() tea.Cmd {
	if m.isWorking {
		return tick()
	}
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		if m.isWorking {
			if m.workTimeLeft > 0 {
				m.workTimeLeft--
				return m, tick()
			}
			m.isWorking = false
			return m, nil
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "space":
			if m.isWorking {
				m.isWorking = false
			} else {
				m.isWorking = true
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	var s string
	s += "Press Space to Start."
	if m.isWorking {
		s += "\n⏳Working\n"
	}
	s += "\nPress q to quit.\n"
	return s
}
