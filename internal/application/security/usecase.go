package security

import (
	"context"

	domain "github.com/local/fgcli/internal/domain/security"
)

type UseCase struct{ repo domain.Repository }

func New(repo domain.Repository) *UseCase { return &UseCase{repo: repo} }

func (u *UseCase) Ips(ctx context.Context, vdom string) ([]domain.IpsAnomaly, error) {
	return u.repo.Ips(ctx, vdom)
}

func (u *UseCase) ListWaf(ctx context.Context, vdom string) ([]domain.WafClass, error) {
	return u.repo.ListWaf(ctx, vdom)
}

func (u *UseCase) GetWaf(ctx context.Context, id int64, vdom string) (domain.WafClass, error) {
	return u.repo.GetWaf(ctx, id, vdom)
}

func (u *UseCase) ListDlp(ctx context.Context, vdom string) ([]domain.DlpPattern, error) {
	return u.repo.ListDlp(ctx, vdom)
}

func (u *UseCase) GetDlp(ctx context.Context, id int64, vdom string) (domain.DlpPattern, error) {
	return u.repo.GetDlp(ctx, id, vdom)
}

func (u *UseCase) DownloadPac(ctx context.Context, vdom string) ([]byte, error) {
	return u.repo.DownloadPac(ctx, vdom)
}
