package matcher

import (
	"net/textproto"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

// HeaderMatcher extends runtime.DefaultHeaderMatcher by also injecting headers.
func HeaderMatcher(key string) (string, bool) {
	// default behavior
	if v, ok := runtime.DefaultHeaderMatcher(key); ok {
		return v, ok
	}

	switch textproto.CanonicalMIMEHeaderKey(key) {
	case "Sso":
		return key, true
	}
	return "", false
}
