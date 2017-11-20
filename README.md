

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
  * [func (g *Goose4) AddTest(t Test)](#Goose4.AddTest)
  * [func (g Goose4) ServeHTTP(w http.ResponseWriter, r *http.Request)](#Goose4.ServeHTTP)
* [type Healthcheck](#Healthcheck)
  * [func NewHealthcheck(t []Test) Healthcheck](#NewHealthcheck)
  * [func (h *Healthcheck) ASG() (output []byte, errors bool, err error)](#Healthcheck.ASG)
  * [func (h *Healthcheck) All() (output []byte, errors bool, err error)](#Healthcheck.All)
  * [func (h *Healthcheck) GTG() (output []byte, errors bool, err error)](#Healthcheck.GTG)
* [type Middleware](#Middleware)
  * [func NewMiddleware(h http.Handler) *Middleware](#NewMiddleware)
  * [func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request)](#Middleware.ServeHTTP)
* [type Status](#Status)
  * [func (s Status) Marshal(boot time.Time) ([]byte, error)](#Status.Marshal)
* [type System](#System)
  * [func NewSystem(boot time.Time) System](#NewSystem)
* [type Test](#Test)


#### <a name="pkg-files">Package files</a>
[config.go](/src/github.com/zeebox/goose4/config.go) [doc.go](/src/github.com/zeebox/goose4/doc.go) [error.go](/src/github.com/zeebox/goose4/error.go) [goose4.go](/src/github.com/zeebox/goose4/goose4.go) [healthcheck.go](/src/github.com/zeebox/goose4/healthcheck.go) [middleware.go](/src/github.com/zeebox/goose4/middleware.go) [status.go](/src/github.com/zeebox/goose4/status.go) 






## <a name="Config">type</a> [Config](/src/target/config.go?s=196:647#L10)
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










### <a name="Config.Marshal">func</a> (Config) [Marshal](/src/target/config.go?s=768:815#L24)
``` go
func (c Config) Marshal() (j []byte, err error)
```
Marshal returns a json document and, potentially, an error in order to
respond with configuration for a service.




## <a name="Error">type</a> [Error](/src/target/error.go?s=118:204#L8)
``` go
type Error struct {
    Status  int    `json:"status"`
    Message string `json:"message"`
}
```
Error is a simple placeholder to store non-2xx, non-3xx response data










### <a name="Error.Marshal">func</a> (Error) [Marshal](/src/target/error.go?s=245:285#L14)
``` go
func (e Error) Marshal() ([]byte, error)
```
Marshal wraps an Error in some json




## <a name="Goose4">type</a> [Goose4](/src/target/goose4.go?s=130:200#L11)
``` go
type Goose4 struct {
    // contains filtered or unexported fields
}
```
Goose4 holds goose4 configuration and provides functions thereon







### <a name="NewGoose4">func</a> [NewGoose4](/src/target/goose4.go?s=270:316#L19)
``` go
func NewGoose4(c Config) (g Goose4, err error)
```
NewGoose4 returns a Goose4 object to be used as net/http handler





### <a name="Goose4.AddTest">func</a> (\*Goose4) [AddTest](/src/target/goose4.go?s=490:522#L28)
``` go
func (g *Goose4) AddTest(t Test)
```
AddTest updates a Goose4 test list for healthchecks. These tests are used
to determine whether a service is up or not




### <a name="Goose4.ServeHTTP">func</a> (Goose4) [ServeHTTP](/src/target/goose4.go?s=612:677#L33)
``` go
func (g Goose4) ServeHTTP(w http.ResponseWriter, r *http.Request)
```
ServeHTTP is an http router to serve se4 endpoints




## <a name="Healthcheck">type</a> [Healthcheck](/src/target/healthcheck.go?s=1546:1701#L51)
``` go
type Healthcheck struct {
    ReportTime time.Time `json:"report_as_of"`
    Duration   string    `json:"report_duration"`
    Tests      []Test    `json:"tests"`
}
```
Healthcheck provides a full view of healthchecks and whether they fail or not







### <a name="NewHealthcheck">func</a> [NewHealthcheck](/src/target/healthcheck.go?s=1703:1744#L57)
``` go
func NewHealthcheck(t []Test) Healthcheck
```




### <a name="Healthcheck.ASG">func</a> (\*Healthcheck) [ASG](/src/target/healthcheck.go?s=2169:2236#L78)
``` go
func (h *Healthcheck) ASG() (output []byte, errors bool, err error)
```
ASG runs critical tests




### <a name="Healthcheck.All">func</a> (\*Healthcheck) [All](/src/target/healthcheck.go?s=1840:1907#L64)
``` go
func (h *Healthcheck) All() (output []byte, errors bool, err error)
```
All runs all tests; both critical and non-critical




### <a name="Healthcheck.GTG">func</a> (\*Healthcheck) [GTG](/src/target/healthcheck.go?s=2013:2080#L71)
``` go
func (h *Healthcheck) GTG() (output []byte, errors bool, err error)
```
GTG runs non-critical tests: "Good to go"




## <a name="Middleware">type</a> [Middleware](/src/target/middleware.go?s=132:196#L10)
``` go
type Middleware struct {
    SE4 Goose4
    // contains filtered or unexported fields
}
```
Middleware handles the "/service" prefix to comply with the SE4 prefix







### <a name="NewMiddleware">func</a> [NewMiddleware](/src/target/middleware.go?s=286:332#L17)
``` go
func NewMiddleware(h http.Handler) *Middleware
```
NewMiddleware takes an http handler
to wrap and returns mutable Middleware object





### <a name="Middleware.ServeHTTP">func</a> (\*Middleware) [ServeHTTP](/src/target/middleware.go?s=446:516#L24)
``` go
func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request)
```
ServeHTTP wraps our requests and handles any calls to `/service*`.




## <a name="Status">type</a> [Status](/src/target/status.go?s=247:285#L15)
``` go
type Status struct {
    Config
    System
}
```
Status embeds Config and System to give a concise system status










### <a name="Status.Marshal">func</a> (Status) [Marshal](/src/target/status.go?s=378:433#L22)
``` go
func (s Status) Marshal(boot time.Time) ([]byte, error)
```
Marshal returns a status doc based on passed in config and up-to-date
system details




## <a name="System">type</a> [System](/src/target/status.go?s=599:942#L32)
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







### <a name="NewSystem">func</a> [NewSystem](/src/target/status.go?s=944:981#L43)
``` go
func NewSystem(boot time.Time) System
```




## <a name="Test">type</a> [Test](/src/target/healthcheck.go?s=351:1253#L12)
``` go
type Test struct {
    // A simple name to help identify tests from one another
    // there is no enforcement of uniqueness- it is left to the developer
    // to ensure these names make sense
    Name string `json:"test_name"`

    // Critical tests will trigger an `asg` failure- these failures mean, essentially:
    // "This instance is broken and must be recycled"
    // The nil value of this is false and may be omited, or set to false explicity, to
    // set it so a failure just means:
    // "This instance cannot accept traffic but is, functionally, fine"
    // - these failures are useful during boot, fo rexample.
    Critical bool `json:"-"`

    // F is a function which returns true for successful or false for a failure
    F func() bool `json:"-"`

    // The following are overwritten on whatsit
    Result   string    `json:"test_result"`
    Duration string    `json:"duration_millis"`
    TestTime time.Time `json:"tested_at"`
}
```
Test provides a way of having an API pass it's own healthcheck tests,
<a href="https://github.com/beamly/SE4/blob/master/SE4.md#healthcheck">https://github.com/beamly/SE4/blob/master/SE4.md#healthcheck</a>)
into goose4 to be run for the `/healthcheck/` endpoints. These are run in parallel
and so tests which rely on one another/ sequentialness are not allowed














- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
