package secrets

import (
	"context"
	"encoding/json"

	awsv2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/pkg/errors"

	"github.com/mauricetjmurphy/ms-common/clients/aws"
)

//go:generate mockery --output secretsmocks --outpkg secretsmocks --name Client
type Client interface {
	// GetValue retrieves the target secret value on given secret identifier.
	GetValue(ctx context.Context, secretID string, target interface{}) error
	// GetSecret retrieves Secret on given secret identifier.
	GetSecret(ctx context.Context, secretID string) (*Secret, error)
	// GetAzureSecret retrieves AzureAPISecret on given secret identifier.
	GetAzureSecret(ctx context.Context, secretID string) (*AzureAPISecret, error)
}

// secretsClient is a wrapper for the aws secrets manager.
type secretsClient struct {
	*secretsmanager.Client
}

func NewFromConfig(cfg awsv2.Config, optFns ...func(*secretsmanager.Options)) Client {
	return &secretsClient{
		Client: secretsmanager.NewFromConfig(cfg, optFns...),
	}
}

func New(ctx context.Context, region string) (Client, error) {
	conf, err := aws.LoadConfig(ctx, region)
	if err != nil {
		return nil, errors.Wrap(err, "secrets : failed to load aws configs")
	}
	return &secretsClient{
		secretsmanager.NewFromConfig(conf),
	}, nil
}

func (c *secretsClient) GetSecret(ctx context.Context, secretID string) (*Secret, error) {
	secret := &Secret{}
	err := c.GetValue(ctx, secretID, secret)
	if err != nil {
		return nil, err
	}
	return secret, nil
}

func (c *secretsClient) GetValue(ctx context.Context, secretID string, target interface{}) error {
	if len(secretID) == 0 {
		return errors.New("secrets: secretId cannot be empty")
	}
	params := &secretsmanager.GetSecretValueInput{
		SecretId: awsv2.String(secretID),
	}
	result, err := c.GetSecretValue(ctx, params)
	if err != nil {
		return errors.Wrap(err, "secrets: failed GetSecretValue")
	}
	if result == nil {
		return errors.Wrap(err, "secrets: secret value is nil")
	}
	if err := json.Unmarshal([]byte(*result.SecretString), target); err != nil {
		return errors.Wrap(err, "secrets: failed to unmarshal a secret value")
	}
	return nil
}

func (c *secretsClient) GetAzureSecret(ctx context.Context, secretID string) (*AzureAPISecret, error) {
	secret := &AzureAPISecret{}
	err := c.GetValue(ctx, secretID, secret)
	if err != nil {
		return nil, err
	}
	return secret, nil
}
