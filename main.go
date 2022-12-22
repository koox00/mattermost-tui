package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/wrap"
)

type model struct {
	linesChan chan string
	width     int
	err       error
}

type lineMsg string

type errMsg struct{ err error }

type Log struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Msg       string                 `json:"msg"`
	Fields    map[string]interface{} `json:"-"`
}

func (e errMsg) Error() string { return e.err.Error() }

func waitForLine(linesChan chan string) tea.Cmd {
	return func() tea.Msg {
		return lineMsg(<-linesChan)
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		runServer(m.linesChan),
		waitForLine(m.linesChan),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width

	case lineMsg:
		line := ""
		ll, err := fromJSON(string(msg))
		if err != nil {
			line = formatText(string(msg))
		} else {
			line = formatLog(ll)
		}

		return m, tea.Batch(
			tea.Println(wrap.String(line, m.width)),
			waitForLine(m.linesChan),
		)

	case errMsg:
		m.err = msg
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "s":
			return m, stopServer(m.linesChan, nil)
		case "d":
			return m, stopServer(m.linesChan, debugServer(m.linesChan))
		case "ctrl+c", "q":
			return m, stopServer(m.linesChan, tea.Quit)
		case "r":
			m.err = nil
			return m, restartServer(m.linesChan)
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err)
	}
	return ""
}

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	m := model{
		linesChan: make(chan string),
	}
	defer close(m.linesChan)
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
		os.Exit(1)
	}
}
