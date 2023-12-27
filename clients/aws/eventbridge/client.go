package eventbridge

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	eb "github.com/aws/aws-sdk-go-v2/service/eventbridge"
)

// Client is the interface to wrap necessary AWS EventBridge APIs.
//
//go:generate mockery --output eventbridgemocks --outpkg eventbridgemocks --name Client
type Client interface {

	// PutEvents sends custom events to Amazon EventBridge matching rules.
	// Return events  output and an error if any.
	PutEvents(context.Context, *eb.PutEventsInput) (*eb.PutEventsOutput, error)
}

// clientImpl to satisfy the Client interface.
type clientImpl struct {
	*eb.Client
}

// New creates the EventBridge Client instance.
func New(cfg aws.Config, optFns ...func(*eb.Options)) Client {
	return &clientImpl{
		Client: eb.NewFromConfig(cfg, optFns...),
	}
}

func (c *clientImpl) PutEvents(ctx context.Context, params *eb.PutEventsInput) (*eb.PutEventsOutput, error) {
	return c.Client.PutEvents(ctx, params)
}
