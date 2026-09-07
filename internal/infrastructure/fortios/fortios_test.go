package fortios

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMoveURLShape(t *testing.T) {
	var path, query string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path, query = r.URL.Path, r.URL.RawQuery
		w.Write([]byte(`{"status":"success"}`))
	}))
	defer srv.Close()
	repo := NewPolicyRepo(NewClient(srv.URL, "tok", false))
	if err := repo.Move(context.Background(), 16, 15, 0, "root"); err != nil {
		t.Fatal(err)
	}
	if path != "/api/v2/cmdb/firewall/policy/16/move" {
		t.Fatalf("wrong move path: %s", path)
	}
	if !strings.Contains(query, "before=15") || !strings.Contains(query, "vdom=root") {
		t.Fatalf("wrong move query: %s", query)
	}
}

func TestDhcp424MeansEmpty(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(424)
		w.Write([]byte(`{"status":"error","http_status":424}`))
	}))
	defer srv.Close()
	repo := NewSysRepo(NewClient(srv.URL, "tok", false))
	list, err := repo.Dhcp(context.Background())
	if err != nil {
		t.Fatalf("424 should decode as empty, got: %v", err)
	}
	if len(list) != 0 {
		t.Fatalf("expected empty list, got %+v", list)
	}
}

func TestPac424MeansEmpty(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(424)
		w.Write([]byte(`{"status":"error","http_status":424}`))
	}))
	defer srv.Close()
	repo := NewSecurityRepo(NewClient(srv.URL, "tok", false))
	b, err := repo.DownloadPac(context.Background(), "root")
	if err != nil {
		t.Fatalf("424 should decode as empty, got: %v", err)
	}
	if len(b) != 0 {
		t.Fatalf("expected empty bytes, got %d", len(b))
	}
}

func TestAllowEmptyBody(t *testing.T) {
	var list []string
	if err := decodeResultsAllowEmpty([]byte(``), &list); err != nil {
		t.Fatalf("empty body should decode empty: %v", err)
	}
	if err := decodeResultsAllowEmpty([]byte("  \n"), &list); err != nil {
		t.Fatalf("blank body should decode empty: %v", err)
	}
}

func TestLogShapes(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/log/memory/dns/raw", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"results":[{"date":"2026-09-07","msg":"query"}]}`))
	})
	mux.HandleFunc("/api/v2/log/disk/ips/raw", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{"date":"2026-09-07"}]`))
	})
	mux.HandleFunc("/api/v2/log/memory/empty/raw", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(``))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	repo := NewLogRepo(NewClient(srv.URL, "tok", false))

	got, err := repo.Fetch(context.Background(), "memory", "dns", "root")
	if err != nil || len(got) != 1 {
		t.Fatalf("envelope shape failed: %+v %v", got, err)
	}
	got, err = repo.Fetch(context.Background(), "disk", "ips", "root")
	if err != nil || len(got) != 1 {
		t.Fatalf("bare-array shape failed: %+v %v", got, err)
	}
	got, err = repo.Fetch(context.Background(), "memory", "empty", "root")
	if err != nil || len(got) != 0 {
		t.Fatalf("empty body failed: %+v %v", got, err)
	}
	if _, err := repo.Fetch(context.Background(), "nope", "dns", "root"); err == nil {
		t.Fatal("bad source should fail")
	}
	if _, err := repo.Fetch(context.Background(), "memory", "", "root"); err == nil {
		t.Fatal("empty type should fail")
	}
}

func TestDecodeResultsAllowEmpty(t *testing.T) {
	var list []string
	if err := decodeResultsAllowEmpty([]byte(`{"results":null}`), &list); err != nil {
		t.Fatalf("null results should decode empty: %v", err)
	}
	if list != nil {
		t.Fatalf("expected nil list, got %+v", list)
	}
	if err := decodeResultsAllowEmpty([]byte(`{"results":[]}`), &list); err != nil {
		t.Fatalf("empty results should decode empty: %v", err)
	}
}

func TestBackupDownloadBytes(t *testing.T) {
	const payload = "#config-version=FGT60F-7.4.9\n"
	var auth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth = r.Header.Get("Authorization")
		if r.URL.Path != "/api/v2/monitor/system/config/backup" {
			t.Errorf("wrong backup path: %s", r.URL.Path)
		}
		w.Write([]byte(payload))
	}))
	defer srv.Close()
	repo := NewBackupRepo(NewClient(srv.URL, "tok", false))
	b, err := repo.Download(context.Background(), "global")
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != payload {
		t.Fatalf("backup bytes mismatch: %q", b)
	}
	if auth != "Bearer tok" {
		t.Fatalf("missing bearer auth: %q", auth)
	}
}
