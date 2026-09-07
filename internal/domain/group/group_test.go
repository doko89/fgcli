package group

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalRealAddrGrp(t *testing.T) {
	raw := `{"name":"BO_RT","member":[{"name":"N_192.168.18.0/24"},{"name":"N_10.20.30.0/24"}],"comment":""}`
	var g Group
	if err := json.Unmarshal([]byte(raw), &g); err != nil {
		t.Fatalf("real addrgrp shape should decode: %v", err)
	}
	if len(g.Member) != 2 || g.Member[0] != "N_192.168.18.0/24" {
		t.Fatalf("unexpected group: %+v", g)
	}
	if err := (Group{}).Validate(); err == nil {
		t.Fatal("empty name should fail")
	}
}
