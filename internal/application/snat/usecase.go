package snat

import (
	"context"

	domain "github.com/local/fgcli/internal/domain/snat"
)

type UseCase struct{ repo domain.Repository }

func New(repo domain.Repository) *UseCase { return &UseCase{repo: repo} }

func (u *UseCase) List(ctx context.Context, vdom string) ([]domain.Map, error) {
	return u.repo.List(ctx, vdom)
}

func (u *UseCase) Get(ctx context.Context, id int64, vdom string) (domain.Map, error) {
	return u.repo.Get(ctx, id, vdom)
}

func (u *UseCase) Create(ctx context.Context, m domain.Map, vdom string) (domain.Map, error) {
	if err := m.Validate(); err != nil {
		return domain.Map{}, err
	}
	return u.repo.Create(ctx, m, vdom)
}

func (u *UseCase) Update(ctx context.Context, m domain.Map, vdom string) (domain.Map, error) {
	if err := m.Validate(); err != nil {
		return domain.Map{}, err
	}
	return u.repo.Update(ctx, m, vdom)
}

func (u *UseCase) Delete(ctx context.Context, id int64, vdom string) error {
	return u.repo.Delete(ctx, id, vdom)
}
