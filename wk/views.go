package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"wk/pkg/ui"
	"wk/pkg/wanikani"

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
	case LessonsView:
		b.WriteString(m.lessonsView())
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
	// Render lessons
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
			b.WriteString("\n  " + m.renderSubjects(lesson.SubjectIDs))
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
			b.WriteString("\n  " + m.renderSubjects(review.SubjectIDs))
		}
	}

	return b.String() + "\n"
}

// Render lessons view
func (m model) lessonsView() string {
	if m.Summary == nil {
		return "loading..."
	}
	var b strings.Builder
	b.WriteString("  ===== Lessons =====\n\n")
	b.WriteString("  Lessons\n")
	for _, lesson := range m.Summary.Data.Lessons {
		b.WriteString(fmt.Sprintf("  Available at %s\n\n", lesson.AvailableAt.Format(time.RFC1123)))
		b.WriteString(fmt.Sprintf("    %d subjects:\n\n", len(lesson.SubjectIDs)))
		for i, subjectID := range lesson.SubjectIDs {
			subject, err := wanikani.GetSubject(context.Background(), m.subjectRepo, subjectID)
			if err != nil {
				b.WriteString("\n  skipped " + string(subject.GetObject().ObjectType) + "due to error: " + err.Error() + "\n")
			} else if subject.KanjiData != nil {
				b.WriteString(fmt.Sprintf("  %s, ", subject.KanjiData.Characters))
			} else if subject.VocabularyData != nil {
				b.WriteString(fmt.Sprintf("  %s, ", subject.VocabularyData.Characters))
			} else {
				b.WriteString("\n  skipped " + string(subject.GetObject().ObjectType) + "\n")
			}
			if i > 0 && i%10 == 0 {
				b.WriteString("\n")
			}
		}
	}
	return b.String() + "\n"
}

// gets and renders a list of subjects
func (m model) renderSubjects(subjectIDs []wanikaniapi.WKID) string {
	var b strings.Builder
	subjectCount := 0
	for _, subjectID := range subjectIDs {
		if subjectCount == 100 {
			b.WriteString(fmt.Sprintf("+ %d more\n", len(subjectIDs)-100))
			break
		}
		subject, err := wanikani.GetSubject(context.Background(), m.subjectRepo, subjectID)
		if err != nil {
			b.WriteString("\n  skipped due to error: " + err.Error() + "\n")
		} else if subject.KanjiData != nil {
			b.WriteString(ui.Kanji(subject.KanjiData.Characters))
			if subjectCount < len(subjectIDs)-1 {
				b.WriteString(", ")
			}
			subjectCount++
		} else if subject.VocabularyData != nil {
			b.WriteString(ui.Vocab(subject.VocabularyData.Characters))
			if subjectCount < len(subjectIDs)-1 {
				b.WriteString(", ")
			}
			subjectCount++
		} else if subject.RadicalData != nil && subject.RadicalData.Characters != nil {
			b.WriteString(ui.Radical(*subject.RadicalData.Characters))
			if subjectCount < len(subjectIDs)-1 {
				b.WriteString(", ")
			}
			subjectCount++
		} else {
			foo, _ := json.Marshal(subject)
			b.WriteString("\n\n" + string(foo) + "\n\n")
		}
		if subjectCount > 0 && subjectCount%10 == 0 { // hopefully can delete this if line wrap can work
			b.WriteString("\n  ")
		}
	}
	return b.String()
}
