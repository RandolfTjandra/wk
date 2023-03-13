package wanikani

import (
	"context"

	"github.com/brandur/wanikaniapi"
	_ "github.com/mattn/go-sqlite3"
)

// Get assignments
func GetAssignments(ctx context.Context, wkClient *wanikaniapi.Client) ([]*wanikaniapi.Assignment, error) {
	isBurned := false
	isStarted := true
	isHidden := false

	var assignments []*wanikaniapi.Assignment
	err := wkClient.PageFully(func(id *wanikaniapi.WKID) (*wanikaniapi.PageObject, error) {
		listParams := wanikaniapi.AssignmentListParams{
			Burned:  &isBurned,
			Started: &isStarted,
			Hidden:  &isHidden,
			ListParams: wanikaniapi.ListParams{
				PageAfterID: id,
			},
		}
		page, err := wkClient.AssignmentList(&listParams)
		if err != nil {
			return nil, err
		}

		assignments = append(assignments, page.Data...)
		return &page.PageObject, nil
	})

	if err != nil {
		return nil, err
	}

	return assignments, nil
}
