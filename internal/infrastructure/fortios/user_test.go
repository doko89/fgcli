package fortios

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	duser "github.com/local/fgcli/internal/domain/user"
)

func TestUserGroupAddMemberRoundtrip(t *testing.T) {
	var putPath, putBody string
	var putQuery string
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/cmdb/user/group/VPN", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Write([]byte(`{"results":[{"name":"VPN","member":[{"name":"alice"},{"name":"bob"}]}]}`))
			return
		}
		if r.Method == "PUT" {
			putPath, putQuery = r.URL.Path, r.URL.RawQuery
			b, _ := io.ReadAll(r.Body)
			putBody = string(b)
			w.Write([]byte(`{"status":"success"}`))
			return
		}
		t.Errorf("unexpected method %s", r.Method)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	repo := NewUserRepo(NewClient(srv.URL, "tok", false))

	gr, err := repo.GroupGet(context.Background(), "VPN", "root")
	if err != nil {
		t.Fatal(err)
	}
	if len(gr.Members) != 2 {
		t.Fatalf("expected 2 members, got %+v", gr)
	}
	members, added := duser.AddUnique(gr.Members, "doko_baru")
	if !added {
		t.Fatal("new member should be added")
	}
	gr.Members = members
	if _, err := repo.GroupUpdate(context.Background(), gr, "root"); err != nil {
		t.Fatal(err)
	}
	if putPath != "/api/v2/cmdb/user/group/VPN" {
		t.Fatalf("wrong PUT path: %s", putPath)
	}
	if !strings.Contains(putQuery, "vdom=root") {
		t.Fatalf("PUT must carry vdom: %s", putQuery)
	}
	var sent duser.UserGroup
	if err := json.Unmarshal([]byte(putBody), &sent); err != nil {
		t.Fatal(err)
	}
	if len(sent.Members) != 3 || sent.Members[2].Name != "doko_baru" {
		t.Fatalf("PUT must carry all 3 members (no overwrite loss): %s", putBody)
	}
}

func TestLocalCreatePathAndRedaction(t *testing.T) {
	var postPath string
	var postBody map[string]any
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/cmdb/user/local", func(w http.ResponseWriter, r *http.Request) {
		postPath = r.URL.Path
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &postBody)
		w.Write([]byte(`{"status":"success"}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	repo := NewUserRepo(NewClient(srv.URL, "tok", false))

	got, err := repo.LocalCreate(context.Background(), duser.LocalUser{Name: "doko_baru", Passwd: "s3cret"}, "root")
	if err != nil {
		t.Fatal(err)
	}
	if postPath != "/api/v2/cmdb/user/local" {
		t.Fatalf("wrong POST path: %s", postPath)
	}
	if postBody["name"] != "doko_baru" || postBody["passwd"] != "s3cret" {
		t.Fatalf("POST must carry name+passwd: %v", postBody)
	}
	if got.Passwd != "" {
		t.Fatal("returned user must not echo the password")
	}
}

func TestVpnSslPath(t *testing.T) {
	var path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		w.Write([]byte(`{"results":{"name":"test","status":"enable"}}`))
	}))
	defer srv.Close()
	repo := NewVpnRepo(NewClient(srv.URL, "tok", false))
	got, err := repo.Ssl(context.Background(), "root")
	if err != nil {
		t.Fatal(err)
	}
	if path != "/api/v2/cmdb/vpn.ssl/settings" {
		t.Fatalf("wrong ssl path: %s", path)
	}
	m, ok := got.(map[string]any)
	if !ok || m["status"] != "enable" {
		t.Fatalf("unexpected ssl settings: %+v", got)
	}
}
