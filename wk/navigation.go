package main

import (
	"github.com/brandur/wanikaniapi"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) handleIndexKeyPress(key tea.KeyMsg) (tea.Model, tea.Cmd) {
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
		m.currentPage = m.navChoices[m.cursors[IndexView]]
		// Set up what has to happen when a new page is selected
		switch m.currentPage {
		case SummaryView:
			return m, m.commander.GetSummary
		case ReviewsView:
			assignmentIDs := []wanikaniapi.WKID{}
			for _, assignment := range m.Assignments {
				assignmentIDs = append(assignmentIDs, assignment.ID)
			}
			return m, m.commander.GetReviews(assignmentIDs...)
		case AssignmentsView:
			return m, m.commander.GetAssignments
		}
	}
	return m, nil
}

func (m model) handleSummaryKeyPress(key tea.KeyMsg) (tea.Model, tea.Cmd) {
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
	}
	return m, nil
}

func (m model) handleDefaultKeyPress(key tea.KeyMsg) (tea.Model, tea.Cmd) {
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
	}
	return m, nil
}
