package circuit

import (
	"time"

	"github.com/cep21/circuit/v3"
	"github.com/cep21/circuit/v3/closers/hystrix"
	log "github.com/mauricetjmurphy/ms-common/logx"
)

const (
	defaultTimeout               = 5 * time.Minute
	defaultMaxConcurrentRequests = 500

	// values defaults for openerFactory configuration
	defaultErrorThresholdPercentage = 50
	defaultRequestVolumeThreshold   = 20
	defaultRollingDuration          = 10 * time.Second
	defaultNumBuckets               = 10

	// values defaults for closerFactory configuration
	defaultSleepWindow                  = 5 * time.Second
	defaultHalfOpenAttempts             = 1
	defaultRequiredConcurrentSuccessful = 1
)

// Config presents the circuit configuration.
type Config struct {
	// The factory builds the circuit configuration.
	CircuitFactory func(circuitName string) circuit.Config

	// the hystrix factory builds hystrix configuration.
	HystrixFactory *hystrix.Factory
}

func applyDefault(conf *Config) {
	if conf.CircuitFactory == nil {
		conf.CircuitFactory = defaultCircuitConfigFn()
	}
	if conf.HystrixFactory == nil {
		openerFactory := defaultConfigureOpenerFn()
		closerFactory := defaultConfigureCloserFn()
		conf.HystrixFactory = &hystrix.Factory{
			CreateConfigureOpener: []func(circuitName string) hystrix.ConfigureOpener{openerFactory},
			CreateConfigureCloser: []func(circuitName string) hystrix.ConfigureCloser{closerFactory},
		}
	}
}

func defaultConfigureOpenerFn() func(circuitName string) hystrix.ConfigureOpener {
	return func(circuitName string) hystrix.ConfigureOpener {
		return hystrix.ConfigureOpener{
			ErrorThresholdPercentage: defaultErrorThresholdPercentage,
			RequestVolumeThreshold:   defaultRequestVolumeThreshold,
			RollingDuration:          defaultRollingDuration,
			NumBuckets:               defaultNumBuckets,
		}
	}
}

func defaultConfigureCloserFn() func(circuitName string) hystrix.ConfigureCloser {
	return func(circuitName string) hystrix.ConfigureCloser {
		return hystrix.ConfigureCloser{
			SleepWindow:                  defaultSleepWindow,
			HalfOpenAttempts:             defaultHalfOpenAttempts,
			RequiredConcurrentSuccessful: defaultRequiredConcurrentSuccessful,
		}
	}
}

func defaultCircuitConfigFn() func(circuitName string) circuit.Config {
	return func(circuitName string) circuit.Config {
		return circuit.Config{
			General: circuit.GeneralConfig{
				GoLostErrors: func(err error, panics interface{}) {
					log.Errorf("circuit_name: "+circuitName, err)
				},
			},
			Execution: circuit.ExecutionConfig{
				Timeout:               defaultTimeout,
				MaxConcurrentRequests: defaultMaxConcurrentRequests,
			},
		}
	}
}
