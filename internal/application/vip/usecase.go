package vip

import (
	"context"

	domain "github.com/local/fgcli/internal/domain/vip"
)

type UseCase struct{ repo domain.Repository }

func New(repo domain.Repository) *UseCase { return &UseCase{repo: repo} }

func (u *UseCase) List(ctx context.Context, vdom string) ([]domain.Vip, error) {
	return u.repo.List(ctx, vdom)
}

func (u *UseCase) Get(ctx context.Context, name, vdom string) (domain.Vip, error) {
	return u.repo.Get(ctx, name, vdom)
}

func (u *UseCase) Create(ctx context.Context, v domain.Vip, vdom string) (domain.Vip, error) {
	if err := v.Validate(); err != nil {
		return domain.Vip{}, err
	}
	return u.repo.Create(ctx, v, vdom)
}

func (u *UseCase) Update(ctx context.Context, v domain.Vip, vdom string) (domain.Vip, error) {
	if err := v.Validate(); err != nil {
		return domain.Vip{}, err
	}
	return u.repo.Update(ctx, v, vdom)
}

func (u *UseCase) Delete(ctx context.Context, name, vdom string) error {
	return u.repo.Delete(ctx, name, vdom)
}
