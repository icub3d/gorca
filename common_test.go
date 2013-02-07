// Copyright 2013 Joshua Marsh. All rights reserved.  Use of this
// source code is governed by a BSD-style license that can be found in
// the LICENSE file.

package gorca

import (
	"appengine"
	"fmt"
	"github.com/icub3d/appenginetesting"
	"github.com/icub3d/testhelper"
	"net/http"
	"net/http/httptest"
	"testing"
)

type LogAndFunc func(c appengine.Context, w http.ResponseWriter,
	r *http.Request, err error)

func TestLogAnds(t *testing.T) {
	h := testhelper.New(t)

	// These are our test cases.
	tests := []struct {
		f      LogAndFunc
		fn     string
		method string
		url    string
		err    error
		ecode  int
		ebody  string
	}{
		// LogAndNotFound
		{
			f:      LogAndNotFound,
			fn:     "LogAndNotFound",
			method: "GET",
			url:    "/",
			err:    fmt.Errorf("no such file or directory"),
			ecode:  http.StatusNotFound,
			ebody:  `{"Type":"error","Message":"Not found."}`,
		},

		// LogAndFailed
		{
			f:      LogAndFailed,
			fn:     "LogAndFailed",
			method: "GET",
			url:    "/",
			err:    fmt.Errorf("oopsie, you failed"),
			ecode:  http.StatusBadRequest,
			ebody:  `{"Type":"error","Message":"Failed."}`,
		},

		// LogAndUnexpected
		{
			f:      LogAndUnexpected,
			fn:     "LogAndUnexpected",
			method: "GET",
			url:    "/",
			err:    fmt.Errorf("oopsie, i failed"),
			ecode:  http.StatusInternalServerError,
			ebody:  `{"Type":"error","Message":"Something unexpected happened."}`,
		},
	}

	// We are going to reuse the context.
	c, err := appenginetesting.NewContext(nil)
	h.FatalNotNil("creating contxt", err)
	defer c.Close()

	for i, test := range tests {
		h.SetIndex(i)

		// Make the request and writer.
		w := httptest.NewRecorder()
		r, err := http.NewRequest(test.method, test.url, nil)
		h.FatalNotNil("creating request", err)

		h.SetFunc(`%s(c, w, r, err("%v"))`, test.fn, test.err)

		// Call the test function
		test.f(c, w, r, test.err)

		// Check the values.
		h.ErrorNotEqual("response code", w.Code, test.ecode)
		h.ErrorNotEqual("response body", w.Body.String(), test.ebody)
	}
}

func TestLogAndMessage(t *testing.T) {
	h := testhelper.New(t)

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
	h.FatalNotNil("creating context", err)
	defer c.Close()

	for i, test := range tests {
		h.SetIndex(i)

		// Make the request and writer.
		w := httptest.NewRecorder()
		r, err := http.NewRequest(test.method, test.url, nil)
		h.FatalNotNil("creating request", err)

		h.SetFunc(`LogAndMessage(c, w, r, err("%v"), "%s, "%s", %d)`,
			test.err, test.mtype, test.msg, test.ecode)

		// Call the test function
		LogAndMessage(c, w, r, test.err, test.mtype, test.msg, test.ecode)

		// Check the values.
		h.ErrorNotEqual("response code", w.Code, test.ecode)
		h.ErrorNotEqual("response body", w.Body.String(), test.ebody)
	}
}
