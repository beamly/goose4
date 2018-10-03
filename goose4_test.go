package goose4

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"testing"
)

type rw struct {
	headers http.Header
	body    string
	status  int
}

func newrw() *rw {
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

		{"/service/healthcheck", "GET", []Test{{F: HealthTestFailure, RequiredForASG: true, RequiredForGTG: true}}, 500, "", "application/json", true},
		{"/service/healthcheck/asg", "GET", []Test{{F: HealthTestFailure, RequiredForASG: true}}, 500, `"Bad"`, "text/plain", false},
		{"/service/healthcheck/gtg", "GET", []Test{{F: HealthTestFailure, RequiredForGTG: true}}, 500, `"Bad"`, "text/plain", false},
	} {
		t.Run(fmt.Sprintf("%s %s", test.method, test.path), func(t *testing.T) {
			g, _ := NewGoose4(Config{})
			g.tests = test.tests
			w := newrw()
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
		title              string
		testName           string
		testRequiredForASG bool
		testRequiredForGTG bool
		testFunc           func() bool
	}{
		{title: "A simple, boring passing test",
			testName:           "a_test",
			testRequiredForASG: true,
			testRequiredForGTG: true,
			testFunc:           func() bool { return true },
		},
		{
			title:              "A simple, boring non-passing test",
			testName:           "another_test",
			testRequiredForASG: true,
			testRequiredForGTG: true,
			testFunc:           func() bool { return false },
		},
	} {
		t.Run(test.title, func(t *testing.T) {
			t0 := Test{
				Name:           test.testName,
				RequiredForASG: test.testRequiredForASG,
				RequiredForGTG: test.testRequiredForGTG,
				F:              test.testFunc,
			}

			g, _ := NewGoose4(Config{})
			g.AddTest(t0)

			t.Run("Test Name", func(t *testing.T) {
				if test.testName != g.tests[0].Name {
					t.Errorf("expected %q, received %q", test.testName, g.tests[0].Name)
				}
			})

			t.Run("Test Required For ASG Value", func(t *testing.T) {
				if test.testRequiredForASG != g.tests[0].RequiredForASG {
					t.Errorf("expected %v, received %v", test.testRequiredForASG, g.tests[0].RequiredForASG)
				}
			})

			t.Run("Test Required For GTG Value", func(t *testing.T) {
				if test.testRequiredForGTG != g.tests[0].RequiredForGTG {
					t.Errorf("expected %v, received %v", test.testRequiredForGTG, g.tests[0].RequiredForGTG)
				}
			})

			// This really checks that g.tests[0] has the same name as testFunc; it doesn't
			// check whether the code is the same.
			//
			// We can, though, be reasonably confident this is good enough: they're both
			// anonymous functions that go internally names.
			t.Run("Test Function", func(t *testing.T) {
				t0Name := runtime.FuncForPC(reflect.ValueOf(test.testFunc).Pointer()).Name()
				t1Name := runtime.FuncForPC(reflect.ValueOf(g.tests[0].F).Pointer()).Name()

				t.Run("Valid Function", func(t *testing.T) {
					if t0Name != t1Name {
						t.Errorf("expected %v, received %v", t0Name, t1Name)
					}
				})

				t.Run("Invalid Function", func(t *testing.T) {
					fF := func() bool { return true }
					fFName := runtime.FuncForPC(reflect.ValueOf(fF).Pointer()).Name()

					if fFName == t1Name {
						t.Errorf("expected %v, received %v", fFName, t1Name)
					}
				})

				t.Run("Mutated Function", func(t *testing.T) {
					g.tests[0].F = func() bool { return false }
					t0Name := runtime.FuncForPC(reflect.ValueOf(test.testFunc).Pointer()).Name()

					// Validate that go isn't assigning the same name to stuff
					if t0Name != t1Name {
						t.Errorf("Received name %v which is a duplicate", t0Name)
					}
				})
			})
		})
	}
}
