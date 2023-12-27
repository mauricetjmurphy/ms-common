package cors

import (
	"net/http"

	"github.com/rs/cors"
)

// AllowCORSHandler allows CORS specification on the request, and add relevant CORS headers request.
func AllowCORSHandler(h http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"HEAD", "GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Accept", "Accept-Language", "Authorization", "Content-Type", "X-CSRF-Token", "SSO"},
	}).Handler(h)
}
