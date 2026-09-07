package snat

import (
	"context"
	"errors"

	"github.com/local/fgcli/internal/domain/common"
)

type Map struct {
	ID        int64         `json:"policyid"`
	Status    string        `json:"status,omitempty"`
	Comments  string        `json:"comments,omitempty"`
	SrcIntf   []common.Name `json:"srcintf,omitempty"`
	DstIntf   []common.Name `json:"dstintf,omitempty"`
	OrigAddr  []common.Name `json:"orig-addr,omitempty"`
	DstAddr   []common.Name `json:"dst-addr,omitempty"`
	NatIPPool []common.Name `json:"nat-ippool,omitempty"`
	Protocol  string        `json:"protocol,omitempty"`
	OrigPort  string        `json:"orig-port,omitempty"`
	NatPort   string        `json:"nat-port,omitempty"`
}

func (m Map) Validate() error {
	if m.ID <= 0 {
		return errors.New("snat: policyid must be > 0")
	}
	return nil
}

type Repository interface {
	List(ctx context.Context, vdom string) ([]Map, error)
	Get(ctx context.Context, id int64, vdom string) (Map, error)
	Create(ctx context.Context, m Map, vdom string) (Map, error)
	Update(ctx context.Context, m Map, vdom string) (Map, error)
	Delete(ctx context.Context, id int64, vdom string) error
}
