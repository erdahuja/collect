package form

import (
	"collect/foundation/event"
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

// EventSource represents the source of the given event.
const EventSource = "form"

// Set of user relatated events.
const (
	EventCreated = "FormCreated"
)

type EventParams struct {
	Form
}

// String returns a string representation of the event parameters.
func (p *EventParams) String() string {
	return fmt.Sprintf("&EventParams{Form:%v}", p.Form)
}

// Marshal returns the event parameters encoded as JSON.
func (p *EventParams) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

// UnmarshalCreated parses the event parameters from JSON.
func UnmarshalCreated(rawParams []byte) (*EventParams, error) {
	var params EventParams
	err := json.Unmarshal(rawParams, &params)
	if err != nil {
		return nil, fmt.Errorf("expected an encoded %T: %w", params, err)
	}

	return &params, nil
}

func (c *Core) registerEventHandlers(evnCore *event.Core) {
	evnCore.AddHandler(EventSource, EventCreated, c.handleFormCreatedEvent)
}

func (c *Core) handleFormCreatedEvent(ev event.Event, publisher *kafka.Writer) error {
	var params EventParams
	err := json.Unmarshal(ev.RawParams, &params)
	if err != nil {
		return fmt.Errorf("expected an encoded %T: %w", params, err)
	}

	c.log.Infow("form create event", "params", params)

	// publish on kafka
	if err := publisher.WriteMessages(context.TODO(), kafka.Message{
		Key:   []byte(fmt.Sprintf("%s:%s", ev.Source, ev.Type)),
		Value: ev.RawParams,
	}); err != nil {
		return fmt.Errorf("error while writing %T: %w", params, err)
	}

	return nil
}
