package user

import (
	"context"
	"errors"
	"strings"

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

func (u *UseCase) GroupList(ctx context.Context, vdom string) ([]domain.UserGroup, error) {
	return u.repo.GroupList(ctx, vdom)
}

func (u *UseCase) GroupGet(ctx context.Context, name, vdom string) (domain.UserGroup, error) {
	if strings.TrimSpace(name) == "" {
		return domain.UserGroup{}, errors.New("user group: name is required")
	}
	return u.repo.GroupGet(ctx, name, vdom)
}

func (u *UseCase) GroupUpdate(ctx context.Context, g domain.UserGroup, vdom string) (domain.UserGroup, error) {
	if err := g.Validate(); err != nil {
		return domain.UserGroup{}, err
	}
	return u.repo.GroupUpdate(ctx, g, vdom)
}

func (u *UseCase) LocalList(ctx context.Context, vdom string) ([]domain.LocalUser, error) {
	return u.repo.LocalList(ctx, vdom)
}

func (u *UseCase) LocalGet(ctx context.Context, name, vdom string) (domain.LocalUser, error) {
	if strings.TrimSpace(name) == "" {
		return domain.LocalUser{}, errors.New("user local: name is required")
	}
	return u.repo.LocalGet(ctx, name, vdom)
}

func (u *UseCase) LocalCreate(ctx context.Context, usr domain.LocalUser, vdom string) (domain.LocalUser, error) {
	if err := usr.Validate(); err != nil {
		return domain.LocalUser{}, err
	}
	if usr.Passwd == "" {
		return domain.LocalUser{}, errors.New("user local: --password is required")
	}
	return u.repo.LocalCreate(ctx, usr, vdom)
}

func (u *UseCase) LocalDelete(ctx context.Context, name, vdom string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("user local: name is required")
	}
	return u.repo.LocalDelete(ctx, name, vdom)
}
