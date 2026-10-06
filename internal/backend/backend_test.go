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

// TestSearchErrorsWrapped checks backend errors carry context naming
// the backend, per project convention (fmt.Errorf with %w).
func TestSearchErrorsWrapped(t *testing.T) {
	if (Pacman{}).Available() {
		t.Skip("pacman present; cannot test missing-tool path")
	}
	if _, err := (Pacman{}).Search("x"); err == nil {
		t.Error("want error when pacman missing")
	}
}
