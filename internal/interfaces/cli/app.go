package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
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
	daddr "github.com/local/fgcli/internal/domain/address"
	dsys "github.com/local/fgcli/internal/domain/svc"
	"github.com/local/fgcli/internal/infrastructure/output"
)

type Deps struct {
	Addr      *au.UseCase
	Pol       *pu.UseCase
	Vip       *vu.UseCase
	Pool      *iu.UseCase
	Snat      *su.UseCase
	Svc       *se.SvcUseCase
	SvcGroup  *se.GroupUseCase
	VipGrp    *gru.UseCase
	AddrGrp   *gru.UseCase
	Iface     *fi.IfaceUseCase
	Zone      *fi.ZoneUseCase
	Vpn       *nu.VpnUseCase
	Routing   *nu.RoutingUseCase
	DnsFilter *nu.DnsFilterUseCase
	Monitor   *mu.UseCase
	Log       *lu.UseCase
	User      *uu.UseCase
	Security  *secu.UseCase
	Voip      *se.VoipUseCase
	Raw       *ru.UseCase
	Backup    *bu.UseCase
	Tracer    *tru.UseCase
	Sys       dsys.Repository
	Version   string
	Insecure  bool
}

type globals struct {
	vdom    string
	vdomSet bool
	dryRun  bool
	fromSrc string
	yes     bool
	filter  string
}

// valueFlags are flags that consume the next arg (their values must not
// be mistaken for commands or --help requests).
var valueFlags = map[string]bool{
	"--profile": true, "--host": true, "--api-key": true, "--token": true,
	"--api-key-env": true, "--token-env": true,
	"--vdom": true, "--from-file": true, "--data": true, "--out": true,
	"--filter": true, "--contains": true,
	"--name": true, "--subnet": true, "--comment": true, "--type": true,
	"--protocol": true, "--tcp": true, "--udp": true, "--sctp": true,
	"--member": true, "--src": true, "--dst": true, "--dport": true,
	"--new-id": true, "--before": true, "--after": true, "--policyid": true,
	"--srcintf": true, "--dstintf": true, "--srcaddr": true, "--dstaddr": true,
	"--service": true, "--action": true, "--extip": true, "--mappedip": true,
	"--extintf": true, "--extport": true, "--mappedport": true,
	"--startip": true, "--endip": true, "--scope": true,
	"--password": true, "--passwd": true, "--group": true, "--status": true,
}

// HelpTarget detects an explicit help request anywhere in argv (help, -h,
// --help), skipping flag values. It returns the resource the help is about,
// or "" for general help. ok=false means no help was requested.
func HelpTarget(argv []string) (target string, ok bool) {
	for i := 0; i < len(argv); i++ {
		a := argv[i]
		if strings.HasPrefix(a, "--") && !strings.Contains(a, "=") && valueFlags[a] {
			i++
			continue
		}
		if len(a) > 0 && a[0] == '-' && a != "-h" && a != "--help" {
			continue
		}
		switch a {
		case "help":
			if i+1 < len(argv) {
				if _, known := HelpFor(argv[i+1]); known {
					return argv[i+1], true
				}
			}
			return "", true
		case "-h", "--help":
			return target, true
		}
		if _, known := HelpFor(a); known {
			target = a
		}
	}
	return "", false
}

// TopCommand returns the resource verb (address, profile, ...) skipping
// global flags, so `fgcli --profile fg1 address list` routes correctly.
func TopCommand(argv []string) string {
	for i := 0; i < len(argv); i++ {
		a := argv[i]
		if len(a) > 2 && a[:2] == "--" {
			if !strings.Contains(a, "=") && valueFlags[a] {
				i++
			}
			continue
		}
		if len(a) > 0 && a[0] == '-' {
			continue
		}
		return a
	}
	return ""
}

func Run(ctx context.Context, d Deps, argv []string) int {
	if len(argv) < 1 {
		usage()
		return 1
	}
	if target, ok := HelpTarget(argv); ok {
		if target == "" {
			usage()
			return 0
		}
		PrintHelp(target)
		return 0
	}
	cmd := TopCommand(argv)
	rest := argv
	for i, a := range argv {
		if a == cmd {
			rest = argv[i+1:]
			break
		}
	}
	switch cmd {
	case "address":
		return runAddress(ctx, d, rest)
	case "policy":
		return runPolicy(ctx, d, rest)
	case "vip":
		return runVip(ctx, d, rest)
	case "ippool":
		return runPool(ctx, d, rest)
	case "snat", "central-snat-map":
		return runSnat(ctx, d, rest)
	case "vipgrp":
		return runGroup(ctx, d.VipGrp, "vipgrp", rest)
	case "addrgrp":
		return runGroup(ctx, d.AddrGrp, "addrgrp", rest)
	case "service":
		return runService(ctx, d, rest)
	case "service-group":
		return runServiceGroup(ctx, d, rest)
	case "interface":
		return runIface(ctx, d, rest)
	case "zone":
		return runZone(ctx, d, rest)
	case "user":
		return runUser(ctx, d, rest)
	case "security":
		return runSecurity(ctx, d, rest)
	case "network":
		return runNetwork(ctx, d, rest)
	case "switch":
		return runSwitch(ctx, d, rest)
	case "wifi":
		return runWifi(ctx, d, rest)
	case "fortiview":
		return runFortiView(ctx, d, rest)
	case "log":
		return runLog(ctx, d, rest)
	case "raw":
		return runRaw(ctx, d, rest)
	case "trace":
		return runTrace(ctx, d, rest)
	case "backup":
		return runBackup(ctx, d, rest)
	case "system":
		return runSystem(ctx, d, rest)
	case "vpn":
		return runVpn(ctx, d, rest)
	case "routing":
		return runRouting(ctx, d, rest)
	case "doctor":
		st, err := d.Sys.Status(ctx)
		if err != nil {
			output.Fail("connect_error", err)
			return 2
		}
		output.Print(map[string]any{"reachable": true, "insecure_skip_verify": d.Insecure, "status": st})
		return 0
	case "version":
		output.Print(map[string]string{"version": d.Version})
		return 0
	case "help", "-h", "--help":
		if len(rest) >= 1 {
			if _, ok := HelpFor(rest[0]); ok {
				PrintHelp(rest[0])
				return 0
			}
		}
		usage()
		return 0
	default:
		output.Fail("usage", fmt.Errorf("unknown command %q", cmd))
		return 1
	}
}

func usage() {
	output.Logf("fgcli [--profile NAME] <address|policy|system|doctor|profile|version> ...")
	output.Logf("  address list|get|create|update|delete [--filter SUB]")
	output.Logf("  policy  list|get|delete [--filter SUB]")
	output.Logf("  vip     list|get|create|update|delete [--filter SUB]")
	output.Logf("  ippool  list|get|create|update|delete [--filter SUB]")
	output.Logf("  snat    list|get|create|update|delete [--filter SUB]")
	output.Logf("  service list|get|create|update|delete [--filter SUB]")
	output.Logf("  service-group list|get|create|update|delete [--filter SUB]")
	output.Logf("  interface list|get [--filter SUB]  (read-only)")
	output.Logf("  zone      list|get [--filter SUB]  (read-only)")
	output.Logf("  system status|license|fortiguard|ntp|dns|dhcp|snmp")
	output.Logf("  vpn ipsec [--filter SUB]  |  vpn ssl get  |  routing status")
	output.Logf("  user firewall|banned  |  user group add-member G U  |  user local create N --password P")
	output.Logf("  security ips|waf|dlp|proxy-pac")
	output.Logf("  network dnsfilter list|get  |  service voip get <name>")
	output.Logf("  switch status  |  wifi ap  |  fortiview")
	output.Logf("  log <fortianalyzer|forticloud|memory|disk> <type> [--filter SUB]")
	output.Logf("  raw get|post|put|delete <path>  (escape hatch)")
	output.Logf("  trace --dst IP [--src IP] [--dport N]")
	output.Logf("  system  status")
	output.Logf("  profile list|show|use|add|rm")
	output.Logf("Global env: FG_HOST FG_API_KEY FG_VDOM FG_INSECURE=1")
	output.Logf("All stdout is JSON {ok,data,error}; logs go to stderr.")
}

func parseGlobals(args []string) (globals, []string, map[string]string) {
	g := globals{vdom: os.Getenv("FG_VDOM")}
	var rest []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--vdom" && i+1 < len(args):
			i++
			g.vdom = args[i]
			g.vdomSet = true
		case len(a) > 7 && a[:7] == "--vdom=":
			g.vdom = a[7:]
			g.vdomSet = true
		case a == "--dry-run":
			g.dryRun = true
		case a == "--pretty" || a == "--json":
			output.Pretty = a == "--pretty"
		case a == "--from-stdin":
			g.fromSrc = "-"
		case a == "--from-file" && i+1 < len(args):
			i++
			g.fromSrc = args[i]
		case len(a) > 12 && a[:12] == "--from-file=":
			g.fromSrc = a[12:]
		case a == "--yes":
			g.yes = true
		case (a == "--filter" || a == "--contains") && i+1 < len(args):
			i++
			g.filter = args[i]
		case len(a) > 9 && a[:9] == "--filter=":
			g.filter = a[9:]
		case len(a) > 11 && a[:11] == "--contains=":
			g.filter = a[11:]
		case a == "--profile" && i+1 < len(args):
			i++ // consumed in main; stripped here so positionals stay clean
		case len(a) > 10 && a[:10] == "--profile=":
		case a == "--host" && i+1 < len(args):
			i++
		case len(a) > 7 && a[:7] == "--host=":
		case a == "--api-key" && i+1 < len(args):
			i++
		case len(a) > 10 && a[:10] == "--api-key=":
		case a == "--token" && i+1 < len(args):
			i++
		case len(a) > 8 && a[:8] == "--token=":
		case a == "--api-key-env" && i+1 < len(args):
			i++
		case len(a) > 14 && a[:14] == "--api-key-env=":
		case a == "--token-env" && i+1 < len(args):
			i++
		case len(a) > 12 && a[:12] == "--token-env=":
		case a == "--insecure":
		default:
			rest = append(rest, a)
		}
	}
	if g.vdom == "" {
		g.vdom = "root"
	}
	return g, rest, map[string]string{}
}

func readJSONInput(fromSrc string, v any) error {
	var r io.Reader = os.Stdin
	if fromSrc != "" && fromSrc != "-" {
		f, err := os.Open(fromSrc)
		if err != nil {
			return err
		}
		defer f.Close()
		r = f
	}
	return json.NewDecoder(r).Decode(v)
}

func flagVal(args []string, name string) string {
	for i := 0; i < len(args); i++ {
		if args[i] == "--"+name && i+1 < len(args) {
			return args[i+1]
		}
		if len(args[i]) > len(name)+2 && args[i][:len(name)+3] == "--"+name+"=" {
			return args[i][len(name)+3:]
		}
	}
	return ""
}

func runAddress(ctx context.Context, d Deps, args []string) int {
	if len(args) < 1 {
		output.Fail("usage", fmt.Errorf("usage: fgcli address list|get|create|update|delete"))
		return 1
	}
	verb, rest := args[0], args[1:]
	g, rest, _ := parseGlobals(rest)
	switch verb {
	case "list":
		list, err := d.Addr.List(ctx, g.vdom)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(filterSlice(list, g.filter))
		return 0
	case "get":
		if len(rest) < 1 {
			output.Fail("usage", fmt.Errorf("usage: fgcli address get <name>"))
			return 1
		}
		a, err := d.Addr.Get(ctx, rest[0], g.vdom)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(a)
		return 0
	case "create", "update":
		var a daddr.Address
		if g.fromSrc != "" {
			if err := readJSONInput(g.fromSrc, &a); err != nil {
				output.Fail("input_error", err)
				return 1
			}
		} else {
			a = daddr.Address{
				Name:    flagVal(rest, "name"),
				Subnet:  flagVal(rest, "subnet"),
				Type:    flagVal(rest, "type"),
				Comment: flagVal(rest, "comment"),
			}
			if len(rest) >= 1 && a.Name == "" && verb == "update" {
				a.Name = rest[0]
			}
		}
		if err := a.Validate(); err != nil {
			output.Fail("validation_error", err)
			return 1
		}
		if g.dryRun {
			output.Print(map[string]any{"dry_run": true, "address": a, "vdom": g.vdom})
			return 0
		}
		var err error
		if verb == "create" {
			a, err = d.Addr.Create(ctx, a, g.vdom)
		} else {
			a, err = d.Addr.Update(ctx, a, g.vdom)
		}
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(a)
		return 0
	case "delete":
		if len(rest) < 1 {
			output.Fail("usage", fmt.Errorf("usage: fgcli address delete <name> --yes"))
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
		if err := d.Addr.Delete(ctx, rest[0], g.vdom); err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(map[string]string{"deleted": rest[0]})
		return 0
	default:
		output.Fail("usage", fmt.Errorf("unknown address verb %q", verb))
		return 1
	}
}

func runPolicy(ctx context.Context, d Deps, args []string) int {
	if len(args) < 1 {
		output.Fail("usage", fmt.Errorf("usage: fgcli policy list|get|create|update|enable|disable|move|clone|delete"))
		return 1
	}
	verb, rest := args[0], args[1:]
	g, rest, _ := parseGlobals(rest)
	switch verb {
	case "create", "update":
		return runPolicyWrite(ctx, d, verb, g, rest)
	case "enable", "disable":
		return runPolicyStatus(ctx, d, verb, g, rest)
	case "move":
		return runPolicyMove(ctx, d, g, rest, args)
	case "clone":
		return runPolicyClone(ctx, d, g, rest, args)
	case "list":
		list, err := d.Pol.List(ctx, g.vdom)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(filterSlice(list, g.filter))
		return 0
	case "get":
		if len(rest) < 1 {
			output.Fail("usage", fmt.Errorf("usage: fgcli policy get <id>"))
			return 1
		}
		id, err := strconv.ParseInt(rest[0], 10, 64)
		if err != nil {
			output.Fail("input_error", err)
			return 1
		}
		p, err := d.Pol.Get(ctx, id, g.vdom)
		if err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(p)
		return 0
	case "delete":
		if len(rest) < 1 {
			output.Fail("usage", fmt.Errorf("usage: fgcli policy delete <id> --yes"))
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
		if err := d.Pol.Delete(ctx, id, g.vdom); err != nil {
			output.Fail("api_error", err)
			return 2
		}
		output.Print(map[string]int64{"deleted": id})
		return 0
	default:
		output.Fail("usage", fmt.Errorf("unknown policy verb %q", verb))
		return 1
	}
}
