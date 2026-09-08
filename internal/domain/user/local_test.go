package user

import (
	"testing"
)

func TestAddUniqueIdempotent(t *testing.T) {
	ms, added := AddUnique(nil, "doko_baru")
	if !added || len(ms) != 1 {
		t.Fatalf("first add should insert: %+v %v", ms, added)
	}
	ms, added = AddUnique(ms, "DOKO_BARU")
	if added || len(ms) != 1 {
		t.Fatalf("re-add must be idempotent (case-insensitive): %+v %v", ms, added)
	}
}

func TestRemoveName(t *testing.T) {
	ms := []GroupMember{{Name: "a"}, {Name: "b"}}
	kept, removed := RemoveName(ms, "B")
	if !removed || len(kept) != 1 || kept[0].Name != "a" {
		t.Fatalf("remove should drop b: %+v %v", kept, removed)
	}
	if _, removed := RemoveName(kept, "zzz"); removed {
		t.Fatal("removing absent member must report removed=false")
	}
}

func TestUserGroupValidate(t *testing.T) {
	if err := (UserGroup{Name: "VPN"}).Validate(); err != nil {
		t.Fatalf("valid group rejected: %v", err)
	}
	if err := (UserGroup{}).Validate(); err == nil {
		t.Fatal("empty group name must fail")
	}
}

func TestLocalUserValidate(t *testing.T) {
	if err := (LocalUser{Name: "n", Passwd: "p"}).Validate(); err != nil {
		t.Fatalf("valid user rejected: %v", err)
	}
	if err := (LocalUser{Passwd: "p"}).Validate(); err == nil {
		t.Fatal("empty user name must fail")
	}
}
