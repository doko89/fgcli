package cli

import (
	"fmt"

	"github.com/local/fgcli/internal/infrastructure/output"
	"github.com/local/fgcli/internal/infrastructure/profile"
)

// RunProfile manages ~/.fgcli/config.yaml. It needs no auth.
func RunProfile(args []string) int {
	if len(args) < 1 {
		output.Fail("usage", fmt.Errorf("usage: fgcli profile list|show|use|add|rm"))
		return 1
	}
	verb, rest := args[0], args[1:]
	f, err := profile.LoadFile()
	if err != nil {
		output.Fail("config_error", err)
		return 1
	}
	switch verb {
	case "list":
		output.Print(map[string]any{"profiles": f.Summaries()})
		return 0
	case "validate":
		warns := profile.ValidateFile(f)
		if len(rest) >= 1 {
			keep := warns[:0]
			for _, w := range warns {
				if w.Profile == rest[0] || w.Profile == "" {
					keep = append(keep, w)
				}
			}
			warns = keep
		}
		valid := true
		for _, w := range warns {
			if w.Severity == "error" {
				valid = false
			}
		}
		output.Print(map[string]any{"valid": valid, "warnings": warns})
		return 0
	case "show":
		if len(rest) < 1 {
			output.Fail("usage", fmt.Errorf("usage: fgcli profile show <name>"))
			return 1
		}
		if _, err := f.Get(rest[0]); err != nil {
			output.Fail("not_found", err)
			return 1
		}
		for _, s := range f.Summaries() {
			if s.Name == rest[0] {
				output.Print(s)
				return 0
			}
		}
		return 0
	case "use", "set":
		if len(rest) < 1 {
			output.Fail("usage", fmt.Errorf("usage: fgcli profile use <name>"))
			return 1
		}
		if _, err := f.Get(rest[0]); err != nil {
			output.Fail("not_found", err)
			return 1
		}
		f.Active = rest[0]
		if err := f.Save(); err != nil {
			output.Fail("config_error", err)
			return 1
		}
		output.Print(map[string]string{"active": f.Active})
		return 0
	case "add":
		return profileAdd(f, rest)
	case "rm", "remove", "delete":
		if len(rest) < 1 {
			output.Fail("usage", fmt.Errorf("usage: fgcli profile rm <name> --yes"))
			return 1
		}
		g, rest, _ := parseGlobals(rest)
		if !g.yes {
			output.Fail("confirm_required", fmt.Errorf("refusing without --yes"))
			return 1
		}
		if _, err := f.Get(rest[0]); err != nil {
			output.Fail("not_found", err)
			return 1
		}
		delete(f.Profiles, rest[0])
		if f.Active == rest[0] {
			f.Active = ""
		}
		if err := f.Save(); err != nil {
			output.Fail("config_error", err)
			return 1
		}
		output.Print(map[string]string{"deleted": rest[0]})
		return 0
	default:
		output.Fail("usage", fmt.Errorf("unknown profile verb %q", verb))
		return 1
	}
}

func profileAdd(f *profile.File, args []string) int {
	_, rest, _ := parseGlobals(args)
	if len(rest) < 1 {
		output.Fail("usage", fmt.Errorf("usage: fgcli profile add <name> --host H --api-key K ..."))
		return 1
	}
	name := rest[0]
	key := flagVal(args, "api-key")
	if key == "" {
		key = flagVal(args, "token")
	}
	keyEnv := flagVal(args, "api-key-env")
	if keyEnv == "" {
		keyEnv = flagVal(args, "token-env")
	}
	pr := profile.Profile{
		Host:      flagVal(args, "host"),
		APIKey:    key,
		APIKeyEnv: keyEnv,
		Vdom:      flagVal(args, "vdom"),
		Insecure:  hasFlag(args, "--insecure"),
	}
	if pr.Host == "" {
		output.Fail("input_error", fmt.Errorf("--host is required"))
		return 1
	}
	if pr.APIKey == "" && pr.APIKeyEnv == "" {
		output.Fail("input_error", fmt.Errorf("auth required: --api-key/--api-key-env"))
		return 1
	}
	if _, exists := f.Profiles[name]; exists && !hasFlag(args, "--force") {
		output.Fail("exists", profile.ErrExists)
		return 1
	}
	if f.Profiles == nil {
		f.Profiles = map[string]profile.Profile{}
	}
	f.Profiles[name] = pr
	if err := f.Save(); err != nil {
		output.Fail("config_error", err)
		return 1
	}
	output.Print(map[string]string{"added": name})
	return 0
}

func hasFlag(args []string, name string) bool {
	for _, a := range args {
		if a == name {
			return true
		}
	}
	return false
}
