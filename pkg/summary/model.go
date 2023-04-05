package summary

import (
	"fmt"
	"strings"
	"time"
	"wk/pkg/ui"

	"github.com/brandur/wanikaniapi"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model interface {
	tea.Model
	HandleSummaryKeyPress(key tea.KeyMsg) (tea.Model, tea.Cmd)
	GetSummary() *wanikaniapi.Summary
	GetSummaryLessons() []*wanikaniapi.SummaryLesson
	GetCurrentReviews() *wanikaniapi.SummaryReview
}

type model struct {
	spinner spinner.Model
	cursor  int

	Summary        *wanikaniapi.Summary
	SummaryLessons []*wanikaniapi.SummaryLesson
	CurrentReviews *wanikaniapi.SummaryReview
	SummaryReviews []*wanikaniapi.SummaryReview
	// The following maps represent SummaryLessons and SummaryReviews stacked on
	// top of each other. This is so that navigation can be managed by a single
	// slice
	SummaryExpansion map[int]bool
	SummarySubjects  map[int][]*wanikaniapi.Subject
}

func New() Model {
	s := spinner.New()
	s.Spinner = spinner.MiniDot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return &model{
		spinner:          s,
		cursor:           0,
		SummaryExpansion: make(map[int]bool),
		SummarySubjects:  make(map[int][]*wanikaniapi.Subject),
	}
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(GetSummary, m.spinner.Tick)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case *wanikaniapi.Summary:
		m.Summary = msg
		// Reset summary lessons
		m.SummaryLessons = []*wanikaniapi.SummaryLesson{}
		for _, lesson := range m.Summary.Data.Lessons {
			if len(lesson.SubjectIDs) == 0 {
				continue
			}
			m.SummaryLessons = append(m.SummaryLessons, lesson)
		}
		// Resest summary reviews
		m.SummaryReviews = []*wanikaniapi.SummaryReview{}
		for _, review := range m.Summary.Data.Reviews {
			if len(review.SubjectIDs) == 0 {
				continue
			}
			m.SummaryReviews = append(m.SummaryReviews, review)
		}
		m.CurrentReviews = m.SummaryReviews[0]
		m.SummaryReviews = m.SummaryReviews[1:]
	case SummaryExpansion:
		m.SummarySubjects[msg.cursor] = msg.subjects
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		//update both spinners
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	if m.Summary == nil {
		return "loading..."
	}
	var b strings.Builder
	b.WriteString("  Lessons:")
	if len(m.SummaryLessons) == 0 {
		b.WriteString("\n    No lessons!")
	}

	for i, lesson := range m.SummaryLessons {
		if m.cursor == i {
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
				b.WriteString(" ")
			} else {
				b.WriteString("\n  " + ui.RenderSubjects(m.SummarySubjects[i]))
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
	if len(m.SummaryReviews) == 0 {
		b.WriteString("\n    No Reviews! Good job!")
	}
	for i, review := range m.SummaryReviews {
		// the first review just contains current reviews
		if m.cursor == i+len(m.SummaryLessons) {
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
				b.WriteString("\n  " + ui.RenderSubjects(m.SummarySubjects[i+len(m.SummaryLessons)]))
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
	b.WriteString(fmt.Sprintf("\ncursor pos: %d", m.cursor))

	return b.String() + "\n"
}

func (m *model) HandleSummaryKeyPress(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "q":
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.SummaryReviews)+len(m.SummaryLessons)-1 {
			m.cursor++
		}
	case "enter":
		cursor := m.cursor
		m.SummaryExpansion[cursor] = !m.SummaryExpansion[cursor]
		if m.SummaryExpansion[cursor] {
			var subjectIDs []wanikaniapi.WKID
			if cursor < len(m.SummaryLessons) {
				subjectIDs = m.SummaryLessons[cursor].SubjectIDs
			} else {
				localCursor := cursor - len(m.SummaryLessons)
				subjectIDs = m.SummaryReviews[localCursor].SubjectIDs
			}
			return m, GetSubjects(cursor, subjectIDs)
		}
	}
	return m, nil
}

func (m *model) GetSummary() *wanikaniapi.Summary {
	return m.Summary
}

func (m *model) GetSummaryLessons() []*wanikaniapi.SummaryLesson {
	return m.SummaryLessons
}

func (m *model) GetCurrentReviews() *wanikaniapi.SummaryReview {
	return m.CurrentReviews
}
