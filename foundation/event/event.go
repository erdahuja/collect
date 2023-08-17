// Package event provides business access to events in the system.
package event

import (
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// Core manages the set of APIs for event access.
type Core struct {
	log       *zap.SugaredLogger
	publisher *kafka.Writer
	handlers  map[string]map[string][]HandleFunc
}

// NewCore constructs a core for event api access.
func NewCore(log *zap.SugaredLogger, publisher *kafka.Writer) *Core {
	return &Core{
		log:       log,
		publisher: publisher,
		handlers:  map[string]map[string][]HandleFunc{},
	}
}

// SendEvent sends event to all handlers registered for the specified event.
func (c *Core) SendEvent(event Event) error {
	c.log.Infow("sendevent", "status", "started", "event", event.String())
	defer c.log.Infow("sendevent", "status", "completed")

	if m, ok := c.handlers[event.Source]; ok {
		if hfs, ok := m[event.Type]; ok {
			for _, hf := range hfs {
				c.log.Infow("sendevent", "status", "sending")

				if err := hf(event, c.publisher); err != nil {
					c.log.Infof("sendevent", "ERROR", err)
				}
			}
		}
	}
	return nil
}

// AddHandler add handler to specific event from specific source.
func (c *Core) AddHandler(source, t string, f HandleFunc) {
	ss, ok := c.handlers[source]
	if !ok {
		ss = map[string][]HandleFunc{}
	}

	ss[t] = append(ss[t], f)
	c.handlers[source] = ss
}
