package service

import "context"

type VoipGetter interface {
	Get(ctx context.Context, name, vdom string) (VoipProfile, error)
}

type SipSettings struct {
	Status string `json:"status,omitempty"`
}

type VoipProfile struct {
	Name       string      `json:"name"`
	FeatureSet string      `json:"feature-set,omitempty"`
	Comment    string      `json:"comment,omitempty"`
	Sip        SipSettings `json:"sip,omitempty"`
}
