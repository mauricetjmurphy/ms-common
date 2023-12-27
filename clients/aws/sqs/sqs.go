package sqs

import (
	"context"
	"fmt"

	awsv2 "github.com/aws/aws-sdk-go-v2/aws"
	awssqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/pkg/errors"

	"github.com/mauricetjmurphy/ms-common/clients/aws"
)

//go:generate mockery --output sqsmocks --outpkg sqsmocks --name Client
type Client interface {
	// ReceiveMessages retrieves one or more messages from the specified queue.
	ReceiveMessages(ctx context.Context, maxReceivedMessages int) (*awssqs.ReceiveMessageOutput, error)
	// DeleteMessage deletes message that were received from the queue
	DeleteMessage(ctx context.Context, message types.Message) error
	// SendMessage sends message the specified queue
	SendMessage(ctx context.Context, message *awssqs.SendMessageInput) (*awssqs.SendMessageOutput, error)
}

// sqsClient is a wrapper for the aws sqs client.
type sqsClient struct {
	*awssqs.Client
	*ClientParams
}

type ClientParams struct {
	// Region is the AWS region where the SQS queue is hosted.
	Region string
	// QueueUrl is the name of an SQS queue.
	QueueURL string
	// MaxReceivedMessage is the max number of messaged to be polled.
	MaxReceivedMessage int
}

func New(ctx context.Context, params *ClientParams) (Client, error) {
	if params == nil {
		return nil, errors.New("sqs : failed to miss client params configs")
	}
	if params.QueueURL == "" {
		return nil, errors.New("sqs : failed to miss QueueUrl params configs")
	}

	conf, err := aws.LoadConfig(ctx, params.Region)
	if err != nil {
		return nil, errors.Wrap(err, "sqs : failed to load aws configs")
	}
	return &sqsClient{
		Client:       awssqs.NewFromConfig(conf),
		ClientParams: params,
	}, nil
}

func (c *sqsClient) ReceiveMessages(ctx context.Context, maxReceivedMessages int) (*awssqs.ReceiveMessageOutput, error) {
	if maxReceivedMessages < 1 || maxReceivedMessages > c.MaxReceivedMessage {
		return nil, fmt.Errorf("maxReceivedMessages must be from 1 to %d", c.MaxReceivedMessage)
	}
	input := &awssqs.ReceiveMessageInput{
		QueueUrl:            awsv2.String(c.QueueURL),
		MaxNumberOfMessages: int32(maxReceivedMessages),
	}
	output, err := c.Client.ReceiveMessage(ctx, input)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (c *sqsClient) DeleteMessage(ctx context.Context, message types.Message) error {
	params := &awssqs.DeleteMessageInput{
		ReceiptHandle: awsv2.String(*message.ReceiptHandle),
		QueueUrl:      awsv2.String(c.QueueURL),
	}
	_, err := c.Client.DeleteMessage(ctx, params)
	return err
}

func (c *sqsClient) SendMessage(ctx context.Context, params *awssqs.SendMessageInput) (*awssqs.SendMessageOutput, error) {
	return c.Client.SendMessage(ctx, params)
}
