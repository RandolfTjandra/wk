package main

import (
	"github.com/brandur/wanikaniapi"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	navChoices  []PageView
	currentPage PageView
	cursors     map[PageView]int
	response    []byte
	err         error

	User *wanikaniapi.User

	Summary          *wanikaniapi.Summary
	SummaryLessons   []*wanikaniapi.SummaryLesson
	SummaryReviews   []*wanikaniapi.SummaryReview
	SummaryExpansion map[int]bool
}

func (m model) Init() tea.Cmd {
	return tea.Batch(getUser, tea.EnterAltScreen)
}

func initialModel(view PageView) model {
	return model{
		currentPage: view,
		navChoices: []PageView{
			SummaryView,
			LessonsView,
			AssignmentsView,
			ReviewsView,
			SettingsView,
		},
		cursors: map[PageView]int{
			IndexView:   0,
			SummaryView: 0,
		},
		SummaryExpansion: make(map[int]bool),
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case *wanikaniapi.User:
		m.User = msg
		return m, nil
	case *wanikaniapi.Summary:
		m.Summary = msg
		i := 0
		// prep summary lessons
		for _, lesson := range m.Summary.Data.Lessons {
			if len(lesson.SubjectIDs) == 0 {
				continue
			}
			m.SummaryLessons = append(m.SummaryLessons, lesson)
			i++
		}
		// prep summary reviews
		for _, review := range m.Summary.Data.Reviews {
			if len(review.SubjectIDs) == 0 {
				continue
			}
			m.SummaryReviews = append(m.SummaryReviews, review)
			i++
		}
		return m, nil
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
			switch m.currentPage {
			case SummaryView:
				return m, getSummary
			}
		case SummaryView:
			m.SummaryExpansion[m.cursors[SummaryView]] = !m.SummaryExpansion[m.cursors[SummaryView]]
		}
	}
	return m, nil
}
