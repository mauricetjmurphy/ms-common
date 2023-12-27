package auth

import (
	"context"
	"net/http"
	"regexp"

	"github.com/mauricetjmurphy/ms-common/http/utils"
)

type contextKey int

const ssoIDContextKey contextKey = iota

// DefaultPatternExcludedPaths is patterns to exclude authorized check health or swagger.
const DefaultPatternExcludedPaths = "^/api/.*/health\\b|\\bswagger$"

// SSOHandler handle authorized incoming request by SSO token in header
func SSOHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Defaults exclude authorized check request URLs.
		if ok, _ := regexp.MatchString(DefaultPatternExcludedPaths, r.URL.Path); ok {
			h.ServeHTTP(w, r)
			return
		}
		// Authorize Check
		ssoID, err := utils.GetSsoID(r)
		if err != nil {
			utils.Unauthorized(w)
			return
		}
		r = r.WithContext(WithSsoID(r.Context(), ssoID))
		h.ServeHTTP(w, r)
	})
}

func WithSsoID(ctx context.Context, ssoID string) context.Context {
	return context.WithValue(ctx, ssoIDContextKey, ssoID)
}

// GetSsoID retrieves the SSO value from incoming request.
func GetSsoID(ctx context.Context) string {
	if v, ok := ctx.Value(ssoIDContextKey).(string); ok {
		return v
	}
	return ""
}
