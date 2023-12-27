package utils

import (
	"encoding/json"
	"io"
	"net/http"

	chttp "github.com/NBCUniversal/gvs-ms-common/clients/http"
	"github.com/NBCUniversal/gvs-ms-common/logx"
)

const (
	ContentType     = "Content-Type"
	ContentTypeJSON = "application/json"
)

// UnmarshalRequest presents the function to convert to target object
func UnmarshalRequest(body io.ReadCloser, target interface{}) error {
	return json.NewDecoder(body).Decode(&target)
}

// OK presents the function to handle write response body.
func OK(w http.ResponseWriter, body interface{}) {
	responseJSON, err := json.Marshal(body)
	if err != nil {
		logx.Errorf("router : unable to marshal a body %v on err %v", body, err)
		return
	}

	_, err = w.Write(responseJSON)
	if err != nil {
		logx.Errorf("router : unable to write a body %v on err %v", responseJSON, err)
		return
	}
}

// JSON presents the utility function to response json body content
func JSON(w http.ResponseWriter, body interface{}) {
	w.Header().Set(ContentType, ContentTypeJSON)
	OK(w, body)
}

// ServerError presents the function to handle error response.
func ServerError(w http.ResponseWriter, err error) {
	var httpErr = chttp.NewErrWith(err)
	w.WriteHeader(httpErr.StatusCode)
	errJSON, _ := json.Marshal(httpErr)
	_, err = w.Write(errJSON)
	if err != nil {
		logx.Errorf("router : unable to write a error body %v on err %v", errJSON, err)
		return
	}
}

// JSONError presents the utility to handle error json body response
func JSONError(w http.ResponseWriter, err error) {
	w.Header().Set(ContentType, ContentTypeJSON)
	ServerError(w, err)
}

// BadRequest presents the utility to return the bad request (400) error json body.
func BadRequest(w http.ResponseWriter, message string) {
	JSONError(w, chttp.BadRequest(message))
}

// Unauthorized presents the utility to return the un-authorized request (401) error json body.
func Unauthorized(w http.ResponseWriter) {
	JSONError(w, chttp.Unauthorized())
}

// InternalError presents the utility to return the internal error json body.
func InternalError(w http.ResponseWriter, err error) {
	JSONError(w, chttp.InternalServerError(err))
}
