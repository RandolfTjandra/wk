package main

import (
	"github.com/brandur/wanikaniapi"
	tea "github.com/charmbracelet/bubbletea"

	"wk/pkg/db"
)

type model struct {
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
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		m.commander.GetUser,
		m.commander.GetSummary,
		m.commander.GetAssignments,
	)
}

func initialModel(commander Commander, view PageView, subjectRepo db.SubjectRepo) model {
	return model{
		commander:   commander,
		currentPage: view,
		navChoices: []PageView{
			SummaryView,
			LessonsView,
			ReviewsView,
			AssignmentsView,
			SettingsView,
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
		return m, nil

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
	}

	return m, nil
}
