package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"wk/pkg/ui"
	"wk/pkg/wanikani"
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
)

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err)
	}
	var b strings.Builder
	switch m.currentPage {
	case IndexView:
		b.WriteString(indexView(m))
	case SummaryView:
		b.WriteString(summaryView(m))
	case LessonsView:
		b.WriteString(lessonsView(m))
	default:
		b.WriteString("incomplete: work in progress\n")
	}

	// The footer
	b.WriteString(ui.Subtle("  j/k, up/down: select") + ui.Dot +
		ui.Subtle("enter: choose") + ui.Dot + ui.Subtle("esc: back to index") + "\n")
	b.WriteString("\n  Press q to quit.\n")
	return b.String()
}

// Render index view
func indexView(m model) string {
	var b strings.Builder
	b.WriteString("                  ===== Menu =====\n\n")
	if m.User != nil {
		b.WriteString("                       Hello " + m.User.Data.Username + "\n")
		b.WriteString("  Level Progress:\n\n")
		b.WriteString("  " + ui.Progressbar(80, float64(m.User.Data.Level)/50) + "\n\n")
	}
	for i, choice := range m.navChoices {
		if i == m.cursor {
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
func assignmentsView(m model) string {
	if len(m.response) == 0 {
		return fmt.Sprintf("loading...")
	}
	return fmt.Sprintf("%s", string(m.response))
}

// Render summary view
func summaryView(m model) string {
	if m.Summary == nil {
		return "loading..."
	}
	var b strings.Builder
	b.WriteString("  ===== Summary =====\n\n")
	b.WriteString("  Lessons\n")
	for _, lesson := range m.Summary.Data.Lessons {
		b.WriteString(fmt.Sprintf("  Available at %s\n\n", lesson.AvailableAt.Format(time.RFC1123)))
		b.WriteString(fmt.Sprintf("    %d subjects:\n\n", len(lesson.SubjectIDs)))
		for i, subjectID := range lesson.SubjectIDs {
			subject, err := wanikani.GetSubject(context.Background(), subjectID)
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

// Render lessons view
func lessonsView(m model) string {
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
			subject, err := wanikani.GetSubject(context.Background(), subjectID)
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