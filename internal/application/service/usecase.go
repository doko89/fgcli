package service

import (
	"context"

	domain "github.com/local/fgcli/internal/domain/service"
)

type SvcUseCase struct{ repo domain.SvcRepository }

func NewSvc(repo domain.SvcRepository) *SvcUseCase { return &SvcUseCase{repo: repo} }

func (u *SvcUseCase) List(ctx context.Context, vdom string) ([]domain.Svc, error) {
	return u.repo.List(ctx, vdom)
}

func (u *SvcUseCase) Get(ctx context.Context, name, vdom string) (domain.Svc, error) {
	return u.repo.Get(ctx, name, vdom)
}

func (u *SvcUseCase) Create(ctx context.Context, s domain.Svc, vdom string) (domain.Svc, error) {
	if err := s.Validate(); err != nil {
		return domain.Svc{}, err
	}
	return u.repo.Create(ctx, s, vdom)
}

func (u *SvcUseCase) Update(ctx context.Context, s domain.Svc, vdom string) (domain.Svc, error) {
	if err := s.Validate(); err != nil {
		return domain.Svc{}, err
	}
	return u.repo.Update(ctx, s, vdom)
}

func (u *SvcUseCase) Delete(ctx context.Context, name, vdom string) error {
	return u.repo.Delete(ctx, name, vdom)
}

type GroupUseCase struct{ repo domain.GroupRepository }

func NewGroup(repo domain.GroupRepository) *GroupUseCase { return &GroupUseCase{repo: repo} }

func (u *GroupUseCase) List(ctx context.Context, vdom string) ([]domain.Group, error) {
	return u.repo.List(ctx, vdom)
}

func (u *GroupUseCase) Get(ctx context.Context, name, vdom string) (domain.Group, error) {
	return u.repo.Get(ctx, name, vdom)
}

func (u *GroupUseCase) Create(ctx context.Context, g domain.Group, vdom string) (domain.Group, error) {
	if err := g.Validate(); err != nil {
		return domain.Group{}, err
	}
	return u.repo.Create(ctx, g, vdom)
}

func (u *GroupUseCase) Update(ctx context.Context, g domain.Group, vdom string) (domain.Group, error) {
	if err := g.Validate(); err != nil {
		return domain.Group{}, err
	}
	return u.repo.Update(ctx, g, vdom)
}

func (u *GroupUseCase) Delete(ctx context.Context, name, vdom string) error {
	return u.repo.Delete(ctx, name, vdom)
}

type VoipUseCase struct{ repo domain.VoipGetter }

func NewVoip(repo domain.VoipGetter) *VoipUseCase { return &VoipUseCase{repo: repo} }

func (u *VoipUseCase) Get(ctx context.Context, name, vdom string) (domain.VoipProfile, error) {
	return u.repo.Get(ctx, name, vdom)
}
