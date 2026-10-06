//go:build integration && pacman

package backend

import (
	"os/exec"
	"testing"
)

func TestPacmanSearchIntegration(t *testing.T) {
	if _, err := exec.LookPath("pacman"); err != nil {
		t.Skip("pacman not installed")
	}
	pkgs, err := Pacman{}.Search("linux")
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(pkgs) == 0 {
		t.Fatal("no results for 'linux'")
	}
	for _, p := range pkgs {
		if p.Name == "" || p.Version == "" {
			t.Errorf("empty field in %+v", p)
		}
		if p.Source != "pacman" {
			t.Errorf("source = %q, want pacman", p.Source)
		}
	}
}

func TestPacmanNoMatchIntegration(t *testing.T) {
	if _, err := exec.LookPath("pacman"); err != nil {
		t.Skip("pacman not installed")
	}
	pkgs, err := Pacman{}.Search("nova-pkg-no-such-package-zzz")
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(pkgs) != 0 {
		t.Fatalf("got %d results, want 0", len(pkgs))
	}
}

func TestPacmanListIntegration(t *testing.T) {
	if _, err := exec.LookPath("pacman"); err != nil {
		t.Skip("pacman not installed")
	}
	pkgs, err := Pacman{}.List()
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(pkgs) == 0 {
		t.Fatal("no installed packages")
	}
	for _, p := range pkgs {
		if p.Name == "" || p.Version == "" {
			t.Errorf("empty field in %+v", p)
		}
	}
}
