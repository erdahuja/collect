package answer

import (
	"context"
	"encoding/json"
	"fmt"

	"collect/business/core/response"
	"collect/foundation/event"

	"github.com/segmentio/kafka-go"
)

// EventSource represents the source of the given event.
const EventSource = "answer"

// Set of user relatated events.
const (
	EventCreated = "AnswerCreated"
)

// EventParamsCreated is the event parameters for the created event.
type EventParamsCreated struct {
	Answer Answer
	Response response.Response
}

func (c *Core) registerEventHandlers(evnCore *event.Core) {
	evnCore.AddHandler(EventSource, EventCreated, c.handleAnswerCreatedEvent)
}

func (c *Core) handleAnswerCreatedEvent(ev event.Event, publisher *kafka.Writer) error {
	var params EventParamsCreated
	err := json.Unmarshal(ev.RawParams, &params)
	if err != nil {
		return fmt.Errorf("expected an encoded %T: %w", params, err)
	}

	c.log.Infow("answer create event", "id", params.Answer.AnswerID)

	// publish on kafka
	publisher.WriteMessages(context.TODO(), kafka.Message{
		Key:   []byte(fmt.Sprintf("%s:%s", ev.Source, ev.Type)),
		Value: ev.RawParams,
	})

	return nil
}

// Marshal returns the event parameters encoded as JSON.
func (p *EventParamsCreated) Marshal() ([]byte, error) {
	return json.Marshal(p)
}
