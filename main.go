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
	// Colors
	teaColor    = lipgloss.Color("#AF5E1F") // Tea brown
	cupColor    = lipgloss.Color("#D6A278") // Cup color
	steamColor  = lipgloss.Color("#B4C5BC") // Steam gray
	borderColor = lipgloss.Color("#E6CCB2") // Light border
	textColor   = lipgloss.Color("#9C6644") // Text brown

	// Styles
	titleStyle = lipgloss.NewStyle().
			Foreground(teaColor).
			Bold(true).
			Margin(1).
			Align(lipgloss.Center)

	timeStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Align(lipgloss.Center)

	quitStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Align(lipgloss.Center).
			Margin(1)

	// Progress bar characters
	steam = []string{"∿", "≋", "⋮", "⋰", "⋱"} // Steam animation frames
)

type model struct {
	progress    int
	width       int
	initialized bool
	totalTime   time.Duration
	elapsedTime time.Duration
	steamFrame  int
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
		return "Preparing your tea break...\n"
	}

	// Create animated steam
	steamAnim := lipgloss.NewStyle().Foreground(steamColor).Render(steam[m.steamFrame])

	// Tea cup ASCII art with animated steam
	teaCup := fmt.Sprintf(`
    %s %s %s
     )))
    (((
  +-----+
  |     |
        |   Tea Time!
  +-----+
`, steamAnim, steamAnim, steamAnim)

	// Progress bar
	barWidth := m.width - 20
	if barWidth < 10 {
		barWidth = 10
	}

	// Stylized progress bar
	var progressBar strings.Builder
	progressBar.WriteString("╭" + strings.Repeat("━", barWidth+2) + "╮\n│ ")

	cupPosition := int(float64(barWidth) * float64(m.progress) / 100.0)
	for i := 0; i < barWidth; i++ {
		if i == cupPosition {
			progressBar.WriteString(lipgloss.NewStyle().Foreground(cupColor).Render("☕"))
		} else if i < cupPosition {
			progressBar.WriteString(lipgloss.NewStyle().Foreground(teaColor).Render("●"))
		} else {
			progressBar.WriteString(lipgloss.NewStyle().Foreground(steamColor).Render("○"))
		}
	}
	progressBar.WriteString(" │\n")
	progressBar.WriteString("╰" + strings.Repeat("━", barWidth+2) + "╯")

	// Combine all components
	title := titleStyle.Render(teaCup)
	progress := fmt.Sprintf("%d%%", m.progress)
	timeInfo := timeStyle.Render(fmt.Sprintf(
		"Steeping time: %s / %s",
		m.elapsedTime.Truncate(time.Second).String(),
		m.totalTime.String(),
	))
	quit := quitStyle.Render("Press 'q' to cancel your tea break")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		progress,
		progressBar.String(),
		timeInfo,
		quit,
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
