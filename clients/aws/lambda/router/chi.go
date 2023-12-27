package router

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/go-chi/chi/v5"
)

// ChiLambdaV2 makes it easy to send API Gateway proxy events to a Chi
// Mux. The library transforms the proxy event into an HTTP request and then
// creates a proxy response object from the http.ResponseWriter
type ChiLambdaV2 struct {
	core.RequestAccessorV2
	chiMux *chi.Mux
}

// New creates a new instance of the ChiLambdaV2 object.
// Receives an initialized *chi.Mux object - normally created with chi.NewRouter().
// It returns the initialized instance of the ChiLambdaV2 object.
func New(chi *chi.Mux) *ChiLambdaV2 {
	return &ChiLambdaV2{chiMux: chi}
}

// ProxyWithContext receives context and an API Gateway proxy event,
// transforms them into a http.Request object, and sends it to the chi.Mux for routing.
// It returns a proxy response object generated from the http.ResponseWriter.
func (g *ChiLambdaV2) ProxyWithContext(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	chiRequest, err := g.EventToRequestWithContext(ctx, req)
	return g.proxyInternal(chiRequest, err)
}

func (g *ChiLambdaV2) proxyInternal(chiRequest *http.Request, err error) (events.APIGatewayV2HTTPResponse, error) {
	if err != nil {
		return core.GatewayTimeoutV2(), core.NewLoggedError("Could not convert proxy event to request: %v", err)
	}

	respWriter := core.NewProxyResponseWriterV2()
	g.chiMux.ServeHTTP(http.ResponseWriter(respWriter), chiRequest)

	proxyResponse, err := respWriter.GetProxyResponse()
	if err != nil {
		return core.GatewayTimeoutV2(), core.NewLoggedError("Error while generating proxy response: %v", err)
	}

	return proxyResponse, nil
}
