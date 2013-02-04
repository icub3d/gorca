package gorca

import (
	"github.com/icub3d/appenginetesting"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetUserOrUnexpected(t *testing.T) {
	c, err := appenginetesting.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	// Test loggin in.
	c.Login("test@example.com", true)

	// Make the request and writer.
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	u, ok := GetUserOrUnexpected(c, w, r)
	if !ok {
		t.Fatal("getting user")
	}

	e := "test@example.com"
	if u.Email != e {
		t.Errorf("expected '%v' for email, but got: %v", e, u.Email)
	}

	if !u.Admin {
		t.Errorf("expected '' for email, but got: %v", true, u.Admin)
	}

	// Now see what happens when we log out.
	c.Logout()
	_, ok = GetUserOrUnexpected(c, w, r)
	if ok {
		t.Errorf("expected !ok for logged out user, but got: ok")
	}
}

func TestGetUserLogoutURL(t *testing.T) {
	c, err := appenginetesting.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	// Make the request and writer.
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	url, ok := GetUserLogoutURL(c, w, r, "/")
	if !ok {
		t.Fatal("getting url")
	}

	surl := "/_ah/login?continue=http%3A//127.0.0.1%3A"
	eurl := "/&action=Logout"
	if !strings.HasPrefix(url, surl) || !strings.HasSuffix(url, eurl) {
		t.Errorf("Expecting '%v' for url, but got: %v",
			surl+"[PORT]"+eurl, url)
	}
}
