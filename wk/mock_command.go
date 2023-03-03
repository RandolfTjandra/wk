package main

import (
	"encoding/json"

	"github.com/brandur/wanikaniapi"
	tea "github.com/charmbracelet/bubbletea"
)

type mockCommander struct {
}

func (c mockCommander) GetUser() tea.Msg {
	example := []byte(`{
  "object": "user",
  "url": "https://api.wanikani.com/v2/user",
  "data_updated_at": "2018-04-06T14:26:53.022245Z",
  "data": {
    "id": "5a6a5234-a392-4a87-8f3f-33342afe8a42",
    "username": "example_user",
    "level": 5,
    "profile_url": "https://www.wanikani.com/users/example_user",
    "started_at": "2012-05-11T00:52:18.958466Z",
    "current_vacation_started_at": null,
    "subscription": {
      "active": true,
      "type": "recurring",
      "max_level_granted": 60,
      "period_ends_at": "2018-12-11T13:32:19.485748Z"
    },
    "preferences": {
      "default_voice_actor_id": 1,
      "lessons_autoplay_audio": false,
      "lessons_batch_size": 10,
      "lessons_presentation_order": "ascending_level_then_subject",
      "reviews_autoplay_audio": false,
      "reviews_display_srs_indicator": true
    }
  }
}`)
	resource := wanikaniapi.User{}
	err := json.Unmarshal(example, &resource)
	if err == nil {
		return &resource
	} else {
		return errMsg{err}
	}
}

func (c mockCommander) GetSummary() tea.Msg {
	example := []byte(`{
	"object": "report",
	"url": "https://api.wanikani.com/v2/summary",
	"data_updated_at": "2018-04-11T21:00:00.000000Z",
	"data": {
		"lessons": [
			{
				"available_at": "2018-04-11T21:00:00.000000Z",
				"subject_ids": [
					25,
					26
				]
			}
		],
		"next_reviews_at": "2018-04-11T21:00:00.000000Z",
		"reviews": [
			{
				"available_at": "2018-04-11T21:00:00.000000Z",
				"subject_ids": [
					21,
					23,
					24
				]
			},
			{
				"available_at": "2018-04-11T22:00:00.000000Z",
				"subject_ids": []
			}
		]
	}
	}`)
	resource := wanikaniapi.Summary{}
	err := json.Unmarshal(example, &resource)
	if err == nil {
		return &resource
	} else {
		return errMsg{err}
	}
}

func (c mockCommander) GetSubjects(subjectIDs []wanikaniapi.WKID) func() tea.Msg {
	example := []byte(`{
  "id": 1,
  "object": "radical",
  "url": "https://api.wanikani.com/v2/subjects/1",
  "data_updated_at": "2018-03-29T23:13:14.064836Z",
  "data": {
    "amalgamation_subject_ids": [
      5,
      4,
      98
    ],
    "auxiliary_meanings": [
      {
        "meaning": "ground",
        "type": "blacklist"
      }
    ],
    "characters": "一",
    "character_images": [
      {
        "url": "https://cdn.wanikani.com/images/legacy/576-subject-1-without-css-original.svg?1520987227",
        "metadata": {
          "inline_styles": false
        },
        "content_type": "image/svg+xml"
      }
    ],
    "created_at": "2012-02-27T18:08:16.000000Z",
    "document_url": "https://www.wanikani.com/radicals/ground",
    "hidden_at": null,
    "lesson_position": 1,
    "level": 1,
    "meanings": [
      {
        "meaning": "Ground",
        "primary": true,
        "accepted_answer": true
      }
    ],
    "meaning_mnemonic": "This radical consists of a single, horizontal stroke. What's the biggest, single, horizontal stroke? That's the ground. Look at the <radical>ground</radical>, look at this radical, now look at the ground again. Kind of the same, right?",
    "slug": "ground",
    "spaced_repetition_system_id": 2
  }
}`)
	resource := wanikaniapi.Subject{}
	err := json.Unmarshal(example, &resource)

	return func() tea.Msg {
		if err == nil {
			return []*wanikaniapi.Subject{&resource}
		} else {
			return errMsg{err}
		}
	}
}
