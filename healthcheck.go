package goose4

import (
	"time"
)

// Test provides a way of having an API pass it's own healthcheck tests,
// https://github.com/beamly/SE4/blob/master/SE4.md#healthcheck)
// into goose4 to be run for the `/healthcheck/` endpoints. These are run in parallel
// and so tests which rely on one another/ sequentialness are not allowed
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
	Critical bool `json:-`

	// F is a function which returns true for successful or false for a failure
	F func() bool `json:-`

	result   string        `json:"test_result"`
	duration time.Duration `json:"duration_millis"`
	testTime time.Time     `json:"tested_at"`
}

// Healthcheck provides a full view of healthchecks and whether they fail or not
type Healthcheck struct {
	reportTime time.Time     `json:"report_as_of"`
	duration   time.Duration `json:"report_duration"`
	tests      []Test        `json:"tests"`
}
