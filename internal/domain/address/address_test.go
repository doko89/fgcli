package address

import "testing"

func TestValidate(t *testing.T) {
	if err := (Address{Name: "lan-net", Subnet: "10.0.0.0/24"}).Validate(); err != nil {
		t.Fatalf("valid address rejected: %v", err)
	}
	if err := (Address{}).Validate(); err == nil {
		t.Fatal("empty name should fail")
	}
	if err := (Address{Name: "bad name"}).Validate(); err == nil {
		t.Fatal("name with space should fail")
	}
}
