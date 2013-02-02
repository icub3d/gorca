// Package gorca contains common RESTful structures, methods, and
// functions that are useful go appengine applications.
package gorca

import (
	"fmt"
	"net/http"
)

// Message is a basic JSON response.
type Message struct {
	Type    string
	Message string
}

// ErrMsgs contains common JSON response messages.
var ErrMsgs map[string]string = map[string]string{
	"failed":       "Failed.",
	"notfound":     "Not found.",
	"success":      "Success.",
	"unexpected":   "Something unexpected happened.",
	"unauthorized": "You are not authorized to do that.",
}

// NotFoundFunc is a http.HandlerFunc that returns a standard 404 not
// found as well as a JSON response with the error.
func NotFoundFunc(w http.ResponseWriter, r *http.Request) {
	err := fmt.Errorf("request failed: %s %s", r.Method, r.URL.String())
	LogAndNotFound(w, r, err)
}

// LogAndNotFound logs the given error message and returns a not found
// JSON error message as well as a 404.
func LogAndNotFound(w http.ResponseWriter, r *http.Request, err error) {
	err = fmt.Errorf("not found: %s %s: %v", r.Method, r.URL, err)
	LogAndMessage(w, r, err, "error", ErrMsgs["notfound"],
		http.StatusNotFound)
}

// LogAndFailed logs the given error message and returns a failed
// JSON error message as well as a 400.
func LogAndFailed(w http.ResponseWriter, r *http.Request, err error) {
	err = fmt.Errorf("failed: %s %s: %v", r.Method, r.URL, err)
	LogAndMessage(w, r, err, "error", ErrMsgs["failed"],
		http.StatusBadRequest)
}

// LogAndUnexpected logs the given error message and returns an
// internal server error JSON error message as well as a 500.
func LogAndUnexpected(w http.ResponseWriter, r *http.Request, err error) {
	err = fmt.Errorf("unexpected: %s %s: %v", r.Method, r.URL, err)
	LogAndMessage(w, r, err, "error", ErrMsgs["unexpexted"],
		http.StatusInternalServerError)
}

// LogAndMessage logs the given error (if it is not nil) then sends
// the given JSON message and status as the response.
func LogAndMessage(w http.ResponseWriter, r *http.Request, err error,
	mtype, message string, code int) {

	if err != nil {
		Log(r, "error", err.Error())
	}
	Log(r, "info", "sent response (%s): %s", mtype, message)

	WriteMessage(w, r, mtype, message, code)
}
