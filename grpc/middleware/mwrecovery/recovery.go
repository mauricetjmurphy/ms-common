package mwrecovery

import (
	"os"
	"runtime/debug"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/mauricetjmurphy/ms-common/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const AllowPanicEnv = "ALLOW_PANIC"

var DefaultPanicHandler grpc_recovery.RecoveryHandlerFunc = func(p interface{}) (err error) {
	debug.PrintStack()
	if os.Getenv(AllowPanicEnv) == "true" {
		panic(err)
	}
	logx.Errorf("panic on err %v", err)
	return status.Errorf(codes.Unknown, "unexpected failure")
}

func DefaultUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(DefaultPanicHandler),
	}
	return grpc_recovery.UnaryServerInterceptor(opts...)
}
