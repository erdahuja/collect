package user

import (
	"collect/foundation/event"
	"time"
)

// User represents an individual user in db
type User struct {
	ID           int64     `json:"user_id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Roles        []Role    `json:"roles"`
	PasswordHash []byte    `json:"-"`
	Active       bool      `json:"active"`
	DateCreated  time.Time `json:"dateCreated"`
	DateUpdated  time.Time `json:"dateUpdated"`
}

// NewUser contains information needed to create a new User.
type NewUser struct {
	Name            string `json:"name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Roles           []Role `json:"roles" validate:"required"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"passwordConfirm" validate:"eqfield=Password"`
}

// CreatedEvent constructs an event for when a user is created.
func (uu NewUser) CreatedEvent(user User) event.Event {
	params := EventParams{
		User: user,
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

// QueryEvent constructs an event for when a user is created.
func (uu User) QueryEvent() event.Event {
	params := EventParams{
		User: uu,
	}

	rawParams, err := params.Marshal()
	if err != nil {
		panic(err)
	}

	return event.Event{
		Source:    EventSource,
		Type:      EventQuery,
		RawParams: rawParams,
	}
}
