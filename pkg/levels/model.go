package levels

import (
	"fmt"
	"strings"

	"github.com/brandur/wanikaniapi"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
		//update both spinners
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	if m.Levels == nil {
		return "loading..." + m.spinner.View()
	}
	var b strings.Builder
	for _, level := range m.Levels {
		var passedString string
		pass := level.Data.PassedAt
		if pass != nil {
			passedString = pass.String()
		} else {
			passedString = "in progress"
		}
		b.WriteString(fmt.Sprintf("%d: %s\n", level.Data.Level, passedString))
	}

	return b.String() + "\n"
}

func (m *model) GetLevels() []*wanikaniapi.LevelProgression {
	return m.Levels
}
