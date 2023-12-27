package logx

import (
	"context"

	"github.com/sirupsen/logrus"
)

// Fields is alias logrus.Fields
type Fields logrus.Fields

type Logger interface {
	Infof(format string, args ...interface{})
	Infoln(args ...interface{})
	Info(ctx context.Context, msg interface{}, fields ...Fields)

	Debugf(format string, args ...interface{})
	Debugln(args ...interface{})
	Debug(ctx context.Context, msg interface{}, fields ...Fields)

	Warnf(format string, args ...interface{})
	Warnln(args ...interface{})
	Warn(ctx context.Context, msg interface{}, fields ...Fields)

	Errorf(format string, args ...interface{})
	Errorln(args ...interface{})
	Error(ctx context.Context, msg interface{}, fields ...Fields)

	Panicf(format string, args ...interface{})
	Panicln(args ...interface{})
	Panic(ctx context.Context, msg interface{}, fields ...Fields)

	WithContext(ctx context.Context) *logrus.Entry
}

var _ Logger = (*defaultLogger)(nil)

type contextKey int

const contextFieldsKey contextKey = iota

type defaultLogger struct {
	Level     logrus.Level
	Formatter logrus.Formatter
}

func (l *defaultLogger) Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

func (l *defaultLogger) Infoln(args ...interface{}) {
	logrus.Infoln(args...)
}

func (l *defaultLogger) Info(ctx context.Context, msg interface{}, fields ...Fields) {
	l.entryFromContext(ctx, msg, fields...).Info(msg)
}

func (l *defaultLogger) Warnf(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}

func (l *defaultLogger) Warnln(args ...interface{}) {
	logrus.Warnln(args...)
}

func (l *defaultLogger) Warn(ctx context.Context, msg interface{}, fields ...Fields) {
	l.entryFromContext(ctx, msg, fields...).Warn(msg)
}

func (l *defaultLogger) Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

func (l *defaultLogger) Debugln(args ...interface{}) {
	logrus.Debugln(args...)
}

func (l *defaultLogger) Debug(ctx context.Context, msg interface{}, fields ...Fields) {
	l.entryFromContext(ctx, msg, fields...).Debug(msg)
}

func (l *defaultLogger) Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

func (l *defaultLogger) Errorln(args ...interface{}) {
	logrus.Errorln(args...)
}

func (l *defaultLogger) Error(ctx context.Context, msg interface{}, fields ...Fields) {
	l.entryFromContext(ctx, msg, fields...).Error(msg)
}

func (l *defaultLogger) Panicf(format string, args ...interface{}) {
	logrus.Panicf(format, args...)
}

func (l *defaultLogger) Panicln(args ...interface{}) {
	logrus.Panicln(args...)
}

func (l *defaultLogger) Panic(ctx context.Context, msg interface{}, fields ...Fields) {
	l.entryFromContext(ctx, msg, fields...).Panic(msg)
}

func (l *defaultLogger) WithContext(ctx context.Context) *logrus.Entry {
	return logrus.WithContext(ctx)
}

func (l *defaultLogger) entryFromContext(ctx context.Context, msg interface{}, fields ...Fields) *logrus.Entry {
	merged := Fields{}
	// Store the original error
	if err, ok := msg.(error); ok {
		merged["error"] = err
	}
	// Extract fields from context.
	if ctx != nil {
		cf, ok := ctx.Value(contextFieldsKey).(Fields)
		if ok {
			for k, v := range cf {
				merged[k] = v
			}
		}
	}
	// Extract fields from optional passed in fields arg(s).
	for _, f := range fields {
		for k, v := range f {
			merged[k] = v
		}
	}
	return logrus.WithFields(logrus.Fields(merged))
}
