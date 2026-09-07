package fortios

import (
	"context"

	dmon "github.com/local/fgcli/internal/domain/monitor"
)

type MonitorRepo struct{ c *Client }

func NewMonitorRepo(c *Client) *MonitorRepo { return &MonitorRepo{c: c} }

func (s *MonitorRepo) Switches(ctx context.Context, vdom string) ([]dmon.ManagedSwitch, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/monitor/switch-controller/managed-switch/status", nil, nil)
	if err != nil {
		return nil, err
	}
	var out []dmon.ManagedSwitch
	if err := decodeResultsAllowEmpty(b, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *MonitorRepo) Wifi(ctx context.Context, vdom string) (dmon.WifiStatus, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/monitor/wifi/ap_status", vdomQuery(vdom), nil)
	if err != nil {
		return dmon.WifiStatus{}, err
	}
	var out dmon.WifiStatus
	return out, decodeResults(b, &out)
}

func (s *MonitorRepo) FortiView(ctx context.Context, vdom string) (dmon.FortiView, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/monitor/fortiview/statistics", nil, nil)
	if err != nil {
		return dmon.FortiView{}, err
	}
	var out dmon.FortiView
	return out, decodeResults(b, &out)
}
