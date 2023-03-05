package wanikani

import (
	"context"

	"github.com/brandur/wanikaniapi"
	_ "github.com/mattn/go-sqlite3"
)

// Return all voice actors
func GetVoiceActors(ctx context.Context,
	wkClient wanikaniapi.Client,
) (*wanikaniapi.VoiceActorPage, error) {
	res, err := wkClient.VoiceActorList(&wanikaniapi.VoiceActorListParams{})
	return res, err
}
