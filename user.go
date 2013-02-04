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
func GetUserOrUnexpected(c appengine.Context, w http.ResponseWriter,
	r *http.Request) (*user.User, bool) {

	// Get the current user.
	u := user.Current(c)
	if u == nil {
		LogAndUnexpected(c, w, r,
			fmt.Errorf("no user found, but auth is required."))
		return nil, false
	}

	return u, true

}

// GetUserLogoutURL fetches the currently logged in user's LogoutURL
// and returns it. The bool returns determines if the get was
// successful. If not, a JSON "unexpected" message is sent as the
// response. That case should terminate your response processing.
func GetUserLogoutURL(c appengine.Context, w http.ResponseWriter,
	r *http.Request, dest string) (string, bool) {

	// Get their logout URL.
	logout, err := user.LogoutURL(c, dest)
	if err != nil {
		LogAndUnexpected(c, w, r, fmt.Errorf("calling LogoutURL: %s", err))
		return "", false
	}

	return logout, true
}
