package goose4

import (
	"encoding/json"
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
	Critical bool `json:"-"`

	// F is a function which returns true for successful or false for a failure
	F func() bool `json:"-"`

	// The following are overwritten on whatsit
	Result   string    `json:"test_result"`
	Duration string    `json:"duration_millis"`
	TestTime time.Time `json:"tested_at"`
}

func (t *Test) run() bool {
	t.TestTime = time.Now()
	success := t.F()

	if success {
		t.Result = "passed"
	} else {
		t.Result = "failed"
	}

	t.Duration = time.Since(t.TestTime).String()

	return success
}

// Healthcheck provides a full view of healthchecks and whether they fail or not
type Healthcheck struct {
	ReportTime time.Time `json:"report_as_of"`
	Duration   string    `json:"report_duration"`
	Tests      []Test    `json:"tests"`
}

func NewHealthcheck(t []Test) Healthcheck {
	return Healthcheck{
		Tests: t,
	}
}

// All runs all tests; both critical and non-critical
func (h *Healthcheck) All() (output []byte, errors bool, err error) {
	output, errors, err = h.runTests(true, true)

	return
}

// GTG runs non-critical tests: "Good to go"
func (h *Healthcheck) GTG() (output []byte, errors bool, err error) {
	output, errors, err = h.runTests(false, true)

	return
}

// ASG runs critical tests
func (h *Healthcheck) ASG() (output []byte, errors bool, err error) {
	output, errors, err = h.runTests(true, false)

	return
}

func (h *Healthcheck) runTests(critical, noncritical bool) ([]byte, bool, error) {
	h.ReportTime = time.Now()

	var errs bool
	bchan := make(chan Test)

	testList := []Test{}
	if critical {
		testList = testByStatus(h.Tests, testList, true)
	}

	if noncritical {
		testList = testByStatus(h.Tests, testList, false)
	}

	if len(testList) > 0 {
		for _, t := range testList {
			go func(t0 Test) {
				if !t0.run() {
					errs = true
				}

				bchan <- t0
			}(t)
		}

		count := 1
		completedTests := []Test{}
		for t := range bchan {
			completedTests = append(completedTests, t)

			if count == len(testList) {
				break
			}
			count++
		}

		h.Tests = completedTests
	}

	h.Duration = time.Since(h.ReportTime).String()
	j, err := json.Marshal(h)

	return j, errs, err
}

func testByStatus(t1 []Test, t2 []Test, critical bool) []Test {
	for _, t := range t1 {
		if t.Critical == critical {
			t2 = append(t2, t)
		}
	}

	return t2
}
