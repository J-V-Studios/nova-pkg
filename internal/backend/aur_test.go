package backend

import (
	"encoding/json"
	"testing"
)

// Sample AUR RPC v5 reply, trimmed to two results. No network involved:
// the test decodes a fixed literal and checks parseAURSearch.
func TestParseAURSearch(t *testing.T) {
	sample := `{
		"version": 5,
		"type": "search",
		"resultcount": 2,
		"results": [
			{
				"ID": 1462842,
				"Name": "cava",
				"PackageBaseID": 114751,
				"Version": "0.10.4-1",
				"Description": "Cross-platform Audio Visualizer",
				"URL": "https://github.com/karlstav/cava"
			},
			{
				"ID": 990909,
				"Name": "cava-git",
				"PackageBaseID": 114753,
				"Version": "0.8.3.r1.g8a6b6d3-1",
				"Description": "Cross-platform Audio Visualizer (git version)"
			}
		]
	}`

	var ar aurResponse
	if err := json.Unmarshal([]byte(sample), &ar); err != nil {
		t.Fatalf("unmarshal sample: %v", err)
	}

	pkgs := parseAURSearch(ar)
	if len(pkgs) != 2 {
		t.Fatalf("got %d packages, want 2", len(pkgs))
	}
	if pkgs[0].Name != "cava" || pkgs[0].Version != "0.10.4-1" ||
		pkgs[0].Description != "Cross-platform Audio Visualizer" || pkgs[0].Source != "aur" {
		t.Errorf("unexpected first package: %+v", pkgs[0])
	}
	if pkgs[1].Name != "cava-git" || pkgs[1].Version != "0.8.3.r1.g8a6b6d3-1" {
		t.Errorf("unexpected second package: %+v", pkgs[1])
	}
}

func TestParseAURSearchEmpty(t *testing.T) {
	if pkgs := parseAURSearch(aurResponse{}); len(pkgs) != 0 {
		t.Fatalf("got %d packages, want 0", len(pkgs))
	}
}
