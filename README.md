

# goose4
`import "github.com/zeebox/goose4"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Subdirectories](#pkg-subdirectories)

## <a name="pkg-overview">Overview</a>
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
	    RunbookURI: "<a href="https://example.com/goose4_runbook.html">https://example.com/goose4_runbook.html</a>",
	    Version: "v0.0.1",
	}
	se4, err := goose4.NewGoose4(c)

Mounting se4 is just as easy:


	http.Handle("/service/", se4)
	panic(http.ListenAndServe(":80", nil))




## <a name="pkg-index">Index</a>
* [type Config](#Config)
  * [func (c Config) Marshal() (j []byte, err error)](#Config.Marshal)
* [type Error](#Error)
  * [func (e Error) Marshal() ([]byte, error)](#Error.Marshal)
* [type Goose4](#Goose4)
  * [func NewGoose4(c Config) (g Goose4, err error)](#NewGoose4)
  * [func (g Goose4) ServeHTTP(w http.ResponseWriter, r *http.Request)](#Goose4.ServeHTTP)
* [type Status](#Status)
  * [func (s Status) Marshal(boot time.Time) ([]byte, error)](#Status.Marshal)
* [type System](#System)
  * [func NewSystem(boot time.Time) System](#NewSystem)


#### <a name="pkg-files">Package files</a>
[config.go](/src/github.com/zeebox/goose4/config.go) [doc.go](/src/github.com/zeebox/goose4/doc.go) [error.go](/src/github.com/zeebox/goose4/error.go) [goose4.go](/src/github.com/zeebox/goose4/goose4.go) [status.go](/src/github.com/zeebox/goose4/status.go) 






## <a name="Config">type</a> [Config](/src/target/config.go?s=196:647#L1)
``` go
type Config struct {
    ArtifactID      string    `json:"artifact_id"`
    BuildNumber     string    `json:"build_number"`
    BuildMachine    string    `json:"build_machine"`
    BuiltBy         string    `json:"built_by"`
    BuiltWhen       time.Time `json:"built_when"`
    CompilerVersion string    `json:"compiler_version"`
    GitSha          string    `json:"git_sha1"`
    RunbookURI      string    `json:"runbook_uri"`
    Version         string    `json:"version"`
}
```
Config implements a subset of <a href="https://github.com/beamly/SE4/blob/master/SE4.md#status">https://github.com/beamly/SE4/blob/master/SE4.md#status</a>
and is used to configure static values for goose4.










### <a name="Config.Marshal">func</a> (Config) [Marshal](/src/target/config.go?s=768:815#L14)
``` go
func (c Config) Marshal() (j []byte, err error)
```
Marshal returns a json document and, potentially, an error in order to
respond with configuration for a service.




## <a name="Error">type</a> [Error](/src/target/error.go?s=118:204#L1)
``` go
type Error struct {
    Status  int    `json:"status"`
    Message string `json:"message"`
}
```
Error is a simple placeholder to store non-2xx, non-3xx response data










### <a name="Error.Marshal">func</a> (Error) [Marshal](/src/target/error.go?s=245:285#L4)
``` go
func (e Error) Marshal() ([]byte, error)
```
Marshal wraps an Error in some json




## <a name="Goose4">type</a> [Goose4](/src/target/goose4.go?s=130:185#L1)
``` go
type Goose4 struct {
    // contains filtered or unexported fields
}
```
Goose4 holds goose4 configuration and provides functions thereon







### <a name="NewGoose4">func</a> [NewGoose4](/src/target/goose4.go?s=255:301#L7)
``` go
func NewGoose4(c Config) (g Goose4, err error)
```
NewGoose4 returns a Goose4 object to be used as net/http handler





### <a name="Goose4.ServeHTTP">func</a> (Goose4) [ServeHTTP](/src/target/goose4.go?s=405:470#L15)
``` go
func (g Goose4) ServeHTTP(w http.ResponseWriter, r *http.Request)
```
ServeHTTP is an http router to serve se4 endpoints




## <a name="Status">type</a> [Status](/src/target/status.go?s=247:285#L5)
``` go
type Status struct {
    Config
    System
}
```
Status embeds Config and System to give a concise system status










### <a name="Status.Marshal">func</a> (Status) [Marshal](/src/target/status.go?s=378:433#L12)
``` go
func (s Status) Marshal(boot time.Time) ([]byte, error)
```
Marshal returns a status doc based on passed in config and up-to-date
system details




## <a name="System">type</a> [System](/src/target/status.go?s=599:942#L22)
``` go
type System struct {
    MachineName string `json:"machine_name"`
    OSArch      string `json:"os_arch"`
    OSLoad      string `json:"os_avgload"`
    OSName      string `json:"os_name"`
    OSProcs     string `json:"os_numprocessors"`
    OSVersion   string `json:"os_version"`
    UpDuration  string `json:"up_duration"`
    UpSince     string `json:"up_since"`
}
```
System contains system specific data for status responses







### <a name="NewSystem">func</a> [NewSystem](/src/target/status.go?s=944:981#L33)
``` go
func NewSystem(boot time.Time) System
```








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
