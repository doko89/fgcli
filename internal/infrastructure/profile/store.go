package profile

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/local/fgcli/internal/infrastructure/config"
	"gopkg.in/yaml.v3"
)

type Profile struct {
	Host      string `yaml:"host"`
	APIKey    string `yaml:"api_key,omitempty"`
	APIKeyEnv string `yaml:"api_key_env,omitempty"`
	Token     string `yaml:"token,omitempty"`
	TokenEnv  string `yaml:"token_env,omitempty"`
	Vdom      string `yaml:"vdom,omitempty"`
	Insecure  bool   `yaml:"insecure,omitempty"`
}

type File struct {
	Active   string             `yaml:"active,omitempty"`
	Profiles map[string]Profile `yaml:"profiles"`
}

type Summary struct {
	Name      string `json:"name"`
	Host      string `json:"host"`
	Vdom      string `json:"vdom"`
	Insecure  bool   `json:"insecure"`
	HasAPIKey bool   `json:"has_api_key"`
	Active    bool   `json:"active"`
}

func Path() (string, error) {
	if p := os.Getenv("FGCLI_CONFIG"); p != "" {
		return p, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".fgcli", "config.yaml"), nil
}

func LoadFile() (*File, error) {
	p, err := Path()
	if err != nil {
		return nil, err
	}
	b, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return &File{Profiles: map[string]Profile{}}, nil
		}
		return nil, err
	}
	var f File
	if err := yaml.Unmarshal(b, &f); err != nil {
		return nil, fmt.Errorf("profile: parse %s: %w", p, err)
	}
	if f.Profiles == nil {
		f.Profiles = map[string]Profile{}
	}
	return &f, nil
}

func (f *File) Save() error {
	p, err := Path()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o700); err != nil {
		return err
	}
	b, err := yaml.Marshal(f)
	if err != nil {
		return err
	}
	return os.WriteFile(p, b, 0o600)
}

func (f *File) Summaries() []Summary {
	out := make([]Summary, 0, len(f.Profiles))
	for name, pr := range f.Profiles {
		out = append(out, Summary{
			Name: name, Host: pr.Host,
			Vdom: orDefault(pr.Vdom, "root"), Insecure: pr.Insecure,
			HasAPIKey: pr.APIKey != "" || pr.APIKeyEnv != "" || pr.Token != "" || pr.TokenEnv != "",
			Active:    name == f.Active,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func (f *File) Get(name string) (Profile, error) {
	pr, ok := f.Profiles[name]
	if !ok {
		return Profile{}, fmt.Errorf("profile %q not found", name)
	}
	return pr, nil
}

func Resolve(base config.Config, f *File, name string) (config.Config, error) {
	if name == "" {
		name = f.Active
	}
	if name == "" {
		return withDefaults(base), nil
	}
	pr, err := f.Get(name)
	if err != nil {
		return config.Config{}, err
	}
	c := base
	if c.Host == "" {
		c.Host = pr.Host
	}
	if c.APIKey == "" {
		c.APIKey = secret(pr.APIKey, pr.APIKeyEnv)
	}
	if c.APIKey == "" {
		c.APIKey = secret(pr.Token, pr.TokenEnv)
	}
	if c.Vdom == "" {
		c.Vdom = pr.Vdom
	}
	c.Insecure = c.Insecure || pr.Insecure
	return withDefaults(c), nil
}

func secret(inline, env string) string {
	if inline != "" {
		return inline
	}
	if env != "" {
		return os.Getenv(env)
	}
	return ""
}

func withDefaults(c config.Config) config.Config {
	if c.Vdom == "" {
		c.Vdom = "root"
	}
	return c
}

func orDefault(s, d string) string {
	if s == "" {
		return d
	}
	return s
}

var ErrExists = errors.New("profile already exists (use --force to overwrite)")
