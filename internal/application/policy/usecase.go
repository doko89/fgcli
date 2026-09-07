package policy

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/local/fgcli/internal/domain/policy"
)

type UseCase struct{ repo domain.Repository }

func New(repo domain.Repository) *UseCase { return &UseCase{repo: repo} }

func (u *UseCase) List(ctx context.Context, vdom string) ([]domain.Policy, error) {
	return u.repo.List(ctx, vdom)
}

func (u *UseCase) Get(ctx context.Context, id int64, vdom string) (domain.Policy, error) {
	return u.repo.Get(ctx, id, vdom)
}

func (u *UseCase) Delete(ctx context.Context, id int64, vdom string) error {
	return u.repo.Delete(ctx, id, vdom)
}

func (u *UseCase) Create(ctx context.Context, p domain.Policy, vdom string) (domain.Policy, error) {
	if err := p.Validate(); err != nil {
		return domain.Policy{}, err
	}
	return u.repo.Create(ctx, p, vdom)
}

func (u *UseCase) Update(ctx context.Context, p domain.Policy, vdom string) (domain.Policy, error) {
	if err := p.Validate(); err != nil {
		return domain.Policy{}, err
	}
	return u.repo.Update(ctx, p, vdom)
}

func (u *UseCase) SetStatus(ctx context.Context, id int64, status, vdom string) (domain.Policy, error) {
	if status != "enable" && status != "disable" {
		return domain.Policy{}, errors.New("policy: status must be enable|disable")
	}
	p, err := u.repo.Get(ctx, id, vdom)
	if err != nil {
		return domain.Policy{}, err
	}
	p.Status = status
	if err := p.Validate(); err != nil {
		return domain.Policy{}, err
	}
	return u.repo.Update(ctx, p, vdom)
}

func (u *UseCase) Move(ctx context.Context, id, before, after int64, vdom string) error {
	if (before <= 0) == (after <= 0) {
		return errors.New("policy: exactly one of --before/--after is required")
	}
	return u.repo.Move(ctx, id, before, after, vdom)
}

func (u *UseCase) Clone(ctx context.Context, srcID, newID int64, name, vdom string) (domain.Policy, error) {
	src, err := u.repo.Get(ctx, srcID, vdom)
	if err != nil {
		return domain.Policy{}, err
	}
	src.ID = newID
	if name != "" {
		src.Name = name
	} else {
		src.Name = fmt.Sprintf("%s_clone%d", src.Name, newID)
	}
	if err := src.Validate(); err != nil {
		return domain.Policy{}, err
	}
	return u.repo.Create(ctx, src, vdom)
}
