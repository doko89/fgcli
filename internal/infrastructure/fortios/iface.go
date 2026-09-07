package fortios

import (
	"context"
	"net/url"

	dif "github.com/local/fgcli/internal/domain/iface"
)

type IfaceRepo struct{ c *Client }

func NewIfaceRepo(c *Client) *IfaceRepo { return &IfaceRepo{c: c} }

func (s *IfaceRepo) List(ctx context.Context, vdom string) ([]dif.Iface, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/cmdb/system/interface", vdomQuery(vdom), nil)
	if err != nil {
		return nil, err
	}
	var out []dif.Iface
	return out, decodeResults(b, &out)
}

func (s *IfaceRepo) Get(ctx context.Context, name, vdom string) (dif.Iface, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/cmdb/system/interface/"+url.PathEscape(name), vdomQuery(vdom), nil)
	if err != nil {
		return dif.Iface{}, err
	}
	var list []dif.Iface
	if err := decodeResults(b, &list); err != nil {
		return dif.Iface{}, err
	}
	if len(list) == 0 {
		return dif.Iface{}, context.DeadlineExceeded
	}
	return list[0], nil
}

type ZoneRepo struct{ c *Client }

func NewZoneRepo(c *Client) *ZoneRepo { return &ZoneRepo{c: c} }

func (s *ZoneRepo) List(ctx context.Context, vdom string) ([]dif.Zone, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/cmdb/system/zone", vdomQuery(vdom), nil)
	if err != nil {
		return nil, err
	}
	var out []dif.Zone
	return out, decodeResults(b, &out)
}

func (s *ZoneRepo) Get(ctx context.Context, name, vdom string) (dif.Zone, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/cmdb/system/zone/"+url.PathEscape(name), vdomQuery(vdom), nil)
	if err != nil {
		return dif.Zone{}, err
	}
	var list []dif.Zone
	if err := decodeResults(b, &list); err != nil {
		return dif.Zone{}, err
	}
	if len(list) == 0 {
		return dif.Zone{}, context.DeadlineExceeded
	}
	return list[0], nil
}
