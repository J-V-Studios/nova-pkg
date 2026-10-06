package backend

import "testing"

// TestBackendContract checks invariants every backend must satisfy,
// regardless of whether its tool is installed.
func TestBackendContract(t *testing.T) {
	for _, b := range []Backend{Pacman{}, AUR{}, Flatpak{}} {
		t.Run(b.Name(), func(t *testing.T) {
			if b.Name() == "" {
				t.Error("empty Name()")
			}
			// Available must never panic, installed or not.
			_ = b.Available()
		})
	}
}

func TestValidateName(t *testing.T) {
	valid := []string{"cava", "linux", "de.haeckerfelix.Shortwave", "cava-git", "lib32-openal", "foo+bar", "x@y", "a_b"}
	for _, n := range valid {
		if err := ValidateName(n); err != nil {
			t.Errorf("ValidateName(%q) = %v, want nil", n, err)
		}
	}
	// Flag injection, path traversal, shell metachars, empty.
	invalid := []string{"-Syyu", "--upload-pack=evil", "../etc", "a/b", "a;b", "a b", "$(id)", "", "-", "."}
	for _, n := range invalid {
		if err := ValidateName(n); err == nil {
			t.Errorf("ValidateName(%q) = nil, want error", n)
		}
	}
}

// TestSearchMissingToolNoResults: when the underlying tool is absent,
// Search must not crash; empty output maps to no results, no error.
// This matches how pacman/flatpak treat "no matches" (exit 1, no output).
func TestSearchMissingToolNoResults(t *testing.T) {
	if (Pacman{}).Available() {
		t.Skip("pacman present; cannot test missing-tool path")
	}
	pkgs, err := (Pacman{}).Search("x")
	if err != nil {
		t.Errorf("missing tool should not error, got %v", err)
	}
	if len(pkgs) != 0 {
		t.Errorf("got %d packages, want 0", len(pkgs))
	}
}
