package main

import (
	"net/http"
	"time"

	"github.com/zeebox/goose4"
)

func main() {
	// Create a simple config for goose4. In reality, though, this stuff
	// is probably more useful created at compile time. See:
	// https://jamescun.com/golang/compile-ldflags/ for easily consumed docs
	c := goose4.Config{
		ArtifactID:      "some-artifact",
		BuildNumber:     "123",
		BuildMachine:    "localhost",
		BuiltBy:         "ci-user",
		BuiltWhen:       time.Now(),
		CompilerVersion: "go version go1.7.4 darwin/amd64",
		GitSha:          "32b619ba997dfbfafd528ae3fea4e2cba8116be8",
		RunbookURI:      "https://example.com/goose4_runbook.html",
		Version:         "v0.0.1",
	}

	// Create a Goose4 handler
	se4, err := goose4.NewGoose4(c)
	if err != nil {
		panic(err)
	}

	// Create and add a healthcheck test pointing to a real function
	someCheck := goose4.Test{
		Name:     "Check something works or something",
		Critical: true,       // We want to remove this instance/server from the ASG if this fails
		F:        dummyCheck, // Note: no brackets
	}
	se4.AddTest(someCheck)

	// Add a truly anonymous function
	anotherCheck := goose4.Test{
		Name:     "Some important thing",
		Critical: true, // If the test fails then just don't send traffic to it
		F:        func() bool { return true },
	}
	se4.AddTest(anotherCheck)

	// Mount Goose4 handler for all se4 routes
	http.Handle("/service/", se4)
	panic(http.ListenAndServe(":8000", nil))
}

func dummyCheck() bool {
	time.Sleep(1 * time.Second)
	return true
}
