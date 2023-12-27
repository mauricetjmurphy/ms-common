package mwauthz

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Authenticator propagates the authentication process
type Authenticator interface {
	// HandleAuth process authentication on given incoming request.
	HandleAuth(ctx context.Context, request interface{}) (context.Context, error)
}

// AuthInterceptor is the customized GRPC interceptor.
type AuthInterceptor interface {
	UnaryServerInterceptor() grpc.UnaryServerInterceptor
}

// New creates the gRPC authentication interceptor.
func New(authenticator Authenticator, excludeMethodPatterns ...string) AuthInterceptor {
	return &authInterceptorImpl{
		Authenticator:         authenticator,
		ExcludeMethodPatterns: excludeMethodPatterns,
	}
}

type authInterceptorImpl struct {
	Authenticator         Authenticator
	ExcludeMethodPatterns []string
}

func (a *authInterceptorImpl) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {

		if len(a.ExcludeMethodPatterns) > 0 {
			for _, e := range a.ExcludeMethodPatterns {
				if info.FullMethod == e {
					return handler(ctx, req)
				}
			}
		}

		newCtx, err := a.Authenticator.HandleAuth(ctx, req)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		return handler(newCtx, req)
	}
}
