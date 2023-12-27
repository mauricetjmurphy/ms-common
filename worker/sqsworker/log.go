package sqsworker

import "fmt"

// Logger is the minimal interface for logging methods.
type Logger interface {
	Infof(format string, messages ...string)
	Debugf(format string, messages ...string)
	Errorf(format string, messages ...string)
}

type StandardLogger struct {
}

func (l *StandardLogger) Infof(format string, messages ...string) {
	fmt.Printf(format, messages)
}

func (l *StandardLogger) Debugf(format string, messages ...string) {
	fmt.Printf(format, messages)
}

func (l *StandardLogger) Errorf(format string, messages ...string) {
	fmt.Printf(format, messages)
}
