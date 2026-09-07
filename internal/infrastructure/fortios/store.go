package fortios

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
	"strings"

	daddr "github.com/local/fgcli/internal/domain/address"
	dpol "github.com/local/fgcli/internal/domain/policy"
	dsys "github.com/local/fgcli/internal/domain/svc"
)

func vdomQuery(vdom string) url.Values {
	q := url.Values{}
	if vdom != "" && vdom != "*" {
		q.Set("vdom", vdom)
	}
	return q
}

type AddressRepo struct{ c *Client }

func NewAddressRepo(c *Client) *AddressRepo { return &AddressRepo{c: c} }

func (s *AddressRepo) List(ctx context.Context, vdom string) ([]daddr.Address, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/cmdb/firewall/address", vdomQuery(vdom), nil)
	if err != nil {
		return nil, err
	}
	var out []daddr.Address
	return out, decodeResults(b, &out)
}

func (s *AddressRepo) Get(ctx context.Context, name, vdom string) (daddr.Address, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/cmdb/firewall/address/"+url.PathEscape(name), vdomQuery(vdom), nil)
	if err != nil {
		return daddr.Address{}, err
	}
	var list []daddr.Address
	if err := decodeResults(b, &list); err != nil {
		return daddr.Address{}, err
	}
	if len(list) == 0 {
		return daddr.Address{}, context.DeadlineExceeded
	}
	return list[0], nil
}

func (s *AddressRepo) Create(ctx context.Context, a daddr.Address, vdom string) (daddr.Address, error) {
	if _, err := s.c.do(ctx, "POST", "/api/v2/cmdb/firewall/address", vdomQuery(vdom), a); err != nil {
		return daddr.Address{}, err
	}
	return a, nil
}

func (s *AddressRepo) Update(ctx context.Context, a daddr.Address, vdom string) (daddr.Address, error) {
	if _, err := s.c.do(ctx, "PUT", "/api/v2/cmdb/firewall/address/"+url.PathEscape(a.Name), vdomQuery(vdom), a); err != nil {
		return daddr.Address{}, err
	}
	return a, nil
}

func (s *AddressRepo) Delete(ctx context.Context, name, vdom string) error {
	_, err := s.c.do(ctx, "DELETE", "/api/v2/cmdb/firewall/address/"+url.PathEscape(name), vdomQuery(vdom), nil)
	return err
}

type PolicyRepo struct{ c *Client }

func NewPolicyRepo(c *Client) *PolicyRepo { return &PolicyRepo{c: c} }

func (s *PolicyRepo) List(ctx context.Context, vdom string) ([]dpol.Policy, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/cmdb/firewall/policy", vdomQuery(vdom), nil)
	if err != nil {
		return nil, err
	}
	var out []dpol.Policy
	return out, decodeResults(b, &out)
}

func (s *PolicyRepo) Get(ctx context.Context, id int64, vdom string) (dpol.Policy, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/cmdb/firewall/policy/"+strconv.FormatInt(id, 10), vdomQuery(vdom), nil)
	if err != nil {
		return dpol.Policy{}, err
	}
	var list []dpol.Policy
	if err := decodeResults(b, &list); err != nil {
		return dpol.Policy{}, err
	}
	if len(list) == 0 {
		return dpol.Policy{}, context.DeadlineExceeded
	}
	return list[0], nil
}

func (s *PolicyRepo) Delete(ctx context.Context, id int64, vdom string) error {
	_, err := s.c.do(ctx, "DELETE", "/api/v2/cmdb/firewall/policy/"+strconv.FormatInt(id, 10), vdomQuery(vdom), nil)
	return err
}

func (s *PolicyRepo) Create(ctx context.Context, p dpol.Policy, vdom string) (dpol.Policy, error) {
	if _, err := s.c.do(ctx, "POST", "/api/v2/cmdb/firewall/policy", vdomQuery(vdom), p); err != nil {
		return dpol.Policy{}, err
	}
	return p, nil
}

func (s *PolicyRepo) Update(ctx context.Context, p dpol.Policy, vdom string) (dpol.Policy, error) {
	if _, err := s.c.do(ctx, "PUT", "/api/v2/cmdb/firewall/policy/"+strconv.FormatInt(p.ID, 10), vdomQuery(vdom), p); err != nil {
		return dpol.Policy{}, err
	}
	return p, nil
}

func (s *PolicyRepo) Move(ctx context.Context, id, before, after int64, vdom string) error {
	q := vdomQuery(vdom)
	if before > 0 {
		q.Set("before", strconv.FormatInt(before, 10))
	} else {
		q.Set("after", strconv.FormatInt(after, 10))
	}
	_, err := s.c.do(ctx, "POST", "/api/v2/cmdb/firewall/policy/"+strconv.FormatInt(id, 10)+"/move", q, nil)
	return err
}

type SysRepo struct{ c *Client }

func NewSysRepo(c *Client) *SysRepo { return &SysRepo{c: c} }

func (s *SysRepo) Status(ctx context.Context) (dsys.Status, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/monitor/system/status", nil, nil)
	if err != nil {
		return dsys.Status{}, err
	}
	// Hostname/model live in results, serial/version/build in the envelope.
	var out dsys.Status
	if err := decodeResults(b, &out); err != nil {
		return dsys.Status{}, err
	}
	var env struct {
		Serial   string `json:"serial"`
		Version  string `json:"version"`
		Build    int    `json:"build"`
		Hostname string `json:"hostname"`
	}
	if err := json.Unmarshal(b, &env); err != nil {
		return dsys.Status{}, err
	}
	if out.Serial == "" {
		out.Serial = env.Serial
	}
	if out.Version == "" {
		out.Version = env.Version
	}
	if out.Build == 0 {
		out.Build = env.Build
	}
	if out.Hostname == "" {
		out.Hostname = env.Hostname
	}
	return out, nil
}

func (s *SysRepo) License(ctx context.Context) (dsys.License, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/monitor/license/status", nil, nil)
	if err != nil {
		return dsys.License{}, err
	}
	var out dsys.License
	return out, decodeResults(b, &out)
}

func (s *SysRepo) Fortiguard(ctx context.Context) (dsys.Fortiguard, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/monitor/fortiguard/service-communication-stats", nil, nil)
	if err != nil {
		return nil, err
	}
	var out dsys.Fortiguard
	return out, decodeResults(b, &out)
}

func (s *SysRepo) Ntp(ctx context.Context) ([]dsys.NtpServer, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/monitor/system/ntp/status", nil, nil)
	if err != nil {
		return nil, err
	}
	var out []dsys.NtpServer
	return out, decodeResults(b, &out)
}

func (s *SysRepo) Dns(ctx context.Context) ([]dsys.DnsEntry, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/monitor/system/acquired-dns", nil, nil)
	if err != nil {
		return nil, err
	}
	var out []dsys.DnsEntry
	return out, decodeResults(b, &out)
}

func (s *SysRepo) Dhcp(ctx context.Context) ([]dsys.DhcpEntry, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/monitor/system/interface/dhcp-status", nil, nil)
	if err != nil {
		// FortiOS answers HTTP 424 when no interface uses DHCP: that is
		// an empty report, not a failure.
		if strings.Contains(err.Error(), "http 424") {
			return []dsys.DhcpEntry{}, nil
		}
		return nil, err
	}
	var out []dsys.DhcpEntry
	if err := decodeResultsAllowEmpty(b, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *SysRepo) Snmp(ctx context.Context) (dsys.Snmp, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/cmdb/system.snmp/sysinfo", nil, nil)
	if err != nil {
		return dsys.Snmp{}, err
	}
	var out dsys.Snmp
	return out, decodeResults(b, &out)
}
