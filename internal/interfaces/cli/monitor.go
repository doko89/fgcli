package cli

import (
	"context"
	"fmt"

	"github.com/local/fgcli/internal/infrastructure/output"
)

func runSwitch(ctx context.Context, d Deps, args []string) int {
	if len(args) < 1 || args[0] != "status" {
		output.Fail("usage", fmt.Errorf("usage: fgcli switch status [--filter SUB]"))
		return 1
	}
	g, _, _ := parseGlobals(args[1:])
	list, err := d.Monitor.Switches(ctx, g.vdom)
	if err != nil {
		output.Fail("api_error", err)
		return 2
	}
	output.Print(filterSlice(list, g.filter))
	return 0
}

func runWifi(ctx context.Context, d Deps, args []string) int {
	if len(args) < 1 || args[0] != "ap" {
		output.Fail("usage", fmt.Errorf("usage: fgcli wifi ap"))
		return 1
	}
	g, _, _ := parseGlobals(args[1:])
	st, err := d.Monitor.Wifi(ctx, g.vdom)
	if err != nil {
		output.Fail("api_error", err)
		return 2
	}
	output.Print(st)
	return 0
}

func runFortiView(ctx context.Context, d Deps, args []string) int {
	g, _, _ := parseGlobals(args)
	fv, err := d.Monitor.FortiView(ctx, g.vdom)
	if err != nil {
		output.Fail("api_error", err)
		return 2
	}
	output.Print(fv)
	return 0
}

func runLog(ctx context.Context, d Deps, args []string) int {
	if len(args) < 1 {
		output.Fail("usage", fmt.Errorf("usage: fgcli log <fortianalyzer|forticloud|memory|disk> <type> [--filter SUB]"))
		return 1
	}
	g, rest, _ := parseGlobals(args)
	if len(rest) < 2 {
		output.Fail("usage", fmt.Errorf("usage: fgcli log <fortianalyzer|forticloud|memory|disk> <type> [--filter SUB]"))
		return 1
	}
	list, err := d.Log.Fetch(ctx, rest[0], rest[1], g.vdom)
	if err != nil {
		output.Fail("api_error", err)
		return 2
	}
	output.Print(filterSlice(list, g.filter))
	return 0
}
