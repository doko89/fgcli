package ippool

import (
	"context"
	"errors"
	"strings"
)

type Pool struct {
	Name      string `json:"name"`
	Type      string `json:"type,omitempty"`
	StartIP   string `json:"startip,omitempty"`
	EndIP     string `json:"endip,omitempty"`
	AssocIntf string `json:"associated-interface,omitempty"`
	Comments  string `json:"comments,omitempty"`
}

func (p Pool) Validate() error {
	if strings.TrimSpace(p.Name) == "" {
		return errors.New("ippool: name is required")
	}
	return nil
}

type Repository interface {
	List(ctx context.Context, vdom string) ([]Pool, error)
	Get(ctx context.Context, name, vdom string) (Pool, error)
	Create(ctx context.Context, p Pool, vdom string) (Pool, error)
	Update(ctx context.Context, p Pool, vdom string) (Pool, error)
	Delete(ctx context.Context, name, vdom string) error
}
