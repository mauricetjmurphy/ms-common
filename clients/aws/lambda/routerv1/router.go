package routerv1

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	awslambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/go-chi/chi/v5"
	"github.com/mauricetjmurphy/ms-common/clients/aws/lambda/utils"
	"github.com/mauricetjmurphy/ms-common/logx"
)

type Handler func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

// Default initialize the default go-chi router instance.
func Default() *chi.Mux {
	return chi.NewRouter()
}

// StartFunc invokes Lambda when run inside a Lambda function triggered from API Gateway via REST API type,
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

// wrap customs error response handler
func wrap(f Handler) Handler {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		obj, err := f(ctx, req)
		if err != nil {
			logx.Errorf("router err : %v", err)
			return internalServer(NewInternalErr()), err
		}
		return obj, err
	}
}

// ErrorStatus presents the error status model
type ErrorStatus struct {
	// Code indicate the unique error identifier code
	Code string `json:"code"`
	// Message presents the error message description
	Message string `json:"message"`
}

// NewInternalErr returns the internal server error status
func NewInternalErr() ErrorStatus {
	return ErrorStatus{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "Internal Server Error",
	}
}

func internalServer(err ErrorStatus) events.APIGatewayProxyResponse {
	payload, _ := json.Marshal(err)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       string(payload),
	}
}
