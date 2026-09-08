package profile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/local/fgcli/internal/infrastructure/config"
)

func tmpConfig(t *testing.T) {
	t.Helper()
	t.Setenv("FGCLI_CONFIG", filepath.Join(t.TempDir(), "config.yaml"))
}

func TestDefaultPathUsesDotDir(t *testing.T) {
	t.Setenv("FGCLI_CONFIG", "")
	home := t.TempDir()
	t.Setenv("HOME", home)
	p, err := Path()
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(home, ".fgcli", "config.yaml")
	if p != want {
		t.Fatalf("default path = %q, want %q", p, want)
	}
}

func TestSaveLoadRoundtripAndPerms(t *testing.T) {
	tmpConfig(t)
	f, err := LoadFile()
	if err != nil {
		t.Fatal(err)
	}
	if len(f.Profiles) != 0 {
		t.Fatalf("fresh file should be empty: %+v", f)
	}
	f.Profiles["fg1"] = Profile{Host: "https://10.0.0.1", APIKeyEnv: "FG_API_KEY_FG1", Insecure: true}
	f.Active = "fg1"
	if err := f.Save(); err != nil {
		t.Fatal(err)
	}
	p, _ := Path()
	st, err := os.Stat(p)
	if err != nil {
		t.Fatal(err)
	}
	if st.Mode().Perm() != 0o600 {
		t.Fatalf("config must be 0600, got %o", st.Mode().Perm())
	}
	got, err := LoadFile()
	if err != nil {
		t.Fatal(err)
	}
	if got.Active != "fg1" || got.Profiles["fg1"].Host != "https://10.0.0.1" {
		t.Fatalf("roundtrip mismatch: %+v", got)
	}
}

func TestResolvePrecedence(t *testing.T) {
	tmpConfig(t)
	t.Setenv("FG_API_KEY_FG1", "s3cret")
	f := &File{Profiles: map[string]Profile{
		"fg1": {Host: "https://from-profile", APIKeyEnv: "FG_API_KEY_FG1", Vdom: "root"},
	}, Active: "fg1"}

	c, err := Resolve(config.Config{}, f, "")
	if err != nil {
		t.Fatal(err)
	}
	if c.Host != "https://from-profile" || c.APIKey != "s3cret" || c.Vdom != "root" {
		t.Fatalf("profile should fill: %+v", c)
	}

	c2, err := Resolve(config.Config{Host: "https://from-flag"}, f, "")
	if err != nil {
		t.Fatal(err)
	}
	if c2.Host != "https://from-flag" {
		t.Fatalf("flag should win: %+v", c2)
	}

	f.Profiles["fg2"] = Profile{Host: "https://fg2", APIKey: "x"}
	c3, err := Resolve(config.Config{}, f, "fg2")
	if err != nil {
		t.Fatal(err)
	}
	if c3.Host != "https://fg2" {
		t.Fatalf("explicit profile should win over active: %+v", c3)
	}

	if _, err := Resolve(config.Config{}, f, "nope"); err == nil {
		t.Fatal("unknown profile should fail")
	}
}

func TestResolveLegacyTokenFallback(t *testing.T) {
	tmpConfig(t)
	t.Setenv("FG_LEGACY_FG1", "legacy-secret")
	f := &File{Profiles: map[string]Profile{
		"fg1": {Host: "https://legacy", TokenEnv: "FG_LEGACY_FG1"},
	}, Active: "fg1"}
	c, err := Resolve(config.Config{}, f, "")
	if err != nil {
		t.Fatal(err)
	}
	if c.APIKey != "legacy-secret" {
		t.Fatalf("legacy token_env should resolve to APIKey: %+v", c)
	}
}

func TestIsPlaceholder(t *testing.T) {
	for _, k := range []string{"will_be_available", "WILL_BE_AVAILABLE", "changeme", "xxx", "  placeholder  ", "tbd"} {
		if !IsPlaceholder(k) {
			t.Fatalf("expected placeholder: %q", k)
		}
	}
	for _, k := range []string{"xmsjcGbx99txwhjy0NgsNxd0mm5G87", "s3cret", "xxxtra-long-key"} {
		if IsPlaceholder(k) {
			t.Fatalf("real key flagged as placeholder: %q", k)
		}
	}
}

func TestValidateFile(t *testing.T) {
	tmpConfig(t)
	t.Setenv("FG_EMPTY_REF", "")
	f := &File{
		Active: "ghost",
		Profiles: map[string]Profile{
			"ho-readonly": {Host: "https://172.28.29.1:15443", APIKey: "realkey123", Insecure: true},
			"ho-admin":    {Host: "https://172.28.29.1:15443", APIKey: "will_be_available"},
			"env-missing": {Host: "https://10.0.0.1", APIKeyEnv: "FG_EMPTY_REF"},
			"no-secret":   {Host: "https://10.0.0.2"},
			"no-host":     {APIKey: "realkey123"},
			"cyber-admin": {Host: "https://172.28.29.1:15443", APIKey: "alsorealkey"},
		},
	}
	warns := ValidateFile(f)
	byCode := map[string]int{}
	for _, w := range warns {
		byCode[w.Code]++
	}
	for code, want := range map[string]int{
		"unknown_active": 1, "placeholder_key": 1, "missing_secret": 2,
		"missing_host": 1, "insecure": 1, "duplicate_host": 1, "no_active": 0,
	} {
		if byCode[code] != want {
			t.Fatalf("code %q: got %d warnings, want %d (%+v)", code, byCode[code], want, warns)
		}
	}
}

func TestSummariesRedactSecrets(t *testing.T) {
	f := &File{Profiles: map[string]Profile{
		"fg1": {Host: "h", APIKey: "tok"},
		"fg2": {Host: "h2", TokenEnv: "FG_LEGACY"},
	}, Active: "fg1"}
	s := f.Summaries()
	if len(s) != 2 {
		t.Fatalf("unexpected summary: %+v", s)
	}
	for _, sm := range s {
		if !sm.HasAPIKey {
			t.Fatalf("expected HasAPIKey: %+v", sm)
		}
	}
	if !s[0].Active {
		t.Fatalf("expected fg1 active: %+v", s)
	}
}
