package policy

import (
	"encoding/json"
	"testing"
)

func TestNameUnmarshalBothShapes(t *testing.T) {
	var p Policy
	raw := `{"policyid":16,"srcaddr":[{"name":"all","q_origin_key":"all"}],"dstaddr":[{"name":"wan"}],"service":["HTTPS"],"action":"accept","status":"enable"}`
	if err := json.Unmarshal([]byte(raw), &p); err != nil {
		t.Fatalf("object refs should decode: %v", err)
	}
	if len(p.SrcAddr) != 1 || p.SrcAddr[0] != "all" {
		t.Fatalf("unexpected srcaddr: %+v", p.SrcAddr)
	}
	if len(p.Service) != 1 || p.Service[0] != "HTTPS" {
		t.Fatalf("unexpected service: %+v", p.Service)
	}

	var n Name
	if err := json.Unmarshal([]byte(`"lan-net"`), &n); err != nil || n != "lan-net" {
		t.Fatalf("bare string ref should decode: %v %q", err, n)
	}
}
