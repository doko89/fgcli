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
		return runGroupList(ctx, uc, g)
	case "get":
		return runGroupGet(ctx, uc, resource, g, rest)
	case "create", "update":
		return runGroupSave(ctx, uc, resource, verb, g, rest)
	case "delete":
		return runGroupDelete(ctx, uc, resource, g, rest)
	default:
		output.Fail("usage", fmt.Errorf("unknown %s verb %q", resource, verb))
		return 1
	}
}

func runGroupList(ctx context.Context, uc *gu.UseCase, g globals) int {
	list, err := uc.List(ctx, g.vdom)
	if err != nil {
		output.Fail("api_error", err)
		return 2
	}
	output.Print(filterSlice(list, g.filter))
	return 0
}

func runGroupGet(ctx context.Context, uc *gu.UseCase, resource string, g globals, rest []string) int {
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
}

func runGroupSave(ctx context.Context, uc *gu.UseCase, resource string, verb string, g globals, rest []string) int {
	gr, ok := loadGroupInput(g, rest, verb)
	if !ok {
		return 1
	}
	if err := gr.Validate(); err != nil {
		output.Fail("validation_error", err)
		return 1
	}
	if g.dryRun {
		output.Print(map[string]any{"dry_run": true, "group": gr, "vdom": g.vdom})
		return 0
	}
	return persistGroup(ctx, uc, verb, gr, g.vdom)
}

func loadGroupInput(g globals, rest []string, verb string) (dgrp.Group, bool) {
	if g.fromSrc != "" {
		var gr dgrp.Group
		if err := readJSONInput(g.fromSrc, &gr); err != nil {
			output.Fail("input_error", err)
			return dgrp.Group{}, false
		}
		return gr, true
	}
	return buildGroupFromFlags(rest, verb), true
}

func buildGroupFromFlags(rest []string, verb string) dgrp.Group {
	gr := dgrp.Group{Name: flagVal(rest, "name"), Comment: flagVal(rest, "comment")}
	appendGroupMembers(&gr, flagVal(rest, "member"))
	if len(rest) >= 1 && gr.Name == "" && verb == "update" {
		gr.Name = rest[0]
	}
	return gr
}

func appendGroupMembers(gr *dgrp.Group, members string) {
	for _, m := range strings.Split(members, ",") {
		trimmed := strings.TrimSpace(m)
		if trimmed == "" {
			continue
		}
		gr.Member = append(gr.Member, common.Name(trimmed))
	}
}

func persistGroup(ctx context.Context, uc *gu.UseCase, verb string, gr dgrp.Group, vdom string) int {
	var err error
	if verb == "create" {
		gr, err = uc.Create(ctx, gr, vdom)
	} else {
		gr, err = uc.Update(ctx, gr, vdom)
	}
	if err != nil {
		output.Fail("api_error", err)
		return 2
	}
	output.Print(gr)
	return 0
}

func runGroupDelete(ctx context.Context, uc *gu.UseCase, resource string, g globals, rest []string) int {
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
		return runVipList(ctx, d, g)
	case "get":
		return runVipGet(ctx, d, g, rest)
	case "create", "update":
		return runVipSave(ctx, d, verb, g, rest)
	case "delete":
		return runVipDelete(ctx, d, g, rest)
	default:
		output.Fail("usage", fmt.Errorf("unknown vip verb %q", verb))
		return 1
	}
}

func runVipList(ctx context.Context, d Deps, g globals) int {
	list, err := d.Vip.List(ctx, g.vdom)
	if err != nil {
		output.Fail("api_error", err)
		return 2
	}
	output.Print(filterSlice(list, g.filter))
	return 0
}

func runVipGet(ctx context.Context, d Deps, g globals, rest []string) int {
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
}

func runVipSave(ctx context.Context, d Deps, verb string, g globals, rest []string) int {
	v, ok := loadVipInput(g, rest, verb)
	if !ok {
		return 1
	}
	if err := v.Validate(); err != nil {
		output.Fail("validation_error", err)
		return 1
	}
	if g.dryRun {
		output.Print(map[string]any{"dry_run": true, "vip": v, "vdom": g.vdom})
		return 0
	}
	return persistVip(ctx, d, verb, v, g.vdom)
}

func loadVipInput(g globals, rest []string, verb string) (dvip.Vip, bool) {
	if g.fromSrc != "" {
		var v dvip.Vip
		if err := readJSONInput(g.fromSrc, &v); err != nil {
			output.Fail("input_error", err)
			return dvip.Vip{}, false
		}
		return v, true
	}
	return buildVipFromFlags(rest, verb), true
}

func buildVipFromFlags(rest []string, verb string) dvip.Vip {
	v := dvip.Vip{
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
	appendVipMappedIP(&v, flagVal(rest, "mappedip"))
	if len(rest) >= 1 && v.Name == "" && verb == "update" {
		v.Name = rest[0]
	}
	return v
}

func appendVipMappedIP(v *dvip.Vip, mappedip string) {
	for _, r := range strings.Split(mappedip, ",") {
		trimmed := strings.TrimSpace(r)
		if trimmed == "" {
			continue
		}
		v.MappedIP = append(v.MappedIP, dvip.MappedRange{Range: trimmed})
	}
}

func persistVip(ctx context.Context, d Deps, verb string, v dvip.Vip, vdom string) int {
	var err error
	if verb == "create" {
		v, err = d.Vip.Create(ctx, v, vdom)
	} else {
		v, err = d.Vip.Update(ctx, v, vdom)
	}
	if err != nil {
		output.Fail("api_error", err)
		return 2
	}
	output.Print(v)
	return 0
}

func runVipDelete(ctx context.Context, d Deps, g globals, rest []string) int {
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
		return runPoolList(ctx, d, g)
	case "get":
		return runPoolGet(ctx, d, g, rest)
	case "create", "update":
		return runPoolSave(ctx, d, verb, g, rest)
	case "delete":
		return runPoolDelete(ctx, d, g, rest)
	default:
		output.Fail("usage", fmt.Errorf("unknown ippool verb %q", verb))
		return 1
	}
}

func runPoolList(ctx context.Context, d Deps, g globals) int {
	list, err := d.Pool.List(ctx, g.vdom)
	if err != nil {
		output.Fail("api_error", err)
		return 2
	}
	output.Print(filterSlice(list, g.filter))
	return 0
}

func runPoolGet(ctx context.Context, d Deps, g globals, rest []string) int {
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
}

func runPoolSave(ctx context.Context, d Deps, verb string, g globals, rest []string) int {
	p, ok := loadPoolInput(g, rest, verb)
	if !ok {
		return 1
	}
	if err := p.Validate(); err != nil {
		output.Fail("validation_error", err)
		return 1
	}
	if g.dryRun {
		output.Print(map[string]any{"dry_run": true, "ippool": p, "vdom": g.vdom})
		return 0
	}
	return persistPool(ctx, d, verb, p, g.vdom)
}

func loadPoolInput(g globals, rest []string, verb string) (dip.Pool, bool) {
	if g.fromSrc != "" {
		var p dip.Pool
		if err := readJSONInput(g.fromSrc, &p); err != nil {
			output.Fail("input_error", err)
			return dip.Pool{}, false
		}
		return p, true
	}
	return buildPoolFromFlags(rest, verb), true
}

func buildPoolFromFlags(rest []string, verb string) dip.Pool {
	p := dip.Pool{
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
	return p
}

func persistPool(ctx context.Context, d Deps, verb string, p dip.Pool, vdom string) int {
	var err error
	if verb == "create" {
		p, err = d.Pool.Create(ctx, p, vdom)
	} else {
		p, err = d.Pool.Update(ctx, p, vdom)
	}
	if err != nil {
		output.Fail("api_error", err)
		return 2
	}
	output.Print(p)
	return 0
}

func runPoolDelete(ctx context.Context, d Deps, g globals, rest []string) int {
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
		return runSnatList(ctx, d, g)
	case "get":
		return runSnatGet(ctx, d, g, rest)
	case "create", "update":
		return runSnatSave(ctx, d, verb, g, rest)
	case "delete":
		return runSnatDelete(ctx, d, g, rest)
	default:
		output.Fail("usage", fmt.Errorf("unknown snat verb %q", verb))
		return 1
	}
}

func runSnatList(ctx context.Context, d Deps, g globals) int {
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
}

func runSnatGet(ctx context.Context, d Deps, g globals, rest []string) int {
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
}

func runSnatSave(ctx context.Context, d Deps, verb string, g globals, rest []string) int {
	m, ok := loadSnatInput(g, rest, verb)
	if !ok {
		return 1
	}
	if err := m.Validate(); err != nil {
		output.Fail("validation_error", err)
		return 1
	}
	if g.dryRun {
		output.Print(map[string]any{"dry_run": true, "snat": m, "vdom": g.vdom})
		return 0
	}
	return persistSnat(ctx, d, verb, m, g.vdom)
}

func loadSnatInput(g globals, rest []string, verb string) (dsn.Map, bool) {
	if g.fromSrc != "" {
		var m dsn.Map
		if err := readJSONInput(g.fromSrc, &m); err != nil {
			output.Fail("input_error", err)
			return dsn.Map{}, false
		}
		return m, true
	}
	return buildSnatFromFlags(rest, verb), true
}

func buildSnatFromFlags(rest []string, verb string) dsn.Map {
	id, _ := strconv.ParseInt(flagVal(rest, "policyid"), 10, 64)
	m := dsn.Map{
		ID:       id,
		Status:   flagVal(rest, "status"),
		Comments: flagVal(rest, "comments"),
		Protocol: flagVal(rest, "protocol"),
	}
	if len(rest) >= 1 && m.ID == 0 && verb == "update" {
		m.ID, _ = strconv.ParseInt(rest[0], 10, 64)
	}
	return m
}

func persistSnat(ctx context.Context, d Deps, verb string, m dsn.Map, vdom string) int {
	var err error
	if verb == "create" {
		m, err = d.Snat.Create(ctx, m, vdom)
	} else {
		m, err = d.Snat.Update(ctx, m, vdom)
	}
	if err != nil {
		output.Fail("api_error", err)
		return 2
	}
	output.Print(m)
	return 0
}

func runSnatDelete(ctx context.Context, d Deps, g globals, rest []string) int {
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
}
