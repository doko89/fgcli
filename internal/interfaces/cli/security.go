package cli

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/local/fgcli/internal/infrastructure/output"
)

func runSecurity(ctx context.Context, d Deps, args []string) int {
	if len(args) < 1 {
		output.Fail("usage", fmt.Errorf("usage: fgcli security ips|waf|dlp|proxy-pac"))
		return 1
	}
	verb, rest := args[0], args[1:]
	g, rest, _ := parseGlobals(rest)
	switch verb {
	case "ips":
		list, err := d.Security.Ips(ctx, g.vdom)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(filterSlice(list, g.filter))
		return 0
	case "waf":
		if len(rest) >= 1 && rest[0] != "list" {
			id, ok := secID(rest)
			if !ok {
				return 1
			}
			v, err := d.Security.GetWaf(ctx, id, g.vdom)
			if err != nil {
				output.Fail("api_error", err)
				return 2
			}
			output.Print(v)
			return 0
		}
		list, err := d.Security.ListWaf(ctx, g.vdom)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(filterSlice(list, g.filter))
		return 0
	case "dlp":
		if len(rest) >= 1 && rest[0] != "list" {
			id, ok := secID(rest)
			if !ok {
				return 1
			}
			v, err := d.Security.GetDlp(ctx, id, g.vdom)
			if err != nil {
				output.Fail("api_error", err)
				return 2
			}
			output.Print(v)
			return 0
		}
		list, err := d.Security.ListDlp(ctx, g.vdom)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(filterSlice(list, g.filter))
		return 0
	case "proxy-pac":
		out := flagVal(args, "out")
		if out == "" {
			out = "proxy.pac"
		}
		b, err := d.Security.DownloadPac(ctx, g.vdom)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		if err := os.WriteFile(out, b, 0o600); err != nil {
			output.Fail("input_error", err)
			return 1
		}
		output.Print(map[string]any{"path": out, "bytes": len(b)})
		return 0
	default:
		output.Fail("usage", fmt.Errorf("unknown security verb %q", verb))
		return 1
	}
}

func secID(rest []string) (int64, bool) {
	arg := rest[0]
	if arg == "get" {
		if len(rest) < 2 {
			output.Fail("usage", fmt.Errorf("usage: fgcli security <waf|dlp> get <id>"))
			return 0, false
		}
		arg = rest[1]
	}
	id, err := strconv.ParseInt(arg, 10, 64)
	if err != nil {
		output.Fail("input_error", err)
		return 0, false
	}
	return id, true
}
