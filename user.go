package gorca

import (
	"appengine"
	"appengine/user"
	"fmt"
	"net/http"
)

// GetUserOrUnexpected fetches the currently logged in user and
// returns it. The bool returns determines if the get was
// successful. If not, a JSON "unexpected" message is sent as the
// response. That case should terminate your response processing.
func GetUserOrUnexpected(w http.ResponseWriter,
	r *http.Request) (*user.User, bool) {

	cxt := appengine.NewContext(r)

	// Get the current user.
	u := user.Current(cxt)
	if u == nil {
		LogAndUnexpected(w, r,
			fmt.Errorf("no user found, but auth is required."))
		return nil, false
	}

	return u, true

}

// GetUserLogoutURL fetches the currently logged in user's LogoutURL
// and returns it. The bool returns determines if the get was
// successful. If not, a JSON "unexpected" message is sent as the
// response. That case should terminate your response processing.
func GetUserLogoutURL(w http.ResponseWriter, r *http.Request,
	dest string) (string, bool) {

	cxt := appengine.NewContext(r)

	// Get their logout URL.
	logout, err := user.LogoutURL(cxt, dest)
	if err != nil {
		LogAndUnexpected(w, r, fmt.Errorf("calling LogoutURL: %s", err))
		return "", false
	}

	return logout, true
}
