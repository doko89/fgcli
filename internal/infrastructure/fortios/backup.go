package fortios

import (
	"context"
	"net/url"
)

type BackupRepo struct{ c *Client }

func NewBackupRepo(c *Client) *BackupRepo { return &BackupRepo{c: c} }

func (s *BackupRepo) Download(ctx context.Context, scope string) ([]byte, error) {
	q := url.Values{"scope": {scope}}
	return s.c.do(ctx, "GET", "/api/v2/monitor/system/config/backup", q, nil)
}
