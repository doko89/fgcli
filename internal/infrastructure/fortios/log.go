package fortios

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	dlog "github.com/local/fgcli/internal/domain/log"
)

type LogRepo struct{ c *Client }

func NewLogRepo(c *Client) *LogRepo { return &LogRepo{c: c} }

func (s *LogRepo) Fetch(ctx context.Context, source, logtype, vdom string) ([]dlog.Entry, error) {
	if err := dlog.ValidSource(source); err != nil {
		return nil, err
	}
	if strings.TrimSpace(logtype) == "" {
		return nil, fmt.Errorf("log: type is required (e.g. dns, traffic, ips)")
	}
	path := "/api/v2/log/" + strings.ToLower(source) + "/" + logtype + "/raw"
	b, err := s.c.do(ctx, "GET", path, nil, nil)
	if err != nil {
		return nil, err
	}
	var out []dlog.Entry
	if err := decodeResultsAllowEmpty(b, &out); err != nil {
		// Fall back to bare-array bodies.
		if trimmed := bytes.TrimSpace(b); len(trimmed) > 0 && trimmed[0] == '[' {
			if err2 := json.Unmarshal(trimmed, &out); err2 == nil {
				return out, nil
			}
		}
		return nil, err
	}
	return out, nil
}
