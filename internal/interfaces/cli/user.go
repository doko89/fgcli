package cli

import (
	"context"
	"fmt"

	"github.com/local/fgcli/internal/infrastructure/output"
)

func runUser(ctx context.Context, d Deps, args []string) int {
	if len(args) < 1 {
		output.Fail("usage", fmt.Errorf("usage: fgcli user firewall|banned"))
		return 1
	}
	verb, rest := args[0], args[1:]
	g, _, _ := parseGlobals(rest)
	switch verb {
	case "firewall":
		list, err := d.User.Firewall(ctx, g.vdom)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(filterSlice(list, g.filter))
		return 0
	case "banned":
		list, err := d.User.Banned(ctx, g.vdom)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(filterSlice(list, g.filter))
		return 0
	default:
		output.Fail("usage", fmt.Errorf("unknown user verb %q", verb))
		return 1
	}
}
