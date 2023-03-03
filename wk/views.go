package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"wk/pkg/ui"

	"github.com/brandur/wanikaniapi"
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
	DebuggingView   PageView = "debugging"
)

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err)
	}
	var b strings.Builder
	switch m.currentPage {
	case IndexView:
		b.WriteString(m.indexView())
	case SummaryView:
		b.WriteString(m.summaryView())
	default:
		b.WriteString("incomplete: work in progress\n")
	}

	// The footer
	b.WriteString(ui.Subtle("\n  j/k, up/down: select") + ui.Dot +
		ui.Subtle("enter: choose") + ui.Dot + ui.Subtle("esc: back to index") + "\n")
	b.WriteString("\n  Press q to quit.\n")
	return b.String()
}

// Render index view
func (m model) indexView() string {
	var b strings.Builder
	b.WriteString(ui.H1Title.Render("Wk"))
	b.WriteString("\n\n")
	if m.User != nil {
		b.WriteString("  Hello " + m.User.Data.Username + "\n")
		b.WriteString("  Level Progress:\n\n")
		b.WriteString("  " + ui.Progressbar(80, float64(m.User.Data.Level), 50) + "\n\n")
	}
	for i, choice := range m.navChoices {
		if i == m.cursors[IndexView] {
			cursor := ui.Keyword(">")
			b.WriteString(fmt.Sprintf("  %s %s\n", cursor, choice))
		} else {
			b.WriteString(fmt.Sprintf("    %s\n", choice))
		}
	}
	b.WriteString("\n\n")
	return b.String()
}

// Render assignments view
func (m model) assignmentsView() string {
	if len(m.response) == 0 {
		return fmt.Sprintf("loading...")
	}
	return fmt.Sprintf("%s", string(m.response))
}

// Render summary view
func (m model) summaryView() string {
	if m.Summary == nil {
		return "loading..."
	}
	var b strings.Builder
	b.WriteString(ui.H1Title.Render("Summary"))
	b.WriteString("\n\n")
	b.WriteString("  Lessons:")
	for i, lesson := range m.SummaryLessons {
		if m.cursors[SummaryView] == i {
			cursor := ui.Keyword(">")
			b.WriteString(fmt.Sprintf("\n  %s ", cursor))
		} else {
			b.WriteString("\n    ")
		}
		b.WriteString(fmt.Sprintf("%d available on %s at %s",
			len(lesson.SubjectIDs),
			lesson.AvailableAt.Weekday().String(),
			lesson.AvailableAt.Local().Format(time.TimeOnly),
		))
		if m.SummaryExpansion[i] {
			b.WriteString("\n  " + m.renderSubjects(m.SummarySubjects[i]))
			renderedCount := len(m.SummarySubjects[i])
			count := len(m.SummaryLessons[i].SubjectIDs)
			missingCount := count - renderedCount
			if missingCount > 0 {
				b.WriteString(fmt.Sprintf("+%d more", missingCount))
			}
		}
	}
	b.WriteString("\n\n")

	// Render reviews
	b.WriteString("  Upcoming Reviews:")
	for i, review := range m.SummaryReviews {
		if m.cursors[SummaryView] == i+len(m.SummaryLessons) {
			cursor := ui.Keyword(">")
			b.WriteString(fmt.Sprintf("\n  %s ", cursor))
		} else {
			b.WriteString("\n    ")
		}
		b.WriteString(fmt.Sprintf("%d available on %s at %s",
			len(review.SubjectIDs),
			review.AvailableAt.Weekday().String(),
			review.AvailableAt.Local().Format(time.TimeOnly),
		))
		if m.SummaryExpansion[i+len(m.SummaryLessons)] {
			b.WriteString("\n  " + m.renderSubjects(m.SummarySubjects[i+len(m.SummaryLessons)]))
			renderedCount := len(m.SummarySubjects[i+len(m.SummaryLessons)])
			count := len(m.SummaryReviews[i].SubjectIDs)
			missingCount := count - renderedCount
			if missingCount > 0 {
				b.WriteString(fmt.Sprintf("+%d more", missingCount))
			}
		}
	}

	return b.String() + "\n"
}

// renders a list of subjects
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
