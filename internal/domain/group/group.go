package group

import (
	"context"
	"errors"
	"strings"

	"github.com/local/fgcli/internal/domain/common"
)

// Group covers vipgrp/addrgrp ({name, member[], comment}).
type Group struct {
	Name    string        `json:"name"`
	Member  []common.Name `json:"member,omitempty"`
	Comment string        `json:"comment,omitempty"`
}

func (g Group) Validate() error {
	if strings.TrimSpace(g.Name) == "" {
		return errors.New("group: name is required")
	}
	return nil
}

type Repository interface {
	List(ctx context.Context, vdom string) ([]Group, error)
	Get(ctx context.Context, name, vdom string) (Group, error)
	Create(ctx context.Context, g Group, vdom string) (Group, error)
	Update(ctx context.Context, g Group, vdom string) (Group, error)
	Delete(ctx context.Context, name, vdom string) error
}
