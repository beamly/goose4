package main

import (
	"net/http"
	"time"

	"github.com/zeebox/goose4"
)

func main() {
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

	se4, err := goose4.NewGoose4(c)
	if err != nil {
		panic(err)
	}

	http.Handle("/service/", se4)
	panic(http.ListenAndServe(":8000", nil))
}
