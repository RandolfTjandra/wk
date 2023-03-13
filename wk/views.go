package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"wk/pkg/ui"

	"github.com/brandur/wanikaniapi"
	"github.com/charmbracelet/lipgloss"
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
		b.WriteString(m.indexView())
	case SummaryView:
		m.content = m.summaryView()
		b.WriteString(m.indexView())
	case ReviewsView:
		m.content = m.reviewsView()
		b.WriteString(m.indexView())
	case AssignmentsView:
		m.content = m.assignmentsView()
		b.WriteString(m.indexView())
	case AccountView:
		m.content = m.accountView()
		b.WriteString(m.indexView())
	default:
		b.WriteString("incomplete: work in progress\n")
	}
	// b.WriteString("debugging: " + m.content)

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
	b.WriteString("\n")

	summaryHeader := ""

	if m.Summary == nil {
		summaryHeader = "loading summary.."
	} else {
		totalLessons := 0
		for _, lesson := range m.SummaryLessons {
			totalLessons += len(lesson.SubjectIDs)
		}
		// The first review contains all the reviews currently available.
		// Subsequent reviews contain reviews that will become available in the future.
		totalReviews := len(m.SummaryReviews[0].SubjectIDs)
		summaryHeader = fmt.Sprintf("%d lessons %d reviews", totalLessons, totalReviews)
	}

	if m.User != nil {
		greeting := fmt.Sprintf("Hello %s..", m.User.Data.Username)
		remainingWidth := ui.UIWidth - len(greeting)
		summaryHeader = lipgloss.NewStyle().Width(remainingWidth).
			AlignHorizontal(lipgloss.Right).
			Render(summaryHeader)

		greeting = lipgloss.JoinHorizontal(0, greeting, summaryHeader)
		greeting = lipgloss.NewStyle().Width(ui.UIWidth).
			Margin(0, ui.UIXMargin).Render(greeting)
		b.WriteString(greeting)
		b.WriteString("\n\n")

		barTitle := lipgloss.NewStyle().Render("Level Progress:")
		bar := lipgloss.NewStyle().Width(ui.UIWidth - len(barTitle) - 2*ui.UIXMargin).Align(lipgloss.Right).Render(
			ui.Progressbar(50, float64(m.User.Data.Level), 50),
		)
		loadingBar := lipgloss.JoinHorizontal(0, barTitle, bar)

		b.WriteString(lipgloss.NewStyle().Margin(0, ui.UIXMargin).Padding(0, ui.UIXMargin).Render(loadingBar))
		b.WriteString("\n\n")
	}

	var choices strings.Builder
	for i, choice := range m.navChoices {
		if i == m.cursors[IndexView] {
			cursor := ui.Keyword(">")
			choices.WriteString(fmt.Sprintf("%s %s", cursor, choice))
		} else {
			choices.WriteString(fmt.Sprintf("  %s", choice))
		}
		if i < len(m.navChoices)-1 {
			choices.WriteString("\n")
		}
	}
	renderedChoices := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("87")).
		MarginLeft(ui.UIXMargin).Padding(0, 2, 0, 1).
		Render(choices.String())

	main := lipgloss.JoinHorizontal(0, renderedChoices, m.content)
	b.WriteString(main)
	b.WriteString("\n\n")
	return b.String()
}

// Render assignments view
func (m model) assignmentsView() string {
	if m.Assignments == nil {
		return fmt.Sprintf("loading...")
	}
	var b strings.Builder
	b.WriteString("  Assignments:\n")
	tally := map[string]int{}
	for _, assignment := range m.Assignments {
		var classification string
		switch stage := assignment.Data.SRSStage; {
		case stage < 5:
			classification = "apprentice"
		case stage < 7:
			classification = "guru"
		case stage < 8:
			classification = "master"
		case stage < 9:
			classification = "enlightened"
		}
		tally[classification] += 1
	}
	b.WriteString(fmt.Sprintf("Apprentice:  %d\n", tally["apprentice"]))
	b.WriteString(fmt.Sprintf("Guru:        %d\n", tally["guru"]))
	b.WriteString(fmt.Sprintf("Master:      %d\n", tally["master"]))
	b.WriteString(fmt.Sprintf("Enlightened: %d\n", tally["enlightened"]))

	return b.String() + "\n"
}

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
	b.WriteString("  Reviews:")
	for i, review := range m.SummaryReviews {
		if i == 0 {
			continue
		}
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

// Render reviews view
func (m model) reviewsView() string {
	if m.Reviews == nil {
		return "loading..."
	}
	var b strings.Builder
	b.WriteString("  Reviews:")
	for i, review := range m.Reviews.Data {
		b.WriteString(fmt.Sprintf("%d: %d\n", i, review.Data.SubjectID))
	}

	return b.String() + "\n"
}

func (m model) accountView() string {
	var b strings.Builder
	activeStatus := ui.BatsuMark
	if m.User.Data.Subscription.Active {
		activeStatus = ui.CheckMark
	}
	userDataContent := fmt.Sprintf(
		"Subscription\n"+
			"Active: %s\n"+
			"Type:   %s",
		activeStatus,
		m.User.Data.Subscription.Type,
	)
	userData := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("87")).
		MarginLeft(1).Padding(0, 2, 0, 1).
		Render(userDataContent)
	b.WriteString(userData)
	return b.String()
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
