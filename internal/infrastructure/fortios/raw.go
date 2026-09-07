package fortios

import (
	"context"
	"encoding/json"
	"strings"

	draw "github.com/local/fgcli/internal/domain/raw"
)

type RawRepo struct{ c *Client }

func NewRawRepo(c *Client) *RawRepo { return &RawRepo{c: c} }

func (s *RawRepo) Do(ctx context.Context, method, path, vdom string, body json.RawMessage) (json.RawMessage, error) {
	if err := draw.ValidMethod(method); err != nil {
		return nil, err
	}
	if !strings.HasPrefix(path, "/api/") {
		path = "/api/v2/" + strings.TrimPrefix(path, "/")
	}
	var payload any
	if len(body) > 0 {
		if err := json.Unmarshal(body, &payload); err != nil {
			return nil, err
		}
	}
	b, err := s.c.do(ctx, strings.ToUpper(method), path, vdomQuery(vdom), payload)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(b), nil
}
