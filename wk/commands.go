package main

import (
	"context"

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

// return *wanikaniapi.ReviewPage
func GetReviews(reviewIDs ...wanikaniapi.WKID) func() tea.Msg {
	return func() tea.Msg {
		reviews, err := wanikani.GetReviews(context.Background(), reviewIDs...)
		if err != nil {
			return errMsg{err}
		}
		return reviews
	}
}

// TODO: should be moved to assignment commands when that's created
// return []*wanikaniapi.Assignment
func GetAssignments() tea.Msg {
	assignments, err := wanikani.GetAssignments(context.Background())
	if err != nil {
		return errMsg{err}
	}

	return assignments
}

// // return *wanikaniapi.VoiceActorPage
// func (c *commander) GetVoiceActors() tea.Msg {
// 	return nil
// }
