package main

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/acarl005/stripansi"
	"github.com/charmbracelet/lipgloss"

	"wk/pkg/ui"
	"wk/pkg/wanikani"
)

// Render index view
func (m mainModel) indexView() string {
	var b strings.Builder
	b.WriteString(ui.H1Title.Render("~Wk~"))
	b.WriteString("\n")

	// Summary header
	summaryHeader := m.renderSummaryHeader()

	// User greeting
	if m.account.GetUser() != nil {
		greeting := fmt.Sprintf("Hello %s..", m.account.GetUser().Data.Username)
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
		bar := lipgloss.NewStyle().Width(ui.UIWidth - len(barTitle) - 2*ui.UIXMargin).
			Align(lipgloss.Right).Render(
			ui.Progressbar(50, float64(m.account.GetUser().Data.Level), 50),
		)
		loadingBar := lipgloss.JoinHorizontal(0, barTitle, bar)

		b.WriteString(lipgloss.NewStyle().Margin(0, ui.UIXMargin).Padding(0, ui.UIXMargin).Render(loadingBar))
		b.WriteString("\n\n")
	}

	// Assignments
	assignments := m.renderAssignments()
	b.WriteString(lipgloss.NewStyle().Width(ui.UIWidth).Align(lipgloss.Center).Render(assignments))
	b.WriteString("\n")

	// Navigation
	renderedChoices := m.renderNavigation()
	choicesLen := utf8.RuneCountInString(stripansi.Strip(strings.Split(renderedChoices, "\n")[0])) + 1
	renderedContent := lipgloss.NewStyle().
		Width(ui.UIWidth-choicesLen).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("87")).
		MarginLeft(1).Padding(0, 2, 0, 1).
		Render(m.content)

	main := lipgloss.JoinHorizontal(0, renderedChoices, renderedContent)
	b.WriteString(main)
	b.WriteString("\n\n")
	if m.err != nil {
		b.WriteString(fmt.Sprintf("%#v\n", m.err))
	}
	return b.String()
}

func (m mainModel) renderNavigation() string {
	// Render navigation
	var choices strings.Builder
	for i, choice := range m.navChoices {
		if i == m.cursors[IndexView] {
			cursor := ">"
			if m.currentPage == IndexView {
				cursor = ui.Keyword(">")
			}
			choices.WriteString(fmt.Sprintf("%s %s", cursor, choice))
		} else {
			choices.WriteString(fmt.Sprintf("  %s", choice))
		}
		if i < len(m.navChoices)-1 {
			choices.WriteString("\n")
		}
	}
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("87")).
		MarginLeft(ui.UIXMargin).Padding(0, 2, 0, 1).
		Render(choices.String())
}

// Render summary which contains remaining lessons and reviews
func (m mainModel) renderSummaryHeader() string {
	if m.summary.GetSummary() == nil {
		return "loading summary.."
	} else {
		totalLessons := 0
		for _, lesson := range m.summary.GetSummaryLessons() {
			totalLessons += len(lesson.SubjectIDs)
		}
		// The first review contains all the reviews currently available.
		// Subsequent reviews contain reviews that will become available in the future.
		totalReviews := len(m.summary.GetCurrentReviews().SubjectIDs)
		return fmt.Sprintf("%d lessons %d reviews", totalLessons, totalReviews)
	}
}

func (m mainModel) renderAssignments() string {
	if len(m.Assignments) == 0 {
		return "loading assignments\n"
	}
	var b strings.Builder
	assignmentsSorted := wanikani.ClassifyAssignments(m.Assignments)
	b.WriteString(string(wanikani.Apprentice) + ": " + ui.Apprentice(fmt.Sprintf("%d ", len(assignmentsSorted[wanikani.Apprentice]))))
	b.WriteString(string(wanikani.Guru) + ": " + ui.Guru(fmt.Sprintf("%d ", len(assignmentsSorted[wanikani.Guru]))))
	b.WriteString(string(wanikani.Master) + ": " + ui.Master(fmt.Sprintf("%d ", len(assignmentsSorted[wanikani.Master]))))
	b.WriteString(string(wanikani.Enlightened) + ": " + ui.Enlightened(fmt.Sprintf("%d\n", len(assignmentsSorted[wanikani.Enlightened]))))

	return b.String()
}
