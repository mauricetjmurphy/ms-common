package recovery

import (
	"encoding/json"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/NBCUniversal/gvs-ms-common/logx"
)

const AllowPanicEnv = "ALLOW_PANIC"

// PanicRecoveryHandle is an HTTP middleware that recovers from a panic.
func PanicRecoveryHandle(inner http.Handler) http.Handler {
	middleware := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
				if os.Getenv(AllowPanicEnv) == "true" {
					panic(err)
				}
				logx.Errorf("panic request %v on err %v", r.URL.Path, err)
				status := http.StatusInternalServerError
				writeJSONErr(w, http.StatusText(status), status)
			}
		}()

		inner.ServeHTTP(w, r)
	}

	return http.HandlerFunc(middleware)
}

func writeJSONErr(w http.ResponseWriter, message string, status int) {
	response := map[string]interface{}{
		"message": []string{message},
		"status":  status,
	}
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	if e := json.NewEncoder(w).Encode(response); e != nil {
		http.Error(w, message, status)
	}
}
