package backend

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

// withAURServer points the RPC URL at a test server for the duration
// of one test.
func withAURServer(t *testing.T, handler http.HandlerFunc) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	old := aurRPCURL
	aurRPCURL = srv.URL + "/?v=5"
	t.Cleanup(func() { aurRPCURL = old })
}

func TestAURSearchHTTP(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		withAURServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("arg") != "cava" {
				t.Errorf("query arg = %q, want cava", r.URL.Query().Get("arg"))
			}
			w.Write([]byte(`{"version":5,"type":"search","resultcount":1,
				"results":[{"Name":"cava","Version":"0.10.4-1","Description":"viz"}]}`))
		})
		pkgs, err := AUR{}.Search("cava")
		if err != nil {
			t.Fatalf("search: %v", err)
		}
		if len(pkgs) != 1 || pkgs[0].Name != "cava" {
			t.Fatalf("unexpected packages: %+v", pkgs)
		}
	})

	t.Run("query escaped", func(t *testing.T) {
		withAURServer(t, func(w http.ResponseWriter, r *http.Request) {
			if got := r.URL.Query().Get("arg"); got != "a b&c" {
				t.Errorf("arg = %q, want %q", got, "a b&c")
			}
			w.Write([]byte(`{"version":5,"results":[]}`))
		})
		if _, err := (AUR{}).Search("a b&c"); err != nil {
			t.Fatalf("search: %v", err)
		}
	})

	t.Run("http 500", func(t *testing.T) {
		withAURServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})
		if _, err := (AUR{}).Search("x"); err == nil {
			t.Fatal("want error on 500, got nil")
		}
	})

	t.Run("malformed json", func(t *testing.T) {
		withAURServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{not json`))
		})
		if _, err := (AUR{}).Search("x"); err == nil {
			t.Fatal("want decode error, got nil")
		}
	})

	t.Run("unreachable host", func(t *testing.T) {
		old := aurRPCURL
		aurRPCURL = "http://127.0.0.1:1/?v=5" // nothing listens on port 1
		t.Cleanup(func() { aurRPCURL = old })
		if _, err := (AUR{}).Search("x"); err == nil {
			t.Fatal("want connection error, got nil")
		}
	})
}

func TestOutdatedAUR(t *testing.T) {
	foreign := []Package{
		{Name: "old", Version: "1.0-1", Source: "aur"},
		{Name: "fresh", Version: "2.0-1", Source: "aur"},
		{Name: "gone", Version: "1.0-1", Source: "aur"}, // dropped from AUR
	}
	latest := map[string]string{"old": "1.1-1", "fresh": "2.0-1"}
	out := outdatedAUR(foreign, latest)
	if len(out) != 1 || out[0].Name != "old" {
		t.Fatalf("outdated = %+v, want only 'old'", out)
	}
}
