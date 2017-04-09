package goose4

import (
	"testing"
	"time"
)

var (
	TestStatus = Status{Config: TestConfig}
)

func TestStatusMarshal(t *testing.T) {
	for _, test := range []struct {
		title       string
		status      Status
		expectError bool
	}{
		{"simple status", TestStatus, false},
	} {
		t.Run(test.title, func(t *testing.T) {
			_, err := test.status.Marshal(time.Now())
			if err == nil {
				if test.expectError {
					t.Errorf("Status.Marshal(): expected error, received none")
				}
			} else {
				if !test.expectError {
					t.Errorf("Status.Marshal(): unexpected error %v", err)
				}
			}
		})
	}
}
