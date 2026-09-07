package policy

import (
	"context"
	"testing"

	domain "github.com/local/fgcli/internal/domain/policy"
)

type fakeRepo struct {
	policies map[int64]domain.Policy
	moved    [3]int64
}

func (f *fakeRepo) List(ctx context.Context, vdom string) ([]domain.Policy, error) {
	return nil, nil
}
func (f *fakeRepo) Get(ctx context.Context, id int64, vdom string) (domain.Policy, error) {
	return f.policies[id], nil
}
func (f *fakeRepo) Create(ctx context.Context, p domain.Policy, vdom string) (domain.Policy, error) {
	f.policies[p.ID] = p
	return p, nil
}
func (f *fakeRepo) Update(ctx context.Context, p domain.Policy, vdom string) (domain.Policy, error) {
	f.policies[p.ID] = p
	return p, nil
}
func (f *fakeRepo) Move(ctx context.Context, id, before, after int64, vdom string) error {
	f.moved = [3]int64{id, before, after}
	return nil
}
func (f *fakeRepo) Delete(ctx context.Context, id int64, vdom string) error { return nil }

func TestCloneRenamesAndReIDs(t *testing.T) {
	uc := New(&fakeRepo{policies: map[int64]domain.Policy{
		7: {ID: 7, Name: "MAIL_IN", Action: "accept", Status: "enable"},
	}})
	got, err := uc.Clone(context.Background(), 7, 70, "", "root")
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != 70 || got.Name != "MAIL_IN_clone70" {
		t.Fatalf("unexpected clone: %+v", got)
	}
	if _, err := uc.Clone(context.Background(), 7, 0, "", "root"); err == nil {
		t.Fatal("clone with id 0 should fail validation")
	}
}

func TestMoveRequiresExactlyOneAnchor(t *testing.T) {
	uc := New(&fakeRepo{policies: map[int64]domain.Policy{}})
	if err := uc.Move(context.Background(), 1, 0, 0, "root"); err == nil {
		t.Fatal("move without anchor should fail")
	}
	if err := uc.Move(context.Background(), 1, 2, 3, "root"); err == nil {
		t.Fatal("move with both anchors should fail")
	}
}

func TestSetStatusRoundtrip(t *testing.T) {
	repo := &fakeRepo{policies: map[int64]domain.Policy{2: {ID: 2, Status: "disable"}}}
	uc := New(repo)
	got, err := uc.SetStatus(context.Background(), 2, "enable", "root")
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != "enable" {
		t.Fatalf("status not flipped: %+v", got)
	}
	if _, err := uc.SetStatus(context.Background(), 2, "bogus", "root"); err == nil {
		t.Fatal("bogus status should fail")
	}
}
