package account

import (
	"fmt"
	"strings"
	"wk/pkg/ui"

	// "wk/pkg/ui"

	"github.com/brandur/wanikaniapi"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	// "github.com/charmbracelet/log"
)

type Model interface {
	tea.Model
	GetUser() *wanikaniapi.User
}

type model struct {
	spinner spinner.Model

	User *wanikaniapi.User
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
	case *wanikaniapi.User:
		m.User = msg
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	if m.User == nil {
		return m.spinner.View() + "loading..."
	}
	var b strings.Builder
	activeStatus := ui.BatsuMark
	if m.User.Data.Subscription.Active {
		activeStatus = ui.CheckMark
	}
	userDataContent := fmt.Sprintf(
		"Subscription\n"+
			"Active: %s\n"+
			"Type:   %s",
		activeStatus,
		m.User.Data.Subscription.Type,
	)
	userData := lipgloss.NewStyle().Render(userDataContent)
	b.WriteString(userData)
	return b.String()
}

func (m *model) GetUser() *wanikaniapi.User {
	return m.User
}
