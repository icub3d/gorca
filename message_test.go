package gorca

import (
	"fmt"
	"github.com/icub3d/appenginetesting"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteJSON(t *testing.T) {
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
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	for i, test := range tests {

		// Make the request and writer.
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		WriteJSON(c, w, r, test.data)

		// Check the status
		if w.Code != test.code {
			t.Errorf("(%v) expexted %v as response code. Got: %v",
				i, test.code, w.Code)
		}

		body := w.Body.String()
		if body != test.response {
			t.Errorf("(%v) expexted %v as response body. Got: %v",
				i, test.response, body)
		}

	}

}

func TestWriteMessage(t *testing.T) {
	// This is already heavily testing in common_test.go
}

func TestWriteSuccessMessage(t *testing.T) {

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
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	for i, test := range tests {
		// Make the request and writer.
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Call the test function
		WriteSuccessMessage(c, w, r)

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

func TestWriteResponse(t *testing.T) {

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
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	for i, test := range tests {
		// Make the request and writer.
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Call the test function
		WriteResponse(c, w, r, []byte(test.body))

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

func TestUnmarshalOrFail(t *testing.T) {

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
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	for i, test := range tests {
		// Make the request and writer.
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Call the test function
		results := UnmarshalOrFail(c, w, r, test.bytes, test.where)

		if results != test.result {
			t.Fatalf("(%v) expexted %v as result. Got: %v",
				i, test.result, results)
		}

		// We don't get anything back on the wire if it succeeded.
		if test.result == true {
			continue
		}

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

func TestGetBodyOrFail(t *testing.T) {

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
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	for i, test := range tests {
		// Make the request and writer.
		w := httptest.NewRecorder()

		// Call the test function
		body, results := GetBodyOrFail(c, w, test.request)

		if results != test.result {
			t.Fatalf("(%v) expexted %v as result. Got: %v",
				i, test.result, results)
		}

		// If we failed, we should test the wire.
		if test.result == false {
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
		} else {
			// We should test the body.
			if string(body) != test.body {
				t.Errorf("(%v) expexted %v as response body. Got: %v",
					i, test.body, body)
			}
		}
	}
}

func TestUnmarshalFromBodyOrFail(t *testing.T) {

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
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	for i, test := range tests {
		// Make the request and writer.
		w := httptest.NewRecorder()

		// Call the test function
		results := UnmarshalFromBodyOrFail(c, w, test.request, test.where)

		if results != test.result {
			t.Fatalf("(%v) expexted %v as result. Got: %v",
				i, test.result, results)
		}

		// If we failed, we should test the wire.
		if test.result == false {
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
