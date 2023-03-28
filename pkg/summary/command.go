package summary

import (
	"context"
	"log"

	"github.com/brandur/wanikaniapi"
	tea "github.com/charmbracelet/bubbletea"

	"wk/pkg/db"
	"wk/pkg/wanikani"
)

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

type wkRes struct {
	httpStatus int
	resBody    int
}

type Commander interface {
	GetSummary() tea.Msg
	GetSubjects(cursor int, subjectIDs []wanikaniapi.WKID) func() tea.Msg
}

type commander struct {
	subjectRepo    db.SubjectRepo
	wanikaniClient *wanikaniapi.Client
}

func NewCommander(live bool, subjectRepo db.SubjectRepo, wanikaniClient *wanikaniapi.Client) Commander {
	if live {
		return &commander{
			subjectRepo:    subjectRepo,
			wanikaniClient: wanikaniClient,
		}
	}
	return nil
}

// return *wanikaniapi.Summary
func (c *commander) GetSummary() tea.Msg {
	summary, err := wanikani.GetSummary(context.Background(), c.wanikaniClient)
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
func (c *commander) GetSubjects(cursor int, subjectIDs []wanikaniapi.WKID) func() tea.Msg {
	return func() tea.Msg {
		subjectCount := 0
		subjects := []*wanikaniapi.Subject{}
		for _, subjectID := range subjectIDs {
			if subjectCount == 100 { // cap to avoid rate limiting
				break
			}
			subject, err := wanikani.GetSubject(context.Background(), c.subjectRepo, *c.wanikaniClient, subjectID)
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
