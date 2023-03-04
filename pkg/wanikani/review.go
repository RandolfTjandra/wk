package wanikani

import (
	"context"

	"github.com/brandur/wanikaniapi"
)

func GetReviews(ctx context.Context, wkClient *wanikaniapi.Client) (*wanikaniapi.ReviewPage, error) {
	res, err := wkClient.ReviewList(&wanikaniapi.ReviewListParams{})
	return res, err
}
