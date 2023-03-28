package wanikani

import (
	"context"

	"github.com/brandur/wanikaniapi"
	_ "github.com/mattn/go-sqlite3"
)

type SRSStageName string

const (
	Apprentice  SRSStageName = "apprentice"
	Guru        SRSStageName = "guru"
	Master      SRSStageName = "master"
	Enlightened SRSStageName = "enlightened"
	Burned      SRSStageName = "burned"
)

// Get assignments
func GetAssignments(ctx context.Context) ([]*wanikaniapi.Assignment, error) {
	isBurned := false
	isStarted := true
	isHidden := false

	var assignments []*wanikaniapi.Assignment
	err := Client.PageFully(func(id *wanikaniapi.WKID) (*wanikaniapi.PageObject, error) {
		listParams := wanikaniapi.AssignmentListParams{
			Burned:  &isBurned,
			Started: &isStarted,
			Hidden:  &isHidden,
			ListParams: wanikaniapi.ListParams{
				PageAfterID: id,
			},
		}
		page, err := Client.AssignmentList(&listParams)
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

func ClassifyAssignments(assignments []*wanikaniapi.Assignment) map[SRSStageName][]*wanikaniapi.Assignment {
	assignmentMap := make(map[SRSStageName][]*wanikaniapi.Assignment)
	for _, assignment := range assignments {
		switch stage := assignment.Data.SRSStage; {
		case stage < 5:
			assignmentMap[Apprentice] = append(assignmentMap[Apprentice], assignment)
		case stage < 7:
			assignmentMap[Guru] = append(assignmentMap[Guru], assignment)
		case stage < 8:
			assignmentMap[Master] = append(assignmentMap[Master], assignment)
		case stage < 9:
			assignmentMap[Enlightened] = append(assignmentMap[Enlightened], assignment)
		}
	}
	return assignmentMap
}
