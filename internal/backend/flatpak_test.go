package backend

import "testing"

func TestParseFlatpakSearch(t *testing.T) {
	// flatpak omits the header when stdout is not a terminal.
	sample := "Cava\tConsole-based Audio Visualizer\torg.cava.Cava\t0.10.4\n" +
		"Shortwave\tListen to internet radio\tde.haeckerfelix.Shortwave\t3.2.0\n"
	pkgs := parseFlatpakSearch(sample)
	if len(pkgs) != 2 {
		t.Fatalf("got %d packages, want 2", len(pkgs))
	}
	if pkgs[0].Name != "org.cava.Cava" || pkgs[0].Version != "0.10.4" ||
		pkgs[0].Description != "Cava - Console-based Audio Visualizer" || pkgs[0].Source != "flatpak" {
		t.Errorf("unexpected first package: %+v", pkgs[0])
	}
	if pkgs[1].Name != "de.haeckerfelix.Shortwave" {
		t.Errorf("unexpected second package: %+v", pkgs[1])
	}
}

func TestParseFlatpakSearchWithHeader(t *testing.T) {
	// On a terminal flatpak prints a header row; it must be skipped.
	sample := "Name\tDescription\tApplication ID\tVersion\n" +
		"Cava\tConsole-based Audio Visualizer\torg.cava.Cava\t0.10.4\n"
	pkgs := parseFlatpakSearch(sample)
	if len(pkgs) != 1 {
		t.Fatalf("got %d packages, want 1", len(pkgs))
	}
	if pkgs[0].Name != "org.cava.Cava" {
		t.Errorf("unexpected package: %+v", pkgs[0])
	}
}

func TestParseFlatpakSearchEmpty(t *testing.T) {
	if pkgs := parseFlatpakSearch(""); len(pkgs) != 0 {
		t.Fatalf("got %d packages, want 0", len(pkgs))
	}
}

func TestParseFlatpakList(t *testing.T) {
	sample := "Shortwave\tde.haeckerfelix.Shortwave\t5.1.0\n"
	pkgs := parseFlatpakList(sample)
	if len(pkgs) != 1 {
		t.Fatalf("got %d packages, want 1", len(pkgs))
	}
	if pkgs[0].Name != "de.haeckerfelix.Shortwave" || pkgs[0].Version != "5.1.0" ||
		pkgs[0].Description != "Shortwave" || pkgs[0].Source != "flatpak" {
		t.Errorf("unexpected package: %+v", pkgs[0])
	}
}

func TestParseFlatpakListWithHeader(t *testing.T) {
	sample := "Name\tApplication ID\tVersion\n" +
		"Shortwave\tde.haeckerfelix.Shortwave\t5.1.0\n"
	pkgs := parseFlatpakList(sample)
	if len(pkgs) != 1 {
		t.Fatalf("got %d packages, want 1", len(pkgs))
	}
}
