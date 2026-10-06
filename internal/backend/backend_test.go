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
