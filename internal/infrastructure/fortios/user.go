package fortios

import (
	"context"

	duser "github.com/local/fgcli/internal/domain/user"
)

type UserRepo struct{ c *Client }

func NewUserRepo(c *Client) *UserRepo { return &UserRepo{c: c} }

func (s *UserRepo) Firewall(ctx context.Context, vdom string) ([]duser.FwUser, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/monitor/user/firewall", nil, nil)
	if err != nil {
		return nil, err
	}
	var out []duser.FwUser
	if err := decodeResultsAllowEmpty(b, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *UserRepo) Banned(ctx context.Context, vdom string) ([]duser.BannedUser, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/monitor/user/banned", nil, nil)
	if err != nil {
		return nil, err
	}
	var out []duser.BannedUser
	if err := decodeResultsAllowEmpty(b, &out); err != nil {
		return nil, err
	}
	return out, nil
}
