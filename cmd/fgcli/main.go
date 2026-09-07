package main

import (
	"context"
	"os"
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

var version = "v0.2.0"

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(argv []string) int {
	// version/help/profile must work without auth (AI probe friendly)
	cmd := cli.TopCommand(argv)
	if len(argv) == 0 || cmd == "version" || cmd == "help" || cmd == "-h" || cmd == "--help" {
		if len(argv) == 0 {
			output.Fail("usage", errUsage())
			return 1
		}
		if cmd == "version" {
			output.Print(map[string]string{"version": version})
			return 0
		}
		if cmd == "help" {
			for i, a := range argv {
				if a == "help" && i+1 < len(argv) {
					if _, known := cli.HelpFor(argv[i+1]); known {
						cli.PrintHelp(argv[i+1])
						return 0
					}
				}
			}
		}
		output.Logf("fgcli %s — FortiGate CLI for humans and AI agents", version)
		output.Logf("Usage: FG_HOST=... FG_API_KEY=... fgcli <address|policy|system|doctor|profile|version> ...")
		output.Logf("   or: fgcli --profile fg1 ...  (see fgcli profile add)")
		output.Logf("Docs: README.md (contract: stdout = single JSON envelope)")
		return 0
	}
	if cmd == "profile" {
		for i, a := range argv {
			if a == "profile" {
				return cli.RunProfile(argv[i+1:])
			}
		}
	}
	// Per-resource help works without auth.
	for i, a := range argv {
		if a == cmd {
			if _, known := cli.HelpFor(cmd); known && cli.WantHelp(argv[i+1:]) {
				cli.PrintHelp(cmd)
				return 0
			}
			break
		}
	}
	// Precedence: flags > env > profile file.
	cfg := overrideFromArgs(config.LoadEnv(), argv)
	pf, err := profile.LoadFile()
	if err != nil {
		output.Fail("config_error", err)
		return 1
	}
	cfg, err = profile.Resolve(cfg, pf, profileName(argv))
	if err != nil {
		output.Fail("config_error", err)
		return 1
	}
	if err := cfg.Validate(); err != nil {
		output.Fail("config_error", err)
		return 1
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
		Sys:     fortios.NewSysRepo(client),
		Version: version,
	}
	// default vdom from env when flag absent
	return cli.Run(ctx, deps, argv)
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
