package iface

import (
	"context"

	"github.com/local/fgcli/internal/domain/common"
)

// Iface is a slim view of system/interface (the API object has 200+
// fields; only the audit-relevant ones are kept).
type Iface struct {
	Name   string `json:"name"`
	IP     string `json:"ip,omitempty"`
	Status string `json:"status,omitempty"`
	Type   string `json:"type,omitempty"`
	Alias  string `json:"alias,omitempty"`
	Vdom   string `json:"vdom,omitempty"`
}

type Zone struct {
	Name        string        `json:"name"`
	Interface   []common.Name `json:"interface,omitempty"`
	Intrazone   string        `json:"intrazone,omitempty"`
	Description string        `json:"description,omitempty"`
}

// Read-only: interfaces/zones are platform config, provisioned outside
// this tool. List/Get is enough to validate policy/vip references.
type IfaceRepository interface {
	List(ctx context.Context, vdom string) ([]Iface, error)
	Get(ctx context.Context, name, vdom string) (Iface, error)
}

type ZoneRepository interface {
	List(ctx context.Context, vdom string) ([]Zone, error)
	Get(ctx context.Context, name, vdom string) (Zone, error)
}
