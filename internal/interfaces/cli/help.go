package cli

import (
	"github.com/local/fgcli/internal/infrastructure/output"
)

// Help texts are per-resource (not per-verb); shown without auth.
var helpTexts = map[string]string{
	"address":       "fgcli address list [--filter SUB] [--vdom V]\nfgcli address get <name>\nfgcli address create --name N --subnet S [--comment C] | --from-stdin\nfgcli address update <name>|--name N [--subnet S] | --from-stdin\nfgcli address delete <name> --yes [--dry-run]",
	"policy":        "fgcli policy list [--filter SUB] [--vdom V]\nfgcli policy get <id>\nfgcli policy create --policyid N --name M --srcintf a,b --dstintf c --srcaddr A --dstaddr B --service S [--action accept] | --from-stdin\nfgcli policy update <id> [same flags] | --from-stdin\nfgcli policy enable|disable <id> [--dry-run]\nfgcli policy move <id> --before X | --after Y [--dry-run]\nfgcli policy clone <src-id> --new-id N [--name M] [--dry-run]\nfgcli policy delete <id> --yes",
	"vip":           "fgcli vip list [--filter SUB] [--vdom V]\nfgcli vip get <name>\nfgcli vip create --name N --extip IP --mappedip IP[,IP] [--extintf wan1] [--protocol tcp] [--extport P --mappedport P] | --from-stdin\nfgcli vip update <name>|--name N [same flags] | --from-stdin\nfgcli vip delete <name> --yes [--dry-run]",
	"vipgrp":        "fgcli vipgrp list [--filter SUB] [--vdom V]\nfgcli vipgrp get <name>\nfgcli vipgrp create --name N --member VIP1,VIP2 [--comment C] | --from-stdin\nfgcli vipgrp update <name>|--name N [same flags] | --from-stdin\nfgcli vipgrp delete <name> --yes [--dry-run]",
	"addrgrp":       "fgcli addrgrp list [--filter SUB] [--vdom V]\nfgcli addrgrp get <name>\nfgcli addrgrp create --name N --member A1,A2 [--comment C] | --from-stdin\nfgcli addrgrp update <name>|--name N [same flags] | --from-stdin\nfgcli addrgrp delete <name> --yes [--dry-run]",
	"ippool":        "fgcli ippool list [--filter SUB] [--vdom V]\nfgcli ippool get <name>\nfgcli ippool create --name N --type overload --startip A --endip B | --from-stdin\nfgcli ippool update <name>|--name N [same flags] | --from-stdin\nfgcli ippool delete <name> --yes [--dry-run]",
	"snat":          "fgcli snat list [--filter SUB] [--vdom V]\nfgcli snat get <policyid>\nfgcli snat create --policyid N [flags] | --from-stdin (JSON for addrs/pools)\nfgcli snat update <policyid>|--policyid N [flags] | --from-stdin\nfgcli snat delete <policyid> --yes  (alias: central-snat-map)",
	"service":       "fgcli service list [--filter SUB] [--vdom V]\nfgcli service get <name>\nfgcli service create --name N [--protocol TCP] [--tcp 443] [--udp U] | --from-stdin\nfgcli service update <name>|--name N [same flags] | --from-stdin\nfgcli service delete <name> --yes [--dry-run]\nfgcli service voip get <name>",
	"service-group": "fgcli service-group list [--filter SUB] [--vdom V]\nfgcli service-group get <name>\nfgcli service-group create --name N --member S1,S2 | --from-stdin\nfgcli service-group update <name>|--name N [same flags] | --from-stdin\nfgcli service-group delete <name> --yes [--dry-run]",
	"interface":     "fgcli interface list [--filter SUB] [--vdom V]  (read-only)\nfgcli interface get <name>",
	"zone":          "fgcli zone list [--filter SUB] [--vdom V]  (read-only)\nfgcli zone get <name>",
	"raw":           "fgcli raw get|post|put|delete <path> [--vdom V] [--data JSON|--from-stdin|--from-file] [--dry-run]\npath may omit the /api/v2 prefix; --vdom sent only when explicit",
	"trace":         "fgcli trace --dst IP [--src IP] [--dport N] [--vdom V]\ncorrelates vip (extip+dport) with matching policies (src/dst/service)",
	"system":        "fgcli system status|license|fortiguard|ntp|dns|dhcp|snmp",
	"vpn":           "fgcli vpn ipsec [--filter SUB]",
	"routing":       "fgcli routing status",
	"user":          "fgcli user firewall|banned [--filter SUB]",
	"security":      "fgcli security ips [--filter SUB]\nfgcli security waf [list|get <id>]\nfgcli security dlp [list|get <id>]\nfgcli security proxy-pac [--out PATH]",
	"network":       "fgcli network dnsfilter list|get <id>",
	"switch":        "fgcli switch status [--filter SUB]",
	"wifi":          "fgcli wifi ap",
	"fortiview":     "fgcli fortiview",
	"log":           "fgcli log <fortianalyzer|forticloud|memory|disk> <type> [--filter SUB]",
	"doctor":        "fgcli doctor",
	"profile":       "fgcli profile list|show <name>|use <name>|add <name> ...|rm <name> --yes\nfgcli profile set <name>  (alias of use)",
	"backup":        "fgcli backup download [--scope global] [--out PATH]  (needs RW perms)",
}

// HelpFor returns the help text and whether the resource is known.
func HelpFor(resource string) (string, bool) {
	t, ok := helpTexts[resource]
	return t, ok
}

// WantHelp reports explicit help requests: help/-h/--help anywhere.
func WantHelp(args []string) bool {
	for _, a := range args {
		if a == "help" || a == "-h" || a == "--help" {
			return true
		}
	}
	return false
}

// PrintHelp writes human text to stderr and keeps the JSON envelope on stdout.
func PrintHelp(resource string) {
	t, ok := HelpFor(resource)
	if !ok {
		usage()
		return
	}
	output.Logf("%s", t)
	output.Print(map[string]string{"help": t})
}
