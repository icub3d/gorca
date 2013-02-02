package gorca

import (
	"appengine"
	"fmt"
	"net/http"
)

// Log is a helper function that logs the given message to appenging
// with the given priority. Accepted priorities are "debug", "info",
// "warn", "error", and "crit". Other values default to "error".
func Log(r *http.Request, priority string, message string,
	params ...interface{}) {
	// Get the context for logging.
	c := appengine.NewContext(r)

	message = fmt.Sprintf("[%s] %s: %s", r.RemoteAddr, r.URL, message)

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
