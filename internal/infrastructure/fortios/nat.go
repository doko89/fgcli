package fortios

import (
	"context"
	"net/url"
	"strconv"

	dip "github.com/local/fgcli/internal/domain/ippool"
	dsn "github.com/local/fgcli/internal/domain/snat"
	dvip "github.com/local/fgcli/internal/domain/vip"
)

const (
	apiFirewallVip         = "/api/v2/cmdb/firewall/vip"
	apiFirewallIppool      = "/api/v2/cmdb/firewall/ippool"
	apiFirewallCentralSnat = "/api/v2/cmdb/firewall/central-snat-map"
)

type VipRepo struct{ c *Client }

func NewVipRepo(c *Client) *VipRepo { return &VipRepo{c: c} }

func (s *VipRepo) List(ctx context.Context, vdom string) ([]dvip.Vip, error) {
	b, err := s.c.do(ctx, "GET", apiFirewallVip, vdomQuery(vdom), nil)
	if err != nil {
		return nil, err
	}
	var out []dvip.Vip
	return out, decodeResults(b, &out)
}

func (s *VipRepo) Get(ctx context.Context, name, vdom string) (dvip.Vip, error) {
	b, err := s.c.do(ctx, "GET", apiFirewallVip+"/"+url.PathEscape(name), vdomQuery(vdom), nil)
	if err != nil {
		return dvip.Vip{}, err
	}
	var list []dvip.Vip
	if err := decodeResults(b, &list); err != nil {
		return dvip.Vip{}, err
	}
	if len(list) == 0 {
		return dvip.Vip{}, context.DeadlineExceeded
	}
	return list[0], nil
}

func (s *VipRepo) Create(ctx context.Context, v dvip.Vip, vdom string) (dvip.Vip, error) {
	if _, err := s.c.do(ctx, "POST", apiFirewallVip, vdomQuery(vdom), v); err != nil {
		return dvip.Vip{}, err
	}
	return v, nil
}

func (s *VipRepo) Update(ctx context.Context, v dvip.Vip, vdom string) (dvip.Vip, error) {
	if _, err := s.c.do(ctx, "PUT", apiFirewallVip+"/"+url.PathEscape(v.Name), vdomQuery(vdom), v); err != nil {
		return dvip.Vip{}, err
	}
	return v, nil
}

func (s *VipRepo) Delete(ctx context.Context, name, vdom string) error {
	_, err := s.c.do(ctx, "DELETE", apiFirewallVip+"/"+url.PathEscape(name), vdomQuery(vdom), nil)
	return err
}

type PoolRepo struct{ c *Client }

func NewPoolRepo(c *Client) *PoolRepo { return &PoolRepo{c: c} }

func (s *PoolRepo) List(ctx context.Context, vdom string) ([]dip.Pool, error) {
	b, err := s.c.do(ctx, "GET", apiFirewallIppool, vdomQuery(vdom), nil)
	if err != nil {
		return nil, err
	}
	var out []dip.Pool
	return out, decodeResults(b, &out)
}

func (s *PoolRepo) Get(ctx context.Context, name, vdom string) (dip.Pool, error) {
	b, err := s.c.do(ctx, "GET", apiFirewallIppool+"/"+url.PathEscape(name), vdomQuery(vdom), nil)
	if err != nil {
		return dip.Pool{}, err
	}
	var list []dip.Pool
	if err := decodeResults(b, &list); err != nil {
		return dip.Pool{}, err
	}
	if len(list) == 0 {
		return dip.Pool{}, context.DeadlineExceeded
	}
	return list[0], nil
}

func (s *PoolRepo) Create(ctx context.Context, p dip.Pool, vdom string) (dip.Pool, error) {
	if _, err := s.c.do(ctx, "POST", apiFirewallIppool, vdomQuery(vdom), p); err != nil {
		return dip.Pool{}, err
	}
	return p, nil
}

func (s *PoolRepo) Update(ctx context.Context, p dip.Pool, vdom string) (dip.Pool, error) {
	if _, err := s.c.do(ctx, "PUT", apiFirewallIppool+"/"+url.PathEscape(p.Name), vdomQuery(vdom), p); err != nil {
		return dip.Pool{}, err
	}
	return p, nil
}

func (s *PoolRepo) Delete(ctx context.Context, name, vdom string) error {
	_, err := s.c.do(ctx, "DELETE", apiFirewallIppool+"/"+url.PathEscape(name), vdomQuery(vdom), nil)
	return err
}

type SnatRepo struct{ c *Client }

func NewSnatRepo(c *Client) *SnatRepo { return &SnatRepo{c: c} }

func (s *SnatRepo) List(ctx context.Context, vdom string) ([]dsn.Map, error) {
	b, err := s.c.do(ctx, "GET", apiFirewallCentralSnat, vdomQuery(vdom), nil)
	if err != nil {
		return nil, err
	}
	var out []dsn.Map
	return out, decodeResults(b, &out)
}

func (s *SnatRepo) Get(ctx context.Context, id int64, vdom string) (dsn.Map, error) {
	b, err := s.c.do(ctx, "GET", apiFirewallCentralSnat+"/"+strconv.FormatInt(id, 10), vdomQuery(vdom), nil)
	if err != nil {
		return dsn.Map{}, err
	}
	var list []dsn.Map
	if err := decodeResults(b, &list); err != nil {
		return dsn.Map{}, err
	}
	if len(list) == 0 {
		return dsn.Map{}, context.DeadlineExceeded
	}
	return list[0], nil
}

func (s *SnatRepo) Create(ctx context.Context, m dsn.Map, vdom string) (dsn.Map, error) {
	if _, err := s.c.do(ctx, "POST", apiFirewallCentralSnat, vdomQuery(vdom), m); err != nil {
		return dsn.Map{}, err
	}
	return m, nil
}

func (s *SnatRepo) Update(ctx context.Context, m dsn.Map, vdom string) (dsn.Map, error) {
	if _, err := s.c.do(ctx, "PUT", apiFirewallCentralSnat+"/"+strconv.FormatInt(m.ID, 10), vdomQuery(vdom), m); err != nil {
		return dsn.Map{}, err
	}
	return m, nil
}

func (s *SnatRepo) Delete(ctx context.Context, id int64, vdom string) error {
	_, err := s.c.do(ctx, "DELETE", apiFirewallCentralSnat+"/"+strconv.FormatInt(id, 10), vdomQuery(vdom), nil)
	return err
}
