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

func TestHelpTarget(t *testing.T) {
	cases := []struct {
		argv   []string
		target string
		ok     bool
	}{
		{[]string{"--help"}, "", true},
		{[]string{"-h"}, "", true},
		{[]string{"help"}, "", true},
		{[]string{"help", "vip"}, "vip", true},
		{[]string{"help", "nope"}, "", true},
		{[]string{"vip", "--help"}, "vip", true},
		{[]string{"--profile", "fg1", "--help"}, "", true},
		{[]string{"--profile", "fg1", "address", "list"}, "", false},
		{[]string{"address", "list"}, "", false},
		// flag values must not trigger help
		{[]string{"raw", "post", "x", "--data", "--help"}, "", false},
		{[]string{"user", "local", "create", "n", "--password", "--help"}, "", false},
	}
	for _, c := range cases {
		target, ok := HelpTarget(c.argv)
		if target != c.target || ok != c.ok {
			t.Fatalf("HelpTarget(%v) = (%q,%v), want (%q,%v)", c.argv, target, ok, c.target, c.ok)
		}
	}
}

func TestTopCommandSkipsValuedFlags(t *testing.T) {
	if got := TopCommand([]string{"--profile", "fg1", "address", "list"}); got != "address" {
		t.Fatalf("TopCommand with --profile = %q", got)
	}
	if got := TopCommand([]string{"--data", "--help", "raw", "get"}); got != "raw" {
		t.Fatalf("TopCommand with --data value = %q", got)
	}
}
