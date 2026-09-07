package log

import (
	"context"

	domain "github.com/local/fgcli/internal/domain/log"
)

type UseCase struct{ repo domain.Repository }

func New(repo domain.Repository) *UseCase { return &UseCase{repo: repo} }

func (u *UseCase) Fetch(ctx context.Context, source, logtype, vdom string) ([]domain.Entry, error) {
	if err := domain.ValidSource(source); err != nil {
		return nil, err
	}
	return u.repo.Fetch(ctx, source, logtype, vdom)
}
