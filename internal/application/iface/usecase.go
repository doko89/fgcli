package iface

import (
	"context"

	domain "github.com/local/fgcli/internal/domain/iface"
)

type IfaceUseCase struct{ repo domain.IfaceRepository }

func NewIface(repo domain.IfaceRepository) *IfaceUseCase { return &IfaceUseCase{repo: repo} }

func (u *IfaceUseCase) List(ctx context.Context, vdom string) ([]domain.Iface, error) {
	return u.repo.List(ctx, vdom)
}

func (u *IfaceUseCase) Get(ctx context.Context, name, vdom string) (domain.Iface, error) {
	return u.repo.Get(ctx, name, vdom)
}

type ZoneUseCase struct{ repo domain.ZoneRepository }

func NewZone(repo domain.ZoneRepository) *ZoneUseCase { return &ZoneUseCase{repo: repo} }

func (u *ZoneUseCase) List(ctx context.Context, vdom string) ([]domain.Zone, error) {
	return u.repo.List(ctx, vdom)
}

func (u *ZoneUseCase) Get(ctx context.Context, name, vdom string) (domain.Zone, error) {
	return u.repo.Get(ctx, name, vdom)
}
