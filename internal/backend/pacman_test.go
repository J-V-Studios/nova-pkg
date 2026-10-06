package backend

import "testing"

func TestParseSearch(t *testing.T) {
	sample := `core/linux 6.10.2.arch1-1 [installed]
    The Linux kernel and modules
extra/firefox 130.0-1
    Fast, Private & Safe Web Browser
`
	pkgs := parseSearch(sample)
	if len(pkgs) != 2 {
		t.Fatalf("got %d packages, want 2", len(pkgs))
	}

	if pkgs[0].Name != "linux" || pkgs[0].Version != "6.10.2.arch1-1" ||
		pkgs[0].Description != "The Linux kernel and modules" || pkgs[0].Source != "pacman" {
		t.Errorf("unexpected first package: %+v", pkgs[0])
	}
	if pkgs[1].Name != "firefox" || pkgs[1].Version != "130.0-1" ||
		pkgs[1].Description != "Fast, Private & Safe Web Browser" {
		t.Errorf("unexpected second package: %+v", pkgs[1])
	}
}

func TestParseSearchEmpty(t *testing.T) {
	if pkgs := parseSearch(""); len(pkgs) != 0 {
		t.Fatalf("got %d packages, want 0", len(pkgs))
	}
}

func TestParsePacmanQ(t *testing.T) {
	sample := "linux 6.10.2.arch1-1\nfirefox 130.0-1\n"
	pkgs := parsePacmanQ(sample, "pacman")
	if len(pkgs) != 2 {
		t.Fatalf("got %d packages, want 2", len(pkgs))
	}
	if pkgs[0].Name != "linux" || pkgs[0].Version != "6.10.2.arch1-1" || pkgs[0].Source != "pacman" {
		t.Errorf("unexpected first package: %+v", pkgs[0])
	}
	if pkgs[1].Name != "firefox" || pkgs[1].Version != "130.0-1" {
		t.Errorf("unexpected second package: %+v", pkgs[1])
	}
}

func TestParsePacmanQEmpty(t *testing.T) {
	if pkgs := parsePacmanQ("", "pacman"); len(pkgs) != 0 {
		t.Fatalf("got %d packages, want 0", len(pkgs))
	}
}
