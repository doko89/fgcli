package user

import (
	"encoding/json"
	"testing"
)

func TestFwUserDecode(t *testing.T) {
	raw := `{"username":"FZ_14A","ipaddr":"10.212.134.200",
		"usergroup":[{"type":"user","name":"FZ_14A"},{"type":"group","name":"VPN Vendor"}],
		"duration_secs":1507,"method":"Firewall"}`
	var u FwUser
	if err := json.Unmarshal([]byte(raw), &u); err != nil {
		t.Fatal(err)
	}
	if u.Username != "FZ_14A" || u.IP != "10.212.134.200" || len(u.Groups) != 2 || u.Groups[1] != "VPN Vendor" {
		t.Fatalf("unexpected user: %+v", u)
	}
}
