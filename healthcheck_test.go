package goose4

import (
	"reflect"
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

func TestNewHealthcheck(t *testing.T) {
	var testSuccess = Test{F: HealthTestSuccess}
	var testFailure = Test{F: HealthTestFailure}

	for _, test := range []struct {
		title string
		tests []Test
	}{
		{"Multiple tests", []Test{testSuccess, testFailure}},
		{"Single test", []Test{testSuccess}},
		{"No tests", []Test{}},
	} {
		t.Run(test.title, func(t *testing.T) {
			h := NewHealthcheck(test.tests)

			t.Run("Tests are set correctly", func(t *testing.T) {
				if !reflect.DeepEqual(test.tests, h.Tests) {
					t.Errorf("expected %v, received %v", test.tests, h.Tests)
				}
			})
		})
	}
}

func TestHealthcheckAll(t *testing.T) {
	for _, test := range []struct {
		title          string
		f              func() bool
		silent         bool
		expectedErrors bool
		expectError    bool
	}{
		{"A simple successful healthcheck", HealthTestSuccess, false, false, false},
		{"A simple failing healthcheck", HealthTestFailure, false, true, false},
		{"A silent successful healthcheck", HealthTestSuccess, true, false, false},
		{"A silent failing healthcheck", HealthTestFailure, true, false, false},
	} {
		t.Run(test.title, func(t *testing.T) {
			t0 := Test{F: test.f, Critical: true, Silent: test.silent}
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

func TestHealthcheckASG(t *testing.T) {
	for _, test := range []struct {
		title          string
		f              func() bool
		silent         bool
		expectedErrors bool
		expectError    bool
	}{
		{"A simple successful healthcheck", HealthTestSuccess, false, false, false},
		{"A simple failing healthcheck", HealthTestFailure, false, true, false},
		{"A silent successful healthcheck", HealthTestSuccess, true, false, false},
		{"A silent failing healthcheck", HealthTestFailure, true, false, false},
	} {
		t.Run(test.title, func(t *testing.T) {
			t0 := Test{F: test.f, Critical: true, Silent: test.silent}
			h := Healthcheck{Tests: []Test{t0}}

			_, errs, err := h.ASG()

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

func TestHealthcheckGTG(t *testing.T) {
	for _, test := range []struct {
		title          string
		f              func() bool
		silent         bool
		expectedErrors bool
		expectError    bool
	}{
		{"A simple successful healthcheck", HealthTestSuccess, false, false, false},
		{"A simple failing healthcheck", HealthTestFailure, false, true, false},
		{"A silent successful healthcheck", HealthTestSuccess, true, false, false},
		{"A silent failing healthcheck", HealthTestFailure, true, false, false},
	} {
		t.Run(test.title, func(t *testing.T) {
			t0 := Test{F: test.f, Critical: false, Silent: test.silent}
			h := Healthcheck{Tests: []Test{t0}}

			_, errs, err := h.GTG()

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
