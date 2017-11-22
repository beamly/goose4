package goose4

import (
	"testing"
	"net/http"
	"fmt"
	"net/url"
	"net/http/httptest"
)

const (
	TestResponseBody = "hello, world!"
)

var (
	TestURLCatch, _ = url.Parse("https://user:pass@example.com/service/")
)

type TestAPI struct{}

func (ta TestAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, TestResponseBody)
}

func TestMiddleware(t *testing.T) {
	t.Run("Interface", testMiddlewareInterface)
	t.Run("RouteCatch", testMiddlewareRouteCatch)
}

func testMiddlewareInterface(t *testing.T) {
	var mw http.Handler
	mw = NewMiddleware(nil)
	if mw == nil {
		t.Errorf("NewMiddleWare: Interface not correctly implemented")
	}
}

func testMiddlewareRouteCatch(t *testing.T) {
	mw := NewMiddleware(TestAPI{})
	r := &http.Request{URL: TestURLCatch}

	rec := httptest.NewRecorder()
	mw.ServeHTTP(rec, r)

	fmt.Println(mw)
}
