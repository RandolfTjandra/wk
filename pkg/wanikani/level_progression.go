package wanikani

import (
	"context"

	"github.com/brandur/wanikaniapi"
)

func GetLevelProgressions(
	ctx context.Context,
	wkClient *wanikaniapi.Client,
) ([]*wanikaniapi.LevelProgression, error) {
	var levelProgressions []*wanikaniapi.LevelProgression
	err := wkClient.PageFully(func(id *wanikaniapi.WKID) (*wanikaniapi.PageObject, error) {
		listParams := wanikaniapi.LevelProgressionListParams{
			ListParams: wanikaniapi.ListParams{
				PageAfterID: id,
			},
		}
		page, err := wkClient.LevelProgressionList(&listParams)
		if err != nil {
			return nil, err
		}

		levelProgressions = append(levelProgressions, page.Data...)
		return &page.PageObject, nil
	})

	if err != nil {
		return nil, err
	}

	return levelProgressions, nil
}
