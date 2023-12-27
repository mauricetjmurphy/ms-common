package publisher

import "time"

// Config presents eventbus configuration
type Config struct {
	encoder      Encoder
	eventBusName string
	timeout      time.Duration
}

type Opts func(config *Config)

// WithEncoder provide the encoder message details.
// Default with JSON encoder
func WithEncoder(encoder Encoder) Opts {
	return func(config *Config) {
		config.encoder = encoder
	}
}

// WithEventBusName provides the eventbus name.
func WithEventBusName(eventBusName string) Opts {
	return func(config *Config) {
		config.eventBusName = eventBusName
	}
}

func WithTimeout(timeout time.Duration) Opts {
	return func(config *Config) {
		config.timeout = timeout
	}
}
