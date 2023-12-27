package eventbridge

import (
	"context"
)

// Publisher is the interface to send messages event into the event bus.
type Publisher interface {
	// PublishEvent a message event to the event bus
	// It will block process until the message is queued into the bus
	// It returns an error if any.
	PublishEvent(context.Context, *Message) error
}

type PublisherFunc func(context.Context, *Message) error

func (f PublisherFunc) PublishEvent(ctx context.Context, message *Message) error {
	return f(ctx, message)
}

func NopePublisher() PublisherFunc {
	return func(context.Context, *Message) error {
		return nil
	}
}
