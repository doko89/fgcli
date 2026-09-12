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
		return runSecurityIps(ctx, d, g)
	case "waf":
		return runSecurityWaf(ctx, d, g, rest)
	case "dlp":
		return runSecurityDlp(ctx, d, g, rest)
	case "proxy-pac":
		return runSecurityPac(ctx, d, args, g)
	default:
		output.Fail("usage", fmt.Errorf("unknown security verb %q", verb))
		return 1
	}
}

func runSecurityIps(ctx context.Context, d Deps, g globals) int {
	list, err := d.Security.Ips(ctx, g.vdom)
	if err != nil {
		output.Fail("api_error", err)
		return 2
	}
	output.Print(filterSlice(list, g.filter))
	return 0
}

func runSecurityWaf(ctx context.Context, d Deps, g globals, rest []string) int {
	if len(rest) >= 1 && rest[0] != "list" {
		return runSecurityGetWaf(ctx, d, g, rest)
	}
	return runSecurityListWaf(ctx, d, g)
}

func runSecurityGetWaf(ctx context.Context, d Deps, g globals, rest []string) int {
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

func runSecurityListWaf(ctx context.Context, d Deps, g globals) int {
	list, err := d.Security.ListWaf(ctx, g.vdom)
	if err != nil {
		output.Fail("api_error", err)
		return 2
	}
	output.Print(filterSlice(list, g.filter))
	return 0
}

func runSecurityDlp(ctx context.Context, d Deps, g globals, rest []string) int {
	if len(rest) >= 1 && rest[0] != "list" {
		return runSecurityGetDlp(ctx, d, g, rest)
	}
	return runSecurityListDlp(ctx, d, g)
}

func runSecurityGetDlp(ctx context.Context, d Deps, g globals, rest []string) int {
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

func runSecurityListDlp(ctx context.Context, d Deps, g globals) int {
	list, err := d.Security.ListDlp(ctx, g.vdom)
	if err != nil {
		output.Fail("api_error", err)
		return 2
	}
	output.Print(filterSlice(list, g.filter))
	return 0
}

func runSecurityPac(ctx context.Context, d Deps, args []string, g globals) int {
	out := flagVal(args, "out")
	if out == "" {
		out = "proxy.pac"
	}
	return downloadPac(ctx, d, g, out)
}

func downloadPac(ctx context.Context, d Deps, g globals, out string) int {
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
