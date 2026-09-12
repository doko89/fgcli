package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	au "github.com/local/fgcli/internal/application/address"
	bu "github.com/local/fgcli/internal/application/backup"
	gru "github.com/local/fgcli/internal/application/group"
	fi "github.com/local/fgcli/internal/application/iface"
	iu "github.com/local/fgcli/internal/application/ippool"
	lu "github.com/local/fgcli/internal/application/log"
	mu "github.com/local/fgcli/internal/application/monitor"
	nu "github.com/local/fgcli/internal/application/network"
	pu "github.com/local/fgcli/internal/application/policy"
	ru "github.com/local/fgcli/internal/application/raw"
	secu "github.com/local/fgcli/internal/application/security"
	se "github.com/local/fgcli/internal/application/service"
	su "github.com/local/fgcli/internal/application/snat"
	tru "github.com/local/fgcli/internal/application/trace"
	uu "github.com/local/fgcli/internal/application/user"
	vu "github.com/local/fgcli/internal/application/vip"
	"github.com/local/fgcli/internal/infrastructure/config"
	"github.com/local/fgcli/internal/infrastructure/fortios"
	"github.com/local/fgcli/internal/infrastructure/output"
	"github.com/local/fgcli/internal/infrastructure/profile"
	"github.com/local/fgcli/internal/interfaces/cli"
)

var version = "v0.2.1"

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(argv []string) int {
	if len(argv) == 0 {
		output.Fail("usage", errUsage())
		return 1
	}
	applyPrettyFlag(argv)
	// version/help/profile must work without auth (AI probe friendly)
	if code, done := runHelpCommand(argv); done {
		return code
	}
	cmd := cli.TopCommand(argv)
	if code, done := runVersionCommand(cmd); done {
		return code
	}
	if code, done := runProfileCommand(cmd, argv); done {
		return code
	}
	// Precedence: flags > env > profile file.
	cfg, code, done := resolveConfig(argv)
	if done {
		return code
	}
	if cfg.Insecure {
		output.Logf("warn: TLS verification disabled (insecure:true / FG_INSECURE=1) — do not use in prod")
	}
	ctx := context.Background()
	client := fortios.NewClient(cfg.Host, cfg.APIKey, cfg.Insecure)
	deps := cli.Deps{
		Addr:      au.New(fortios.NewAddressRepo(client)),
		Pol:       pu.New(fortios.NewPolicyRepo(client)),
		Vip:       vu.New(fortios.NewVipRepo(client)),
		Pool:      iu.New(fortios.NewPoolRepo(client)),
		Snat:      su.New(fortios.NewSnatRepo(client)),
		Svc:       se.NewSvc(fortios.NewSvcRepo(client)),
		SvcGroup:  se.NewGroup(fortios.NewSvcGroupRepo(client)),
		Iface:     fi.NewIface(fortios.NewIfaceRepo(client)),
		Zone:      fi.NewZone(fortios.NewZoneRepo(client)),
		VipGrp:    gru.New(fortios.NewVipGrpRepo(client)),
		AddrGrp:   gru.New(fortios.NewAddrGrpRepo(client)),
		Vpn:       nu.NewVpn(fortios.NewVpnRepo(client)),
		Routing:   nu.NewRouting(fortios.NewRoutingRepo(client)),
		DnsFilter: nu.NewDnsFilter(fortios.NewDnsFilterRepo(client)),
		User:      uu.New(fortios.NewUserRepo(client)),
		Security:  secu.New(fortios.NewSecurityRepo(client)),
		Monitor:   mu.New(fortios.NewMonitorRepo(client)),
		Log:       lu.New(fortios.NewLogRepo(client)),
		Voip:      se.NewVoip(fortios.NewVoipRepo(client)),
		Raw:       ru.New(fortios.NewRawRepo(client)),
		Backup:    bu.New(fortios.NewBackupRepo(client)),
		Tracer: &tru.UseCase{
			Vips:  fortios.NewVipRepo(client),
			Pols:  fortios.NewPolicyRepo(client),
			Addrs: fortios.NewAddressRepo(client),
			Svcs:  fortios.NewSvcRepo(client),
		},
		Sys:      fortios.NewSysRepo(client),
		Version:  version,
		Insecure: cfg.Insecure,
	}
	// default vdom from env when flag absent
	return cli.Run(ctx, deps, argv)
}

func applyPrettyFlag(argv []string) {
	for _, a := range argv {
		if a == "--pretty" {
			output.Pretty = true
		}
	}
}

func runHelpCommand(argv []string) (int, bool) {
	target, ok := cli.HelpTarget(argv)
	if !ok {
		return 0, false
	}
	if target == "" {
		printGeneralHelp()
		return 0, true
	}
	cli.PrintHelp(target)
	return 0, true
}

func runVersionCommand(cmd string) (int, bool) {
	if cmd != "version" {
		return 0, false
	}
	output.Print(map[string]string{"version": version})
	return 0, true
}

func runProfileCommand(cmd string, argv []string) (int, bool) {
	if cmd != "profile" {
		return 0, false
	}
	warnLegacyConfig()
	for i, a := range argv {
		if a == "profile" {
			return cli.RunProfile(argv[i+1:]), true
		}
	}
	return 0, false
}

func resolveConfig(argv []string) (config.Config, int, bool) {
	cfg := overrideFromArgs(config.LoadEnv(), argv)
	pf, err := profile.LoadFile()
	if err != nil {
		output.Fail("config_error", err)
		return config.Config{}, 1, true
	}
	cfg, err = profile.Resolve(cfg, pf, profileName(argv))
	if err != nil {
		output.Fail("config_error", err)
		return config.Config{}, 1, true
	}
	if err := cfg.Validate(); err != nil {
		output.Fail("config_error", err)
		return config.Config{}, 1, true
	}
	if profile.IsPlaceholder(cfg.APIKey) {
		output.Fail("config_error", errPlaceholderKey())
		return config.Config{}, 1, true
	}
	return cfg, 0, false
}

func printGeneralHelp() {
	lines := []string{
		"fgcli " + version + " — FortiGate CLI for humans and AI agents",
		"Usage: FG_HOST=... FG_API_KEY=... fgcli <address|policy|system|doctor|profile|version> ...",
		"   or: fgcli --profile fg1 ...  (see fgcli profile add)",
		"   or: fgcli help <resource> | fgcli <resource> --help  (also: --help, -h)",
		"Config: single file ~/.fgcli/config.yaml (override: FGCLI_CONFIG)",
		"Docs: README.md (contract: stdout = single JSON envelope)",
	}
	for _, l := range lines {
		output.Logf("%s", l)
	}
	output.Print(map[string]string{"help": strings.Join(lines, "\n")})
}

// warnLegacyConfig points at files from older setups that fgcli ignores.
func warnLegacyConfig() {
	home, err := os.UserHomeDir()
	if err != nil || os.Getenv("FGCLI_CONFIG") != "" {
		return
	}
	for _, p := range []string{
		filepath.Join(home, "fgcli", "config.yaml"),
		filepath.Join(home, "tmp", "FG-CLI", "config.ini"),
	} {
		if _, err := os.Stat(p); err == nil {
			output.Logf("hint: legacy config %s is ignored; only ~/.fgcli/config.yaml is used", p)
		}
	}
}

func errPlaceholderKey() error {
	return usageErr{"auth looks like a placeholder api_key (run: fgcli profile validate)"}
}

func overrideFromArgs(c config.Config, argv []string) config.Config {
	for i := 0; i < len(argv); i++ {
		switch {
		case argv[i] == "--host" && i+1 < len(argv):
			c.Host = argv[i+1]
		case strings.HasPrefix(argv[i], "--host="):
			c.Host = strings.TrimPrefix(argv[i], "--host=")
		case argv[i] == "--api-key" && i+1 < len(argv):
			c.APIKey = argv[i+1]
		case strings.HasPrefix(argv[i], "--api-key="):
			c.APIKey = strings.TrimPrefix(argv[i], "--api-key=")
		case argv[i] == "--token" && i+1 < len(argv):
			c.APIKey = argv[i+1]
		case strings.HasPrefix(argv[i], "--token="):
			c.APIKey = strings.TrimPrefix(argv[i], "--token=")
		case argv[i] == "--insecure":
			c.Insecure = true
		}
	}
	return c
}

func profileName(argv []string) string {
	for i := 0; i < len(argv); i++ {
		if argv[i] == "--profile" && i+1 < len(argv) {
			return argv[i+1]
		}
		if strings.HasPrefix(argv[i], "--profile=") {
			return strings.TrimPrefix(argv[i], "--profile=")
		}
	}
	return ""
}

type usageErr struct{ s string }

func (e usageErr) Error() string { return e.s }
func errUsage() error {
	return usageErr{"usage: fgcli <address|policy|system|doctor|profile|version> ..."}
}
