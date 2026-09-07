package user

import (
	"context"

	domain "github.com/local/fgcli/internal/domain/user"
)

type UseCase struct{ repo domain.Repository }

func New(repo domain.Repository) *UseCase { return &UseCase{repo: repo} }

func (u *UseCase) Firewall(ctx context.Context, vdom string) ([]domain.FwUser, error) {
	return u.repo.Firewall(ctx, vdom)
}

func (u *UseCase) Banned(ctx context.Context, vdom string) ([]domain.BannedUser, error) {
	return u.repo.Banned(ctx, vdom)
}
