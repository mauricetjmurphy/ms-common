package server

import (
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// Option is an optional setting applied to the Server.
type Option func(*serverOpts)

// serverOpts presents the server options.
type serverOpts struct {
	GRPCAddr string
	GRPCOpts []grpc.ServerOption

	HTTPAddr        string
	HTTPMuxOpts     []runtime.ServeMuxOption
	HTTPMiddlewares []func(http.Handler) http.Handler

	Logger Logger
}

func defaultServerOpts(addr string) *serverOpts {
	return &serverOpts{
		GRPCAddr: addr,
		HTTPAddr: addr,
		Logger:   &defaultLogger{},
	}
}

// WithGRPCAddr sets gRPC server endpoint options.
func WithGRPCAddr(addr string) Option {
	return func(o *serverOpts) {
		o.GRPCAddr = addr
	}
}

// WithGRPCOpts sets gRPC server options.
func WithGRPCOpts(opts ...grpc.ServerOption) Option {
	return func(o *serverOpts) {
		o.GRPCOpts = opts
	}
}

// WithHTTPAddr sets Http server endpoint options.
func WithHTTPAddr(addr string) Option {
	return func(o *serverOpts) {
		o.HTTPAddr = addr
	}
}

// WithHTTPMuxOpts sets Http server options.
func WithHTTPMuxOpts(opts ...runtime.ServeMuxOption) Option {
	return func(o *serverOpts) {
		o.HTTPMuxOpts = opts
	}
}

func WithHTTPMiddlewares(middleware ...func(http.Handler) http.Handler) Option {
	return func(o *serverOpts) {
		o.HTTPMiddlewares = append(o.HTTPMiddlewares, middleware...)
	}
}

// WithLoggerOpts sets Http server options.
func WithLoggerOpts(logger Logger) Option {
	return func(o *serverOpts) {
		o.Logger = logger
	}
}
