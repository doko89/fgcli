package cli

import (
	"testing"

	daddr "github.com/local/fgcli/internal/domain/address"
)

func TestMatchFilter(t *testing.T) {
	a := daddr.Address{Name: "MAIL_SRV", Subnet: "192.168.1.43/32"}
	if !matchFilter(a, "") {
		t.Fatal("empty filter must match")
	}
	if !matchFilter(a, "192.168.1.43") {
		t.Fatal("subnet should match")
	}
	if !matchFilter(a, "mail_srv") {
		t.Fatal("match should be case-insensitive")
	}
	if matchFilter(a, "10.9.9.9") {
		t.Fatal("non-matching subnet must not match")
	}
}

func TestFilterSlice(t *testing.T) {
	items := []daddr.Address{{Name: "a"}, {Name: "b"}}
	if got := filterSlice(items, ""); len(got) != 2 {
		t.Fatalf("empty filter keeps all: %d", len(got))
	}
	if got := filterSlice(items, "b"); len(got) != 1 || got[0].Name != "b" {
		t.Fatalf("unexpected filter result: %+v", got)
	}
}
