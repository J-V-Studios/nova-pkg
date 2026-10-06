package backend

import (
	"encoding/json"
	"testing"
)

// TestParseAURSearch decodes the golden RPC v5 reply from testdata.
// No network involved.
func TestParseAURSearch(t *testing.T) {
	var ar aurResponse
	if err := json.Unmarshal([]byte(readTestdata(t, "aur_search.json")), &ar); err != nil {
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

// TestAURResponseUnknownFields ensures extra RPC fields (votes,
// popularity, timestamps, ...) do not break decoding.
func TestAURResponseUnknownFields(t *testing.T) {
	sample := `{"version":5,"type":"search","resultcount":1,"results":[
		{"ID":1,"Name":"x","Version":"1.0-1","Description":"d",
		 "NumVotes":42,"Popularity":1.5,"LastModified":1700000000}]}`
	var ar aurResponse
	if err := json.Unmarshal([]byte(sample), &ar); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	pkgs := parseAURSearch(ar)
	if len(pkgs) != 1 || pkgs[0].Name != "x" {
		t.Fatalf("unexpected packages: %+v", pkgs)
	}
}
