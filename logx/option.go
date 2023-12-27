package logx

import "github.com/sirupsen/logrus"

const TextFormat = "text"

// LoggerOpt are options on the logging.
type LoggerOpt func(*defaultLogger)

func JSONFormatter() logrus.Formatter {
	return &logrus.JSONFormatter{}
}

func TextFormatter() logrus.Formatter {
	return &logrus.TextFormatter{FullTimestamp: true}
}

// WithLevel set the level for the logger
func WithLevel(level string) LoggerOpt {
	return func(opt *defaultLogger) {
		if v, err := logrus.ParseLevel(level); err != nil {
			opt.Level = logrus.InfoLevel
		} else {
			opt.Level = v
		}
	}
}

// WithFormat set the format for the logger
func WithFormat(format string) LoggerOpt {
	return func(opt *defaultLogger) {
		switch format {
		case TextFormat:
			opt.Formatter = TextFormatter()
		default:
			opt.Formatter = JSONFormatter()
		}
	}
}
