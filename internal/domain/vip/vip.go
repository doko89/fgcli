package vip

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
)

type MappedRange struct {
	Range string `json:"range"`
}

func (m *MappedRange) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		m.Range = s
		return nil
	}
	type raw MappedRange
	return json.Unmarshal(b, (*raw)(m))
}

type Vip struct {
	Name        string        `json:"name"`
	UUID        string        `json:"uuid,omitempty"`
	Type        string        `json:"type,omitempty"`
	ExtIP       string        `json:"extip,omitempty"`
	ExtIntf     string        `json:"extintf,omitempty"`
	MappedIP    []MappedRange `json:"mappedip,omitempty"`
	PortForward string        `json:"portforward,omitempty"`
	Protocol    string        `json:"protocol,omitempty"`
	ExtPort     string        `json:"extport,omitempty"`
	MappedPort  string        `json:"mappedport,omitempty"`
	Comment     string        `json:"comment,omitempty"`
	Status      string        `json:"status,omitempty"`
}

func (v Vip) Validate() error {
	if strings.TrimSpace(v.Name) == "" {
		return errors.New("vip: name is required")
	}
	return nil
}

type Repository interface {
	List(ctx context.Context, vdom string) ([]Vip, error)
	Get(ctx context.Context, name, vdom string) (Vip, error)
	Create(ctx context.Context, v Vip, vdom string) (Vip, error)
	Update(ctx context.Context, v Vip, vdom string) (Vip, error)
	Delete(ctx context.Context, name, vdom string) error
}
