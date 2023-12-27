package mwlog

import (
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func DefaultUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	logrusEntry := logrus.NewEntry(logrus.StandardLogger())
	return grpc_logrus.UnaryClientInterceptor(logrusEntry, defaultLogrusOpts()...)
}

func DefaultUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	logrusEntry := logrus.NewEntry(logrus.StandardLogger())
	return grpc_logrus.UnaryServerInterceptor(logrusEntry, defaultLogrusOpts()...)
}

func defaultLogrusOpts() []grpc_logrus.Option {
	levelFunc := func(c codes.Code) logrus.Level {
		switch c {
		case codes.OK:
			return logrus.DebugLevel
		default:
			return grpc_logrus.DefaultClientCodeToLevel(c)
		}
	}
	return []grpc_logrus.Option{
		grpc_logrus.WithLevels(levelFunc),
	}
}
