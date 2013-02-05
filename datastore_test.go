// Copyright 2013 Joshua Marsh. All rights reserved.  Use of this
// source code is governed by a BSD-style license that can be found in
// the LICENSE file.

package gorca

import (
	"appengine"
	"appengine/datastore"
	"github.com/icub3d/appenginetesting"
	"github.com/icub3d/testhelper"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewKey(t *testing.T) {
	h := testhelper.New(t)

	// We are going to reuse the context.
	c, err := appenginetesting.NewContext(nil)
	h.FatalNotNil("creating contxt", err)
	defer c.Close()

	// Make the request and writer.
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/newkey", nil)
	h.FatalNotNil("creating request", err)

	skey, key, ok := NewKey(c, w, r, "Item", nil)

	h.ErrorNotEqual("new key", ok, true)
	h.ErrorNotEqual("keys", skey, key.Encode())
}

// This is merely used to hold the strings.
type stringer struct {
	s string
}

func TestPutStringKeys(t *testing.T) {
	// This also tests PutKeys.

	h := testhelper.New(t)

	// We are going to reuse the context.
	c, err := appenginetesting.NewContext(nil)
	h.FatalNotNil("creating contxt", err)
	defer c.Close()

	// These are the tests
	tests := []struct {
		keys   []string
		values []stringer
		expect bool
	}{
		// An empty list.
		{
			keys:   []string{},
			values: []stringer{},
			expect: true,
		},

		// A normal list.
		{
			keys:   []string{makeKey(c).Encode(), makeKey(c).Encode()},
			values: []stringer{stringer{"one"}, stringer{"two"}},
			expect: true,
		},

		// More keys than values.
		{
			keys:   []string{makeKey(c).Encode(), makeKey(c).Encode()},
			values: []stringer{stringer{"one"}},
			expect: false,
		},

		// More values than keys.
		{
			keys:   []string{makeKey(c).Encode(), makeKey(c).Encode()},
			values: []stringer{stringer{"one"}, stringer{"two"}, stringer{"three"}},
			expect: false,
		},

		// Invalid key.
		{
			keys:   []string{makeKey(c).Encode(), makeKey(c).Encode(), "hahaha"},
			values: []stringer{stringer{"one"}, stringer{"two"}, stringer{"three"}},
			expect: false,
		},
	}

	for i, test := range tests {
		h.SetIndex(i)

		// Make the request and writer.
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "/datastore", nil)
		h.FatalNotNil("creating request", err)

		result := PutStringKeys(c, w, r, test.keys, test.values)
		h.FatalNotEqual("put keys", result, test.expect)

		if test.expect == false {
			// Test the output.
			h.ErrorNotEqual("response code", w.Code, http.StatusInternalServerError)
			h.ErrorNotEqual("response body", w.Body.String(),
				`{"Type":"error","Message":"Something unexpected happened."}`)
		} else {
			for _, key := range test.keys {
				// Make sure each of the keys persisted.
				var value stringer
				k, err := datastore.DecodeKey(key)
				h.FatalNotNil("decoding key", err)

				err = datastore.Get(c, k, &value)
				h.FatalNotNil("datastore get", err)

				// High replication seems to make this impossible.
				// h.ErrorNotEqual("datastore value", test.values[j].s, value.s)
			}
		}
	}
}

func TestDeleteStringKeys(t *testing.T) {
	// This also tests DeleteKeys.

	h := testhelper.New(t)

	// We are going to reuse the context.
	c, err := appenginetesting.NewContext(nil)
	h.FatalNotNil("creating contxt", err)
	defer c.Close()

	// These are the tests
	tests := []struct {
		keys    []string
		values  []stringer
		inserts []bool
		deletes []bool
		expects []bool
	}{
		// A normal list.
		{
			keys: []string{
				makeKey(c).Encode(),
				makeKey(c).Encode(),
				makeKey(c).Encode(),
				makeKey(c).Encode(),
				makeKey(c).Encode(),
				makeKey(c).Encode(),
			},
			values: []stringer{
				stringer{"one"},
				stringer{"two"},
				stringer{"three"},
				stringer{"four"},
				stringer{"five"},
				stringer{"six"},
			},
			deletes: []bool{false, false, true, false, true, false},
		},
	}

	for i, test := range tests {
		h.SetIndex(i)

		// Make the request and writer.
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "/datastore", nil)
		h.FatalNotNil("creating request", err)

		// insert the items.
		ok := PutStringKeys(c, w, r, test.keys, test.values)
		h.FatalNotEqual("put keys", ok, true)

		// make the items to delete
		dkeys := make([]string, 0, 0)
		for j, key := range test.keys {
			if test.deletes[j] {
				dkeys = append(dkeys, key)
			}
		}

		// delete the items
		ok = DeleteStringKeys(c, w, r, dkeys)

		// Check each of the items.
		for j, key := range test.keys {
			var value stringer
			k, err := datastore.DecodeKey(key)
			h.FatalNotNil("decoding key", err)

			err = datastore.Get(c, k, &value)
			if test.deletes[j] {
				h.ErrorNil("deleted item", err)
			} else {
				h.ErrorNotNil("deleted item", err)
			}
		}
	}
}

func makeKey(c appengine.Context) *datastore.Key {
	id, _, _ := datastore.AllocateIDs(c, "Item", nil, 1)
	return datastore.NewKey(c, "Item", "", id, nil)
}
