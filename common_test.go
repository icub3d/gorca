package gorca

import (
	"appengine"
	"fmt"
	"github.com/icub3d/appenginetesting"
	"net/http"
	"net/http/httptest"
	"testing"
)

type LogAndFunc func(c appengine.Context, w http.ResponseWriter,
	r *http.Request, err error)

func TestLogAnds(t *testing.T) {

	tests := []struct {
		f      LogAndFunc
		method string
		url    string
		err    error
		ecode  int
		ebody  string
	}{
		// LogAndNotFound
		{
			f:      LogAndNotFound,
			method: "GET",
			url:    "/",
			err:    fmt.Errorf("no such file or directory"),
			ecode:  http.StatusNotFound,
			ebody:  `{"Type":"error","Message":"Not found."}`,
		},

		// LogAndFailed
		{
			f:      LogAndFailed,
			method: "GET",
			url:    "/",
			err:    fmt.Errorf("oopsie, you failed"),
			ecode:  http.StatusBadRequest,
			ebody:  `{"Type":"error","Message":"Failed."}`,
		},

		// LogAndUnexpected
		{
			f:      LogAndUnexpected,
			method: "GET",
			url:    "/",
			err:    fmt.Errorf("oopsie, i failed"),
			ecode:  http.StatusInternalServerError,
			ebody:  `{"Type":"error","Message":"Something unexpected happened."}`,
		},
	}

	// We are going to reuse the context.
	c, err := appenginetesting.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	for i, test := range tests {
		// Make the request and writer.
		w := httptest.NewRecorder()
		r, err := http.NewRequest(test.method, test.url, nil)
		if err != nil {
			t.Fatal(err)
		}

		// Call the test function
		test.f(c, w, r, test.err)

		// Check the status
		if w.Code != test.ecode {
			t.Errorf("(%v) expexted %v as response code. Got: %v",
				i, test.ecode, w.Code)
		}

		body := w.Body.String()
		if body != test.ebody {
			t.Errorf("(%v) expexted %v as response body. Got: %v",
				i, test.ebody, body)
		}
	}
}

func TestLogAndMessage(t *testing.T) {

	tests := []struct {
		method string
		url    string
		err    error
		ecode  int
		ebody  string
		mtype  string
		msg    string
	}{
		// Test A success message. The rest are tested above.
		{
			method: "GET",
			url:    "/",
			err:    nil,
			ecode:  http.StatusOK,
			ebody:  `{"Type":"success","Message":"Success."}`,
			mtype:  "success",
			msg:    "Success.",
		},
	}

	// We are going to reuse the context.
	c, err := appenginetesting.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	for i, test := range tests {
		// Make the request and writer.
		w := httptest.NewRecorder()
		r, err := http.NewRequest(test.method, test.url, nil)
		if err != nil {
			t.Fatal(err)
		}

		// Call the test function
		LogAndMessage(c, w, r, test.err, test.mtype, test.msg, test.ecode)

		// Check the status
		if w.Code != test.ecode {
			t.Errorf("(%v) expexted %v as response code. Got: %v",
				i, test.ecode, w.Code)
		}

		body := w.Body.String()
		if body != test.ebody {
			t.Errorf("(%v) expexted %v as response body. Got: %v",
				i, test.ebody, body)
		}
	}
}
