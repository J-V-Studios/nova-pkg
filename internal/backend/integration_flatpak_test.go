//go:build integration && flatpak

package backend

import (
	"os/exec"
	"testing"
)

func TestFlatpakSearchIntegration(t *testing.T) {
	if _, err := exec.LookPath("flatpak"); err != nil {
		t.Skip("flatpak not installed")
	}
	pkgs, err := Flatpak{}.Search("shortwave")
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(pkgs) == 0 {
		t.Fatal("no results for 'shortwave'")
	}
	for _, p := range pkgs {
		if p.Name == "" {
			t.Errorf("empty app ID in %+v", p)
		}
		if p.Source != "flatpak" {
			t.Errorf("source = %q, want flatpak", p.Source)
		}
	}
}

func TestFlatpakRemoteInfoIntegration(t *testing.T) {
	if _, err := exec.LookPath("flatpak"); err != nil {
		t.Skip("flatpak not installed")
	}
	// Read-only check that a known app ID resolves on flathub.
	out, err := exec.Command("flatpak", "remote-info", "flathub",
		"de.haeckerfelix.Shortwave").CombinedOutput()
	if err != nil {
		t.Fatalf("remote-info: %v\n%s", err, out)
	}
}

func TestFlatpakListIntegration(t *testing.T) {
	if _, err := exec.LookPath("flatpak"); err != nil {
		t.Skip("flatpak not installed")
	}
	// Empty is fine: a clean container has no flatpak apps installed.
	pkgs, err := Flatpak{}.List()
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	for _, p := range pkgs {
		if p.Name == "" || p.Source != "flatpak" {
			t.Errorf("bad package: %+v", p)
		}
	}
}
