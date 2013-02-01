package gorca

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// WriteJSON transforms the given data into JSON and sends it as a
// response. If an error occurs, that will be returned instead.
func WriteJSON(w http.ResponseWriter, r *http.Request, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		LogAndUnexpected(w, r, fmt.Errorf("writing json: %s", err))
		return
	}

	WriteResponse(w, r, b)
}

// WriteMessage prints a standard JSON message to the given writer.
func WriteMessage(w http.ResponseWriter, r *http.Request, mtype, message string, code int) {
	// Make the JSON response.
	m := Message{Type: mtype, Message: message}
	b, err := json.Marshal(m)
	if err != nil {
		// Eeek! just return the message itself.
		http.Error(w, message, code)
		return
	}

	w.WriteHeader(code)
	WriteResponse(w, r, b)
}

// WriteSuccessMessage prints a JSON response of success to the given writer.
func WriteSuccessMessage(w http.ResponseWriter, r *http.Request) {
	WriteMessage(w, r, "success", ErrMsgs["success"], http.StatusOK)
}

// WriteAndLog writes the given data to the given response write. If
// an error occurs, it is logged.
func WriteResponse(w http.ResponseWriter, r *http.Request, bytes []byte) {
	_, err := w.Write(bytes)

	if err != nil {
		Log(r, "error", "Failed to write: %v", err.Error())
		Log(r, "info", "Message was: %v", string(bytes))
	}
}
