package main

import (
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	blackColor  = lipgloss.Color("#000000")
	whiteColor  = lipgloss.Color("#FFFFFF")
	debugColor  = lipgloss.Color("#828BB8")
	infoColor   = lipgloss.Color("#82AAFF")
	warnColor   = lipgloss.Color("#FFC777")
	dangerColor = lipgloss.Color("#FF757F")
)

func formatLabel(label string) string {
	return lipgloss.NewStyle().
		Faint(true).
		Render(label)
}

func formatTime(timestamp string) string {
	DefTimestampFormat := "2006-01-02 15:04:05.000 -07:00"

	style := lipgloss.NewStyle().Align(lipgloss.Center)

	tm, err := time.Parse(DefTimestampFormat, timestamp)

	if err != nil {
		panic(err)
	}

	return style.Render(tm.Format(time.StampMilli))
}

func formatLevel(level string) string {
	style := lipgloss.NewStyle().Bold(true)

	if level == "debug" {
		style = style.Foreground(debugColor)
	}

	if level == "info" {
		style = style.Background(infoColor).Foreground(blackColor)
	}

	if level == "error" {
		style = style.Background(dangerColor).Foreground(blackColor)
	}

	if level == "warn" {
		style = style.Background(warnColor).Foreground(blackColor)
	}

	return style.
		Width(8).
		Align(lipgloss.Center).
		Render(strings.ToUpper(level))
}

func formatCell(text string) string {
	return lipgloss.NewStyle().
		Align(lipgloss.Left).
		Render(text + "\n")
}

func formatMsg(msg string) string {
	return lipgloss.NewStyle().Render(msg)
}

func formatText(line string) string {
	return lipgloss.NewStyle().
		Italic(true).
		Faint(true).
		Render(line)
}
