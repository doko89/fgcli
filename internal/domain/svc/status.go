package svc

import "context"

type Status struct {
	Hostname string `json:"hostname"`
	Version  string `json:"version"`
	Serial   string `json:"serial"`
	Build    int    `json:"build,omitempty"`
	Model    string `json:"model,omitempty"`
}

type Repository interface {
	Status(ctx context.Context) (Status, error)
	License(ctx context.Context) (License, error)
	Fortiguard(ctx context.Context) (Fortiguard, error)
	Ntp(ctx context.Context) ([]NtpServer, error)
	Dns(ctx context.Context) ([]DnsEntry, error)
	Dhcp(ctx context.Context) ([]DhcpEntry, error)
	Snmp(ctx context.Context) (Snmp, error)
}
