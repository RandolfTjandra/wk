package wanikani

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/brandur/wanikaniapi"
	redis "github.com/redis/go-redis/v9"
)

// Return a single subject cached
func GetSubject(ctx context.Context, subjectID wanikaniapi.WKID) (*wanikaniapi.Subject, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     RedisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	val, err := rdb.Get(ctx, strconv.Itoa(int(subjectID))).Result()
	if err != nil { // get from api
		wkClient := wanikaniapi.NewClient(&wanikaniapi.ClientConfig{
			APIToken: ApiKey,
		})
		log.Println("get subject from api: " + strconv.Itoa(int(subjectID)))

		res, err := wkClient.SubjectGet(&wanikaniapi.SubjectGetParams{ID: &subjectID})
		marshalled, _ := json.Marshal(res)
		rdb.Set(ctx, strconv.Itoa(int(subjectID)), marshalled, 0)
		return res, err
	}
	subject := wanikaniapi.Subject{}
	json.Unmarshal([]byte(val), &subject)

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
