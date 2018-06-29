package goose4

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Goose4 holds goose4 configuration and provides functions thereon
type Goose4 struct {
	config Config
	boot   time.Time

	tests []Test
}

// NewGoose4 returns a Goose4 object to be used as net/http handler
func NewGoose4(c Config) (g Goose4, err error) {
	g.config = c
	g.boot = time.Now()

	return
}

// AddTest updates a Goose4 test list for healthchecks. These tests are used
// to determine whether a service is up or not
func (g *Goose4) AddTest(t Test) {
	g.tests = append(g.tests, t)
}

// ServeHTTP is an http router to serve se4 endpoints
func (g Goose4) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body []byte
	var errs bool
	var err error

	w.Header().Set("access-control-allow-origin", "*")
	w.Header().Set("access-control-allow-headers", "origin, content-type, accept")
	w.Header().Set("access-control-allow-methods", "GET")

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		body, err = Error{http.StatusMethodNotAllowed, fmt.Sprintf("Method %q not allowed", r.Method)}.Marshal()
	} else {
		switch r.URL.Path {
		case "/service/config":
			body, err = g.config.Marshal()
		case "/service/status":
			body, err = Status{Config: g.config}.Marshal(g.boot)
		case "/service/healthcheck":
			h := NewHealthcheck(g.tests)
			body, errs, err = h.All()

			if errs {
				w.WriteHeader(http.StatusInternalServerError)
			}
		case "/service/healthcheck/gtg":
			w.Header().Set("Content-Type", "text/plain")

			h := NewHealthcheck(g.tests)
			_, errs, err = h.GTG()

			if errs {
				w.WriteHeader(http.StatusInternalServerError)
				body = []byte(`"Bad"`)
			} else {
				body = []byte(`"OK"`)
			}

		case "/service/healthcheck/asg":
			w.Header().Set("Content-Type", "text/plain")
			h := NewHealthcheck(g.tests)
			_, errs, err = h.ASG()

			if errs {
				w.WriteHeader(http.StatusInternalServerError)
				body = []byte(`"Bad"`)
			} else {
				body = []byte(`"OK"`)
			}

		default:
			w.WriteHeader(http.StatusNotFound)
			body, err = Error{http.StatusNotFound, fmt.Sprintf("No such route %q", r.URL.Path)}.Marshal()
		}
	}

	if err != nil {
		log.Print(err)

		// This will nuke the original error; this is acceptable due to the risk of leaking
		// potentially sensitive information otherwise
		w.WriteHeader(http.StatusInternalServerError)
		body, err = Error{http.StatusInternalServerError, fmt.Sprint("Internal error")}.Marshal()
	}

	fmt.Fprintf(w, string(body))
}
