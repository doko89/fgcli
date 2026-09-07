package ippool

import (
	"context"

	domain "github.com/local/fgcli/internal/domain/ippool"
)

type UseCase struct{ repo domain.Repository }

func New(repo domain.Repository) *UseCase { return &UseCase{repo: repo} }

func (u *UseCase) List(ctx context.Context, vdom string) ([]domain.Pool, error) {
	return u.repo.List(ctx, vdom)
}

func (u *UseCase) Get(ctx context.Context, name, vdom string) (domain.Pool, error) {
	return u.repo.Get(ctx, name, vdom)
}

func (u *UseCase) Create(ctx context.Context, p domain.Pool, vdom string) (domain.Pool, error) {
	if err := p.Validate(); err != nil {
		return domain.Pool{}, err
	}
	return u.repo.Create(ctx, p, vdom)
}

func (u *UseCase) Update(ctx context.Context, p domain.Pool, vdom string) (domain.Pool, error) {
	if err := p.Validate(); err != nil {
		return domain.Pool{}, err
	}
	return u.repo.Update(ctx, p, vdom)
}

func (u *UseCase) Delete(ctx context.Context, name, vdom string) error {
	return u.repo.Delete(ctx, name, vdom)
}
