package goose4

import (
	"encoding/json"
	"time"
)

// Config implements a subset of https://github.com/beamly/SE4/blob/master/SE4.md#status
// and is used to configure static values for goose4.
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

// Marshal returns a json document and, potentially, an error in order to
// respond with configuration for a service.
func (c Config) Marshal() (j []byte, err error) {
	return json.Marshal(c)
}
