package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Coffee-themed colors
	coffeeColor = lipgloss.Color("#4A2C2A") // Dark coffee
	creamColor  = lipgloss.Color("#C4A69D") // Coffee with cream
	beanColor   = lipgloss.Color("#6F4E37") // Coffee bean brown
	steamColor  = lipgloss.Color("#9B9B9B") // Steam gray
	borderColor = lipgloss.Color("#8B4513") // Dark roast
	accentColor = lipgloss.Color("#D2691E") // Light roast
	foamColor   = lipgloss.Color("#E6BE8A") // Cappuccino foam

	// Base styles
	baseStyle = lipgloss.NewStyle().
			Align(lipgloss.Center)

	titleStyle = lipgloss.NewStyle().
			Foreground(foamColor).
			Bold(true).
			Padding(1).
			Align(lipgloss.Center)

	timeStyle = lipgloss.NewStyle().
			Foreground(creamColor).
			Align(lipgloss.Center).
			Padding(0, 1)

	progressStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Padding(0, 2)

	quitStyle = lipgloss.NewStyle().
			Foreground(creamColor).
			Align(lipgloss.Center).
			Margin(1)

	containerStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			AlignVertical(lipgloss.Center)

	// Coffee steam animation frames
	steam = []string{"░", "▒", "▓", "▒", "░"}
)

type model struct {
	progress    int
	width       int
	height      int
	initialized bool
	totalTime   time.Duration
	elapsedTime time.Duration
	steamFrame  int
	taskName    string
}

func (m model) Init() tea.Cmd {
	return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.initialized = true
		return m, nil

	case tickMsg:
		m.elapsedTime += time.Second / 10
		m.progress = int((float64(m.elapsedTime) / float64(m.totalTime)) * 100)
		m.steamFrame = (m.steamFrame + 1) % len(steam)
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
		return "Brewing your coffee...\n"
	}

	// Create animated steam
	steamAnim := lipgloss.NewStyle().Foreground(steamColor).Render(steam[m.steamFrame])

	// Coffee cup ASCII art with proper alignment
	coffeeCup := fmt.Sprintf(`
     %s %s %s
      )  )  )
     (  (  (
    .-========-.
    | ________ |
    | Break!!  |
    | -------- |
    '========='

    %s`, steamAnim, steamAnim, steamAnim, m.taskName)
	// Progress bar
	barWidth := m.width/2 - 10 // Adjusted for better centering
	if barWidth < 10 {
		barWidth = 10
	}

	// Stylized progress bar with coffee beans
	var progressBar strings.Builder
	progressBar.WriteString("╭" + strings.Repeat("━", barWidth+2) + "╮\n")
	progressBar.WriteString("│ ")

	cupPosition := int(float64(barWidth) * float64(m.progress) / 100.0)
	for i := 0; i < barWidth; i++ {
		if i == cupPosition {
			progressBar.WriteString(lipgloss.NewStyle().Foreground(accentColor).Render("☕"))
		} else if i < cupPosition {
			progressBar.WriteString(lipgloss.NewStyle().Foreground(beanColor).Render("♨"))
		} else {
			progressBar.WriteString(lipgloss.NewStyle().Foreground(creamColor).Render("○"))
		}
	}
	progressBar.WriteString(" │\n")
	progressBar.WriteString("╰" + strings.Repeat("━", barWidth+2) + "╯")

	// Status message based on progress
	var status string
	if m.progress < 33 {
		status = "Grinding beans..."
	} else if m.progress < 66 {
		status = "Brewing..."
	} else {
		status = "Almost ready!"
	}

	// Style and center all components
	title := titleStyle.Render(coffeeCup)
	progress := progressStyle.Render(fmt.Sprintf("%d%%", m.progress))
	progressBarStr := progressBar.String() // Removed extra styling
	timeInfo := timeStyle.Render(fmt.Sprintf(
		"%s\nBrew time: %s / %s",
		status,
		m.elapsedTime.Truncate(time.Second).String(),
		m.totalTime.String(),
	))
	quit := quitStyle.Render("Press 'q' to cancel your coffee break")

	// Combine all elements with vertical centering
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		progress,
		progressBarStr,
		timeInfo,
		quit,
	)

	// Center everything in the terminal
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
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
	if flag.NArg() < 1 || flag.NArg() > 2 {
		fmt.Println("Usage: go run main.go <duration> <task_name> (e.g., '30s' 'Making Coffee')")
		os.Exit(1)
	}

	duration, err := parseDuration(flag.Arg(0))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	taskName := "Coffee Brewing"
	if flag.NArg() == 2 {
		taskName = flag.Arg(1)
	}

	p := tea.NewProgram(
		model{totalTime: duration, taskName: taskName},
		tea.WithAltScreen(),
	)

	if err := p.Start(); err != nil {
		fmt.Printf("Error starting app: %v\n", err)
		os.Exit(1)
	}
}
