package utils

import (
	"net/http"

	"github.com/mauricetjmurphy/ms-common/logx"
)

func HealthCheck(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set(ContentType, ContentTypeJSON)

	_, err := w.Write([]byte("{\"status\": \"OK\"}"))
	if err != nil {
		logx.Errorf("middleware : failed to write message %v", err)
	}
}
