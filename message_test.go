package gorca

import (
	"fmt"
	"github.com/icub3d/appenginetesting"
	"github.com/icub3d/testhelper"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	h := testhelper.New(t)

	tests := []struct {
		data     interface{}
		code     int
		response string
	}{
		// Test the error case.
		{
			data: struct {
				C complex128
			}{
				C: complex(1, 1),
			},
			code:     http.StatusInternalServerError,
			response: `{"Type":"error","Message":"Something unexpected happened."}`,
		},

		// Test a normal case.
		{
			data: struct {
				C int64
			}{
				C: 123,
			},
			code:     http.StatusOK,
			response: `{"C":123}`,
		},
	}

	// We can use the same context for all tests.
	c, err := appenginetesting.NewContext(nil)
	h.FatalNotNil("creating context", err)
	defer c.Close()

	for i, test := range tests {
		h.SetIndex(i)

		// Make the request and writer.
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "/", nil)
		h.FatalNotNil("creating request", err)

		WriteJSON(c, w, r, test.data)

		// Check the values.
		h.ErrorNotEqual("response code", w.Code, test.code)
		h.ErrorNotEqual("response body", w.Body.String(), test.response)
	}
}

func TestWriteMessage(t *testing.T) {
	// This is already heavily testing in common_test.go
}

func TestWriteSuccessMessage(t *testing.T) {
	h := testhelper.New(t)

	tests := []struct {
		ecode int
		ebody string
	}{
		// Test A success message. The rest are tested above.
		{
			ecode: http.StatusOK,
			ebody: `{"Type":"success","Message":"Success."}`,
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
		r, err := http.NewRequest("GET", "/", nil)
		h.FatalNotNil("creating request", err)

		// Call the test function
		WriteSuccessMessage(c, w, r)

		// Check the values.
		h.ErrorNotEqual("response code", w.Code, test.ecode)
		h.ErrorNotEqual("response body", w.Body.String(), test.ebody)
	}
}

func TestWriteResponse(t *testing.T) {
	h := testhelper.New(t)

	tests := []struct {
		ecode int
		ebody string
		body  string
	}{
		// Test A success message. The rest are tested above.
		{
			ecode: http.StatusOK,
			ebody: `hello world`,
			body:  `hello world`,
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
		r, err := http.NewRequest("GET", "/", nil)
		h.FatalNotNil("creating request", err)

		// Call the test function
		WriteResponse(c, w, r, []byte(test.body))

		// Check the values.
		h.ErrorNotEqual("response code", w.Code, test.ecode)
		h.ErrorNotEqual("response body", w.Body.String(), test.ebody)
	}
}

func TestUnmarshalOrFail(t *testing.T) {
	h := testhelper.New(t)

	tests := []struct {
		result bool
		ecode  int
		ebody  string
		bytes  []byte
		where  interface{}
	}{
		// Test a good unmarshal.
		{
			result: true,
			ecode:  0,
			ebody:  "",
			bytes:  []byte(`{"C":123}`),
			where:  &struct{ C int64 }{},
		},

		// Test a bad unmarshal.
		{
			result: false,
			ecode:  http.StatusBadRequest,
			ebody:  `{"Type":"error","Message":"Failed."}`,
			bytes:  []byte(`{"C":12`),
			where:  &struct{ C int64 }{},
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
		r, err := http.NewRequest("GET", "/", nil)
		h.FatalNotNil("creating request", err)

		// Call the test function
		results := UnmarshalOrFail(c, w, r, test.bytes, test.where)

		h.FatalNotEqual("umarshal results", results, test.result)

		// We don't get anything back on the wire if it succeeded.
		if test.result == true {
			continue
		}

		// Check the values.
		h.ErrorNotEqual("response code", w.Code, test.ecode)
		h.ErrorNotEqual("response body", w.Body.String(), test.ebody)
	}
}

func TestGetBodyOrFail(t *testing.T) {
	h := testhelper.New(t)

	tests := []struct {
		body    string
		result  bool
		ecode   int
		ebody   string
		request *http.Request
	}{
		// Test a nil body.
		{
			body:    "",
			result:  false,
			ecode:   http.StatusBadRequest,
			ebody:   `{"Type":"error","Message":"Failed."}`,
			request: newRequest("GET", "/", nil),
		},

		// Test a body that will fail reading.
		{
			body:    "",
			result:  false,
			ecode:   http.StatusInternalServerError,
			ebody:   `{"Type":"error","Message":"Something unexpected happened."}`,
			request: newRequest("GET", "/", ErrorReader{}),
		},

		// Test a normal body.
		{
			body:    "hello world",
			result:  true,
			ecode:   0,
			ebody:   "",
			request: newRequest("GET", "/", strings.NewReader("hello world")),
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

		// Call the test function
		body, results := GetBodyOrFail(c, w, test.request)

		h.FatalNotEqual("umarshal results", results, test.result)

		// If we failed, we should test the wire.
		if test.result == false {
			// Check the values.
			h.ErrorNotEqual("response code", w.Code, test.ecode)
			h.ErrorNotEqual("response body", w.Body.String(), test.ebody)
		} else {
			// We should test the body.
			h.ErrorNotEqual("body", string(body), test.body)
		}
	}
}

func TestUnmarshalFromBodyOrFail(t *testing.T) {
	h := testhelper.New(t)

	tests := []struct {
		body    string
		result  bool
		ecode   int
		ebody   string
		request *http.Request
		where   interface{}
	}{
		// Test a nil body.
		{
			result:  false,
			ecode:   http.StatusBadRequest,
			ebody:   `{"Type":"error","Message":"Failed."}`,
			request: newRequest("GET", "/", nil),
			where:   nil,
		},

		// Test a body that will fail reading.
		{
			result:  false,
			ecode:   http.StatusInternalServerError,
			ebody:   `{"Type":"error","Message":"Something unexpected happened."}`,
			request: newRequest("GET", "/", ErrorReader{}),
			where:   nil,
		},

		// Test a normal body.
		{
			result:  true,
			ecode:   0,
			ebody:   "",
			request: newRequest("GET", "/", strings.NewReader(`{"C":123}`)),
			where:   &struct{ C int64 }{},
		},

		// Test a JSON decoding error.
		{
			result:  false,
			ecode:   http.StatusBadRequest,
			ebody:   `{"Type":"error","Message":"Failed."}`,
			request: newRequest("GET", "/", strings.NewReader(`{"C":12`)),
			where:   &struct{ C int64 }{},
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

		// Call the test function
		results := UnmarshalFromBodyOrFail(c, w, test.request, test.where)

		h.FatalNotEqual("umarshal results", results, test.result)

		// If we failed, we should test the wire.
		if test.result == false {
			// Check the values.
			h.ErrorNotEqual("response code", w.Code, test.ecode)
			h.ErrorNotEqual("response body", w.Body.String(), test.ebody)
		}
	}
}

// newRequest is a helper function that wraps away they returned
// error.
func newRequest(method, url string, body io.Reader) *http.Request {
	r, _ := http.NewRequest(method, url, body)

	return r
}

// ErrorRedader implements the Reader interface and always errors our
// on the first read.
type ErrorReader struct{}

// Read returns an error when called.
func (e ErrorReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("I blew up on purpose!")
}
