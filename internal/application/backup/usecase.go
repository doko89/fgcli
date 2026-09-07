package backup

import (
	"context"

	domain "github.com/local/fgcli/internal/domain/backup"
)

type UseCase struct{ repo domain.Repository }

func New(repo domain.Repository) *UseCase { return &UseCase{repo: repo} }

func (u *UseCase) Download(ctx context.Context, scope string) ([]byte, error) {
	if scope == "" {
		scope = "global"
	}
	return u.repo.Download(ctx, scope)
}
