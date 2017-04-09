package goose4

import (
	"testing"
)

func TestNewGoose4(t *testing.T) {
	for _, test := range []struct {
		title       string
		config      Config
		expectError bool
	}{
		{"simple goose4 object", Config{}, false},
	} {
		t.Run(test.title, func(t *testing.T) {
			_, err := NewGoose4(test.config)

			if err == nil {
				if test.expectError {
					t.Errorf("NewGoose4: expected error, received none")
				}
			} else {
				if !test.expectError {
					t.Errorf("NewGoose4: %v", err)
				}
			}
		})
	}
}
