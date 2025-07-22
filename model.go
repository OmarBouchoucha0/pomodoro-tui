package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	workTime      int
	isWorking     bool
	workTimeLeft  int
	breakTime     int
	isOnBreak     bool
	breakTimeLeft int
	cycles        int
	cyclesLeft    int
}
type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func formatTime(time int) string {
	minutes := time / 60
	seconds := time % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func initialModel() model {
	workTotal := 1 * 6
	breakTotal := 1 * 3
	cycles := 3
	return model{
		workTime:      workTotal,
		isWorking:     false,
		workTimeLeft:  workTotal,
		breakTime:     breakTotal,
		isOnBreak:     false,
		breakTimeLeft: breakTotal,
		cycles:        cycles,
		cyclesLeft:    cycles,
	}
}

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("Pomodoro")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		if m.isWorking {
			if m.workTimeLeft > 0 {
				m.workTimeLeft--
				return m, tick() // Continue ticking
			}
			// Work timer finished - switch to break
			m.isWorking = false
			m.isOnBreak = true
			return m, tick() // Start break timer immediately
		}
		if m.isOnBreak {
			if m.breakTimeLeft > 0 {
				m.breakTimeLeft--
				return m, tick()
			}
			// Break finished - reset for next cycle
			m.isOnBreak = false
			m.workTimeLeft = m.workTime
			m.breakTimeLeft = m.breakTime
			m.cyclesLeft--
			if m.cyclesLeft > 0 {
				m.isWorking = true
			} else {
				m.isWorking = false
			}
			return m, tick()
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case " ":
			if m.isWorking {
				m.isWorking = false // Pause work
			} else if m.isOnBreak {
				m.isOnBreak = false // Pause break
			} else {
				// Start work (from paused or initial state)
				m.isWorking = true
				return m, tick()
			}
		case "r":
			// Reset timer
			m.workTimeLeft = m.workTime
			m.isWorking = false
			m.breakTimeLeft = m.breakTime
			m.isOnBreak = false
			m.cyclesLeft = m.cycles
		}
	}
	return m, nil
}

func (m model) View() string {
	var s string

	if m.isWorking {
		// Working state
		s += fmt.Sprintf("Work Time: %s\n\n", formatTime(m.workTimeLeft))
		s += "Working... (Press space to pause)\n"
	} else if m.isOnBreak {
		// Break state
		s += fmt.Sprintf("☕ Break Time: %s\n\n", formatTime(m.breakTimeLeft))
		s += "Break - Go Get Some Rest! (Press space to pause)\n"
	} else {
		// Paused or initial state
		if m.workTimeLeft == 0 && m.breakTimeLeft == m.breakTime {
			// Just finished a full cycle
			s += "🎉 Cycle Complete! Ready for next work session.\n"
			s += fmt.Sprintf("⏰ Work Time: %s\n\n", formatTime(m.workTime))
		} else if m.workTimeLeft < m.workTime {
			// Paused during work
			s += fmt.Sprintf("⏰ Work Time: %s\n\n", formatTime(m.workTimeLeft))
			s += "⏸️ Work Paused (Press space to resume)\n"
		} else if m.breakTimeLeft < m.breakTime {
			// Paused during break
			s += fmt.Sprintf("⏰ Break Time: %s\n\n", formatTime(m.breakTimeLeft))
			s += "⏸️ Break Paused (Press space to resume)\n"
		} else {
			// Initial state
			s += fmt.Sprintf("⏰ Work Time: %s\n\n", formatTime(m.workTimeLeft))
			s += "Press SpaceBar to Start Your Pomodoro!\n"
		}
	}

	s += "\nControls:\n"
	s += "• Space: Start/Pause\n"
	s += "• r: Reset timer\n"
	s += "• q: Quit\n"

	return s
}
