// Package config loads ~/.config/nova-pkg/config.toml.
package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config is the parsed config file.
type Config struct {
	// Priority orders backends for search output and install/remove
	// preference. Unlisted backends append after in default order.
	Priority []string `toml:"priority"`
}

// Load reads the config file. A missing file is not an error: it
// returns the zero Config.
func Load() (Config, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return Config{}, fmt.Errorf("config: %w", err)
	}
	f, err := os.Open(filepath.Join(dir, "nova-pkg", "config.toml"))
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, nil
		}
		return Config{}, fmt.Errorf("config: %w", err)
	}
	defer f.Close()
	return Parse(f)
}

// Parse decodes config from r. Pure function for tests.
func Parse(r io.Reader) (Config, error) {
	var c Config
	if _, err := toml.NewDecoder(r).Decode(&c); err != nil {
		return Config{}, fmt.Errorf("config parse: %w", err)
	}
	return c, nil
}
