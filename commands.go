package main

import (
	"bufio"
	"os/exec"
	"strings"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
)

func runCommand(linesChan chan string, command string, fn func() tea.Msg) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("make", command)
		cmd.Dir = "../mattermost-server"

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return errMsg{err}
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			return errMsg{err}
		}

		err = cmd.Start()
		if err != nil {
			return errMsg{err}
		}

		r := bufio.NewScanner(stdout)
		e := bufio.NewScanner(stderr)

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			for r.Scan() {
				text := r.Text()
				if strings.TrimSpace(text) != "" {
					linesChan <- text
				}
			}
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			for e.Scan() {
				text := e.Text()
				if strings.TrimSpace(text) != "" {
					linesChan <- text
				}
			}
			wg.Done()
		}()

		cmd.Wait()
		wg.Wait()

		if fn != nil {
			return fn()
		}
		return ""
	}
}

func stopGoServer(linesChan chan string, fn func() tea.Msg) tea.Cmd {
	return runCommand(linesChan, "stop-server", fn)
}

func stopDocker(linesChan chan string, fn func() tea.Msg) tea.Cmd {
	return runCommand(linesChan, "stop-docker", fn)
}

func stopServer(linesChan chan string, fn func() tea.Msg) tea.Cmd {
	return stopGoServer(linesChan, stopDocker(linesChan, fn))
}

func runServer(linesChan chan string) tea.Cmd {
	return runCommand(linesChan, "run-server", nil)
}

func restartServer(linesChan chan string) tea.Cmd {
	return runCommand(linesChan, "restart-server", nil)
}

func debugServer(linesChan chan string) tea.Cmd {
	return runCommand(linesChan, "debug-server-headless", nil)
}
