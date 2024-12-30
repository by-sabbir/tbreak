package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	progress    int
	width       int
	initialized bool
	totalTime   time.Duration
	elapsedTime time.Duration
}

func (m model) Init() tea.Cmd {
	return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.initialized = true
		return m, nil

	case tickMsg:
		m.elapsedTime += time.Second / 10
		m.progress = int((float64(m.elapsedTime) / float64(m.totalTime)) * 100)
		if m.elapsedTime >= m.totalTime {
			return m, tea.Quit
		}
		return m, tick()

	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	if !m.initialized {
		return "Initializing...\n"
	}

	barWidth := m.width - 20 // Adjust for padding and labels
	if barWidth < 10 {
		barWidth = 10 // Minimum width for the progress bar
	}

	filledWidth := m.progress * barWidth / 100
	emptyWidth := barWidth - filledWidth

	filled := strings.Repeat("█", filledWidth)
	empty := strings.Repeat(" ", emptyWidth)

	bar := fmt.Sprintf("[%s%s]", filled, empty)
	return fmt.Sprintf(
		"Progress: %d%% %s\nElapsed: %s / %s\nPress 'q' to quit.",
		m.progress,
		bar,
		m.elapsedTime.Truncate(time.Second).String(),
		m.totalTime.String(),
	)
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second/10, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func parseDuration(arg string) (time.Duration, error) {
	if strings.HasSuffix(arg, "m") {
		minutes := strings.TrimSuffix(arg, "m")
		min, err := strconv.Atoi(minutes)
		if err != nil {
			return 0, err
		}
		return time.Duration(min) * time.Minute, nil
	} else if strings.HasSuffix(arg, "s") {
		seconds := strings.TrimSuffix(arg, "s")
		sec, err := strconv.Atoi(seconds)
		if err != nil {
			return 0, err
		}
		return time.Duration(sec) * time.Second, nil
	}
	return 0, fmt.Errorf("invalid duration format: %s", arg)
}

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Println("Usage: go run main.go <duration> (e.g., '30s', '2m')")
		os.Exit(1)
	}

	duration, err := parseDuration(flag.Arg(0))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(model{totalTime: duration}, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Printf("Error starting app: %v\n", err)
	}
}
