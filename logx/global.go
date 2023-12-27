package logx

import (
	"context"

	"github.com/sirupsen/logrus"
)

var globalLogger = &defaultLogger{
	Level:     logrus.InfoLevel,
	Formatter: JSONFormatter(),
}

// New creates the Logger instance.
func New(opts ...LoggerOpt) {
	for _, o := range opts {
		o(globalLogger)
	}
	logrus.SetFormatter(globalLogger.Formatter)
	logrus.SetLevel(globalLogger.Level)
}

func Infof(format string, args ...interface{}) {
	globalLogger.Infof(format, args...)
}

func Infoln(args ...interface{}) {
	globalLogger.Infoln(args...)
}

func Info(ctx context.Context, msg interface{}, fields ...Fields) {
	globalLogger.Info(ctx, msg, fields...)
}

func Warnf(format string, args ...interface{}) {
	globalLogger.Warnf(format, args...)
}

func Warnln(args ...interface{}) {
	globalLogger.Warnln(args...)
}

func Warn(ctx context.Context, msg interface{}, fields ...Fields) {
	globalLogger.Warn(ctx, msg, fields...)
}

func Debugf(format string, args ...interface{}) {
	globalLogger.Debugf(format, args...)
}

func Debugln(args ...interface{}) {
	globalLogger.Debugln(args...)
}

func Debug(ctx context.Context, msg interface{}, fields ...Fields) {
	globalLogger.Debug(ctx, msg, fields...)
}

func Errorf(format string, args ...interface{}) {
	globalLogger.Errorf(format, args...)
}

func Errorln(args ...interface{}) {
	globalLogger.Errorln(args...)
}

func Error(ctx context.Context, msg interface{}, fields ...Fields) {
	globalLogger.Error(ctx, msg, fields...)
}

func Panicf(format string, args ...interface{}) {
	globalLogger.Panicf(format, args...)
}

func Panicln(args ...interface{}) {
	globalLogger.Panicln(args...)
}

func Panic(ctx context.Context, msg interface{}, fields ...Fields) {
	globalLogger.Panic(ctx, msg, fields...)
}

func WithContext(ctx context.Context) *logrus.Entry {
	return globalLogger.WithContext(ctx)
}
