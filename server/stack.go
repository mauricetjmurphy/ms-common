package server

import (
	"net"
	"time"

	http "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
)

const (
	retryWait     = 500 * time.Millisecond
	retryDuration = 10 * time.Second
)

type stack struct {
	grpc *grpc.Server
	mux  *http.ServeMux
}

type listeners struct {
	mainListener cmux.CMux
	grpc         net.Listener
	mux          net.Listener
}

func newStack(opts *serverOpts) *stack {
	return &stack{
		grpc: grpc.NewServer(opts.GRPCOpts...),
		mux:  http.NewServeMux(opts.HTTPMuxOpts...),
	}
}

func newListenerSet(opts *serverOpts) (*listeners, error) {
	lis := &listeners{}
	var err error
	lis.grpc, err = newListener(opts.GRPCAddr)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create main listener")
	}
	if opts.GRPCAddr == opts.HTTPAddr {
		mux := cmux.New(lis.grpc)
		lis.grpc = mux.Match(cmux.HTTP2())
		lis.mux = mux.Match(cmux.Any())
		lis.mainListener = mux
	} else {
		lis.mux, err = newListener(opts.HTTPAddr)
	}
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create HTTP listener")
	}
	return lis, nil
}

// newListener start net.Listener on an addr. Keeps retrying if address is already in use.
func newListener(address string) (net.Listener, error) {
	var listener net.Listener
	var err error
	start := time.Now()
	for time.Since(start) < retryDuration {
		listener, err = net.Listen("tcp", address)
		if err == nil {
			return listener, nil
		}
		time.Sleep(retryWait)
	}
	return nil, err
}
