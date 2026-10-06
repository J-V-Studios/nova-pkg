package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteDefault(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nova-pkg", "config.toml")
	if err := writeDefault(path); err != nil {
		t.Fatalf("writeDefault: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	// Generated file must parse and yield the default order.
	c, err := Parse(strings.NewReader(string(data)))
	if err != nil {
		t.Fatalf("generated config does not parse: %v", err)
	}
	if len(c.Priority) != 3 || c.Priority[0] != "pacman" {
		t.Fatalf("priority = %v, want [pacman aur flatpak]", c.Priority)
	}
}

func TestLoadFrom(t *testing.T) {
	t.Run("missing file generates default", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "nova-pkg", "config.toml")
		c, err := loadFrom(path)
		if err != nil {
			t.Fatalf("loadFrom: %v", err)
		}
		if len(c.Priority) != 3 || c.Priority[0] != "pacman" {
			t.Fatalf("priority = %v, want [pacman aur flatpak]", c.Priority)
		}
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("config file not generated: %v", err)
		}
	})

	t.Run("existing file wins", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "config.toml")
		if err := os.WriteFile(path, []byte(`priority = ["aur"]`), 0o644); err != nil {
			t.Fatal(err)
		}
		c, err := loadFrom(path)
		if err != nil {
			t.Fatalf("loadFrom: %v", err)
		}
		if len(c.Priority) != 1 || c.Priority[0] != "aur" {
			t.Fatalf("priority = %v, want [aur]", c.Priority)
		}
	})

	t.Run("corrupt file errors", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "config.toml")
		if err := os.WriteFile(path, []byte("priority = ["), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := loadFrom(path); err == nil {
			t.Fatal("want parse error, got nil")
		}
	})

	t.Run("unwritable dir errors", func(t *testing.T) {
		if os.Getuid() == 0 {
			t.Skip("root ignores permission bits")
		}
		dir := t.TempDir()
		if err := os.Chmod(dir, 0o555); err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(dir, "sub", "config.toml")
		if _, err := loadFrom(path); err == nil {
			t.Fatal("want write error, got nil")
		}
	})
}

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []string
		wantErr bool
	}{
		{name: "empty", input: "", want: nil},
		{name: "priority", input: `priority = ["flatpak", "pacman"]`, want: []string{"flatpak", "pacman"}},
		{name: "bad toml", input: "priority = [", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := Parse(strings.NewReader(tt.input))
			if tt.wantErr {
				if err == nil {
					t.Fatal("want error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			if len(c.Priority) != len(tt.want) {
				t.Fatalf("priority = %v, want %v", c.Priority, tt.want)
			}
			for i := range tt.want {
				if c.Priority[i] != tt.want[i] {
					t.Errorf("priority[%d] = %q, want %q", i, c.Priority[i], tt.want[i])
				}
			}
		})
	}
}
