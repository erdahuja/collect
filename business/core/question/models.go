package question

import (
	"collect/foundation/event"
	"time"
)

// Question represents an individual user in db
type Question struct {
	ID           int64     `json:"question_id"`
	FormID       int64     `json:"form_id"`
	QuestionType string    `json:"question_type"`
	QuestionText string    `json:"question_text"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

// NewQuestion contains information needed to create a new Question.
type NewQuestion struct {
	FormID       int64  `json:"form_id"`
	QuestionType string `json:"question_type"`
	QuestionText string `json:"question_text"`
}

// CreatedEvent constructs an event for when a response is created.
func (uu NewQuestion) CreatedEvent(res Question) event.Event {
	params := EventParams{
		Question: res,
	}

	rawParams, err := params.Marshal()
	if err != nil {
		panic(err)
	}

	return event.Event{
		Source:    EventSource,
		Type:      EventCreated,
		RawParams: rawParams,
	}
}

type DeleteQuestion struct {
	ID int64 `json:"question_id"`
}
