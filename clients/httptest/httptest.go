package httptest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

// Handler is alias http.Handler
type Handler http.Handler

// Request is alias http.Request
type Request *http.Request

// RequestMatcher presents the matching request.
type RequestMatcher struct {
	// HTTP method. eg : POST, GET, PUT, DELETE
	Method string
	// URL path request
	Path string
	// Response raw data. eg: JSON
	RawResponse string
	// Assert request query or body.
	AssertRequest func(request Request)
}

// NewServer starts and returns a new Server on given matcher request.
func NewServer(requestMatcher RequestMatcher) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == requestMatcher.Method && r.URL.Path == requestMatcher.Path {
			if requestMatcher.AssertRequest != nil {
				requestMatcher.AssertRequest(r)
			}
			if len(requestMatcher.RawResponse) > 0 {
				fmt.Fprintln(w, requestMatcher.RawResponse)
			}
		}
	}))
}

// NewServerWith starts and returns a new Server on given custom handler.
func NewServerWith(handler Handler) *httptest.Server {
	return httptest.NewServer(handler)
}

// NewMatcher makes matcher request.
func NewMatcher(method, path, rawResponse string) RequestMatcher {
	return RequestMatcher{
		Method:      method,
		Path:        path,
		RawResponse: rawResponse,
	}
}
