package wanikani

import (
	"context"
	"encoding/json"

	"github.com/brandur/wanikaniapi"
	_ "github.com/mattn/go-sqlite3"

	"wk/pkg/db"
)

// Return a single subject saved to db
func GetSubject(ctx context.Context,
	subjectRepo db.SubjectRepo,
	wkClient wanikaniapi.Client,
	subjectID wanikaniapi.WKID,
) (*wanikaniapi.Subject, error) {
	subjectRaw, err := subjectRepo.GetByID(ctx, int(subjectID))
	if err != nil { // get from api
		res, err := wkClient.SubjectGet(&wanikaniapi.SubjectGetParams{ID: &subjectID})
		if err != nil {
			return &wanikaniapi.Subject{}, err
		}
		marshalled, _ := json.Marshal(res)
		subjectRepo.Insert(int(subjectID), string(marshalled))
		return res, err
	}
	subject := wanikaniapi.Subject{}
	json.Unmarshal([]byte(subjectRaw), &subject)

	return &subject, nil
}
