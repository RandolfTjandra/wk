package wanikani

import (
	"context"
	"encoding/json"
	"time"

	"github.com/brandur/wanikaniapi"
	redis "github.com/redis/go-redis/v9"
)

func GetSummary(ctx context.Context, wkClient *wanikaniapi.Client) (*wanikaniapi.Summary, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     RedisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	val, err := rdb.Get(ctx, "summary").Result()
	if err != nil { // get from api
		res, err := Client.SummaryGet(&wanikaniapi.SummaryGetParams{})
		marshalled, _ := json.Marshal(res)
		rdb.Set(ctx, "summary", marshalled, 1*time.Hour)
		return res, err
	}
	resource := wanikaniapi.Summary{}
	json.Unmarshal([]byte(val), &resource)

	return &resource, nil
}
