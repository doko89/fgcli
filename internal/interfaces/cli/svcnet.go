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
		return runServiceVoip(ctx, d, g, rest)
	case "list":
		return runServiceList(ctx, d, g)
	case "get":
		return runServiceGet(ctx, d, g, rest)
	case "create", "update":
		return runServiceSave(ctx, d, verb, g, rest)
	case "delete":
		return runServiceDelete(ctx, d, g, rest)
	default:
		output.Fail("usage", fmt.Errorf("unknown service verb %q", verb))
		return 1
	}
}

func runServiceVoip(ctx context.Context, d Deps, g globals, rest []string) int {
	if len(rest) < 1 {
		output.Fail("usage", fmt.Errorf("usage: fgcli service voip get <name>"))
		return 1
	}
	name := resolveVoipName(rest)
	if name == "" {
		output.Fail("usage", fmt.Errorf("usage: fgcli service voip get <name>"))
		return 1
	}
	v, err := d.Voip.Get(ctx, name, g.vdom)
	if err != nil {
		output.Fail("api_error", err)
		return 2
	}
	output.Print(v)
	return 0
}

func resolveVoipName(rest []string) string {
	if rest[0] != "get" {
		return rest[0]
	}
	if len(rest) < 2 {
		return ""
	}
	return rest[1]
}

func runServiceList(ctx context.Context, d Deps, g globals) int {
	list, err := d.Svc.List(ctx, g.vdom)
	if err != nil {
		output.Fail("api_error", err)
		return 2
	}
	output.Print(filterSlice(list, g.filter))
	return 0
}

func runServiceGet(ctx context.Context, d Deps, g globals, rest []string) int {
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
}

func runServiceSave(ctx context.Context, d Deps, verb string, g globals, rest []string) int {
	s, ok := loadServiceInput(g, rest, verb)
	if !ok {
		return 1
	}
	if err := s.Validate(); err != nil {
		output.Fail("validation_error", err)
		return 1
	}
	if g.dryRun {
		output.Print(map[string]any{"dry_run": true, "service": s, "vdom": g.vdom})
		return 0
	}
	return persistService(ctx, d, verb, s, g.vdom)
}

func loadServiceInput(g globals, rest []string, verb string) (dsvc.Svc, bool) {
	if g.fromSrc != "" {
		var s dsvc.Svc
		if err := readJSONInput(g.fromSrc, &s); err != nil {
			output.Fail("input_error", err)
			return dsvc.Svc{}, false
		}
		return s, true
	}
	return buildServiceFromFlags(rest, verb), true
}

func buildServiceFromFlags(rest []string, verb string) dsvc.Svc {
	s := dsvc.Svc{
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
	return s
}

func persistService(ctx context.Context, d Deps, verb string, s dsvc.Svc, vdom string) int {
	var err error
	if verb == "create" {
		s, err = d.Svc.Create(ctx, s, vdom)
	} else {
		s, err = d.Svc.Update(ctx, s, vdom)
	}
	if err != nil {
		output.Fail("api_error", err)
		return 2
	}
	output.Print(s)
	return 0
}

func runServiceDelete(ctx context.Context, d Deps, g globals, rest []string) int {
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
		return runServiceGroupList(ctx, d, g)
	case "get":
		return runServiceGroupGet(ctx, d, g, rest)
	case "create", "update":
		return runServiceGroupSave(ctx, d, verb, g, rest)
	case "delete":
		return runServiceGroupDelete(ctx, d, g, rest)
	default:
		output.Fail("usage", fmt.Errorf("unknown service-group verb %q", verb))
		return 1
	}
}

func runServiceGroupList(ctx context.Context, d Deps, g globals) int {
	list, err := d.SvcGroup.List(ctx, g.vdom)
	if err != nil {
		output.Fail("api_error", err)
		return 2
	}
	output.Print(filterSlice(list, g.filter))
	return 0
}

func runServiceGroupGet(ctx context.Context, d Deps, g globals, rest []string) int {
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
}

func runServiceGroupSave(ctx context.Context, d Deps, verb string, g globals, rest []string) int {
	gr, ok := loadServiceGroupInput(g, rest, verb)
	if !ok {
		return 1
	}
	if err := gr.Validate(); err != nil {
		output.Fail("validation_error", err)
		return 1
	}
	if g.dryRun {
		output.Print(map[string]any{"dry_run": true, "service_group": gr, "vdom": g.vdom})
		return 0
	}
	return persistServiceGroup(ctx, d, verb, gr, g.vdom)
}

func loadServiceGroupInput(g globals, rest []string, verb string) (dsvc.Group, bool) {
	if g.fromSrc != "" {
		var gr dsvc.Group
		if err := readJSONInput(g.fromSrc, &gr); err != nil {
			output.Fail("input_error", err)
			return dsvc.Group{}, false
		}
		return gr, true
	}
	return buildServiceGroupFromFlags(rest, verb), true
}

func buildServiceGroupFromFlags(rest []string, verb string) dsvc.Group {
	gr := dsvc.Group{Name: flagVal(rest, "name"), Comment: flagVal(rest, "comment")}
	appendServiceGroupMembers(&gr, flagVal(rest, "member"))
	if len(rest) >= 1 && gr.Name == "" && verb == "update" {
		gr.Name = rest[0]
	}
	return gr
}

func appendServiceGroupMembers(gr *dsvc.Group, members string) {
	for _, m := range strings.Split(members, ",") {
		trimmed := strings.TrimSpace(m)
		if trimmed == "" {
			continue
		}
		gr.Member = append(gr.Member, common.Name(trimmed))
	}
}

func persistServiceGroup(ctx context.Context, d Deps, verb string, gr dsvc.Group, vdom string) int {
	var err error
	if verb == "create" {
		gr, err = d.SvcGroup.Create(ctx, gr, vdom)
	} else {
		gr, err = d.SvcGroup.Update(ctx, gr, vdom)
	}
	if err != nil {
		output.Fail("api_error", err)
		return 2
	}
	output.Print(gr)
	return 0
}

func runServiceGroupDelete(ctx context.Context, d Deps, g globals, rest []string) int {
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
