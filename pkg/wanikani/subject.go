package wanikani

import (
	"context"
	"encoding/json"

	"github.com/brandur/wanikaniapi"
	_ "github.com/mattn/go-sqlite3"

	"wk/pkg/db"
)

// Return a single subject cached
func GetSubject(ctx context.Context, subjectRepo db.SubjectRepo, subjectID wanikaniapi.WKID) (*wanikaniapi.Subject, error) {
	subjectRaw, err := subjectRepo.GetByID(ctx, int(subjectID))
	if err != nil { // get from api
		wkClient := wanikaniapi.NewClient(&wanikaniapi.ClientConfig{
			APIToken: ApiKey,
		})

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

// Get assignments
func GetAssignments(ctx context.Context) (*wanikaniapi.AssignmentPage, error) {
	wkClient := wanikaniapi.NewClient(&wanikaniapi.ClientConfig{
		APIToken: ApiKey,
	})
	res, err := wkClient.AssignmentList(&wanikaniapi.AssignmentListParams{})
	if err != nil {
		return nil, err
	}

	return res, nil
}
