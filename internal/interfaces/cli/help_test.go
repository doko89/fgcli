package cli

import "testing"

func TestHelpForKnownResources(t *testing.T) {
	for _, r := range []string{"vip", "vipgrp", "addrgrp", "policy", "backup", "trace", "raw"} {
		if txt, ok := HelpFor(r); !ok || txt == "" {
			t.Fatalf("missing help for %q", r)
		}
	}
	if _, ok := HelpFor("nope"); ok {
		t.Fatal("unknown resource should not have help")
	}
}

func TestWantHelp(t *testing.T) {
	if !WantHelp([]string{"--help"}) || !WantHelp([]string{"list", "--filter", "x", "-h"}) {
		t.Fatal("help flags should be detected")
	}
	if WantHelp([]string{"list", "--filter", "x"}) {
		t.Fatal("normal args must not trigger help")
	}
}
