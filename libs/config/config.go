package config

import (
	"context"

	"github.com/kelseyhightower/envconfig"
)

// Value represents a generic configuration value or set of values.
type Value interface{}

// Configure is the interface for any configuration which will provide a configuration.
type Configure interface {
	// Load loads the value configuration
	Load(ctx context.Context, value Value) error
}

type ConfigureFn func(ctx context.Context, value Value) error

func (fn ConfigureFn) Load(ctx context.Context, value Value) error {
	return fn(ctx, value)
}

func WithEnv(prefix string) ConfigureFn {
	return func(_ context.Context, value Value) error {
		return envconfig.Process(prefix, value)
	}
}

func WithYaml(path string) ConfigureFn {
	return func(ctx context.Context, value Value) error {
		return NewYaml(path).Load(ctx, value)
	}
}

func Load(ctx context.Context, value Value, opts ...ConfigureFn) error {
	for _, cf := range opts {
		if err := cf(ctx, value); err != nil {
			return err
		}
	}
	return nil
}
