# ☕ TBreak - Terminal Coffee Break Timer

A delightful, coffee-themed terminal break timer built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Features

- 🎨 Beautiful coffee-themed ASCII art animation
- ⏲️ Customizable break duration
- 🎯 Progress tracking with animated progress bar
- ☁️ Animated steam effects
- 📝 Custom task/break naming
- 🎭 Dynamic status messages

## Installation

### Prerequisites

- Go 1.16 or higher

### Building from source

```bash
git clone https://github.com/yourusername/tbreak
cd tbreak
go build
```

## Usage

Run the timer with a specified duration:

```bash
# Run with duration in seconds
./tbreak 30s "Coffee Break"

# Run with duration in minutes
./tbreak 5m "Quick Break"
```

### Command-line Arguments

- First argument: Duration (required)
  - Format: `<number>s` for seconds or `<number>m` for minutes
  - Example: `30s`, `5m`
- Second argument: Task name (optional)
  - Default: "Coffee Brewing"
  - Example: "Coffee Break", "Quick Rest"

### Controls

- Press `q` to quit the timer at any time

## Example

```bash
./tbreak 5m "Afternoon Coffee"
```

This will start a 5-minute timer with "Afternoon Coffee" as the break name, displaying a charming coffee-themed interface with:
- Animated coffee cup with steam
- Progress bar showing completion
- Time remaining
- Current brewing status
