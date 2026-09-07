package address

import (
	"context"

	domain "github.com/local/fgcli/internal/domain/address"
)

type UseCase struct{ repo domain.Repository }

func New(repo domain.Repository) *UseCase { return &UseCase{repo: repo} }

func (u *UseCase) List(ctx context.Context, vdom string) ([]domain.Address, error) {
	return u.repo.List(ctx, vdom)
}

func (u *UseCase) Get(ctx context.Context, name, vdom string) (domain.Address, error) {
	return u.repo.Get(ctx, name, vdom)
}

func (u *UseCase) Create(ctx context.Context, a domain.Address, vdom string) (domain.Address, error) {
	if err := a.Validate(); err != nil {
		return domain.Address{}, err
	}
	return u.repo.Create(ctx, a, vdom)
}

func (u *UseCase) Update(ctx context.Context, a domain.Address, vdom string) (domain.Address, error) {
	if err := a.Validate(); err != nil {
		return domain.Address{}, err
	}
	return u.repo.Update(ctx, a, vdom)
}

func (u *UseCase) Delete(ctx context.Context, name, vdom string) error {
	return u.repo.Delete(ctx, name, vdom)
}
