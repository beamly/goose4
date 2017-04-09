package goose4

import (
	"fmt"
	"log"
	"net/http"
)

// Goose4 holds goose4 configuration and provides functions thereon
type Goose4 struct {
	c Config
}

// NewGoose4 returns a Goose4 object to be used as net/http handler
func NewGoose4(c Config) (g Goose4, err error) {
	g.c = c

	return
}

// ServeHTTP is an http router to serve se4 endpoints
func (g Goose4) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body []byte
	var err error

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		body, err = Error{http.StatusMethodNotAllowed, fmt.Sprintf("Method %q not allowed", r.Method)}.Marshal()
	} else {
		switch r.URL.Path {
		case "/service/config":
			body, err = g.c.Marshal()
		default:
			w.WriteHeader(http.StatusNotFound)
			body, err = Error{http.StatusNotFound, fmt.Sprintf("No such route %q", r.URL.Path)}.Marshal()
		}
	}

	if err != nil {
		// This will nuke the original error; this is acceptable due to the risk of leaking
		// potentially sensitive information otherwise
		w.WriteHeader(http.StatusInternalServerError)
		body, err = Error{http.StatusInternalServerError, fmt.Sprint("Internal error")}.Marshal()

		log.Print(err)
	}

	fmt.Fprintf(w, string(body))
}
