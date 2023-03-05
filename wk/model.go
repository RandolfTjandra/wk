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

	User *wanikaniapi.User

	Summary        *wanikaniapi.Summary
	SummaryLessons []*wanikaniapi.SummaryLesson
	SummaryReviews []*wanikaniapi.SummaryReview

	// The following maps represent SummaryLessons and SummaryReviews stacked on
	// top of each other. This is so that navigation can be managed by a single
	// slice
	SummaryExpansion map[int]bool
	SummarySubjects  map[int][]*wanikaniapi.Subject
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.commander.GetUser, tea.EnterAltScreen)
}

func initialModel(commander Commander, view PageView, subjectRepo db.SubjectRepo) model {
	return model{
		commander:   commander,
		currentPage: view,
		navChoices: []PageView{
			SummaryView,
			LessonsView,
			AssignmentsView,
			ReviewsView,
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
	case *wanikaniapi.User:
		m.User = msg
		return m, nil
	case *wanikaniapi.Summary:
		m.Summary = msg
		m.SummaryLessons = []*wanikaniapi.SummaryLesson{}
		// prep summary lessons
		for _, lesson := range m.Summary.Data.Lessons {
			if len(lesson.SubjectIDs) == 0 {
				continue
			}
			m.SummaryLessons = append(m.SummaryLessons, lesson)
		}
		// prep summary reviews
		m.SummaryReviews = []*wanikaniapi.SummaryReview{}
		for _, review := range m.Summary.Data.Reviews {
			if len(review.SubjectIDs) == 0 {
				continue
			}
			m.SummaryReviews = append(m.SummaryReviews, review)
		}
	case []*wanikaniapi.Subject:
		m.SummarySubjects[m.cursors[SummaryView]] = msg
	case errMsg:
		m.err = msg
		return m, tea.Quit
	// key press
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}

	return m, nil
}

func (m model) handleKeyPress(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "esc":
		m.currentPage = IndexView
		return m, nil
	case "q":
		return m, tea.Quit
	case "up", "k":
		if m.cursors[m.currentPage] > 0 {
			m.cursors[m.currentPage]--
		}
	case "down", "j":
		switch m.currentPage {
		case IndexView:
			if m.cursors[m.currentPage] < len(m.navChoices)-1 {
				m.cursors[m.currentPage]++
			}
		case SummaryView:
			if m.cursors[m.currentPage] < len(m.SummaryReviews)+len(m.SummaryLessons)-1 {
				m.cursors[m.currentPage]++
			}
		}
	case "enter":
		switch m.currentPage {
		case IndexView:
			m.currentPage = m.navChoices[m.cursors[IndexView]]
			// Set up what has to happen when a new page is selected
			switch m.currentPage {
			case SummaryView:
				return m, m.commander.GetSummary
			}
		case SummaryView:
			m.SummaryExpansion[m.cursors[SummaryView]] = !m.SummaryExpansion[m.cursors[SummaryView]]
			if m.SummaryExpansion[m.cursors[SummaryView]] {
				cursor := m.cursors[SummaryView]
				var subjectIDs []wanikaniapi.WKID
				if cursor < len(m.SummaryLessons) {
					subjectIDs = m.SummaryLessons[cursor].SubjectIDs
				} else {
					cursor = cursor - len(m.SummaryLessons)
					subjectIDs = m.SummaryReviews[cursor].SubjectIDs
				}
				return m, m.commander.GetSubjects(subjectIDs)
			}
		case AccountView:
			return m, nil
		}
	}
	return m, nil
}
