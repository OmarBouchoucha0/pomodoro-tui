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
	timeStr := fmt.Sprintf("%02d:%02d", minutes, seconds)
	return fmt.Sprintf("⏰ Time: %s\n\n", timeStr)
}

func initialModel() model {
	workTotal := 1 * 6
	breakTotal := 1 * 3
	return model{
		workTime:      workTotal,
		isWorking:     false,
		workTimeLeft:  workTotal,
		breakTime:     breakTotal,
		isOnBreak:     false,
		breakTimeLeft: breakTotal,
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
			// Timer finished
			m.isWorking = false
			m.isOnBreak = true
			fmt.printf("breaking")
			return m, nil
		}
		if m.isOnBreak {
			if m.breakTimeLeft > 0 {
				m.breakTimeLeft--
				return m, tick()
			}
			m.isOnBreak = false
			return m, nil
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case " ":
			if m.isWorking {
				m.isWorking = false // Pause
			} else {
				m.isWorking = true // Start/Resume
				return m, tick()   // Start the ticker
			}
		case "r":
			// Reset timer
			m.workTimeLeft = m.workTime
			m.isWorking = false
			m.breakTimeLeft = m.breakTime
			m.isOnBreak = false
		}
	}
	return m, nil
}

func (m model) View() string {
	var s string
	if m.isWorking {
		s = formatTime(m.workTimeLeft)
		s += "Working... (Press space to pause)\n"
	} else {
		if m.workTimeLeft == 0 {
			s = "✅ Work session complete!\n"
		} else {
			s = "⏸️  Paused (Press space to start)\n"
		}
	}
	if m.isOnBreak {
		s = formatTime(m.breakTimeLeft)
		s += "✅ Work session complete!\n"
		s += "Go Get Some Rest\n"
	}

	s += "\nControls:\n"
	s += "• Space: Start/Pause\n"
	s += "• r: Reset timer\n"
	s += "• q: Quit\n"

	return s
}
