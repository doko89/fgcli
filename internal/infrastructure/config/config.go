package config

import (
	"errors"
	"os"
)

type Config struct {
	Host     string
	APIKey   string
	Insecure bool
	Vdom     string
}

func LoadEnv() Config {
	key := os.Getenv("FG_API_KEY")
	if key == "" {
		key = os.Getenv("FG_TOKEN")
	}
	c := Config{
		Host:   os.Getenv("FG_HOST"),
		APIKey: key,
		Vdom:   os.Getenv("FG_VDOM"),
	}
	if v := os.Getenv("FG_INSECURE"); v == "1" || v == "true" {
		c.Insecure = true
	}
	return c
}

func (c Config) Validate() error {
	if c.Host == "" {
		return errors.New("FG_HOST is required (e.g. https://192.0.2.1)")
	}
	if c.APIKey == "" {
		return errors.New("auth required: FG_API_KEY (or --api-key)")
	}
	return nil
}

func Load() (Config, error) {
	c := LoadEnv()
	if c.Vdom == "" {
		c.Vdom = "root"
	}
	return c, c.Validate()
}
