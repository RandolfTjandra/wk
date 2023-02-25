package main

import (
	"context"
	"encoding/json"

	"github.com/brandur/wanikaniapi"
	tea "github.com/charmbracelet/bubbletea"

	"wk/pkg/wanikani"
)

type wkRes struct {
	httpStatus int
	resBody    int
}

func getUser() tea.Msg {
	user, err := wanikani.GetUser(context.Background())
	if err != nil {
		return errMsg{err}
	}
	return user
}

func getAssignments() tea.Msg {
	assignmentPage, err := wanikani.GetAssignments(context.Background())
	if err != nil {
		// There was an error making our request. Wrap the error we received
		// in a message and return it.
		return errMsg{err}
	}

	return assignmentPage
}

func getSummary() tea.Msg {
	if false {
		example := []byte(`{
	"object": "report",
	"url": "https://api.wanikani.com/v2/summary",
	"data_updated_at": "2018-04-11T21:00:00.000000Z",
	"data": {
		"lessons": [
			{
				"available_at": "2018-04-11T21:00:00.000000Z",
				"subject_ids": [
					25,
					26
				]
			}
		],
		"next_reviews_at": "2018-04-11T21:00:00.000000Z",
		"reviews": [
			{
				"available_at": "2018-04-11T21:00:00.000000Z",
				"subject_ids": [
					21,
					23,
					24
				]
			},
			{
				"available_at": "2018-04-11T22:00:00.000000Z",
				"subject_ids": []
			}
		]
	}
	}`)
		foo := wanikaniapi.Summary{}
		err := json.Unmarshal(example, &foo)
		if err == nil {
			return &foo
		} else {
			return errMsg{err}
		}
	} else {
		// real. above is debugging with sample
		wkClient := wanikaniapi.NewClient(&wanikaniapi.ClientConfig{
			APIToken: wanikani.ApiKey,
		})

		res, err := wkClient.SummaryGet(&wanikaniapi.SummaryGetParams{})
		if err != nil {
			// There was an error making our request. Wrap the error we received
			// in a message and return it.
			return errMsg{err}
		}

		return res
	}
}

type errMsg struct{ err error }

// For messages that contain errors it's often handy to also implement the
// error interface on the message.
func (e errMsg) Error() string { return e.err.Error() }
