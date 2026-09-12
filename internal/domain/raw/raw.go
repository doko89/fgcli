package raw

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
)

// Doer is the authenticated escape hatch: any API path while more
// specific verbs are still missing. Output stays in the JSON envelope.
type Doer interface {
	Do(ctx context.Context, method, path, vdom string, body json.RawMessage) (json.RawMessage, error)
}

func ValidMethod(m string) error {
	switch strings.ToUpper(m) {
	case "GET", "POST", "PUT", "DELETE":
		return nil
	}
	return errors.New("raw: method must be get|post|put|delete")
}
