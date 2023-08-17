package response

import (
	"collect/foundation/event"
	"time"
)

// Response represents an individual form in the response.
type Response struct {
	ResponseID   int64     `json:"response_id,omitempty"`
	FormID       int64     `json:"form_id,omitempty"`
	RespondentID string    `json:"respondent_id,omitempty"` // some unique identifier so that one respondent can have one response in one form, this can be some id also
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

// NewResponse is what we require from clients when adding a Response to db.
type NewResponse struct {
	FormID       int64  `json:"form_id" validate:"required"`
	RespondentID string `json:"respondent_id" validate:"required"`
}

// DeleteResponse is what we require from clients when deleting a response.
type DeleteResponse struct {
	DeleteResponse int64 `json:"response_id"`
}

// CreatedEvent constructs an event for when a response is created.
func (uu NewResponse) CreatedEvent(res Response) event.Event {
	params := EventParams{
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
