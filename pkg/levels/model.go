package levels

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/brandur/wanikaniapi"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"wk/pkg/ui"
)

type Model interface {
	tea.Model
	GetLevels() []*wanikaniapi.LevelProgression
}

type model struct {
	spinner spinner.Model

	Levels []*wanikaniapi.LevelProgression
}

func New() Model {
	s := spinner.New()
	s.Spinner = spinner.MiniDot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return &model{
		spinner: s,
	}
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case []*wanikaniapi.LevelProgression:
		m.Levels = msg
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	if m.Levels == nil {
		return m.spinner.View() + "loading..."
	}
	var b strings.Builder
	var days []float64
	for _, level := range m.Levels {
		pass := level.Data.PassedAt
		if pass != nil {
			days = append(days, math.Round(level.Data.PassedAt.Sub(*level.Data.StartedAt).Hours()/24))
		} else {
			if level.Data.UnlockedAt != nil {
				days = append(days, math.Round(time.Now().Sub(*level.Data.UnlockedAt).Hours()/24))
			}
		}
	}
	lengths := createBarLengths(days)
	for i, length := range lengths {
		// b.WriteString(fmt.Sprintf("%d: %s\n", level.Data.Level, strings.Repeat("*", int(length))))
		p := ui.LevelProgressBar(45, length, 1, int(days[i]))
		b.WriteString(fmt.Sprintf("%2d: %s\n", i+1, p))
	}

	return b.String() + "\n"
}

func createBarLengths(numbers []float64) []float64 {
	if len(numbers) == 0 {
		return []float64{}
	}
	var max_value, min_value float64
	max_value = numbers[0]
	min_value = numbers[0]

	// Find the minimum and maximum values in the input array
	for _, num := range numbers {
		if num > max_value {
			max_value = num
		}
		if num < min_value {
			min_value = num
		}
	}

	var bar_lengths []float64
	for _, num := range numbers {
		if max_value == min_value {
			bar_lengths = append(bar_lengths, 1)
		} else {
			log_value := math.Log(num - min_value + 1)
			log_range := math.Log(max_value - min_value + 1)
			normalized_value := log_value / log_range
			bar_length := normalized_value
			bar_lengths = append(bar_lengths, bar_length)
		}
	}
	return bar_lengths
}

func (m *model) GetLevels() []*wanikaniapi.LevelProgression {
	return m.Levels
}
