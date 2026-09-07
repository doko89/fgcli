package address

import (
	"context"
	"errors"
	"testing"

	domain "github.com/local/fgcli/internal/domain/address"
)

type fakeRepo struct {
	created domain.Address
	err     error
}

func (f *fakeRepo) List(ctx context.Context, vdom string) ([]domain.Address, error) {
	return []domain.Address{f.created}, f.err
}
func (f *fakeRepo) Get(ctx context.Context, name, vdom string) (domain.Address, error) {
	return f.created, f.err
}
func (f *fakeRepo) Create(ctx context.Context, a domain.Address, vdom string) (domain.Address, error) {
	return a, f.err
}
func (f *fakeRepo) Update(ctx context.Context, a domain.Address, vdom string) (domain.Address, error) {
	return a, f.err
}
func (f *fakeRepo) Delete(ctx context.Context, name, vdom string) error { return f.err }

func TestCreateValidatesBeforeRepo(t *testing.T) {
	uc := New(&fakeRepo{})
	if _, err := uc.Create(context.Background(), domain.Address{}, "root"); err == nil {
		t.Fatal("expected validation error, repo must not be called")
	}
	got, err := uc.Create(context.Background(), domain.Address{Name: "web-srv"}, "root")
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "web-srv" {
		t.Fatalf("unexpected result: %+v", got)
	}
}

func TestRepoErrorPropagates(t *testing.T) {
	uc := New(&fakeRepo{err: errors.New("boom")})
	if _, err := uc.List(context.Background(), "root"); err == nil {
		t.Fatal("expected repo error")
	}
}
