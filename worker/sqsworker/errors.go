package sqsworker

import "github.com/pkg/errors"

var (
	ErrNoMessageInQueue = errors.New("no message in queue")
)
