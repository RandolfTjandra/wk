package wanikani

import "github.com/brandur/wanikaniapi"

var (
	Client *wanikaniapi.Client
)

func Init(apiKey string) {
	Client = wanikaniapi.NewClient(&wanikaniapi.ClientConfig{
		APIToken: apiKey,
	})
}
