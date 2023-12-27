package router

import (
	"context"
	"encoding/json"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/NBCUniversal/gvs-ms-common/clients/aws/lambda/utils"
	"github.com/NBCUniversal/gvs-ms-common/logx"
	"github.com/aws/aws-lambda-go/events"
	awslambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/go-chi/chi/v5"
)

type Handler func(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error)

// Default initialize the default go-chi router instance.
func Default() *chi.Mux {
	return chi.NewRouter()
}

// StartFunc invokes Lambda when run inside a Lambda function triggered from API Gateway via HTTP API type/ Lambda Function URL,
// otherwise, it will listen and server on local mode for development purpose on given address.
// Usage:
//
//	r := router.Default()
//	r.Get("/health", HealthCheck)
//	if err := router.StartFunc(":80", r)); err != nil { .... }
func StartFunc(addr string, mux *chi.Mux) error {
	if !utils.IsLambdaEnv() {
		logx.Infof("router : run the server in localhost mode [%v]", addr)
		return http.ListenAndServe(addr, mux) //#nosec G114
	}
	awslambda.Start(wrap(New(mux).ProxyWithContext))
	return nil
}

// FileServer conveniently sets up a http.FileServer handler to serve static files from a http.FileSystem.
// Usage:
//
//	  r := router.Default()
//		 FileServer(r, "/", http.Dir("swagger/ui"))
func FileServer(r chi.Router, path string, root http.FileSystem) {
	fs := http.StripPrefix(path, http.FileServer(root))
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"
	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})
}

// wrap customs error response handler
func wrap(f Handler) Handler {
	return func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		obj, err := f(ctx, req)

		if err != nil {
			logx.Errorf("router err : %v", err)
			// Verify gRPC Resource
			if status.Code(err) == codes.ResourceExhausted {
				return internalServer(NewgRPCResourceExhaustedErr()), nil
			}
			return internalServer(NewInternalErr()), err
		}
		// Handle error lambda limitations
		payload, _ := json.Marshal(obj)
		if cap(payload) >= AWSLambdaSyncInvocationPayloadMax {
			return internalServer(NewExceededInvocationPayloadErr()), nil
		}
		return obj, err
	}
}

func internalServer(err ErrorStatus) events.APIGatewayV2HTTPResponse {
	payload, _ := json.Marshal(err)
	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       string(payload),
	}
}
