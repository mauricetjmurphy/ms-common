package publisher

import (
	"encoding/json"
)

type Encoder interface {
	Marshal(payload interface{}) (string, error)
}

func DefaultEncode() Encoder {
	return &defaultEncoder{}
}

type defaultEncoder struct{}

func (e *defaultEncoder) Marshal(payload interface{}) (string, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return "nil", err
	}
	return string(b), nil
}
