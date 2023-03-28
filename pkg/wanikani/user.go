package wanikani

import (
	"context"
	"encoding/json"

	"github.com/brandur/wanikaniapi"
	redis "github.com/redis/go-redis/v9"
)

const UserKey = "user"

func GetUser(ctx context.Context) (*wanikaniapi.User, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     RedisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	val, err := rdb.Get(ctx, UserKey).Result()
	if err != nil { // get from api
		res, err := Client.UserGet(&wanikaniapi.UserGetParams{})
		marshalled, _ := json.Marshal(res)
		rdb.Set(ctx, UserKey, marshalled, 0)
		return res, err
	}
	resource := wanikaniapi.User{}
	json.Unmarshal([]byte(val), &resource)

	return &resource, nil
}
