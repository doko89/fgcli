package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/local/fgcli/internal/application/trace"
	"github.com/local/fgcli/internal/infrastructure/output"
)

func runBackup(ctx context.Context, d Deps, args []string) int {
	if len(args) < 1 || args[0] != "download" {
		output.Fail("usage", fmt.Errorf("usage: fgcli backup download [--scope global] [--out PATH]"))
		return 1
	}
	scope := flagVal(args, "scope")
	if scope == "" {
		scope = "global"
	}
	out := flagVal(args, "out")
	if out == "" {
		out = fmt.Sprintf("fgcli-backup-%s.conf", time.Now().Format("20060102-150405"))
	}
	b, err := d.Backup.Download(ctx, scope)
	if err != nil {
		output.Fail("api_error", err)
		return 2
	}
	if err := os.WriteFile(out, b, 0o600); err != nil {
		output.Fail("input_error", err)
		return 1
	}
	output.Print(map[string]any{"path": out, "bytes": len(b), "scope": scope})
	return 0
}

func runRaw(ctx context.Context, d Deps, args []string) int {
	if len(args) < 1 {
		output.Fail("usage", fmt.Errorf("usage: fgcli raw get|post|put|delete <path> [--vdom V] [--data JSON|--from-stdin|--from-file]"))
		return 1
	}
	verb, rest := args[0], args[1:]
	g, rest, _ := parseGlobals(rest)
	if len(rest) < 1 {
		output.Fail("usage", fmt.Errorf("usage: fgcli raw %s <path>", verb))
		return 1
	}
	path, vdom := rest[0], ""
	if g.vdomSet {
		vdom = g.vdom
	}
	var body json.RawMessage
	if g.fromSrc != "" {
		var v any
		if err := readJSONInput(g.fromSrc, &v); err != nil {
			output.Fail("input_error", err)
			return 1
		}
		b, _ := json.Marshal(v)
		body = b
	} else if data := flagVal(args, "data"); data != "" {
		body = json.RawMessage(data)
	}
	if g.dryRun {
		output.Print(map[string]any{"dry_run": true, "method": verb, "path": path, "vdom": vdom})
		return 0
	}
	out, err := d.Raw.Do(ctx, verb, path, vdom, body)
	if err != nil {
		output.Fail("api_error", err)
		return 2
	}
	var v any
	if err := json.Unmarshal(out, &v); err != nil {
		output.Fail("api_error", err)
		return 2
	}
	output.Print(v)
	return 0
}

func runTrace(ctx context.Context, d Deps, args []string) int {
	g, _, _ := parseGlobals(args)
	in := trace.Input{
		Src: flagVal(args, "src"),
		Dst: flagVal(args, "dst"),
	}
	if in.Dst == "" {
		output.Fail("usage", fmt.Errorf("usage: fgcli trace --dst <ip> [--src <ip>] [--dport <n>]"))
		return 1
	}
	if p := flagVal(args, "dport"); p != "" {
		n, err := strconv.Atoi(p)
		if err != nil || n <= 0 || n > 65535 {
			output.Fail("input_error", fmt.Errorf("--dport must be 1-65535"))
			return 1
		}
		in.DPort = n
	}
	res, err := d.Tracer.Trace(ctx, g.vdom, in)
	if err != nil {
		output.Fail("api_error", err)
		return 2
	}
	output.Print(res)
	return 0
}
