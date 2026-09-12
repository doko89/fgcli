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
		return profileList(f)
	case "validate":
		return profileValidate(f, rest)
	case "show":
		return profileShow(f, rest)
	case "use", "set":
		return profileUse(f, rest)
	case "add":
		return profileAdd(f, rest)
	case "rm", "remove", "delete":
		return profileRemove(f, rest)
	default:
		output.Fail("usage", fmt.Errorf("unknown profile verb %q", verb))
		return 1
	}
}

func profileList(f *profile.File) int {
	output.Print(map[string]any{"profiles": f.Summaries()})
	return 0
}

func profileValidate(f *profile.File, rest []string) int {
	warns := profile.ValidateFile(f)
	warns = filterWarnings(warns, rest)
	output.Print(map[string]any{"valid": warningsValid(warns), "warnings": warns})
	return 0
}

func filterWarnings(warns []profile.Warning, rest []string) []profile.Warning {
	if len(rest) < 1 {
		return warns
	}
	keep := warns[:0]
	for _, w := range warns {
		if w.Profile == rest[0] || w.Profile == "" {
			keep = append(keep, w)
		}
	}
	return keep
}

func warningsValid(warns []profile.Warning) bool {
	for _, w := range warns {
		if w.Severity == "error" {
			return false
		}
	}
	return true
}

func profileShow(f *profile.File, rest []string) int {
	if len(rest) < 1 {
		output.Fail("usage", fmt.Errorf("usage: fgcli profile show <name>"))
		return 1
	}
	if _, err := f.Get(rest[0]); err != nil {
		output.Fail("not_found", err)
		return 1
	}
	return printProfileSummary(f, rest[0])
}

func printProfileSummary(f *profile.File, name string) int {
	for _, s := range f.Summaries() {
		if s.Name == name {
			output.Print(s)
			return 0
		}
	}
	return 0
}

func profileUse(f *profile.File, rest []string) int {
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
}

func profileRemove(f *profile.File, args []string) int {
	rest := args
	if len(rest) < 1 {
		output.Fail("usage", fmt.Errorf("usage: fgcli profile rm <name> --yes"))
		return 1
	}
	g, rest, _ := parseGlobals(rest)
	if !g.yes {
		output.Fail("confirm_required", fmt.Errorf("refusing without --yes"))
		return 1
	}
	return deleteProfile(f, rest[0])
}

func deleteProfile(f *profile.File, name string) int {
	if _, err := f.Get(name); err != nil {
		output.Fail("not_found", err)
		return 1
	}
	delete(f.Profiles, name)
	if f.Active == name {
		f.Active = ""
	}
	if err := f.Save(); err != nil {
		output.Fail("config_error", err)
		return 1
	}
	output.Print(map[string]string{"deleted": name})
	return 0
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
