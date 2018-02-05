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
		requiredForASG bool
		requiredForGTG bool
		expectedErrors bool
		expectError    bool
	}{
		{"A successful healthcheck, requires ASG", HealthTestSuccess, true, false, false, false},
		{"A successful healthcheck, requires GTG", HealthTestSuccess, false, true, false, false},
		{"A failing healthcheck, requires ASG", HealthTestFailure, true, false, true, false},
		{"A failing healthcheck, requires GTG", HealthTestFailure, false, true, true, false},
		{"A successful healthcheck, requires ASG and GTG", HealthTestSuccess, true, true, false, false},
		{"A successful healthcheck", HealthTestSuccess, false, false, false, false},
		{"A failing healthcheck, requires ASG and GTG", HealthTestFailure, true, true, false, false},
		{"A failing healthcheck", HealthTestFailure, true, true, false, false},
	} {
		t.Run(test.title, func(t *testing.T) {
			t0 := Test{F: test.f, RequiredForASG: test.requiredForASG, RequiredForGTG: test.requiredForGTG}
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
		requiredForASG bool
		expectedErrors bool
		expectError    bool
	}{
		{"A critical successful healthcheck", HealthTestSuccess, true, false, false},
		{"A critical failing healthcheck", HealthTestFailure, true, true, false},
		{"A non-critical successful healthcheck", HealthTestSuccess, false, false, false},
		{"A non-critical failing healthcheck", HealthTestFailure, false, false, false},
	} {
		t.Run(test.title, func(t *testing.T) {
			t0 := Test{F: test.f, RequiredForASG: test.requiredForASG, RequiredForGTG: true}
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
		requiredForGTG bool
		expectedErrors bool
		expectError    bool
	}{
		{"A critical successful healthcheck", HealthTestSuccess, true, false, false},
		{"A critical failing healthcheck", HealthTestFailure, true, true, false},
		{"A non-critical successful healthcheck", HealthTestSuccess, false, false, false},
		{"A non-critical failing healthcheck", HealthTestFailure, false, false, false},
	} {
		t.Run(test.title, func(t *testing.T) {
			t0 := Test{F: test.f, RequiredForASG: true, RequiredForGTG: test.requiredForGTG}
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
