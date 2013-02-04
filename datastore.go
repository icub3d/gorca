package gorca

import (
	"appengine"
	"appengine/datastore"
	"net/http"
)

// NewKey is a helper function that allocates a new id and uses it to
// make a new key. It returns both the string and struct version fo
// the key. If a failure occured, false is returned and a response was
// returned to the request. This case should be terminal.
func NewKey(c appengine.Context, w http.ResponseWriter, r *http.Request,
	kind string, parent *datastore.Key) (string, *datastore.Key, bool) {

	// Generate a new key for this kind.
	id, _, err := datastore.AllocateIDs(c, kind, parent, 1)
	if err != nil {
		LogAndUnexpected(c, w, r, err)
		return "", nil, false
	}
	key := datastore.NewKey(c, kind, "", id, parent)

	return key.Encode(), key, true
}

// PutStringKeys is a helper function that performs a PutMulti on the
// set of keys and values. If a failure occured, false is returned and
// a response was returned to the request. This case should be
// terminal.
func PutStringKeys(c appengine.Context, w http.ResponseWriter,
	r *http.Request, keys []string, values interface{}) bool {

	dkeys, ok := StringsToKeys(c, w, r, keys)
	if !ok {
		return false
	}

	return PutKeys(c, w, r, dkeys, values)
}

// PutKeys is a helper function the performs a PutMulti on the set of
// keys and values. If a failure occured, false is returned and a
// response was returned to the request. This case should be terminal.
func PutKeys(c appengine.Context, w http.ResponseWriter, r *http.Request,
	keys []*datastore.Key, values interface{}) bool {

	if _, err := datastore.PutMulti(c, keys, values); err != nil {
		LogAndUnexpected(c, w, r, err)
		return false
	}

	return true
}

// DeleteStringKeyAndAncestors is a helper function that remove the given
// key from the datastore as well as all of it's ancestors of the
// given kind. If a failure occured, false is returned and a response
// was returned to the request. This case should be terminal.
func DeleteStringKeyAndAncestors(c appengine.Context, w http.ResponseWriter,
	r *http.Request, kind string, key string) bool {

	// Decode the string version of the key.
	k, err := datastore.DecodeKey(key)
	if err != nil {
		LogAndUnexpected(c, w, r, err)
		return false
	}

	// Call the helper to do the deletions.
	DeleteKeyAndAncestors(c, w, r, kind, k)

	return true
}

// DeleteKeyAndAncestors is a helper function that remove the given
// key from the datastore as well as all of it's ancestors of the
// given kind. If a failure occured, false is returned and a response
// was returned to the request. This case should be terminal.
func DeleteKeyAndAncestors(c appengine.Context, w http.ResponseWriter,
	r *http.Request, kind string, key *datastore.Key) bool {

	// Get all of the ancestors.
	q := datastore.NewQuery(kind).Ancestor(key).KeysOnly()
	keys, err := q.GetAll(c, nil)
	if err != nil {
		LogAndUnexpected(c, w, r, err)
		return false
	}

	// Delete all the items and the list.
	keys = append(keys, key)
	if !DeleteKeys(c, w, r, keys) {
		return false
	}

	return true
}

// DeleteKeys is a helper function that removes all of the given
// keys from the datastore. If a failure occured, false is returned
// and a response was returned to the request. This case should be
// terminal.
func DeleteKeys(c appengine.Context, w http.ResponseWriter,
	r *http.Request, keys []*datastore.Key) bool {

	// Delete all the removed items.
	if err := datastore.DeleteMulti(c, keys); err != nil {
		LogAndUnexpected(c, w, r, err)
		return false
	}

	return true
}

// DeleteStringKeys is a helper function that converts the given
// strings into datastore keys and then calls DeleteKeyHelper on
// them. If a failure occured, false is returned and a response was
// returned to the request. This case should be terminal.
func DeleteStringKeys(c appengine.Context, w http.ResponseWriter, r *http.Request,
	keys []string) bool {

	dkeys, ok := StringsToKeys(c, w, r, keys)
	if !ok {
		return false
	}

	return DeleteKeys(c, w, r, dkeys)
}

// StringToKey is a helper function the turns a string into a
// datastore key. If a failure occured, false is returned and a
// response was returned to the request. This case should be terminal.
func StringToKey(c appengine.Context, w http.ResponseWriter,
	r *http.Request, key string) (*datastore.Key, bool) {

	k, err := datastore.DecodeKey(key)
	if err != nil {
		LogAndUnexpected(c, w, r, err)
		return nil, false
	}

	return k, true
}

// StringsToKeys is a helper function that turns a list of strings
// into a list of datastore keys. If a failure occured, false is
// returned and a response was returned to the request. This case
// should be terminal.
func StringsToKeys(c appengine.Context, w http.ResponseWriter,
	r *http.Request, keys []string) ([]*datastore.Key, bool) {

	dkeys := make([]*datastore.Key, 0, len(keys))
	for _, k := range keys {
		key, ok := StringToKey(c, w, r, k)
		if !ok {
			return nil, false
		}

		dkeys = append(dkeys, key)
	}

	return dkeys, true
}
