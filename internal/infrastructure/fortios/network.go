package fortios

import (
	"context"
	"strconv"

	dnet "github.com/local/fgcli/internal/domain/network"
)

type VpnRepo struct{ c *Client }

func NewVpnRepo(c *Client) *VpnRepo { return &VpnRepo{c: c} }

func (s *VpnRepo) List(ctx context.Context, vdom string) ([]dnet.Tunnel, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/monitor/vpn/ipsec", nil, nil)
	if err != nil {
		return nil, err
	}
	var out []dnet.Tunnel
	if err := decodeResultsAllowEmpty(b, &out); err != nil {
		return nil, err
	}
	return out, nil
}

type RoutingRepo struct{ c *Client }

func NewRoutingRepo(c *Client) *RoutingRepo { return &RoutingRepo{c: c} }

func (s *RoutingRepo) Stats(ctx context.Context, vdom string) (dnet.RouteStats, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/monitor/router/statistics", nil, nil)
	if err != nil {
		return dnet.RouteStats{}, err
	}
	var out dnet.RouteStats
	return out, decodeResults(b, &out)
}

type DnsFilterRepo struct{ c *Client }

func NewDnsFilterRepo(c *Client) *DnsFilterRepo { return &DnsFilterRepo{c: c} }

func (s *DnsFilterRepo) List(ctx context.Context, vdom string) ([]dnet.DnsFilter, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/cmdb/dnsfilter/domain-filter", vdomQuery(vdom), nil)
	if err != nil {
		return nil, err
	}
	var out []dnet.DnsFilter
	return out, decodeResults(b, &out)
}

func (s *DnsFilterRepo) Get(ctx context.Context, id int64, vdom string) (dnet.DnsFilter, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/cmdb/dnsfilter/domain-filter/"+strconv.FormatInt(id, 10), vdomQuery(vdom), nil)
	if err != nil {
		return dnet.DnsFilter{}, err
	}
	var list []dnet.DnsFilter
	if err := decodeResults(b, &list); err != nil {
		return dnet.DnsFilter{}, err
	}
	if len(list) == 0 {
		return dnet.DnsFilter{}, context.DeadlineExceeded
	}
	return list[0], nil
}
