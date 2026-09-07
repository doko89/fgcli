package iface

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalSlimView(t *testing.T) {
	var i Iface
	raw := `{"name":"wan1","ip":"103.118.129.220 255.255.255.240","status":"up","type":"physical","alias":"TO BIZNET","mode":"static","allowaccess":"ping"}`
	if err := json.Unmarshal([]byte(raw), &i); err != nil {
		t.Fatalf("interface shape should decode: %v", err)
	}
	if i.IP != "103.118.129.220 255.255.255.240" || i.Alias != "TO BIZNET" {
		t.Fatalf("unexpected interface: %+v", i)
	}
}
