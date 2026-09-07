package service

import (
	"context"
	"errors"
	"strings"

	"github.com/local/fgcli/internal/domain/common"
)

type Svc struct {
	Name      string `json:"name"`
	Protocol  string `json:"protocol,omitempty"`
	TCPPorts  string `json:"tcp-portrange,omitempty"`
	UDPPorts  string `json:"udp-portrange,omitempty"`
	SCTPPorts string `json:"sctp-portrange,omitempty"`
	Comment   string `json:"comment,omitempty"`
	Category  string `json:"category,omitempty"`
}

func (s Svc) Validate() error {
	if strings.TrimSpace(s.Name) == "" {
		return errors.New("service: name is required")
	}
	return nil
}

type Group struct {
	Name    string        `json:"name"`
	Member  []common.Name `json:"member,omitempty"`
	Comment string        `json:"comment,omitempty"`
}

func (g Group) Validate() error {
	if strings.TrimSpace(g.Name) == "" {
		return errors.New("service-group: name is required")
	}
	return nil
}

type SvcRepository interface {
	List(ctx context.Context, vdom string) ([]Svc, error)
	Get(ctx context.Context, name, vdom string) (Svc, error)
	Create(ctx context.Context, s Svc, vdom string) (Svc, error)
	Update(ctx context.Context, s Svc, vdom string) (Svc, error)
	Delete(ctx context.Context, name, vdom string) error
}

type GroupRepository interface {
	List(ctx context.Context, vdom string) ([]Group, error)
	Get(ctx context.Context, name, vdom string) (Group, error)
	Create(ctx context.Context, g Group, vdom string) (Group, error)
	Update(ctx context.Context, g Group, vdom string) (Group, error)
	Delete(ctx context.Context, name, vdom string) error
}
