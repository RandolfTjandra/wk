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

func (c mockCommander) GetReviews() tea.Msg {
	example := []byte(`{
  "object": "collection",
  "url": "https://api.wanikani.com/v2/reviews",
  "pages": {
    "per_page": 1000,
    "next_url": "https://api.wanikani.com/v2/reviews?page_after_id=534345",
    "previous_url": null
  },
  "total_count": 19201,
  "data_updated_at": "2017-12-20T01:10:17.578705Z",
  "data": [
    {
      "id": 534342,
      "object": "review",
      "url": "https://api.wanikani.com/v2/reviews/534342",
      "data_updated_at": "2017-12-20T01:00:59.255427Z",
      "data": {
        "created_at": "2017-12-20T01:00:59.255427Z",
        "assignment_id": 32132,
        "spaced_repetition_system_id": 1,
        "subject_id": 8,
        "starting_srs_stage": 4,
        "ending_srs_stage": 2,
        "incorrect_meaning_answers": 1,
        "incorrect_reading_answers": 0
      }
    }
  ]
}`)
	resource := wanikaniapi.ReviewPage{}
	err := json.Unmarshal(example, &resource)

	return func() tea.Msg {
		if err == nil {
			return &resource
		} else {
			return errMsg{err}
		}
	}
}

func (c mockCommander) GetAssignments() tea.Msg {
	example := []byte(`{
  "object": "collection",
  "url": "https://api.wanikani.com/v2/assignments",
  "pages": {
    "per_page": 500,
    "next_url": "https://api.wanikani.com/v2/assignments?page_after_id=80469434",
    "previous_url": null
  },
  "total_count": 1600,
  "data_updated_at": "2017-11-29T19:37:03.571377Z",
  "data": [
    {
      "id": 80463006,
      "object": "assignment",
      "url": "https://api.wanikani.com/v2/assignments/80463006",
      "data_updated_at": "2017-10-30T01:51:10.438432Z",
      "data": {
        "created_at": "2017-09-05T23:38:10.695133Z",
        "subject_id": 8761,
        "subject_type": "radical",
        "srs_stage": 8,
        "unlocked_at": "2017-09-05T23:38:10.695133Z",
        "started_at": "2017-09-05T23:41:28.980679Z",
        "passed_at": "2017-09-07T17:14:14.491889Z",
        "burned_at": null,
        "available_at": "2018-02-27T00:00:00.000000Z",
        "resurrected_at": null
      }
    }
  ]
}`)
	resource := wanikaniapi.AssignmentPage{}
	err := json.Unmarshal(example, &resource)

	return func() tea.Msg {
		if err == nil {
			return &resource
		} else {
			return errMsg{err}
		}
	}
}

func (c mockCommander) GetVoiceActors() tea.Msg {
	example := []byte(`{
  "object": "collection",
  "url": "https://api.wanikani.com/v2/voice_actors",
  "pages": {
    "per_page": 500,
    "next_url": null,
    "previous_url": null
  },
  "total_count": 2,
  "data_updated_at": "2017-11-29T19:37:03.571377Z",
  "data": [
    {
      "id": 234,
      "object": "voice_actor",
      "url": "https://api.wanikani.com/v2/voice_actors/1",
      "data_updated_at": "2017-12-20T00:24:47.048380Z",
      "data": {
        "created_at": "2017-12-20T00:03:56.642838Z",
        "name": "Kyoko",
        "gender": "female",
        "description": "Tokyo accent"
      }
    }
  ]
}`)
	resource := wanikaniapi.VoiceActorPage{}
	err := json.Unmarshal(example, &resource)

	return func() tea.Msg {
		if err == nil {
			return &resource
		} else {
			return errMsg{err}
		}
	}
}
