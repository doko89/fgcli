package fortios

import (
	"context"
	"net/url"

	dsvc "github.com/local/fgcli/internal/domain/service"
)

type SvcRepo struct{ c *Client }

func NewSvcRepo(c *Client) *SvcRepo { return &SvcRepo{c: c} }

func (s *SvcRepo) List(ctx context.Context, vdom string) ([]dsvc.Svc, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/cmdb/firewall.service/custom", vdomQuery(vdom), nil)
	if err != nil {
		return nil, err
	}
	var out []dsvc.Svc
	return out, decodeResults(b, &out)
}

func (s *SvcRepo) Get(ctx context.Context, name, vdom string) (dsvc.Svc, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/cmdb/firewall.service/custom/"+url.PathEscape(name), vdomQuery(vdom), nil)
	if err != nil {
		return dsvc.Svc{}, err
	}
	var list []dsvc.Svc
	if err := decodeResults(b, &list); err != nil {
		return dsvc.Svc{}, err
	}
	if len(list) == 0 {
		return dsvc.Svc{}, context.DeadlineExceeded
	}
	return list[0], nil
}

func (s *SvcRepo) Create(ctx context.Context, v dsvc.Svc, vdom string) (dsvc.Svc, error) {
	if _, err := s.c.do(ctx, "POST", "/api/v2/cmdb/firewall.service/custom", vdomQuery(vdom), v); err != nil {
		return dsvc.Svc{}, err
	}
	return v, nil
}

func (s *SvcRepo) Update(ctx context.Context, v dsvc.Svc, vdom string) (dsvc.Svc, error) {
	if _, err := s.c.do(ctx, "PUT", "/api/v2/cmdb/firewall.service/custom/"+url.PathEscape(v.Name), vdomQuery(vdom), v); err != nil {
		return dsvc.Svc{}, err
	}
	return v, nil
}

func (s *SvcRepo) Delete(ctx context.Context, name, vdom string) error {
	_, err := s.c.do(ctx, "DELETE", "/api/v2/cmdb/firewall.service/custom/"+url.PathEscape(name), vdomQuery(vdom), nil)
	return err
}

type SvcGroupRepo struct{ c *Client }

func NewSvcGroupRepo(c *Client) *SvcGroupRepo { return &SvcGroupRepo{c: c} }

func (s *SvcGroupRepo) List(ctx context.Context, vdom string) ([]dsvc.Group, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/cmdb/firewall.service/group", vdomQuery(vdom), nil)
	if err != nil {
		return nil, err
	}
	var out []dsvc.Group
	return out, decodeResults(b, &out)
}

func (s *SvcGroupRepo) Get(ctx context.Context, name, vdom string) (dsvc.Group, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/cmdb/firewall.service/group/"+url.PathEscape(name), vdomQuery(vdom), nil)
	if err != nil {
		return dsvc.Group{}, err
	}
	var list []dsvc.Group
	if err := decodeResults(b, &list); err != nil {
		return dsvc.Group{}, err
	}
	if len(list) == 0 {
		return dsvc.Group{}, context.DeadlineExceeded
	}
	return list[0], nil
}

func (s *SvcGroupRepo) Create(ctx context.Context, g dsvc.Group, vdom string) (dsvc.Group, error) {
	if _, err := s.c.do(ctx, "POST", "/api/v2/cmdb/firewall.service/group", vdomQuery(vdom), g); err != nil {
		return dsvc.Group{}, err
	}
	return g, nil
}

func (s *SvcGroupRepo) Update(ctx context.Context, g dsvc.Group, vdom string) (dsvc.Group, error) {
	if _, err := s.c.do(ctx, "PUT", "/api/v2/cmdb/firewall.service/group/"+url.PathEscape(g.Name), vdomQuery(vdom), g); err != nil {
		return dsvc.Group{}, err
	}
	return g, nil
}

func (s *SvcGroupRepo) Delete(ctx context.Context, name, vdom string) error {
	_, err := s.c.do(ctx, "DELETE", "/api/v2/cmdb/firewall.service/group/"+url.PathEscape(name), vdomQuery(vdom), nil)
	return err
}

type VoipRepo struct{ c *Client }

func NewVoipRepo(c *Client) *VoipRepo { return &VoipRepo{c: c} }

func (s *VoipRepo) Get(ctx context.Context, name, vdom string) (dsvc.VoipProfile, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/cmdb/voip/profile/"+url.PathEscape(name), vdomQuery(vdom), nil)
	if err != nil {
		return dsvc.VoipProfile{}, err
	}
	var list []dsvc.VoipProfile
	if err := decodeResults(b, &list); err != nil {
		return dsvc.VoipProfile{}, err
	}
	if len(list) == 0 {
		return dsvc.VoipProfile{}, context.DeadlineExceeded
	}
	return list[0], nil
}
