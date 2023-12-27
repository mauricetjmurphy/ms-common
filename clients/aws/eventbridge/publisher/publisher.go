package publisher

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	eb "github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"

	"github.com/mauricetjmurphy/ms-common/clients/aws/eventbridge"
)

// Publisher is satisfied Publisher to send custom message to eventbus
type Publisher struct {
	eventbridge.Client
	*Config
}

// New creates the eventbus Publisher instance.
func New(client eventbridge.Client, opts ...Opts) eventbridge.Publisher {
	configs := &Config{
		encoder: DefaultEncode(),
		timeout: 30 * time.Second,
	}
	for _, opt := range opts {
		opt(configs)
	}
	return &Publisher{Client: client, Config: configs}
}

func (p *Publisher) PublishEvent(ctx context.Context, message *eventbridge.Message) error {
	newCtx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()
	details, err := p.encoder.Marshal(message.Detail)
	if err != nil {
		return err
	}
	_, err = p.Client.PutEvents(newCtx, &eb.PutEventsInput{
		Entries: []types.PutEventsRequestEntry{
			{
				Source:       aws.String(message.EventSource),
				DetailType:   aws.String(message.EventDetailType),
				EventBusName: aws.String(p.eventBusName),
				Detail:       aws.String(details),
			},
		},
	})
	return err
}
