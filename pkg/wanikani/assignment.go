package wanikani

import (
	"context"

	"github.com/brandur/wanikaniapi"
	_ "github.com/mattn/go-sqlite3"
)

// Get assignments
func GetAssignments(ctx context.Context, wkClient *wanikaniapi.Client) (*wanikaniapi.AssignmentPage, error) {
	res, err := wkClient.AssignmentList(&wanikaniapi.AssignmentListParams{})
	if err != nil {
		return nil, err
	}

	return res, nil
}
