/*
Package goose4 provides a golang implemenation of the se4 spec. It is a pure golang implementation with no/few extra dependencies.

It provides a `net/http` drop in HTTP Function which will route and provide bits of information.

Initialisation is reasonably simple:

    import (
        "net/http"
        "github.com/zeebox/goose4"
    )

    c := goose4.Config{
        ArtifactID: "some-artifact",
        BuildNumber: "123",
        BuildMachine: "localhost",
        BuiltBy: "ci-user",
        BuiltWhen: Time.now(),
        CompilerVersion: "go version go1.7.4 darwin/amd64",
        GitSha: "32b619ba997dfbfafd528ae3fea4e2cba8116be8",
        RunbookURI: "https://example.com/goose4_runbook.html",
        Version: "v0.0.1",
    }
    se4, err := goose4.NewGoose4(c)


Mounting se4 is just as easy:

    http.Handle("/service/", se4)
    panic(http.ListenAndServe(":80", nil))

*/
package goose4
