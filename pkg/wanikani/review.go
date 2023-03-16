package wanikani

import (
	"context"

	"github.com/brandur/wanikaniapi"
)

func GetReviews(ctx context.Context,
	wkClient *wanikaniapi.Client,
	assignmentIDs ...wanikaniapi.WKID,
) (*wanikaniapi.ReviewPage, error) {
	params := wanikaniapi.ReviewListParams{}
	if len(assignmentIDs) > 0 {
		params.AssignmentIDs = assignmentIDs
	}
	res, err := wkClient.ReviewList(&params)
	return res, err
}
