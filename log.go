// Copyright 2013 Joshua Marsh. All rights reserved.  Use of this
// source code is governed by a BSD-style license that can be found in
// the LICENSE file.

package gorca

import (
	"appengine"
	"fmt"
	"net/http"
)

// Log is a helper function that logs the given message to appenging
// with the given priority. Accepted priorities are "debug", "info",
// "warn", "error", and "crit". Other values default to "error".
func Log(c appengine.Context, r *http.Request, priority string,
	message string, params ...interface{}) {

	message = fmt.Sprintf("[%s] [%s] [%s]: %s", r.RemoteAddr, r.Method,
		r.URL, message)

	switch priority {
	case "debug":
		c.Debugf(message, params...)

	case "info":
		c.Infof(message, params...)

	case "warn":
		c.Warningf(message, params...)

	case "error":
		c.Errorf(message, params...)

	case "crit":
		c.Criticalf(message, params...)

	default:
		c.Errorf(message, params...)
	}
}
