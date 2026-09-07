package vip

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalRealShape(t *testing.T) {
	raw := `{"name":"OXYGEN_172.23.9.2_20","type":"static-nat","extip":"103.83.94.250",
		"mappedip":[{"range":"172.27.28.2","q_origin_key":"172.27.28.2"}],
		"extintf":"wan2","portforward":"enable","status":"enable",
		"protocol":"tcp","extport":"20","mappedport":"20"}`
	var v Vip
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		t.Fatalf("real vip shape should decode: %v", err)
	}
	if v.ExtIP != "103.83.94.250" || len(v.MappedIP) != 1 || v.MappedIP[0].Range != "172.27.28.2" {
		t.Fatalf("unexpected vip: %+v", v)
	}
	if err := v.Validate(); err != nil {
		t.Fatal(err)
	}
	if err := (Vip{}).Validate(); err == nil {
		t.Fatal("empty name should fail")
	}
}
