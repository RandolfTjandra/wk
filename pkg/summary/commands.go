package summary

import (
	"context"
	"log"

	"github.com/brandur/wanikaniapi"
	tea "github.com/charmbracelet/bubbletea"

	"wk/pkg/wanikani"
)

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

type wkRes struct {
	httpStatus int
	resBody    int
}

// return *wanikaniapi.Summary
func GetSummary() tea.Msg {
	summary, err := wanikani.GetSummary(context.Background())
	if err != nil {
		return errMsg{err}
	}

	return summary
}

type SummaryExpansion struct {
	cursor   int
	subjects []*wanikaniapi.Subject
}

// return []*wanikaniapi.Subject
func GetSubjects(cursor int, subjectIDs []wanikaniapi.WKID) func() tea.Msg {
	return func() tea.Msg {
		subjectCount := 0
		subjects := []*wanikaniapi.Subject{}
		for _, subjectID := range subjectIDs {
			if subjectCount == 100 { // cap to avoid rate limiting
				break
			}
			subject, err := wanikani.GetSubject(context.Background(), *wanikani.Client, subjectID)
			if err != nil {
				log.Print("\n  skipped due to error: " + err.Error() + "\n")
				continue
			}
			if subject.KanjiData != nil ||
				subject.VocabularyData != nil ||
				(subject.RadicalData != nil && subject.RadicalData.Characters != nil) {
				subjects = append(subjects, subject)
				subjectCount++
			}
		}

		return SummaryExpansion{
			cursor:   cursor,
			subjects: subjects,
		}
	}
}
