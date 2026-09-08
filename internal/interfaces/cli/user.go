package cli

import (
	"context"
	"fmt"
	"strings"

	duser "github.com/local/fgcli/internal/domain/user"
	"github.com/local/fgcli/internal/infrastructure/output"
)

func runUser(ctx context.Context, d Deps, args []string) int {
	if len(args) < 1 {
		output.Fail("usage", fmt.Errorf("usage: fgcli user firewall|banned|group|local"))
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
	case "group":
		return runUserGroup(ctx, d, rest)
	case "local":
		return runUserLocal(ctx, d, rest)
	default:
		output.Fail("usage", fmt.Errorf("unknown user verb %q", verb))
		return 1
	}
}

func runUserGroup(ctx context.Context, d Deps, args []string) int {
	g, rest, _ := parseGlobals(args)
	if len(rest) < 1 {
		rest = []string{"list"}
	}
	verb, rest := rest[0], rest[1:]
	switch verb {
	case "list":
		list, err := d.User.GroupList(ctx, g.vdom)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(filterSlice(list, g.filter))
		return 0
	case "get":
		if len(rest) < 1 {
			output.Fail("usage", fmt.Errorf("usage: fgcli user group get <GROUP>"))
			return 1
		}
		gr, err := d.User.GroupGet(ctx, rest[0], g.vdom)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(gr)
		return 0
	case "add-member":
		if len(rest) < 2 {
			output.Fail("usage", fmt.Errorf("usage: fgcli user group add-member <GROUP> <USER> [--dry-run]"))
			return 1
		}
		gr, err := d.User.GroupGet(ctx, rest[0], g.vdom)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		members, added := duser.AddUnique(gr.Members, rest[1])
		gr.Members = members
		if g.dryRun {
			output.Print(map[string]any{"dry_run": true, "method": "PUT",
				"path": "cmdb/user/group/" + rest[0], "vdom": g.vdom,
				"member": gr.MemberNames(), "added": added})
			return 0
		}
		gr, err = d.User.GroupUpdate(ctx, gr, g.vdom)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(map[string]any{"group": gr.Name, "member": gr.MemberNames(), "added": added})
		return 0
	case "rm-member", "remove-member", "del-member":
		if len(rest) < 2 {
			output.Fail("usage", fmt.Errorf("usage: fgcli user group rm-member <GROUP> <USER> --yes [--dry-run]"))
			return 1
		}
		if !g.yes && !g.dryRun {
			output.Fail("confirm_required", fmt.Errorf("refusing without --yes (or use --dry-run)"))
			return 1
		}
		gr, err := d.User.GroupGet(ctx, rest[0], g.vdom)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		members, removed := duser.RemoveName(gr.Members, rest[1])
		gr.Members = members
		if g.dryRun {
			output.Print(map[string]any{"dry_run": true, "method": "PUT",
				"path": "cmdb/user/group/" + rest[0], "vdom": g.vdom,
				"member": gr.MemberNames(), "removed": removed})
			return 0
		}
		gr, err = d.User.GroupUpdate(ctx, gr, g.vdom)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(map[string]any{"group": gr.Name, "member": gr.MemberNames(), "removed": removed})
		return 0
	default:
		output.Fail("usage", fmt.Errorf("unknown user group verb %q", verb))
		return 1
	}
}

func runUserLocal(ctx context.Context, d Deps, args []string) int {
	g, rest, _ := parseGlobals(args)
	if len(rest) < 1 {
		rest = []string{"list"}
	}
	verb, rest := rest[0], rest[1:]
	switch verb {
	case "list":
		list, err := d.User.LocalList(ctx, g.vdom)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(filterSlice(list, g.filter))
		return 0
	case "get":
		if len(rest) < 1 {
			output.Fail("usage", fmt.Errorf("usage: fgcli user local get <NAME>"))
			return 1
		}
		u, err := d.User.LocalGet(ctx, rest[0], g.vdom)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(u)
		return 0
	case "create":
		var u duser.LocalUser
		if g.fromSrc != "" {
			if err := readJSONInput(g.fromSrc, &u); err != nil {
				output.Fail("input_error", err)
				return 1
			}
		} else {
			if len(rest) < 1 {
				output.Fail("usage", fmt.Errorf("usage: fgcli user local create <NAME> --password P [--group G,...] [--dry-run]"))
				return 1
			}
			u = duser.LocalUser{
				Name:   rest[0],
				Passwd: flagVal(args, "password"),
				Status: flagVal(args, "status"),
				Type:   flagVal(args, "type"),
			}
			if u.Passwd == "" {
				u.Passwd = flagVal(args, "passwd")
			}
		}
		if err := u.Validate(); err != nil {
			output.Fail("validation_error", err)
			return 1
		}
		if u.Passwd == "" {
			output.Fail("validation_error", fmt.Errorf("user local: --password is required"))
			return 1
		}
		groups := splitCSV(flagVal(args, "group"))
		if g.dryRun {
			output.Print(map[string]any{"dry_run": true, "method": "POST",
				"path": "cmdb/user/local", "vdom": g.vdom,
				"user":          map[string]string{"name": u.Name, "status": u.Status, "type": u.Type},
				"add_to_groups": groups})
			return 0
		}
		created, err := d.User.LocalCreate(ctx, u, g.vdom)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		joined := []string{}
		for _, grp := range groups {
			gr, err := d.User.GroupGet(ctx, grp, g.vdom)
			if err != nil {
				output.Fail("api_error", fmt.Errorf("user %q created but group %q lookup failed: %w", created.Name, grp, err))
				return 2
			}
			members, _ := duser.AddUnique(gr.Members, created.Name)
			gr.Members = members
			if _, err := d.User.GroupUpdate(ctx, gr, g.vdom); err != nil {
				output.Fail("api_error", fmt.Errorf("user %q created but adding to group %q failed: %w", created.Name, grp, err))
				return 2
			}
			joined = append(joined, grp)
		}
		output.Print(map[string]any{"user": created, "added_to_groups": joined})
		return 0
	case "delete":
		if len(rest) < 1 {
			output.Fail("usage", fmt.Errorf("usage: fgcli user local delete <NAME> --yes"))
			return 1
		}
		if !g.yes && !g.dryRun {
			output.Fail("confirm_required", fmt.Errorf("refusing without --yes (or use --dry-run)"))
			return 1
		}
		if g.dryRun {
			output.Print(map[string]any{"dry_run": true, "method": "DELETE",
				"path": "cmdb/user/local/" + rest[0], "vdom": g.vdom})
			return 0
		}
		if err := d.User.LocalDelete(ctx, rest[0], g.vdom); err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(map[string]string{"deleted": rest[0]})
		return 0
	default:
		output.Fail("usage", fmt.Errorf("unknown user local verb %q", verb))
		return 1
	}
}

func splitCSV(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
