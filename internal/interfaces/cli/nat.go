package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	gu "github.com/local/fgcli/internal/application/group"
	"github.com/local/fgcli/internal/domain/common"
	dgrp "github.com/local/fgcli/internal/domain/group"
	dip "github.com/local/fgcli/internal/domain/ippool"
	dsn "github.com/local/fgcli/internal/domain/snat"
	dvip "github.com/local/fgcli/internal/domain/vip"
	"github.com/local/fgcli/internal/infrastructure/output"
)

func runGroup(ctx context.Context, uc *gu.UseCase, resource string, args []string) int {
	if len(args) < 1 {
		output.Fail("usage", fmt.Errorf("usage: fgcli %s list|get|create|update|delete", resource))
		return 1
	}
	verb, rest := args[0], args[1:]
	g, rest, _ := parseGlobals(rest)
	switch verb {
	case "list":
		list, err := uc.List(ctx, g.vdom)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(filterSlice(list, g.filter))
		return 0
	case "get":
		if len(rest) < 1 {
			output.Fail("usage", fmt.Errorf("usage: fgcli %s get <name>", resource))
			return 1
		}
		gr, err := uc.Get(ctx, rest[0], g.vdom)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(gr)
		return 0
	case "create", "update":
		var gr dgrp.Group
		if g.fromSrc != "" {
			if err := readJSONInput(g.fromSrc, &gr); err != nil {
				output.Fail("input_error", err)
				return 1
			}
		} else {
			gr = dgrp.Group{Name: flagVal(rest, "name"), Comment: flagVal(rest, "comment")}
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
			output.Print(map[string]any{"dry_run": true, "group": gr, "vdom": g.vdom})
			return 0
		}
		var err error
		if verb == "create" {
			gr, err = uc.Create(ctx, gr, g.vdom)
		} else {
			gr, err = uc.Update(ctx, gr, g.vdom)
		}
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(gr)
		return 0
	case "delete":
		if len(rest) < 1 {
			output.Fail("usage", fmt.Errorf("usage: fgcli %s delete <name> --yes", resource))
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
		if err := uc.Delete(ctx, rest[0], g.vdom); err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(map[string]string{"deleted": rest[0]})
		return 0
	default:
		output.Fail("usage", fmt.Errorf("unknown %s verb %q", resource, verb))
		return 1
	}
}

// matchFilter reports whether sub appears anywhere in the JSON rendering
// of v (case-insensitive). Empty sub matches everything.
func matchFilter(v any, sub string) bool {
	if sub == "" {
		return true
	}
	b, err := json.Marshal(v)
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(b)), strings.ToLower(sub))
}

func filterSlice[T any](items []T, sub string) []T {
	out := make([]T, 0, len(items))
	for _, it := range items {
		if matchFilter(it, sub) {
			out = append(out, it)
		}
	}
	return out
}

func runVip(ctx context.Context, d Deps, args []string) int {
	if len(args) < 1 {
		output.Fail("usage", fmt.Errorf("usage: fgcli vip list|get|create|update|delete"))
		return 1
	}
	verb, rest := args[0], args[1:]
	g, rest, _ := parseGlobals(rest)
	switch verb {
	case "list":
		list, err := d.Vip.List(ctx, g.vdom)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(filterSlice(list, g.filter))
		return 0
	case "get":
		if len(rest) < 1 {
			output.Fail("usage", fmt.Errorf("usage: fgcli vip get <name>"))
			return 1
		}
		v, err := d.Vip.Get(ctx, rest[0], g.vdom)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(v)
		return 0
	case "create", "update":
		var v dvip.Vip
		if g.fromSrc != "" {
			if err := readJSONInput(g.fromSrc, &v); err != nil {
				output.Fail("input_error", err)
				return 1
			}
		} else {
			v = dvip.Vip{
				Name:        flagVal(rest, "name"),
				Type:        flagVal(rest, "type"),
				ExtIP:       flagVal(rest, "extip"),
				ExtIntf:     flagVal(rest, "extintf"),
				Protocol:    flagVal(rest, "protocol"),
				ExtPort:     flagVal(rest, "extport"),
				MappedPort:  flagVal(rest, "mappedport"),
				PortForward: flagVal(rest, "portforward"),
				Comment:     flagVal(rest, "comment"),
				Status:      flagVal(rest, "status"),
			}
			for _, r := range strings.Split(flagVal(rest, "mappedip"), ",") {
				if strings.TrimSpace(r) != "" {
					v.MappedIP = append(v.MappedIP, dvip.MappedRange{Range: strings.TrimSpace(r)})
				}
			}
			if len(rest) >= 1 && v.Name == "" && verb == "update" {
				v.Name = rest[0]
			}
		}
		if err := v.Validate(); err != nil {
			output.Fail("validation_error", err)
			return 1
		}
		if g.dryRun {
			output.Print(map[string]any{"dry_run": true, "vip": v, "vdom": g.vdom})
			return 0
		}
		var err error
		if verb == "create" {
			v, err = d.Vip.Create(ctx, v, g.vdom)
		} else {
			v, err = d.Vip.Update(ctx, v, g.vdom)
		}
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(v)
		return 0
	case "delete":
		if len(rest) < 1 {
			output.Fail("usage", fmt.Errorf("usage: fgcli vip delete <name> --yes"))
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
		if err := d.Vip.Delete(ctx, rest[0], g.vdom); err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(map[string]string{"deleted": rest[0]})
		return 0
	default:
		output.Fail("usage", fmt.Errorf("unknown vip verb %q", verb))
		return 1
	}
}

func runPool(ctx context.Context, d Deps, args []string) int {
	if len(args) < 1 {
		output.Fail("usage", fmt.Errorf("usage: fgcli ippool list|get|create|update|delete"))
		return 1
	}
	verb, rest := args[0], args[1:]
	g, rest, _ := parseGlobals(rest)
	switch verb {
	case "list":
		list, err := d.Pool.List(ctx, g.vdom)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(filterSlice(list, g.filter))
		return 0
	case "get":
		if len(rest) < 1 {
			output.Fail("usage", fmt.Errorf("usage: fgcli ippool get <name>"))
			return 1
		}
		p, err := d.Pool.Get(ctx, rest[0], g.vdom)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(p)
		return 0
	case "create", "update":
		var p dip.Pool
		if g.fromSrc != "" {
			if err := readJSONInput(g.fromSrc, &p); err != nil {
				output.Fail("input_error", err)
				return 1
			}
		} else {
			p = dip.Pool{
				Name:      flagVal(rest, "name"),
				Type:      flagVal(rest, "type"),
				StartIP:   flagVal(rest, "startip"),
				EndIP:     flagVal(rest, "endip"),
				AssocIntf: flagVal(rest, "assoc-intf"),
				Comments:  flagVal(rest, "comments"),
			}
			if len(rest) >= 1 && p.Name == "" && verb == "update" {
				p.Name = rest[0]
			}
		}
		if err := p.Validate(); err != nil {
			output.Fail("validation_error", err)
			return 1
		}
		if g.dryRun {
			output.Print(map[string]any{"dry_run": true, "ippool": p, "vdom": g.vdom})
			return 0
		}
		var err error
		if verb == "create" {
			p, err = d.Pool.Create(ctx, p, g.vdom)
		} else {
			p, err = d.Pool.Update(ctx, p, g.vdom)
		}
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(p)
		return 0
	case "delete":
		if len(rest) < 1 {
			output.Fail("usage", fmt.Errorf("usage: fgcli ippool delete <name> --yes"))
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
		if err := d.Pool.Delete(ctx, rest[0], g.vdom); err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(map[string]string{"deleted": rest[0]})
		return 0
	default:
		output.Fail("usage", fmt.Errorf("unknown ippool verb %q", verb))
		return 1
	}
}

func runSnat(ctx context.Context, d Deps, args []string) int {
	if len(args) < 1 {
		output.Fail("usage", fmt.Errorf("usage: fgcli snat list|get|create|update|delete"))
		return 1
	}
	verb, rest := args[0], args[1:]
	g, rest, _ := parseGlobals(rest)
	switch verb {
	case "list":
		list, err := d.Snat.List(ctx, g.vdom)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		if list == nil {
			list = []dsn.Map{}
		}
		output.Print(filterSlice(list, g.filter))
		return 0
	case "get":
		if len(rest) < 1 {
			output.Fail("usage", fmt.Errorf("usage: fgcli snat get <policyid>"))
			return 1
		}
		id, err := strconv.ParseInt(rest[0], 10, 64)
		if err != nil {
			output.Fail("input_error", err)
			return 1
		}
		m, err := d.Snat.Get(ctx, id, g.vdom)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(m)
		return 0
	case "create", "update":
		var m dsn.Map
		if g.fromSrc != "" {
			if err := readJSONInput(g.fromSrc, &m); err != nil {
				output.Fail("input_error", err)
				return 1
			}
		} else {
			id, _ := strconv.ParseInt(flagVal(rest, "policyid"), 10, 64)
			m = dsn.Map{
				ID:       id,
				Status:   flagVal(rest, "status"),
				Comments: flagVal(rest, "comments"),
				Protocol: flagVal(rest, "protocol"),
			}
			if len(rest) >= 1 && m.ID == 0 && verb == "update" {
				m.ID, _ = strconv.ParseInt(rest[0], 10, 64)
			}
		}
		if err := m.Validate(); err != nil {
			output.Fail("validation_error", err)
			return 1
		}
		if g.dryRun {
			output.Print(map[string]any{"dry_run": true, "snat": m, "vdom": g.vdom})
			return 0
		}
		var err error
		if verb == "create" {
			m, err = d.Snat.Create(ctx, m, g.vdom)
		} else {
			m, err = d.Snat.Update(ctx, m, g.vdom)
		}
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(m)
		return 0
	case "delete":
		if len(rest) < 1 {
			output.Fail("usage", fmt.Errorf("usage: fgcli snat delete <policyid> --yes"))
			return 1
		}
		if !g.yes && !g.dryRun {
			output.Fail("confirm_required", fmt.Errorf("refusing without --yes"))
			return 1
		}
		id, err := strconv.ParseInt(rest[0], 10, 64)
		if err != nil {
			output.Fail("input_error", err)
			return 1
		}
		if g.dryRun {
			output.Print(map[string]any{"dry_run": true, "delete": id})
			return 0
		}
		if err := d.Snat.Delete(ctx, id, g.vdom); err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(map[string]int64{"deleted": id})
		return 0
	default:
		output.Fail("usage", fmt.Errorf("unknown snat verb %q", verb))
		return 1
	}
}
