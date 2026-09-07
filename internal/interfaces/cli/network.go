package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/local/fgcli/internal/infrastructure/output"
)

func runSystem(ctx context.Context, d Deps, args []string) int {
	if len(args) < 1 {
		output.Fail("usage", fmt.Errorf("usage: fgcli system status|license|fortiguard|ntp|dns|dhcp|snmp"))
		return 1
	}
	verb, rest := args[0], args[1:]
	g, _, _ := parseGlobals(rest)
	switch verb {
	case "status":
		st, err := d.Sys.Status(ctx)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(st)
		return 0
	case "license":
		lic, err := d.Sys.License(ctx)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(lic)
		return 0
	case "fortiguard":
		fg, err := d.Sys.Fortiguard(ctx)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(fg)
		return 0
	case "ntp":
		list, err := d.Sys.Ntp(ctx)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(filterSlice(list, g.filter))
		return 0
	case "dns":
		list, err := d.Sys.Dns(ctx)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(filterSlice(list, g.filter))
		return 0
	case "dhcp":
		list, err := d.Sys.Dhcp(ctx)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(filterSlice(list, g.filter))
		return 0
	case "snmp":
		s, err := d.Sys.Snmp(ctx)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(s)
		return 0
	default:
		output.Fail("usage", fmt.Errorf("unknown system verb %q", verb))
		return 1
	}
}

func runVpn(ctx context.Context, d Deps, args []string) int {
	if len(args) < 1 || args[0] != "ipsec" {
		output.Fail("usage", fmt.Errorf("usage: fgcli vpn ipsec [--filter SUB]"))
		return 1
	}
	g, _, _ := parseGlobals(args[1:])
	list, err := d.Vpn.List(ctx, g.vdom)
	if err != nil {
		output.Fail("api_error", err)
		return 2
	}
	output.Print(filterSlice(list, g.filter))
	return 0
}

func runNetwork(ctx context.Context, d Deps, args []string) int {
	if len(args) < 1 || args[0] != "dnsfilter" {
		output.Fail("usage", fmt.Errorf("usage: fgcli network dnsfilter list|get <id>"))
		return 1
	}
	rest := args[1:]
	g, rest, _ := parseGlobals(rest)
	if len(rest) < 1 || rest[0] == "list" {
		list, err := d.DnsFilter.List(ctx, g.vdom)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(filterSlice(list, g.filter))
		return 0
	}
	if rest[0] == "get" {
		rest = rest[1:]
	}
	if len(rest) < 1 {
		output.Fail("usage", fmt.Errorf("usage: fgcli network dnsfilter get <id>"))
		return 1
	}
	id, err := strconv.ParseInt(rest[0], 10, 64)
	if err != nil {
		output.Fail("input_error", err)
		return 1
	}
	v, err := d.DnsFilter.Get(ctx, id, g.vdom)
	if err != nil {
		output.Fail("api_error", err)
		return 2
	}
	output.Print(v)
	return 0
}

func runRouting(ctx context.Context, d Deps, args []string) int {
	if len(args) < 1 || args[0] != "status" {
		output.Fail("usage", fmt.Errorf("usage: fgcli routing status"))
		return 1
	}
	g, _, _ := parseGlobals(args[1:])
	st, err := d.Routing.Stats(ctx, g.vdom)
	if err != nil {
		output.Fail("api_error", err)
		return 2
	}
	output.Print(st)
	return 0
}
