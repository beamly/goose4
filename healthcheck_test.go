package goose4

import (
	"testing"
)

var (
	HealthTestSuccess = func() bool { return true }
	HealthTestFailure = func() bool { return false }
)

func TestTest_Run(t *testing.T) {
	for _, test := range []struct {
		title          string
		f              func() bool
		expectedResult string
	}{
		{"A successful healthcheck", HealthTestSuccess, "passed"},
		{"An unsuccessful healthcheck", HealthTestFailure, "failed"},
	} {
		t.Run(test.title, func(t *testing.T) {
			t0 := Test{F: test.f}

			_ = t0.run()
			t.Run("Result", func(t *testing.T) {
				if test.expectedResult != t0.Result {
					t.Errorf("expected %q, received %q", test.expectedResult, t0.Result)
				}
			})

		})
	}
}

func TestHealthcheckAll(t *testing.T) {
	for _, test := range []struct {
		title          string
		f              func() bool
		expectedErrors bool
		expectError    bool
	}{
		{"A simple successful healthcheck", HealthTestSuccess, false, false},
		{"A simple failing healthcheck", HealthTestFailure, true, false},
	} {
		t.Run(test.title, func(t *testing.T) {
			t0 := Test{F: test.f, Critical: true}
			h := Healthcheck{Tests: []Test{t0}}

			_, errs, err := h.All()

			t.Run("Test errors", func(t *testing.T) {
				if test.expectedErrors != errs {
					t.Errorf("expected %v, received %v", test.expectedErrors, errs)
				}
			})

			t.Run("Test returns error", func(t *testing.T) {
				if test.expectError == (err == nil) {
					t.Errorf("expected %v, received %v", test.expectedErrors, errs)
				}
			})
		})

	}
}
