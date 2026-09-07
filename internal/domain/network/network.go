package network

import "context"

type Tunnel struct {
	Name      string `json:"name"`
	Comments  string `json:"comments,omitempty"`
	ConnCount int    `json:"connection_count,omitempty"`
	InBytes   int64  `json:"incoming_bytes,omitempty"`
	OutBytes  int64  `json:"outgoing_bytes,omitempty"`
	Rgwy      string `json:"rgwy,omitempty"`
}

type VpnRepository interface {
	List(ctx context.Context, vdom string) ([]Tunnel, error)
}

type RouteStats struct {
	Total int `json:"total_lines"`
	IPv4  int `json:"total_lines_ipv4"`
	IPv6  int `json:"total_lines_ipv6"`
}

// NOTE: scalar status structs above keep zero values on marshal
// (no omitempty): 0 routes / 0 clients is data, not absence.

type RoutingRepository interface {
	Stats(ctx context.Context, vdom string) (RouteStats, error)
}
