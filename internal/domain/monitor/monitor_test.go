package monitor

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestWifiDecode(t *testing.T) {
	var w WifiStatus
	raw := `{"wtp_active":0,"wtp_down":0,"wtp_rebooted":0,"client_count":0,"client_count_max":0}`
	if err := json.Unmarshal([]byte(raw), &w); err != nil {
		t.Fatal(err)
	}
	if w.WtpActive != 0 || w.ClientCount != 0 {
		t.Fatalf("unexpected wifi: %+v", w)
	}
}

func TestZeroCountersSurviveMarshal(t *testing.T) {
	b, err := json.Marshal(WifiStatus{})
	if err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"wtp_active", "client_count"} {
		if !strings.Contains(string(b), k) {
			t.Fatalf("zero counter %q must be kept, got %s", k, b)
		}
	}
}

func TestFortiViewDecode(t *testing.T) {
	var f FortiView
	raw := `{"summary":{"start":1788149423,"end":1788754223},"details":[]}`
	if err := json.Unmarshal([]byte(raw), &f); err != nil {
		t.Fatal(err)
	}
	if f.Summary.Start != 1788149423 || len(f.Details) != 0 {
		t.Fatalf("unexpected fortiview: %+v", f)
	}
}
