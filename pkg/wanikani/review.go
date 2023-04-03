package wanikani

import (
	"context"

	"github.com/brandur/wanikaniapi"
)

func GetReviews(ctx context.Context,
	assignmentIDs ...wanikaniapi.WKID,
) (*wanikaniapi.ReviewPage, error) {
	params := wanikaniapi.ReviewListParams{}
	if len(assignmentIDs) > 0 {
		params.AssignmentIDs = assignmentIDs
	}
	res, err := Client.ReviewList(&params)
	return res, err
}
