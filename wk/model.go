package main

import (
	"github.com/brandur/wanikaniapi"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"wk/pkg/db"
)

type model struct {
	spinner spinner.Model

	commander Commander

	navChoices  []PageView
	currentPage PageView
	cursors     map[PageView]int
	response    []byte
	err         error
	content     string

	User *wanikaniapi.User

	Summary        *wanikaniapi.Summary
	SummaryLessons []*wanikaniapi.SummaryLesson
	SummaryReviews []*wanikaniapi.SummaryReview
	// The following maps represent SummaryLessons and SummaryReviews stacked on
	// top of each other. This is so that navigation can be managed by a single
	// slice
	SummaryExpansion map[int]bool
	SummarySubjects  map[int][]*wanikaniapi.Subject

	CurrentReviews *wanikaniapi.SummaryReview
	Reviews        *wanikaniapi.ReviewPage

	Assignments []*wanikaniapi.Assignment

	Levels []*wanikaniapi.LevelProgression
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		m.commander.GetUser,
		m.commander.GetSummary,
		m.commander.GetAssignments,
		m.spinner.Tick,
	)
}

func initialModel(commander Commander, view PageView, subjectRepo db.SubjectRepo) model {
	s := spinner.New()
	s.Spinner = spinner.MiniDot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{
		spinner:     s,
		commander:   commander,
		currentPage: view,
		navChoices: []PageView{
			SummaryView,
			LessonsView,
			ReviewsView,
			AssignmentsView,
			SettingsView,
			LevelsView,
			AccountView,
		},
		cursors: map[PageView]int{
			IndexView:   0,
			SummaryView: 0,
		},
		SummaryExpansion: make(map[int]bool),
		SummarySubjects:  make(map[int][]*wanikaniapi.Subject),
	}
}

// Receives information and updates the model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errMsg:
		m.err = msg
		return m, nil

	// key press
	case tea.KeyMsg:
		switch m.currentPage {
		case IndexView:
			return m.handleIndexKeyPress(msg)
		case SummaryView:
			return m.handleSummaryKeyPress(msg)
		default:
			return m.handleDefaultKeyPress(msg)
		}

	case *wanikaniapi.User:
		m.User = msg
		return m, m.spinner.Tick

	case *wanikaniapi.Summary:
		m.Summary = msg
		// Reset summary lessons
		m.SummaryLessons = []*wanikaniapi.SummaryLesson{}
		for _, lesson := range m.Summary.Data.Lessons {
			if len(lesson.SubjectIDs) == 0 {
				continue
			}
			m.SummaryLessons = append(m.SummaryLessons, lesson)
		}
		// Resest summary reviews
		m.SummaryReviews = []*wanikaniapi.SummaryReview{}
		for _, review := range m.Summary.Data.Reviews {
			if len(review.SubjectIDs) == 0 {
				continue
			}
			m.SummaryReviews = append(m.SummaryReviews, review)
		}
		m.CurrentReviews = m.SummaryReviews[0]
		m.SummaryReviews = m.SummaryReviews[1:]

	case *wanikaniapi.ReviewPage:
		m.Reviews = msg
	case []*wanikaniapi.Assignment:
		m.Assignments = msg
	case []*wanikaniapi.Subject:
		m.SummarySubjects[m.cursors[SummaryView]] = msg
	case []*wanikaniapi.LevelProgression:
		m.Levels = msg
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}
