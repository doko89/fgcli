package user

import (
	"context"

	"github.com/local/fgcli/internal/domain/common"
)

type FwUser struct {
	Username string        `json:"username,omitempty"`
	IP       string        `json:"ipaddr,omitempty"`
	Groups   []common.Name `json:"usergroup,omitempty"`
	Method   string        `json:"method,omitempty"`
	Duration int64         `json:"duration_secs,omitempty"`
}

// BannedUser shape is partial: the box reports none, fields are best-effort.
type BannedUser struct {
	Username string `json:"username,omitempty"`
	IP       string `json:"ipaddr,omitempty"`
}

type Repository interface {
	Firewall(ctx context.Context, vdom string) ([]FwUser, error)
	Banned(ctx context.Context, vdom string) ([]BannedUser, error)
}
