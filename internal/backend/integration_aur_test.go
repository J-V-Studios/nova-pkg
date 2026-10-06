//go:build integration && aur

package backend

import "testing"

func TestAURSearchIntegration(t *testing.T) {
	pkgs, err := AUR{}.Search("cava")
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(pkgs) == 0 {
		t.Fatal("no results for 'cava'")
	}
	for _, p := range pkgs {
		if p.Name == "" || p.Version == "" {
			t.Errorf("empty field in %+v", p)
		}
		if p.Source != "aur" {
			t.Errorf("source = %q, want aur", p.Source)
		}
	}
}

func TestAURNoMatchIntegration(t *testing.T) {
	pkgs, err := AUR{}.Search("nova-pkg-no-such-package-zzz")
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(pkgs) != 0 {
		t.Fatalf("got %d results, want 0", len(pkgs))
	}
}
