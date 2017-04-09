package goose4

import (
	"testing"
	"time"
)

var (
	TestConfig = Config{"artifact",
		"123",
		"localhost",
		"root",
		time.Time{},
		"go version go1.7.4 darwin/amd64",
		"32b619ba997dfbfafd528ae3fea4e2cba8116be8",
		"https://runbooks.example.com/goose4.md",
		"1.0.0",
	}
	TestJSON = `{"artifact_id":"artifact","build_number":"123","build_machine":"localhost","built_by":"root","built_when":"0001-01-01T00:00:00Z","compiler_version":"go version go1.7.4 darwin/amd64","git_sha1":"32b619ba997dfbfafd528ae3fea4e2cba8116be8","runbook_uri":"https://runbooks.example.com/goose4.md","version":"1.0.0"}`
)

func TestConfigMarshal(t *testing.T) {
	for _, test := range []struct {
		title       string
		config      Config
		output      string
		expectError bool
	}{
		{"simple, valid config", TestConfig, TestJSON, false},
	} {
		t.Run(test.title, func(t *testing.T) {
			output, err := test.config.Marshal()
			if err == nil {
				if test.expectError {
					t.Errorf("Config.Marshal(): expected error, received none")
				}

				if test.output != string(output) {
					t.Errorf("expected %q, received %q", test.output, string(output))
				}
			} else {
				if !test.expectError {
					t.Errorf("Config.Marshal(): received error %v", err)
				}
			}
		})
	}
}
