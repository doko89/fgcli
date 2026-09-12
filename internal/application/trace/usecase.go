package trace

import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"sort"
	"strconv"
	"strings"

	daddr "github.com/local/fgcli/internal/domain/address"
	dpol "github.com/local/fgcli/internal/domain/policy"
	dsvc "github.com/local/fgcli/internal/domain/service"
	dvip "github.com/local/fgcli/internal/domain/vip"
)

type Input struct {
	Src   string `json:"src"`
	Dst   string `json:"dst"`
	DPort int    `json:"dport"`
}

type PolicyHit struct {
	Policy dpol.Policy `json:"policy"`
	Via    []string    `json:"via"`
}

type Result struct {
	Input    Input       `json:"input"`
	Vip      *dvip.Vip   `json:"vip_match"`
	Policies []PolicyHit `json:"policies"`
	Note     string      `json:"note,omitempty"`
}

type UseCase struct {
	Vips  dvip.Repository
	Pols  dpol.Repository
	Addrs daddr.Repository
	Svcs  dsvc.SvcRepository
}

func (u *UseCase) Trace(ctx context.Context, vdom string, in Input) (Result, error) {
	res := Result{Input: in, Policies: []PolicyHit{}}
	vips, err := u.Vips.List(ctx, vdom)
	if err != nil {
		return res, err
	}
	res.Vip = findVip(vips, in)
	pols, err := u.Pols.List(ctx, vdom)
	if err != nil {
		return res, err
	}
	addrs, err := u.Addrs.List(ctx, vdom)
	if err != nil {
		return res, err
	}
	subnets := buildSubnetMap(addrs)
	svcs, err := u.Svcs.List(ctx, vdom)
	if err != nil {
		return res, err
	}
	ports := buildSvcMap(svcs)
	res.Policies = matchPolicies(pols, vipNameOf(res.Vip), in, subnets, ports)
	sortHits(res.Policies)
	if res.Vip == nil {
		res.Note = fmt.Sprintf("no vip with extip %s (and dport %d); inbound is not DNATed on this box", in.Dst, in.DPort)
	}
	return res, nil
}

func findVip(vips []dvip.Vip, in Input) *dvip.Vip {
	for i := range vips {
		if vips[i].ExtIP == in.Dst && portOK(vips[i].PortForward, vips[i].ExtPort, in.DPort) {
			return &vips[i]
		}
	}
	return nil
}

func buildSubnetMap(addrs []daddr.Address) map[string]string {
	subnets := map[string]string{}
	for _, a := range addrs {
		subnets[strings.ToLower(a.Name)] = a.Subnet
	}
	return subnets
}

func buildSvcMap(svcs []dsvc.Svc) map[string]dsvc.Svc {
	ports := map[string]dsvc.Svc{}
	for _, s := range svcs {
		ports[strings.ToLower(s.Name)] = s
	}
	return ports
}

func vipNameOf(v *dvip.Vip) string {
	if v == nil {
		return ""
	}
	return strings.ToLower(v.Name)
}

func matchPolicies(pols []dpol.Policy, vipName string, in Input, subnets map[string]string, ports map[string]dsvc.Svc) []PolicyHit {
	out := []PolicyHit{}
	for _, p := range pols {
		hit, ok := matchOnePolicy(p, vipName, in, subnets, ports)
		if !ok {
			continue
		}
		out = append(out, hit)
	}
	return out
}

func matchOnePolicy(p dpol.Policy, vipName string, in Input, subnets map[string]string, ports map[string]dsvc.Svc) (PolicyHit, bool) {
	via := []string{}
	ok, why := dstLeg(p, vipName, in.Dst, subnets)
	if !ok {
		return PolicyHit{}, false
	}
	via = append(via, why)
	ok, why = srcLeg(p, in.Src, subnets)
	if !ok {
		return PolicyHit{}, false
	}
	via = append(via, why)
	ok, why = svcLeg(p, in.DPort, ports)
	if !ok {
		return PolicyHit{}, false
	}
	via = append(via, why)
	if strings.EqualFold(p.Status, "disable") {
		via = append(via, "status:DISABLED")
	}
	return PolicyHit{Policy: p, Via: via}, true
}

// sortHits orders specific matches first: vip-name > concrete subnet >
// catch-all, disabled policies last. Stable, so policy order is preserved
// within the same specificity.
func sortHits(hits []PolicyHit) {
	score := func(h PolicyHit) int {
		s := 0
		for _, v := range h.Via {
			switch {
			case strings.HasPrefix(v, "dstaddr:vip="):
				s += 4
			case strings.HasPrefix(v, "dstaddr-subnet:") && !strings.Contains(v, "0.0.0.0"):
				s += 2
			case v == "status:DISABLED":
				s -= 8
			}
		}
		return s
	}
	sort.SliceStable(hits, func(i, j int) bool { return score(hits[i]) > score(hits[j]) })
}

func portOK(portforward, extport string, dport int) bool {
	if dport <= 0 {
		return true
	}
	if !strings.EqualFold(portforward, "enable") {
		return true
	}
	return rangeOK(extport, dport)
}

// rangeOK matches FortiOS portrange text: "443", "80-90", "80 443".
func rangeOK(spec string, port int) bool {
	for _, f := range strings.Fields(strings.ReplaceAll(spec, ",", " ")) {
		if strings.Contains(f, "-") {
			lo, hi, ok := splitRange(f)
			if ok && port >= lo && port <= hi {
				return true
			}
			continue
		}
		if n, err := strconv.Atoi(f); err == nil && n == port {
			return true
		}
	}
	return false
}

func splitRange(f string) (int, int, bool) {
	lo, hi, found := strings.Cut(f, "-")
	a, err1 := strconv.Atoi(strings.TrimSpace(lo))
	b, err2 := strconv.Atoi(strings.TrimSpace(hi))
	if !found || err1 != nil || err2 != nil {
		return 0, 0, false
	}
	return a, b, true
}

func dstLeg(p dpol.Policy, vipName, dst string, subnets map[string]string) (bool, string) {
	for _, d := range p.DstAddr {
		n := strings.ToLower(string(d))
		if vipName != "" && n == vipName {
			return true, "dstaddr:vip=" + string(d)
		}
		if sub, ok := subnets[n]; ok && ipInSubnet(sub, dst) {
			return true, "dstaddr-subnet:" + string(d) + "=" + sub
		}
		if n == "all" {
			return true, "dstaddr:all"
		}
	}
	return false, ""
}

func srcLeg(p dpol.Policy, src string, subnets map[string]string) (bool, string) {
	if src == "" {
		return true, "src:any"
	}
	for _, s := range p.SrcAddr {
		n := strings.ToLower(string(s))
		if n == "all" {
			return true, "srcaddr:all"
		}
		if sub, ok := subnets[n]; ok && ipInSubnet(sub, src) {
			return true, "srcaddr-subnet:" + string(s) + "=" + sub
		}
	}
	return false, ""
}

func svcLeg(p dpol.Policy, dport int, ports map[string]dsvc.Svc) (bool, string) {
	if dport <= 0 {
		return true, "service:any-port"
	}
	for _, s := range p.Service {
		n := strings.ToLower(string(s))
		if n == "all" {
			return true, "service:ALL"
		}
		svc, ok := ports[n]
		if !ok {
			return true, "service:" + string(s) + ":unresolved(assume-match)"
		}
		if rangeOK(svc.TCPPorts, dport) || rangeOK(svc.UDPPorts, dport) || rangeOK(svc.SCTPPorts, dport) {
			return true, "service:" + string(s) + ":port-ok"
		}
	}
	return false, ""
}

// ipInSubnet handles FortiOS subnet text: "IP MASK", "IP/CIDR", "IP", "".
func ipInSubnet(subnet, ip string) bool {
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		return false
	}
	subnet = strings.TrimSpace(subnet)
	if subnet == "" {
		return false
	}
	if strings.Contains(subnet, "/") {
		pfx, err := netip.ParsePrefix(subnet)
		if err != nil {
			return false
		}
		return pfx.Contains(addr)
	}
	f := strings.Fields(subnet)
	if len(f) == 1 {
		a, err := netip.ParseAddr(f[0])
		return err == nil && a == addr
	}
	mask := net.ParseIP(f[1]).To4()
	base := net.ParseIP(f[0]).To4()
	raw := net.ParseIP(ip).To4()
	if mask == nil || base == nil || raw == nil {
		return false
	}
	for i := 0; i < 4; i++ {
		if base[i]&mask[i] != raw[i]&mask[i] {
			return false
		}
	}
	return true
}
