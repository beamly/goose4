package goose4

import (
	"encoding/json"
)

// Error is a simple placeholder to store non-2xx, non-3xx response data
type Error struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// Marshal wraps an Error in some json
func (e Error) Marshal() ([]byte, error) {
	return json.Marshal(e)
}
