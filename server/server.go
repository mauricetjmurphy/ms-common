package server

import (
	nett "net/http"
	"os"
	"os/signal"
	"syscall"

	http "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

var stdSignals = []os.Signal{syscall.SIGINT}

// Server presents the server service.
type Server struct {
	stack     *stack
	opts      *serverOpts
	listeners *listeners
}

// Service register the HTTP and GRPC connections.
type Service interface {
	RegisterGRPC(*grpc.Server)
	RegisterHTTP(*http.ServeMux)
}

// NewServer creates a gRPC Server and HTTP server on same port with given server options.
func NewServer(addr string, opts ...Option) *Server {
	serverOpts := defaultServerOpts(addr)
	for _, opt := range opts {
		opt(serverOpts)
	}
	return &Server{opts: serverOpts}
}

// Run starts the service.
func (s *Server) Run(service Service) error {
	var err error
	s.listeners, err = newListenerSet(s.opts)
	if err != nil {
		return errors.Wrap(err, "failed to initialize listeners")
	}

	s.stack = newStack(s.opts)
	service.RegisterGRPC(s.stack.grpc)
	service.RegisterHTTP(s.stack.mux)

	var g errgroup.Group
	g.Go(s.waitForShutdown())
	if err = s.run(); err != nil {
		return err
	}
	return g.Wait()
}

func (s *Server) run() error {
	errChan := make(chan error, 5)
	if s.listeners.mainListener != nil {
		go func() {
			err := s.listeners.mainListener.Serve()
			errChan <- err
		}()
	}
	go func() {
		var handler nett.Handler = s.stack.mux
		for i := len(s.opts.HTTPMiddlewares) - 1; i >= 0; i-- {
			handler = s.opts.HTTPMiddlewares[i](handler)
		}
		err := nett.Serve(s.listeners.mux, handler)
		errChan <- err
	}()
	go func() {
		err := s.stack.grpc.Serve(s.listeners.grpc)
		errChan <- err
	}()
	return <-errChan
}

func (s *Server) waitForShutdown() func() error {
	return func() error {
		sig := make(chan os.Signal, 1)
		if len(stdSignals) > 0 {
			signal.Notify(sig, stdSignals...)
		}

		// Wait for a shutdown signal.
		select {
		case <-sig:
		}

		return s.Shutdown()
	}
}

// Shutdown stops the server gracefully.
func (s *Server) Shutdown() error {
	s.opts.Logger.Log("server : shutting down ...")
	if svc := s.stack; svc != nil {
		svc.grpc.GracefulStop()
		err := s.listeners.mux.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
