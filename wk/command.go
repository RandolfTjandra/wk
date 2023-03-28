package main

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
	GetUser() tea.Msg
	GetSubjects(subjectIDs []wanikaniapi.WKID) func() tea.Msg
	GetReviews(reviewIDs ...wanikaniapi.WKID) func() tea.Msg
	GetAssignments() tea.Msg
	GetVoiceActors() tea.Msg
	GetLevelProgressions() tea.Msg
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
	} else {
		return mockCommander{}
	}
}

// return *wanikaniapi.User
func (c *commander) GetUser() tea.Msg {
	user, err := wanikani.GetUser(context.Background(), c.wanikaniClient)
	if err != nil {
		return errMsg{err}
	}
	return user
}

// return *wanikaniapi.ReviewPage
func (c *commander) GetReviews(reviewIDs ...wanikaniapi.WKID) func() tea.Msg {
	return func() tea.Msg {
		reviews, err := wanikani.GetReviews(context.Background(), c.wanikaniClient, reviewIDs...)
		if err != nil {
			return errMsg{err}
		}
		return reviews
	}
}

// return []*wanikaniapi.Assignment
func (c *commander) GetAssignments() tea.Msg {
	assignments, err := wanikani.GetAssignments(context.Background(), c.wanikaniClient)
	if err != nil {
		return errMsg{err}
	}

	return assignments
}

// return []*wanikaniapi.Subject
func (c *commander) GetSubjects(subjectIDs []wanikaniapi.WKID) func() tea.Msg {
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
		return subjects
	}
}

// return *wanikaniapi.VoiceActorPage
func (c *commander) GetVoiceActors() tea.Msg {
	return nil
}

// return []*wanikaniapi.GetLevelProgression
func (c *commander) GetLevelProgressions() tea.Msg {
	progressions, err := wanikani.GetLevelProgressions(context.Background(), c.wanikaniClient)
	if err != nil {
		return errMsg{err}
	}

	return progressions
}
