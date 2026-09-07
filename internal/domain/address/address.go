package address

import (
	"context"
	"errors"
	"strings"
)

type Address struct {
	Name    string `json:"name"`
	Subnet  string `json:"subnet,omitempty"`
	Type    string `json:"type,omitempty"`
	Comment string `json:"comment,omitempty"`
}

func (a Address) Validate() error {
	if strings.TrimSpace(a.Name) == "" {
		return errors.New("address: name is required")
	}
	if strings.ContainsAny(a.Name, " /\\") {
		return errors.New("address: name must not contain spaces or slashes")
	}
	return nil
}

type Repository interface {
	List(ctx context.Context, vdom string) ([]Address, error)
	Get(ctx context.Context, name, vdom string) (Address, error)
	Create(ctx context.Context, a Address, vdom string) (Address, error)
	Update(ctx context.Context, a Address, vdom string) (Address, error)
	Delete(ctx context.Context, name, vdom string) error
}
