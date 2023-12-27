package circuit

import (
	"github.com/cep21/circuit/v3"
)

// Circuit is alias circuit Circuit.
type Circuit = circuit.Circuit

// ClientManager is a wrapper for the circuit manager.
type ClientManager struct {
	*circuit.Manager
}

// NewDefault creates the client circuit manager on given default configuration.
func NewDefault() *ClientManager {
	return New(&Config{})
}

// New creates the client circuit manager on given configuration.
func New(conf *Config) *ClientManager {
	applyDefault(conf)
	manager := &circuit.Manager{
		DefaultCircuitProperties: []circuit.CommandPropertiesConstructor{
			conf.CircuitFactory,
			conf.HystrixFactory.Configure,
		},
	}
	return &ClientManager{Manager: manager}
}
