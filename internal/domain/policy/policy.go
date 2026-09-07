package policy

import (
	"context"
	"errors"

	"github.com/local/fgcli/internal/domain/common"
)

// Name aliases the shared object-reference type.
type Name = common.Name

type Policy struct {
	ID       int64  `json:"policyid"`
	Name     string `json:"name,omitempty"`
	SrcIntf  []Name `json:"srcintf,omitempty"`
	DstIntf  []Name `json:"dstintf,omitempty"`
	SrcAddr  []Name `json:"srcaddr,omitempty"`
	DstAddr  []Name `json:"dstaddr,omitempty"`
	Service  []Name `json:"service,omitempty"`
	Action   string `json:"action,omitempty"`
	Status   string `json:"status,omitempty"`
	Schedule string `json:"schedule,omitempty"`
	Comments string `json:"comments,omitempty"`
	Log      string `json:"logtraffic,omitempty"`
	NAT      string `json:"nat,omitempty"`
	PoolName []Name `json:"poolname,omitempty"`
}

func (p Policy) Validate() error {
	if p.ID <= 0 {
		return errors.New("policy: policyid must be > 0")
	}
	return nil
}

type Repository interface {
	List(ctx context.Context, vdom string) ([]Policy, error)
	Get(ctx context.Context, id int64, vdom string) (Policy, error)
	Create(ctx context.Context, p Policy, vdom string) (Policy, error)
	Update(ctx context.Context, p Policy, vdom string) (Policy, error)
	// Move reorders id before/after another id (exactly one of before/after > 0).
	Move(ctx context.Context, id, before, after int64, vdom string) error
	Delete(ctx context.Context, id int64, vdom string) error
}
