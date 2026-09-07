package log

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
)

var sources = map[string]bool{
	"fortianalyzer": true, "forticloud": true, "memory": true, "disk": true,
}

// Entry schema varies by log type (15+); passthrough preserves all fields.
type Entry = json.RawMessage

func ValidSource(s string) error {
	if !sources[strings.ToLower(s)] {
		return errors.New("log: source must be fortianalyzer|forticloud|memory|disk")
	}
	return nil
}

type Repository interface {
	Fetch(ctx context.Context, source, logtype, vdom string) ([]Entry, error)
}
