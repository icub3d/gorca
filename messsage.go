package gorca

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
func WriteMessage(w http.ResponseWriter, r *http.Request,
	mtype, message string, code int) {

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

// WriteSuccessMessage prints a JSON response of success to the given
// writer.
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

// UnmarshalOrFail attempts to unmarshal the given bytes as JSON and
// put it in where. if it fails, false is returned and a "failed"
// message is returned. In that case, this should be terminal.
func UnmarshalOrFail(w http.ResponseWriter, r *http.Request, bytes []byte,
	where interface{}) bool {

	err := json.Unmarshal(bytes, where)
	if err != nil {
		LogAndFailed(w, r, err)
		return false
	}

	return true
}

// GetBodyOrFail attempts to read the body from the given request. If
// it succeeds, the body is returned as a string as well as true. If
// it fails, "" and false are returned. The failure is also loged and
// generic error is returned as the response.
func GetBodyOrFail(w http.ResponseWriter, r *http.Request) ([]byte, bool) {
	// Read the body for the JSON.
	if r.Body == nil {
		LogAndFailed(w, r, fmt.Errorf("no JSON found"))
		return nil, false
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		LogAndUnexpected(w, r, err)
		return nil, false
	}

	return body, true
}

// UnmarshalFromBodyOrFail attempts to read the body from the give
// request and umarshal it as JSON into the given interface. If an
// error occurs, the failure is logged and a generic message is
// returned as the response. The boolean value returned signifies the
// success of the operation.
func UnmarshalFromBodyOrFail(w http.ResponseWriter, r *http.Request,
	v interface{}) bool {
	body, success := GetBodyOrFail(w, r)
	if !success {
		return false
	}

	if !UnmarshalOrFail(w, r, body, v) {
		return false
	}

	return true
}
