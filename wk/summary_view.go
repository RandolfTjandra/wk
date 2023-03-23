package main

import (
	"fmt"
	"strings"
	"time"

	"wk/pkg/ui"
)

// Render summary view
func (m model) summaryView() string {
	if m.Summary == nil {
		return "loading..."
	}
	var b strings.Builder
	// b.WriteString(ui.H1Title.Render("Summary"))
	// b.WriteString("\n\n")
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
			if len(m.SummarySubjects[i]) == 0 {
				b.WriteString(" " + m.spinner.View())
			} else {
				b.WriteString("\n  " + m.renderSubjects(m.SummarySubjects[i]))
				renderedCount := len(m.SummarySubjects[i])
				count := len(m.SummaryLessons[i].SubjectIDs)
				missingCount := count - renderedCount
				if missingCount > 0 {
					b.WriteString(fmt.Sprintf("+%d more", missingCount))
				}
			}
		}
	}
	b.WriteString("\n\n")

	// Render reviews
	b.WriteString("  Reviews:")
	for i, review := range m.SummaryReviews {
		// the first review just contains current reviews
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
			if len(m.SummarySubjects[i+len(m.SummaryLessons)]) == 0 {
				b.WriteString(" " + m.spinner.View())
			} else {
				b.WriteString("\n  " + m.renderSubjects(m.SummarySubjects[i+len(m.SummaryLessons)]))
				renderedCount := len(m.SummarySubjects[i+len(m.SummaryLessons)])
				count := len(m.SummaryReviews[i].SubjectIDs)
				missingCount := count - renderedCount
				if missingCount > 0 {
					b.WriteString(fmt.Sprintf("+%d more", missingCount))
				}
			}
		}
	}
	// debugging
	b.WriteString(fmt.Sprintf("\ncursor pos: %d", m.cursors[SummaryView]))

	return b.String() + "\n"
}
