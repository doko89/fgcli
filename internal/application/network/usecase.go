package network

import (
	"context"

	domain "github.com/local/fgcli/internal/domain/network"
)

type VpnUseCase struct{ repo domain.VpnRepository }

func NewVpn(repo domain.VpnRepository) *VpnUseCase { return &VpnUseCase{repo: repo} }

func (u *VpnUseCase) List(ctx context.Context, vdom string) ([]domain.Tunnel, error) {
	return u.repo.List(ctx, vdom)
}

// Ssl returns the SSL-VPN settings object (cmdb/vpn.ssl/settings).
func (u *VpnUseCase) Ssl(ctx context.Context, vdom string) (any, error) {
	return u.repo.Ssl(ctx, vdom)
}

type RoutingUseCase struct{ repo domain.RoutingRepository }

func NewRouting(repo domain.RoutingRepository) *RoutingUseCase { return &RoutingUseCase{repo: repo} }

func (u *RoutingUseCase) Stats(ctx context.Context, vdom string) (domain.RouteStats, error) {
	return u.repo.Stats(ctx, vdom)
}

type DnsFilterUseCase struct{ repo domain.DnsFilterRepository }

func NewDnsFilter(repo domain.DnsFilterRepository) *DnsFilterUseCase {
	return &DnsFilterUseCase{repo: repo}
}

func (u *DnsFilterUseCase) List(ctx context.Context, vdom string) ([]domain.DnsFilter, error) {
	return u.repo.List(ctx, vdom)
}

func (u *DnsFilterUseCase) Get(ctx context.Context, id int64, vdom string) (domain.DnsFilter, error) {
	return u.repo.Get(ctx, id, vdom)
}
