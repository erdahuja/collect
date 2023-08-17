package answer

import (
	"collect/business/core/response"
	"collect/foundation/event"
	"time"
)

// Answer represents an individual answer
type Answer struct {
	AnswerID       int64
	QuestionID     int64
	ResponseID     int64
	AnswerText     *string
	AnswerOptionID *int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// NewAnswer is what we require from clients when adding a Answer to db.
type NewAnswer struct {
	QuestionID     int64  `json:"question_id" validate:"required"`
	AnswerText     string `json:"answer_text,omitempty"`
	AnswerOptionID int    `json:"answer_option_id,omitempty"`
}

// CreatedEvent constructs an event for when a user is created.
func (uu NewAnswer) CreatedEvent(ans Answer, res response.Response) event.Event {
	params := EventParamsCreated{
		Answer:   ans,
		Response: res,
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

// DeleteAnswer is what we require from clients when deleting a response.
type DeleteAnswer struct {
	AnswerID int64 `json:"answer_id"`
}
