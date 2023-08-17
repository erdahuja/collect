package event

import (
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

// HandleFunc represents a function that can receive an event.
type HandleFunc func(Event, *kafka.Writer) error

// Event represents an event between core domains.
type Event struct {
	Source    string
	Type      string
	RawParams []byte
}

func (e Event) Bytes() []byte {
	b, err := json.Marshal(e)
	if err != nil {
		return make([]byte, 0)
	}
	return b
}

// String implements the Stringer interface.
func (e Event) String() string {
	return fmt.Sprintf(
		"Event{Source:%#v, Type:%#v, RawParams:%#v}",
		e.Source, e.Type, string(e.RawParams),
	)
}
