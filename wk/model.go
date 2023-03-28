package main

import (
	"github.com/brandur/wanikaniapi"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"wk/pkg/db"
	"wk/pkg/levels"
	"wk/pkg/summary"
)

type mainModel struct {
	spinner spinner.Model

	commander Commander

	navChoices  []PageView
	currentPage PageView
	cursors     map[PageView]int
	response    []byte
	err         error
	content     string

	User *wanikaniapi.User

	summary summary.Model
	levels  levels.Model

	Reviews *wanikaniapi.ReviewPage

	Assignments []*wanikaniapi.Assignment

	Levels []*wanikaniapi.LevelProgression
}

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		m.commander.GetUser,
		m.commander.GetAssignments,
		m.spinner.Tick,
		m.summary.Init(),
		m.levels.Init(),
	)
}

func initialMainModel(
	commander Commander,
	summaryCommander summary.Commander,
	view PageView,
	subjectRepo db.SubjectRepo,
) mainModel {
	s := spinner.New()
	s.Spinner = spinner.MiniDot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	summaryModel := summary.New(summaryCommander)
	levelsModel := levels.New()
	return mainModel{
		spinner:     s,
		commander:   commander,
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
		summary: summaryModel,
		levels:  levelsModel,
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

	case *wanikaniapi.User:
		m.User = msg
		return m, nil

	case *wanikaniapi.Summary:
		m.summary.Update(msg)
	case summary.SummaryExpansion:
		m.summary.Update(msg)

	case *wanikaniapi.ReviewPage:
		m.Reviews = msg
	case []*wanikaniapi.Assignment:
		m.Assignments = msg
	case []*wanikaniapi.LevelProgression:
		m.Levels = msg
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
