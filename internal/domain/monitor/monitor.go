package monitor

import (
	"context"
	"encoding/json"
)

// ManagedSwitch shape is best-effort: the box reports none.
type ManagedSwitch struct {
	Name    string `json:"name,omitempty"`
	Status  string `json:"status,omitempty"`
	Model   string `json:"model,omitempty"`
	Version string `json:"version,omitempty"`
}

type WifiStatus struct {
	WtpActive   int `json:"wtp_active"`
	WtpDown     int `json:"wtp_down"`
	WtpRebooted int `json:"wtp_rebooted"`
	ClientCount int `json:"client_count"`
	ClientMax   int `json:"client_count_max"`
}

type FortiViewSummary struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

type FortiView struct {
	Summary FortiViewSummary `json:"summary"`
	// Details schema varies by dataset; passthrough preserves all fields.
	Details []json.RawMessage `json:"details,omitempty"`
}

type Repository interface {
	Switches(ctx context.Context, vdom string) ([]ManagedSwitch, error)
	Wifi(ctx context.Context, vdom string) (WifiStatus, error)
	FortiView(ctx context.Context, vdom string) (FortiView, error)
}
