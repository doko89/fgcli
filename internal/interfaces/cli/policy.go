package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/local/fgcli/internal/domain/common"
	dpol "github.com/local/fgcli/internal/domain/policy"
	"github.com/local/fgcli/internal/infrastructure/output"
)

func names(s string) []common.Name {
	var out []common.Name
	for _, m := range strings.Split(s, ",") {
		if strings.TrimSpace(m) != "" {
			out = append(out, common.Name(strings.TrimSpace(m)))
		}
	}
	return out
}

func runPolicyWrite(ctx context.Context, d Deps, verb string, g globals, rest []string) int {
	var p dpol.Policy
	if g.fromSrc != "" {
		if err := readJSONInput(g.fromSrc, &p); err != nil {
			output.Fail("input_error", err)
			return 1
		}
	} else {
		id, _ := strconv.ParseInt(flagVal(rest, "policyid"), 10, 64)
		p = dpol.Policy{
			ID:       id,
			Name:     flagVal(rest, "name"),
			SrcIntf:  names(flagVal(rest, "srcintf")),
			DstIntf:  names(flagVal(rest, "dstintf")),
			SrcAddr:  names(flagVal(rest, "srcaddr")),
			DstAddr:  names(flagVal(rest, "dstaddr")),
			Service:  names(flagVal(rest, "service")),
			Action:   flagVal(rest, "action"),
			Status:   flagVal(rest, "status"),
			Schedule: flagVal(rest, "schedule"),
			Comments: flagVal(rest, "comments"),
			PoolName: names(flagVal(rest, "poolname")),
		}
		if len(rest) >= 1 && p.ID == 0 && verb == "update" {
			p.ID, _ = strconv.ParseInt(rest[0], 10, 64)
		}
	}
	if err := p.Validate(); err != nil {
		output.Fail("validation_error", err)
		return 1
	}
	if g.dryRun {
		output.Print(map[string]any{"dry_run": true, "policy": p, "vdom": g.vdom})
		return 0
	}
	var err error
	if verb == "create" {
		p, err = d.Pol.Create(ctx, p, g.vdom)
	} else {
		p, err = d.Pol.Update(ctx, p, g.vdom)
	}
	if err != nil {
		output.Fail("api_error", err)
		return 2
	}
	output.Print(p)
	return 0
}

func runPolicyStatus(ctx context.Context, d Deps, verb string, g globals, rest []string) int {
	if len(rest) < 1 {
		output.Fail("usage", fmt.Errorf("usage: fgcli policy %s <id>", verb))
		return 1
	}
	id, err := strconv.ParseInt(rest[0], 10, 64)
	if err != nil {
		output.Fail("input_error", err)
		return 1
	}
	if g.dryRun {
		output.Print(map[string]any{"dry_run": true, "status": verb, "policyid": id})
		return 0
	}
	p, err := d.Pol.SetStatus(ctx, id, verb, g.vdom)
	if err != nil {
		output.Fail("api_error", err)
		return 2
	}
	output.Print(p)
	return 0
}

func runPolicyMove(ctx context.Context, d Deps, g globals, rest, args []string) int {
	if len(rest) < 1 {
		output.Fail("usage", fmt.Errorf("usage: fgcli policy move <id> --before X | --after Y"))
		return 1
	}
	id, err := strconv.ParseInt(rest[0], 10, 64)
	if err != nil {
		output.Fail("input_error", err)
		return 1
	}
	before, _ := strconv.ParseInt(flagVal(args, "before"), 10, 64)
	after, _ := strconv.ParseInt(flagVal(args, "after"), 10, 64)
	if (before <= 0) == (after <= 0) {
		output.Fail("input_error", fmt.Errorf("exactly one of --before/--after is required"))
		return 1
	}
	if g.dryRun {
		output.Print(map[string]any{"dry_run": true, "move": id, "before": before, "after": after})
		return 0
	}
	if err := d.Pol.Move(ctx, id, before, after, g.vdom); err != nil {
		output.Fail("api_error", err)
		return 2
	}
	output.Print(map[string]any{"moved": id, "before": before, "after": after})
	return 0
}

func runPolicyClone(ctx context.Context, d Deps, g globals, rest, args []string) int {
	if len(rest) < 1 {
		output.Fail("usage", fmt.Errorf("usage: fgcli policy clone <src-id> --new-id N [--name M]"))
		return 1
	}
	src, err := strconv.ParseInt(rest[0], 10, 64)
	if err != nil {
		output.Fail("input_error", err)
		return 1
	}
	newID, _ := strconv.ParseInt(flagVal(args, "new-id"), 10, 64)
	if newID <= 0 {
		output.Fail("input_error", fmt.Errorf("--new-id is required"))
		return 1
	}
	name := flagVal(args, "name")
	if g.dryRun {
		output.Print(map[string]any{"dry_run": true, "clone": src, "new_id": newID, "name": name})
		return 0
	}
	p, err := d.Pol.Clone(ctx, src, newID, name, g.vdom)
	if err != nil {
		output.Fail("api_error", err)
		return 2
	}
	output.Print(p)
	return 0
}
