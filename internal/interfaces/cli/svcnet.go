package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/local/fgcli/internal/domain/common"
	dsvc "github.com/local/fgcli/internal/domain/service"
	"github.com/local/fgcli/internal/infrastructure/output"
)

func runService(ctx context.Context, d Deps, args []string) int {
	if len(args) < 1 {
		output.Fail("usage", fmt.Errorf("usage: fgcli service list|get|create|update|delete"))
		return 1
	}
	verb, rest := args[0], args[1:]
	g, rest, _ := parseGlobals(rest)
	switch verb {
	case "voip":
		if len(rest) < 1 {
			output.Fail("usage", fmt.Errorf("usage: fgcli service voip get <name>"))
			return 1
		}
		name := rest[0]
		if name == "get" {
			if len(rest) < 2 {
				output.Fail("usage", fmt.Errorf("usage: fgcli service voip get <name>"))
				return 1
			}
			name = rest[1]
		}
		v, err := d.Voip.Get(ctx, name, g.vdom)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(v)
		return 0
	case "list":
		list, err := d.Svc.List(ctx, g.vdom)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(filterSlice(list, g.filter))
		return 0
	case "get":
		if len(rest) < 1 {
			output.Fail("usage", fmt.Errorf("usage: fgcli service get <name>"))
			return 1
		}
		s, err := d.Svc.Get(ctx, rest[0], g.vdom)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(s)
		return 0
	case "create", "update":
		var s dsvc.Svc
		if g.fromSrc != "" {
			if err := readJSONInput(g.fromSrc, &s); err != nil {
				output.Fail("input_error", err)
				return 1
			}
		} else {
			s = dsvc.Svc{
				Name:      flagVal(rest, "name"),
				Protocol:  flagVal(rest, "protocol"),
				TCPPorts:  flagVal(rest, "tcp"),
				UDPPorts:  flagVal(rest, "udp"),
				SCTPPorts: flagVal(rest, "sctp"),
				Comment:   flagVal(rest, "comment"),
			}
			if len(rest) >= 1 && s.Name == "" && verb == "update" {
				s.Name = rest[0]
			}
		}
		if err := s.Validate(); err != nil {
			output.Fail("validation_error", err)
			return 1
		}
		if g.dryRun {
			output.Print(map[string]any{"dry_run": true, "service": s, "vdom": g.vdom})
			return 0
		}
		var err error
		if verb == "create" {
			s, err = d.Svc.Create(ctx, s, g.vdom)
		} else {
			s, err = d.Svc.Update(ctx, s, g.vdom)
		}
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(s)
		return 0
	case "delete":
		if len(rest) < 1 {
			output.Fail("usage", fmt.Errorf("usage: fgcli service delete <name> --yes"))
			return 1
		}
		if !g.yes && !g.dryRun {
			output.Fail("confirm_required", fmt.Errorf("refusing without --yes (or use --dry-run)"))
			return 1
		}
		if g.dryRun {
			output.Print(map[string]any{"dry_run": true, "delete": rest[0]})
			return 0
		}
		if err := d.Svc.Delete(ctx, rest[0], g.vdom); err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(map[string]string{"deleted": rest[0]})
		return 0
	default:
		output.Fail("usage", fmt.Errorf("unknown service verb %q", verb))
		return 1
	}
}

func runServiceGroup(ctx context.Context, d Deps, args []string) int {
	if len(args) < 1 {
		output.Fail("usage", fmt.Errorf("usage: fgcli service-group list|get|create|update|delete"))
		return 1
	}
	verb, rest := args[0], args[1:]
	g, rest, _ := parseGlobals(rest)
	switch verb {
	case "list":
		list, err := d.SvcGroup.List(ctx, g.vdom)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(filterSlice(list, g.filter))
		return 0
	case "get":
		if len(rest) < 1 {
			output.Fail("usage", fmt.Errorf("usage: fgcli service-group get <name>"))
			return 1
		}
		gr, err := d.SvcGroup.Get(ctx, rest[0], g.vdom)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(gr)
		return 0
	case "create", "update":
		var gr dsvc.Group
		if g.fromSrc != "" {
			if err := readJSONInput(g.fromSrc, &gr); err != nil {
				output.Fail("input_error", err)
				return 1
			}
		} else {
			gr = dsvc.Group{Name: flagVal(rest, "name"), Comment: flagVal(rest, "comment")}
			for _, m := range strings.Split(flagVal(rest, "member"), ",") {
				if strings.TrimSpace(m) != "" {
					gr.Member = append(gr.Member, common.Name(strings.TrimSpace(m)))
				}
			}
			if len(rest) >= 1 && gr.Name == "" && verb == "update" {
				gr.Name = rest[0]
			}
		}
		if err := gr.Validate(); err != nil {
			output.Fail("validation_error", err)
			return 1
		}
		if g.dryRun {
			output.Print(map[string]any{"dry_run": true, "service_group": gr, "vdom": g.vdom})
			return 0
		}
		var err error
		if verb == "create" {
			gr, err = d.SvcGroup.Create(ctx, gr, g.vdom)
		} else {
			gr, err = d.SvcGroup.Update(ctx, gr, g.vdom)
		}
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(gr)
		return 0
	case "delete":
		if len(rest) < 1 {
			output.Fail("usage", fmt.Errorf("usage: fgcli service-group delete <name> --yes"))
			return 1
		}
		if !g.yes && !g.dryRun {
			output.Fail("confirm_required", fmt.Errorf("refusing without --yes (or use --dry-run)"))
			return 1
		}
		if g.dryRun {
			output.Print(map[string]any{"dry_run": true, "delete": rest[0]})
			return 0
		}
		if err := d.SvcGroup.Delete(ctx, rest[0], g.vdom); err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(map[string]string{"deleted": rest[0]})
		return 0
	default:
		output.Fail("usage", fmt.Errorf("unknown service-group verb %q", verb))
		return 1
	}
}

func runIface(ctx context.Context, d Deps, args []string) int {
	if len(args) < 1 {
		output.Fail("usage", fmt.Errorf("usage: fgcli interface list|get"))
		return 1
	}
	verb, rest := args[0], args[1:]
	g, rest, _ := parseGlobals(rest)
	switch verb {
	case "list":
		list, err := d.Iface.List(ctx, g.vdom)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(filterSlice(list, g.filter))
		return 0
	case "get":
		if len(rest) < 1 {
			output.Fail("usage", fmt.Errorf("usage: fgcli interface get <name>"))
			return 1
		}
		i, err := d.Iface.Get(ctx, rest[0], g.vdom)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(i)
		return 0
	default:
		output.Fail("usage", fmt.Errorf("unknown interface verb %q (read-only: list|get)", verb))
		return 1
	}
}

func runZone(ctx context.Context, d Deps, args []string) int {
	if len(args) < 1 {
		output.Fail("usage", fmt.Errorf("usage: fgcli zone list|get"))
		return 1
	}
	verb, rest := args[0], args[1:]
	g, rest, _ := parseGlobals(rest)
	switch verb {
	case "list":
		list, err := d.Zone.List(ctx, g.vdom)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(filterSlice(list, g.filter))
		return 0
	case "get":
		if len(rest) < 1 {
			output.Fail("usage", fmt.Errorf("usage: fgcli zone get <name>"))
			return 1
		}
		z, err := d.Zone.Get(ctx, rest[0], g.vdom)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(z)
		return 0
	default:
		output.Fail("usage", fmt.Errorf("unknown zone verb %q (read-only: list|get)", verb))
		return 1
	}
}
