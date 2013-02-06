// Copyright 2013 Joshua Marsh. All rights reserved.  Use of this
// source code is governed by a BSD-style license that can be found in
// the LICENSE file.

// Package gorca contains common RESTful structures, methods, and
// functions that are useful go appengine applications.
//
// If you are testing these functions, there are some steps you need
// to do to setup a proper testing environment. 
//
//    export APPENGINE_SDK=/path/to/google_appengine
//    cd $GOPATH/src
//    ln -s $APPENGINE_SDK/goroot/src/pkg/appengine
//    ln -s $APPENGINE_SDK/goroot/src/pkg/appengine_internal
//    go get github.com/icub3d/appenginetesting
//    cd github.com/icub3d/gorca
//    go test ./...
package gorca

import (
	"appengine"
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

// NotFoundFunc makes a http.HandlerFunc that returns a standard
// 404 not found as well as a JSON response with the error.
func NotFoundFunc(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	LogAndNotFound(c, w, r, fmt.Errorf("not found func"))
}

// LogAndNotFound logs the given error message and returns a not found
// JSON error message as well as a 404.
func LogAndNotFound(c appengine.Context, w http.ResponseWriter,
	r *http.Request, err error) {

	err = fmt.Errorf("not found: %v", err)
	LogAndMessage(c, w, r, err, "error", ErrMsgs["notfound"],
		http.StatusNotFound)
}

// LogAndFailed logs the given error message and returns a failed
// JSON error message as well as a 400.
func LogAndFailed(c appengine.Context, w http.ResponseWriter,
	r *http.Request, err error) {

	err = fmt.Errorf("failed: %v", err)
	LogAndMessage(c, w, r, err, "error", ErrMsgs["failed"],
		http.StatusBadRequest)
}

// LogAndUnexpected logs the given error message and returns an
// internal server error JSON error message as well as a 500.
func LogAndUnexpected(c appengine.Context, w http.ResponseWriter,
	r *http.Request, err error) {

	err = fmt.Errorf("unexpected: %v", err)
	LogAndMessage(c, w, r, err, "error", ErrMsgs["unexpected"],
		http.StatusInternalServerError)
}

// LogAndMessage logs the given error (if it is not nil) then sends
// the given JSON message and status as the response.
func LogAndMessage(c appengine.Context, w http.ResponseWriter,
	r *http.Request, err error, mtype, message string, code int) {

	if err != nil {
		Log(c, r, "error", err.Error())
	}
	Log(c, r, "info", "sent response (%s): %s", mtype, message)

	WriteMessage(c, w, r, mtype, message, code)
}
