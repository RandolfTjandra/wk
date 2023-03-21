package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/brandur/wanikaniapi"
	"github.com/charmbracelet/lipgloss"

	"wk/pkg/ui"
)

// Views
type PageView string

const (
	IndexView       PageView = "index"
	SummaryView     PageView = "summary"
	AssignmentsView PageView = "assignments"
	LevelsView      PageView = "levels"
	LessonsView     PageView = "lessons"
	ReviewsView     PageView = "reviews"
	SettingsView    PageView = "settings"
	AccountView     PageView = "account"
	DebuggingView   PageView = "debugging"
)

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err)
	}
	m.content = ""
	var b strings.Builder
	switch m.currentPage {
	case IndexView:
		m.content = lipgloss.NewStyle().
			Height(5).
			Render("Make selection from side menu")
	case SummaryView:
		m.content = m.summaryView()
	case ReviewsView:
		m.content = m.reviewsView()
	case AssignmentsView:
		m.content = m.assignmentsView()
	case AccountView:
		m.content = m.accountView()
	case LevelsView:
		m.content = m.levelsView()
	default:
		m.content = "incomplete: work in progress"
	}
	b.WriteString(m.indexView())

	// The footer
	b.WriteString(ui.Subtle("\n  j/k, up/down: select") + ui.Dot +
		ui.Subtle("enter: choose") + ui.Dot + ui.Subtle("esc: back to index") + "\n")
	b.WriteString("\n  Press q to quit.\n")
	return b.String()
}

// Renders a list of subjects
func (m model) renderSubjects(subjects []*wanikaniapi.Subject) string {
	var b strings.Builder
	subjectCount := 0
	for _, subject := range subjects {
		if subject.KanjiData != nil {
			b.WriteString(ui.Kanji(subject.KanjiData.Characters))
			if subjectCount < len(subjects)-1 {
				b.WriteString(", ")
			}
			subjectCount++
		} else if subject.VocabularyData != nil {
			b.WriteString(ui.Vocab(subject.VocabularyData.Characters))
			if subjectCount < len(subjects)-1 {
				b.WriteString(", ")
			}
			subjectCount++
		} else if subject.RadicalData != nil && subject.RadicalData.Characters != nil {
			b.WriteString(ui.Radical(*subject.RadicalData.Characters))
			if subjectCount < len(subjects)-1 {
				b.WriteString(", ")
			}
			subjectCount++
		} else {
			foo, _ := json.Marshal(subject)
			b.WriteString("\n\n" + string(foo) + "\n\n")
		}

		// hopefully can delete this if line wrap can work
		if subjectCount > 0 && subjectCount%10 == 0 {
			b.WriteString("\n  ")
		}
	}
	return b.String()
}
