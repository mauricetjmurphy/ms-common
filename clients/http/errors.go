package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Error presents the wrapper error from external response error.
type Error struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
	Cause      string `json:"-"`
}

// NewErr creates the wrapper error on given message and status code.
func NewErr(caused error, message string, statusCode int) *Error {
	err := &Error{Message: message, StatusCode: statusCode}
	if caused != nil {
		err.Cause = caused.Error()
	}
	return err
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	return strconv.Itoa(e.StatusCode) + ": " + e.Message
}

// NewErrWith returns the wrapper client error.
func NewErrWith(err error) *Error {
	if v, ok := err.(*Error); ok {
		return NewErr(err, v.Message, v.StatusCode)
	}
	return InternalServerError(err)
}

// InternalServerError presents the utility to return the internal error json body.
func InternalServerError(err error) *Error {
	return &Error{
		StatusCode: http.StatusInternalServerError,
		Message:    "Internal Server Error",
		Cause:      err.Error(),
	}
}

// BadRequest presents the utility to return the bad request (400) error json body.
func BadRequest(message string) *Error {
	return &Error{
		StatusCode: http.StatusBadRequest,
		Message:    message,
	}
}

// Unauthorized presents the utility to return the un-authorized request (401) error json body.
func Unauthorized() *Error {
	statusUnauthorized := http.StatusUnauthorized
	return &Error{
		StatusCode: statusUnauthorized,
		Message:    http.StatusText(statusUnauthorized),
	}
}

// ParseJSONErr allows parsing JSON error message from response if any.
func ParseJSONErr(resp *http.Response) *Error {
	if resp == nil {
		return &Error{Message: "http : unexpected nil *http.Response"}
	}
	var (
		statusCode = resp.StatusCode
	)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return NewErr(err, fmt.Sprintf("http : unable to read response body %v", resp.Body), statusCode)
	}
	return NewErr(nil, string(body), statusCode)
}

func IsSuccess(response *http.Response) bool {
	return response != nil && response.StatusCode >= 200 && response.StatusCode < 300
}

func IsErrorStatus(response *http.Response) bool {
	return response != nil && !IsSuccess(response)
}
