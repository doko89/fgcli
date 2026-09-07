package svc

import (
	"encoding/json"
	"testing"
)

func TestLicenseSlimDecode(t *testing.T) {
	raw := `{"forticare":{"status":"registered","account":"it@example.id",
		"support":{"hardware":{"status":"expires_soon","support_level":"Advanced HW","expires":1791244800}}},
		"fortiguard":{"supported":true,"connected":true,"server_address":"10.0.0.1:443"}}`
	var l License
	if err := json.Unmarshal([]byte(raw), &l); err != nil {
		t.Fatal(err)
	}
	if l.Forticare.Account != "it@example.id" || l.Forticare.Status != "registered" {
		t.Fatalf("unexpected forticare: %+v", l.Forticare)
	}
	if l.Forticare.Support.Hardware.Status != "expires_soon" {
		t.Fatalf("unexpected support: %+v", l.Forticare.Support)
	}
	if !l.Fortiguard.Connected || l.Fortiguard.Server != "10.0.0.1:443" {
		t.Fatalf("unexpected fortiguard: %+v", l.Fortiguard)
	}
}

func TestFortiguardMapDecode(t *testing.T) {
	var f Fortiguard
	raw := `{"forticare":{"1_hour":[0,1],"24_hour":[2],"1_week":[3]}}`
	if err := json.Unmarshal([]byte(raw), &f); err != nil {
		t.Fatal(err)
	}
	if len(f["forticare"].Hour) != 2 || f["forticare"].Day[0] != 2 {
		t.Fatalf("unexpected buckets: %+v", f)
	}
}

func TestNtpAndSnmpDecode(t *testing.T) {
	var n []NtpServer
	if err := json.Unmarshal([]byte(`[{"server":"ntp2.fortiguard.com","ip":"208.91.112.62","reachable":true,"stratum":4,"offset":33.75}]`), &n); err != nil {
		t.Fatal(err)
	}
	if len(n) != 1 || !n[0].Reachable || n[0].Stratum != 4 {
		t.Fatalf("unexpected ntp: %+v", n)
	}
	var s Snmp
	if err := json.Unmarshal([]byte(`{"status":"disable","description":"","contact-info":"","location":""}`), &s); err != nil {
		t.Fatal(err)
	}
	if s.Status != "disable" {
		t.Fatalf("unexpected snmp: %+v", s)
	}
}
