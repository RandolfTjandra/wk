package main

import (
	"fmt"
	"strings"

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

func (m mainModel) View() string {
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
		// m.content = m.summaryView()
		m.content = m.summary.View()
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
