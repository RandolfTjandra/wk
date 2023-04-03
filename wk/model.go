package main

import (
	"github.com/brandur/wanikaniapi"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"wk/pkg/account"
	"wk/pkg/levels"
	"wk/pkg/summary"
)

type mainModel struct {
	spinner spinner.Model

	navChoices  []PageView
	currentPage PageView
	cursors     map[PageView]int
	response    []byte
	err         error
	content     string

	User *wanikaniapi.User

	summary summary.Model
	levels  levels.Model
	account account.Model

	Reviews *wanikaniapi.ReviewPage

	Assignments []*wanikaniapi.Assignment

	Levels []*wanikaniapi.LevelProgression
}

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		account.GetUser,
		GetAssignments,
		m.spinner.Tick,
		m.summary.Init(),
		m.levels.Init(),
		m.account.Init(),
	)
}

func initialMainModel(
	view PageView,
) mainModel {
	s := spinner.New()
	s.Spinner = spinner.MiniDot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return mainModel{
		spinner:     s,
		currentPage: view,
		navChoices: []PageView{
			SummaryView,
			LevelsView,
			AccountView,
		},
		cursors: map[PageView]int{
			IndexView:   0,
			SummaryView: 0,
		},
		summary: summary.New(),
		levels:  levels.New(),
		account: account.New(),
	}
}

// Receives information and updates the mainModel
func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errMsg:
		m.err = msg
		return m, nil

	// key press
	case tea.KeyMsg:
		if msg.String() == "esc" {
			m.currentPage = IndexView
		}
		switch m.currentPage {
		case IndexView:
			return m.handleIndexKeyPress(msg)
		case SummaryView:
			_, cmd := m.summary.HandleSummaryKeyPress(msg)
			return m, cmd
		default:
			return m.handleDefaultKeyPress(msg)
		}

	case *wanikaniapi.Summary:
		m.summary.Update(msg)
	case summary.SummaryExpansion:
		m.summary.Update(msg)
	case *wanikaniapi.ReviewPage:
		m.Reviews = msg
	case []*wanikaniapi.Assignment:
		m.Assignments = msg
	case []*wanikaniapi.LevelProgression:
		m.levels.Update(msg)
	case *wanikaniapi.User:
		m.account.Update(msg)
	default:
		var cmd tea.Cmd
		var cmds []tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		//update all spinners
		cmds = append(cmds, cmd)

		s, cmd := m.summary.Update(msg)
		m.summary = s.(summary.Model)
		cmds = append(cmds, cmd)

		l, cmd := m.levels.Update(msg)
		m.levels = l.(levels.Model)
		cmds = append(cmds, cmd)

		return m, tea.Batch(cmds...)
	}

	return m, nil
}
