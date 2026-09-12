package raw

import (
	"context"
	"encoding/json"

	domain "github.com/local/fgcli/internal/domain/raw"
)

type UseCase struct{ repo domain.Doer }

func New(repo domain.Doer) *UseCase { return &UseCase{repo: repo} }

func (u *UseCase) Do(ctx context.Context, method, path, vdom string, body json.RawMessage) (json.RawMessage, error) {
	if err := domain.ValidMethod(method); err != nil {
		return nil, err
	}
	return u.repo.Do(ctx, method, path, vdom, body)
}
