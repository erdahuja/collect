package form

import (
	"collect/foundation/event"
	"time"
)

// Form represents an individual form.
type Form struct {
	ID          int64     `json:"form_id"`
	Title       string    `json:"form_title"`
	Description string    `json:"form_description"`
	DateCreated time.Time `json:"created_at"`
	DateUpdated time.Time `json:"updated_at"`
}

// Form is what we require from clients when adding a Form.
type NewForm struct {
	Title       string `json:"form_title"`
	Description string `json:"form_description"`
}

// UpdateForm have fields which clients can send just the fields to update
type UpdateForm struct {
	Title       *string `json:"form_title,omitempty" validate:"omitempty"`
	Description *string `json:"form_description,omitempty" validate:"omitempty"`
}

// CreatedEvent constructs an event for when a user is created.
func (uu NewForm) CreatedEvent(res Form) event.Event {
	params := EventParams{
		Form: res,
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
