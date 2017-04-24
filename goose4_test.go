package goose4

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

type rw struct {
	headers http.Header
	body    string
	status  int
}

func Newrw() *rw {
	return &rw{headers: make(http.Header)}
}

func (r *rw) Header() http.Header {
	return r.headers
}

func (r *rw) Write(b []byte) (int, error) {
	if r.status == 0 {
		r.WriteHeader(200)
	}
	r.body = string(b)
	return len(b), nil
}

func (r *rw) WriteHeader(i int) {
	r.status = i
}

func TestNewGoose4(t *testing.T) {
	for _, test := range []struct {
		title       string
		config      Config
		expectError bool
	}{
		{"simple goose4 object", Config{}, false},
	} {
		t.Run(test.title, func(t *testing.T) {
			_, err := NewGoose4(test.config)

			if err == nil {
				if test.expectError {
					t.Errorf("NewGoose4: expected error, received none")
				}
			} else {
				if !test.expectError {
					t.Errorf("NewGoose4: %v", err)
				}
			}
		})
	}
}

func TestServeHTTP(t *testing.T) {
	var emptyOutput = `{"artifact_id":"","build_number":"","build_machine":"","built_by":"","built_when":"0001-01-01T00:00:00Z","compiler_version":"","git_sha1":"","runbook_uri":"","version":""}`

	for _, test := range []struct {
		path              string
		method            string
		tests             []Test
		expectStatusCode  int
		expectBody        string
		expectContentType string
		ignoreBody        bool // lulz
	}{
		{"/service/config", "GET", []Test{}, 200, emptyOutput, "application/json", false},
		{"/service/status", "GET", []Test{}, 200, "", "application/json", true},
		{"/service/healthcheck", "GET", []Test{}, 200, "", "application/json", true},
		{"/service/healthcheck/asg", "GET", []Test{}, 200, `"OK"`, "text/plain", false},
		{"/service/healthcheck/gtg", "GET", []Test{}, 200, `"OK"`, "text/plain", false},

		{"/service/config", "POST", []Test{}, 405, `{"status":405,"message":"Method \"POST\" not allowed"}`, "application/json", false},
		{"/service/floopydoop", "GET", []Test{}, 404, `{"status":404,"message":"No such route \"/service/floopydoop\""}`, "application/json", false},

		{"/service/healthcheck", "GET", []Test{{F: HealthTestFailure, Critical: true}}, 500, "", "application/json", true},
		{"/service/healthcheck/asg", "GET", []Test{{F: HealthTestFailure, Critical: true}}, 500, `"Bad"`, "text/plain", false},
		{"/service/healthcheck/gtg", "GET", []Test{{F: HealthTestFailure}}, 500, `"Bad"`, "text/plain", false},
	} {
		t.Run(fmt.Sprintf("%s %s", test.method, test.path), func(t *testing.T) {
			g, _ := NewGoose4(Config{})
			g.tests = test.tests
			w := Newrw()
			r := &http.Request{
				Method: test.method,
				URL:    &url.URL{Path: test.path},
			}
			g.ServeHTTP(w, r)

			t.Run("Status code", func(t *testing.T) {
				if w.status != test.expectStatusCode {
					t.Errorf("expected %d, received %d", test.expectStatusCode, w.status)
				}
			})

			t.Run("Body text", func(t *testing.T) {
				if w.body != test.expectBody && !test.ignoreBody {
					t.Errorf("expected %q, received %q", test.expectBody, w.body)
				}
			})

			t.Run("Content Type", func(t *testing.T) {
				ct := w.headers.Get("Content-Type")
				if ct != test.expectContentType {
					t.Errorf("expected %q, received %q", test.expectContentType, ct)
				}
			})
		})
	}
}

func TestAddTest(t *testing.T) {
	for _, test := range []struct {
		title        string
		testName     string
		testCritical bool
		testFunc     func() bool
	}{
		{"A simple, boring test", "a_test", true, func() bool { return true }},
	} {
		t.Run(test.title, func(t *testing.T) {
			t0 := Test{
				Name:     test.testName,
				Critical: test.testCritical,
				F:        test.testFunc,
			}

			g, _ := NewGoose4(Config{})
			g.AddTest(t0)

			t.Run("Test Name", func(t *testing.T) {
				if test.testName != g.tests[0].Name {
					t.Errorf("expected %q, received %q", test.testName, g.tests[0].Name)
				}
			})

			t.Run("Test Critical Value", func(t *testing.T) {
				if test.testCritical != g.tests[0].Critical {
					t.Errorf("expected %v, received %v", test.testCritical, g.tests[0].Critical)
				}
			})

			// t.Run("Test Function", func(t *testing.T) {
			//     if !reflect.DeepEqual(test.testFunc, g.tests[0].F) {
			//         t.Errorf("expected %v, received %v", test.testFunc, g.tests[0].F)
			//     }
			// })
		})
	}

}
