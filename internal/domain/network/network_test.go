package network

import (
	"encoding/json"
	"testing"
)

func TestTunnelSlimDecode(t *testing.T) {
	raw := `{"name":"RT_SitetoIDXREP","comments":"Fiberstar","connection_count":28,
		"incoming_bytes":336185088,"outgoing_bytes":139886108,"rgwy":"139.255.35.134",
		"proxyid":[{"status":"up","p2name":"Tunnel_To_RTReplika"}]}`
	var tun Tunnel
	if err := json.Unmarshal([]byte(raw), &tun); err != nil {
		t.Fatal(err)
	}
	if tun.Name != "RT_SitetoIDXREP" || tun.ConnCount != 28 || tun.Rgwy != "139.255.35.134" {
		t.Fatalf("unexpected tunnel: %+v", tun)
	}
}

func TestRouteStatsDecode(t *testing.T) {
	var st RouteStats
	if err := json.Unmarshal([]byte(`{"total_lines":28,"total_lines_ipv4":28,"total_lines_ipv6":0}`), &st); err != nil {
		t.Fatal(err)
	}
	if st.Total != 28 || st.IPv4 != 28 || st.IPv6 != 0 {
		t.Fatalf("unexpected stats: %+v", st)
	}
}
