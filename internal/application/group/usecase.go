package group

import (
	"context"

	domain "github.com/local/fgcli/internal/domain/group"
)

type UseCase struct{ repo domain.Repository }

func New(repo domain.Repository) *UseCase { return &UseCase{repo: repo} }

func (u *UseCase) List(ctx context.Context, vdom string) ([]domain.Group, error) {
	return u.repo.List(ctx, vdom)
}

func (u *UseCase) Get(ctx context.Context, name, vdom string) (domain.Group, error) {
	return u.repo.Get(ctx, name, vdom)
}

func (u *UseCase) Create(ctx context.Context, g domain.Group, vdom string) (domain.Group, error) {
	if err := g.Validate(); err != nil {
		return domain.Group{}, err
	}
	return u.repo.Create(ctx, g, vdom)
}

func (u *UseCase) Update(ctx context.Context, g domain.Group, vdom string) (domain.Group, error) {
	if err := g.Validate(); err != nil {
		return domain.Group{}, err
	}
	return u.repo.Update(ctx, g, vdom)
}

func (u *UseCase) Delete(ctx context.Context, name, vdom string) error {
	return u.repo.Delete(ctx, name, vdom)
}
