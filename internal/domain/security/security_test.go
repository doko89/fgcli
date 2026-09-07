package security

import (
	"encoding/json"
	"testing"
)

func TestShapesDecode(t *testing.T) {
	var a []IpsAnomaly
	if err := json.Unmarshal([]byte(`[{"id":100663396,"name":"tcp_syn_flood","threshold":2000,"action":7}]`), &a); err != nil {
		t.Fatal(err)
	}
	if len(a) != 1 || a[0].Name != "tcp_syn_flood" || a[0].Threshold != 2000 {
		t.Fatalf("unexpected anomaly: %+v", a)
	}
	var d DlpPattern
	if err := json.Unmarshal([]byte(`{"id":1,"name":"builtin-patterns","entries":[{"pattern":"*.exe","file-type":"unknown"}]}`), &d); err != nil {
		t.Fatal(err)
	}
	if d.Name != "builtin-patterns" || len(d.Entries) != 1 || d.Entries[0].Pattern != "*.exe" {
		t.Fatalf("unexpected dlp: %+v", d)
	}
	var v WafClass
	if err := json.Unmarshal([]byte(`{"id":2,"name":"default","comment":""}`), &v); err != nil {
		t.Fatal(err)
	}
	if v.ID != 2 || v.Name != "default" {
		t.Fatalf("unexpected waf: %+v", v)
	}
}
