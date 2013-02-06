// Copyright 2013 Joshua Marsh. All rights reserved.  Use of this
// source code is governed by a BSD-style license that can be found in
// the LICENSE file.

package gorca

import (
	"fmt"
	"github.com/icub3d/appenginetesting"
	"net/http"
	"net/http/httptest"
)

func ExampleLogAndNotFound() {
	// Note: The LogAnd* functions all work in a similar fashion.

	// Create an appengine context (in this case a mock one).
	c, err := appenginetesting.NewContext(nil)
	if err != nil {
		fmt.Println("creating context", err)
		return
	}
	defer c.Close()

	// Make the request and writer.
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/example/test", nil)
	if err != nil {
		fmt.Println("creating request", err)
		return
	}

	// Simulate an error
	err = fmt.Errorf("WHERE ARE THEY TRYING TO GO? LOL")
	LogAndNotFound(c, w, r, err)

	// Print out what would be returned down the pipe.
	fmt.Println("Response Code:", w.Code)
	fmt.Println("Response Body:", w.Body.String())

	// Output:
	// Response Code: 404
	// Response Body: {"Type":"error","Message":"Not found."}
}

func ExampleGetUserOrUnexpected() {
	// Create an appengine context (in this case a mock one).
	c, err := appenginetesting.NewContext(nil)
	if err != nil {
		fmt.Println("creating context", err)
		return
	}
	defer c.Close()

	// Make the request and writer.
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/example/test", nil)
	if err != nil {
		fmt.Println("creating request", err)
		return
	}

	// Get the user (no one is logged in, so it should return an
	// unexpected).
	u, ok := GetUserOrUnexpected(c, w, r)
	if !ok {
		// Because no one is logged in, we are going to get an unexpected
		// error.Print out what would be returned down the pipe.
		fmt.Println("Response Code:", w.Code)
		fmt.Println("Response Body:", w.Body.String())
	}

	// Now simulate a login and get that user.
	c.Login("test@example.com", true)
	u, ok = GetUserOrUnexpected(c, w, r)
	if !ok {
		fmt.Println("getting user", err)
		return
	}

	fmt.Println("Logged In User:", u.Email)
	fmt.Println("Is Admin:", u.Admin)

	// Output:
	// Response Code: 500
	// Response Body: {"Type":"error","Message":"Something unexpected happened."}
	// Logged In User: test@example.com
	// Is Admin: true
}
