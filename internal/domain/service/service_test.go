package service

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalRealShapes(t *testing.T) {
	var s Svc
	raw := `{"name":"TCP_822","protocol":"TCP","tcp-portrange":"822","comment":""}`
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		t.Fatalf("real service shape should decode: %v", err)
	}
	if s.TCPPorts != "822" || s.Protocol != "TCP" {
		t.Fatalf("unexpected service: %+v", s)
	}
	if err := (Svc{}).Validate(); err == nil {
		t.Fatal("empty name should fail")
	}

	var g Group
	rawg := `{"name":"Email Access","member":[{"name":"SMTP"},{"name":"IMAPS"}]}`
	if err := json.Unmarshal([]byte(rawg), &g); err != nil {
		t.Fatalf("real group shape should decode: %v", err)
	}
	if len(g.Member) != 2 || g.Member[0] != "SMTP" {
		t.Fatalf("unexpected group: %+v", g)
	}
}
