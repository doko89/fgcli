package monitor

import (
	"context"

	domain "github.com/local/fgcli/internal/domain/monitor"
)

type UseCase struct{ repo domain.Repository }

func New(repo domain.Repository) *UseCase { return &UseCase{repo: repo} }

func (u *UseCase) Switches(ctx context.Context, vdom string) ([]domain.ManagedSwitch, error) {
	return u.repo.Switches(ctx, vdom)
}

func (u *UseCase) Wifi(ctx context.Context, vdom string) (domain.WifiStatus, error) {
	return u.repo.Wifi(ctx, vdom)
}

func (u *UseCase) FortiView(ctx context.Context, vdom string) (domain.FortiView, error) {
	return u.repo.FortiView(ctx, vdom)
}
