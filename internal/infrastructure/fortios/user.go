package fortios

import (
	"context"
	"net/url"

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

func (s *UserRepo) GroupList(ctx context.Context, vdom string) ([]duser.UserGroup, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/cmdb/user/group", vdomQuery(vdom), nil)
	if err != nil {
		return nil, err
	}
	var out []duser.UserGroup
	return out, decodeResults(b, &out)
}

func (s *UserRepo) GroupGet(ctx context.Context, name, vdom string) (duser.UserGroup, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/cmdb/user/group/"+url.PathEscape(name), vdomQuery(vdom), nil)
	if err != nil {
		return duser.UserGroup{}, err
	}
	var list []duser.UserGroup
	if err := decodeResults(b, &list); err != nil {
		return duser.UserGroup{}, err
	}
	if len(list) == 0 {
		return duser.UserGroup{}, context.DeadlineExceeded
	}
	return list[0], nil
}

// GroupUpdate replaces the member list (PUT with the full object).
func (s *UserRepo) GroupUpdate(ctx context.Context, g duser.UserGroup, vdom string) (duser.UserGroup, error) {
	if _, err := s.c.do(ctx, "PUT", "/api/v2/cmdb/user/group/"+url.PathEscape(g.Name), vdomQuery(vdom), g); err != nil {
		return duser.UserGroup{}, err
	}
	return g, nil
}

func (s *UserRepo) LocalList(ctx context.Context, vdom string) ([]duser.LocalUser, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/cmdb/user/local", vdomQuery(vdom), nil)
	if err != nil {
		return nil, err
	}
	var out []duser.LocalUser
	return out, decodeResults(b, &out)
}

func (s *UserRepo) LocalGet(ctx context.Context, name, vdom string) (duser.LocalUser, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/cmdb/user/local/"+url.PathEscape(name), vdomQuery(vdom), nil)
	if err != nil {
		return duser.LocalUser{}, err
	}
	var list []duser.LocalUser
	if err := decodeResults(b, &list); err != nil {
		return duser.LocalUser{}, err
	}
	if len(list) == 0 {
		return duser.LocalUser{}, context.DeadlineExceeded
	}
	return list[0], nil
}

func (s *UserRepo) LocalCreate(ctx context.Context, u duser.LocalUser, vdom string) (duser.LocalUser, error) {
	if _, err := s.c.do(ctx, "POST", "/api/v2/cmdb/user/local", vdomQuery(vdom), u); err != nil {
		return duser.LocalUser{}, err
	}
	u.Passwd = ""
	return u, nil
}

func (s *UserRepo) LocalDelete(ctx context.Context, name, vdom string) error {
	_, err := s.c.do(ctx, "DELETE", "/api/v2/cmdb/user/local/"+url.PathEscape(name), vdomQuery(vdom), nil)
	return err
}
