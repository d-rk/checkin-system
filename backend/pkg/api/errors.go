package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

type Error interface {
	// APIError returns an HTTP status code and an API-safe error message.
	APIError() (int, string)
}

var (
	ErrInvalidCredentials = &sentinelAPIError{status: http.StatusUnauthorized, msg: "invalid credentials"}
	ErrInvalidToken       = &sentinelAPIError{status: http.StatusUnauthorized, msg: "invalid token"}
	ErrNotFound           = &sentinelAPIError{status: http.StatusNotFound, msg: "not found"}
	ErrBadRequest         = &sentinelAPIError{status: http.StatusBadRequest, msg: "bad request"}
	ErrConflict           = &sentinelAPIError{status: http.StatusConflict, msg: "conflict"}
)

type sentinelAPIError struct {
	status int
	msg    string
}

func (e sentinelAPIError) Error() string {
	return e.msg
}

func (e sentinelAPIError) APIError() (int, string) {
	return e.status, e.msg
}

func (e sentinelAPIError) Wrap(err error) error {
	return sentinelWrappedError{error: err, sentinel: &e}
}

type sentinelWrappedError struct {
	error

	sentinel *sentinelAPIError
}

func (e sentinelWrappedError) Is(err error) bool {
	return e.sentinel == err
}

func (e sentinelWrappedError) APIError() (int, string) {
	return e.sentinel.APIError()
}

func Wrap(err error, sentinel *sentinelAPIError) error {
	return sentinelWrappedError{error: err, sentinel: sentinel}
}

func ValidateErrorHandlerFunc(w http.ResponseWriter, message string, statusCode int) {
	errorResponse(w, statusCode, message)
}

func handlerError(w http.ResponseWriter, _ *http.Request, err error) {
	w.Header().Add("Content-Type", "application/json")

	var apiErr Error

	if errors.As(err, &apiErr) {
		status, msg := apiErr.APIError()
		errorResponse(w, status, msg)
	} else {
		slog.Error("unexpected error", "error", err)
		errorResponse(w, http.StatusInternalServerError, "internal error")
	}
}

func errorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(ErrorResponse{Message: message})
}
