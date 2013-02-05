// Copyright 2013 Joshua Marsh. All rights reserved.  Use of this
// source code is governed by a BSD-style license that can be found in
// the LICENSE file.

package gorca

import (
	"github.com/icub3d/appenginetesting"
	"github.com/icub3d/testhelper"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetUserOrUnexpected(t *testing.T) {
	h := testhelper.New(t)

	// We are going to reuse the context here.
	c, err := appenginetesting.NewContext(nil)
	h.FatalNotNil("creating context", err)
	defer c.Close()

	tests := []struct {
		f       func()
		success bool
		email   string
		admin   bool
		ecode   int
		ebody   string
	}{
		// Test login for an admin user.
		{
			f: func() {
				c.Login("test@example.com", true)
			},
			success: true,
			email:   "test@example.com",
			admin:   true,
			ecode:   0,
			ebody:   "",
		},

		// Test login for a non admin user.
		{
			f: func() {
				c.Logout()
				c.Login("test@example.com", false)
			},
			success: true,
			email:   "test@example.com",
			admin:   false,
			ecode:   0,
			ebody:   "",
		},

		// Test being logged out.
		{
			f: func() {
				c.Logout()
			},
			success: false,
			email:   "",
			admin:   false,
			ecode:   http.StatusInternalServerError,
			ebody:   `{"Type":"error","Message":"Something unexpected happened."}`,
		},
	}

	for i, test := range tests {
		h.SetIndex(i)

		// Prep the test.
		test.f()

		// Make the request and writer.
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "/", nil)
		h.FatalNotNil("creating request", err)

		// Call the normal function
		u, ok := GetUserOrUnexpected(c, w, r)
		h.FatalNotEqual("success", ok, test.success)

		if test.success {
			// Check the user credentials
			h.ErrorNotEqual("email", u.Email, test.email)
			h.ErrorNotEqual("admin", u.Admin, test.admin)
		} else {
			// Check the values.
			h.ErrorNotEqual("response code", w.Code, test.ecode)
			h.ErrorNotEqual("response body", w.Body.String(), test.ebody)
		}
	}
}

func TestGetUserLogoutURL(t *testing.T) {
	h := testhelper.New(t)

	c, err := appenginetesting.NewContext(nil)
	h.FatalNotNil("creating context", err)
	defer c.Close()

	// Make the request and writer.
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	h.FatalNotNil("creating request", err)

	url, ok := GetUserLogoutURL(c, w, r, "/")
	h.FatalNotEqual("getting url", ok, true)

	surl := "/_ah/login?continue=http%3A//127.0.0.1%3A"
	eurl := "/&action=Logout"
	if !strings.HasPrefix(url, surl) || !strings.HasSuffix(url, eurl) {
		t.Errorf("Expecting '%v' for url, but got: %v",
			surl+"[PORT]"+eurl, url)
	}
}
