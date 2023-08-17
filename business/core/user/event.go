package user

import (
	"encoding/json"
	"fmt"
)

// EventSource represents the source of the given event.
const EventSource = "user"

// Set of user relatated events.
const (
	EventCreated = "UserCreated"
	EventQuery   = "UserQuery"
)

// =============================================================================

// EventParams is the event parameters for the updated event.
type EventParams struct {
	User
}

// String returns a string representation of the event parameters.
func (p *EventParams) String() string {
	return fmt.Sprintf("&EventParams{User:%v}", p.User)
}

// Marshal returns the event parameters encoded as JSON.
func (p *EventParams) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

// UnmarshalUpdated parses the event parameters from JSON.
func UnmarshalUpdated(rawParams []byte) (*EventParams, error) {
	var params EventParams
	err := json.Unmarshal(rawParams, &params)
	if err != nil {
		return nil, fmt.Errorf("expected an encoded %T: %w", params, err)
	}

	return &params, nil
}
