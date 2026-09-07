package trace

import (
	"context"
	"testing"

	daddr "github.com/local/fgcli/internal/domain/address"
	dpol "github.com/local/fgcli/internal/domain/policy"
	dsvc "github.com/local/fgcli/internal/domain/service"
	dvip "github.com/local/fgcli/internal/domain/vip"
)

type stubVip struct{ list []dvip.Vip }

func (s stubVip) List(ctx context.Context, vdom string) ([]dvip.Vip, error) { return s.list, nil }
func (s stubVip) Get(ctx context.Context, n, v string) (dvip.Vip, error)    { return dvip.Vip{}, nil }
func (s stubVip) Create(ctx context.Context, v dvip.Vip, vdom string) (dvip.Vip, error) {
	return v, nil
}
func (s stubVip) Update(ctx context.Context, v dvip.Vip, vdom string) (dvip.Vip, error) {
	return v, nil
}
func (s stubVip) Delete(ctx context.Context, n, v string) error { return nil }

type stubPol struct{ list []dpol.Policy }

func (s stubPol) List(ctx context.Context, vdom string) ([]dpol.Policy, error) { return s.list, nil }
func (s stubPol) Get(ctx context.Context, id int64, v string) (dpol.Policy, error) {
	return dpol.Policy{}, nil
}
func (s stubPol) Create(ctx context.Context, p dpol.Policy, v string) (dpol.Policy, error) {
	return p, nil
}
func (s stubPol) Update(ctx context.Context, p dpol.Policy, v string) (dpol.Policy, error) {
	return p, nil
}
func (s stubPol) Move(ctx context.Context, id, before, after int64, v string) error { return nil }
func (s stubPol) Delete(ctx context.Context, id int64, v string) error              { return nil }

type stubAddr struct{ list []daddr.Address }

func (s stubAddr) List(ctx context.Context, vdom string) ([]daddr.Address, error) {
	return s.list, nil
}
func (s stubAddr) Get(ctx context.Context, n, v string) (daddr.Address, error) {
	return daddr.Address{}, nil
}
func (s stubAddr) Create(ctx context.Context, a daddr.Address, v string) (daddr.Address, error) {
	return a, nil
}
func (s stubAddr) Update(ctx context.Context, a daddr.Address, v string) (daddr.Address, error) {
	return a, nil
}
func (s stubAddr) Delete(ctx context.Context, n, v string) error { return nil }

type stubSvc struct{ list []dsvc.Svc }

func (s stubSvc) List(ctx context.Context, vdom string) ([]dsvc.Svc, error) { return s.list, nil }
func (s stubSvc) Get(ctx context.Context, n, v string) (dsvc.Svc, error)    { return dsvc.Svc{}, nil }
func (s stubSvc) Create(ctx context.Context, v dsvc.Svc, vdom string) (dsvc.Svc, error) {
	return v, nil
}
func (s stubSvc) Update(ctx context.Context, v dsvc.Svc, vdom string) (dsvc.Svc, error) {
	return v, nil
}
func (s stubSvc) Delete(ctx context.Context, n, v string) error { return nil }

func testUC() *UseCase {
	return &UseCase{
		Vips: stubVip{list: []dvip.Vip{{
			Name: "VIP_MAIL", ExtIP: "203.0.113.10", PortForward: "enable",
			Protocol: "tcp", ExtPort: "443", MappedPort: "443",
			MappedIP: []dvip.MappedRange{{Range: "192.168.1.43"}},
		}}},
		Pols: stubPol{list: []dpol.Policy{{
			ID: 7, Name: "MAIL_IN", SrcAddr: []dpol.Name{"all"},
			DstAddr: []dpol.Name{"VIP_MAIL"}, Service: []dpol.Name{"HTTPS"},
			Action: "accept", Status: "enable",
		}, {
			ID: 8, Name: "OTHER", SrcAddr: []dpol.Name{"all"},
			DstAddr: []dpol.Name{"OTHER_NET"}, Service: []dpol.Name{"ALL"},
			Action: "accept", Status: "enable",
		}}},
		Addrs: stubAddr{list: []daddr.Address{
			{Name: "OTHER_NET", Subnet: "10.9.0.0 255.255.0.0"},
		}},
		Svcs: stubSvc{list: []dsvc.Svc{
			{Name: "HTTPS", Protocol: "TCP", TCPPorts: "443"},
		}},
	}
}

func TestTraceVipAndPolicyHit(t *testing.T) {
	res, err := testUC().Trace(context.Background(), "root", Input{Src: "8.8.8.8", Dst: "203.0.113.10", DPort: 443})
	if err != nil {
		t.Fatal(err)
	}
	if res.Vip == nil || res.Vip.Name != "VIP_MAIL" {
		t.Fatalf("expected vip hit: %+v", res.Vip)
	}
	if len(res.Policies) != 1 || res.Policies[0].Policy.ID != 7 {
		t.Fatalf("expected policy 7 only: %+v", res.Policies)
	}
}

func TestTraceNoVipNote(t *testing.T) {
	res, err := testUC().Trace(context.Background(), "root", Input{Dst: "198.51.100.9", DPort: 443})
	if err != nil {
		t.Fatal(err)
	}
	if res.Vip != nil || res.Note == "" {
		t.Fatalf("expected no-vip note: %+v", res)
	}
	if len(res.Policies) != 0 {
		t.Fatalf("no policy should match unknown dst: %+v", res.Policies)
	}
}

func TestTraceWrongPortMisses(t *testing.T) {
	res, err := testUC().Trace(context.Background(), "root", Input{Src: "8.8.8.8", Dst: "203.0.113.10", DPort: 22})
	if err != nil {
		t.Fatal(err)
	}
	if res.Vip != nil {
		t.Fatalf("port 22 must not hit vip 443: %+v", res.Vip)
	}
}

func TestIPInSubnet(t *testing.T) {
	cases := []struct {
		sub, ip string
		want    bool
	}{
		{"192.168.1.0 255.255.255.0", "192.168.1.43", true},
		{"192.168.1.0 255.255.255.0", "192.168.2.1", false},
		{"10.0.0.0/8", "10.9.9.9", true},
		{"203.0.113.10", "203.0.113.10", true},
		{"0.0.0.0 0.0.0.0", "8.8.8.8", true},
		{"", "8.8.8.8", false},
	}
	for _, c := range cases {
		if got := ipInSubnet(c.sub, c.ip); got != c.want {
			t.Errorf("ipInSubnet(%q,%q)=%v want %v", c.sub, c.ip, got, c.want)
		}
	}
}

func TestRangeOK(t *testing.T) {
	if !rangeOK("443", 443) || !rangeOK("80-90", 85) || rangeOK("80-90", 91) {
		t.Fatal("range matching broken")
	}
	if rangeOK("443", 22) {
		t.Fatal("wrong port must not match")
	}
}
