package main

import (
	"errors"
	"testing"

	"github.com/J-V-Studios/nova-pkg/internal/backend"
)

// fakeBackend records calls and returns canned results.
type fakeBackend struct {
	name      string
	available bool
	pkgs      []backend.Package
	searchErr error
}

func (f fakeBackend) Name() string    { return f.name }
func (f fakeBackend) Available() bool { return f.available }
func (f fakeBackend) Search(string) ([]backend.Package, error) {
	return f.pkgs, f.searchErr
}
func (f fakeBackend) Install(string) error             { return nil }
func (f fakeBackend) Remove(string) error              { return nil }
func (f fakeBackend) Update() error                    { return nil }
func (f fakeBackend) List() ([]backend.Package, error) { return f.pkgs, nil }

// withBackends swaps the global backend list for a test and restores it.
func withBackends(t *testing.T, bs []backend.Backend) {
	t.Helper()
	old := backends
	backends = bs
	t.Cleanup(func() { backends = old })
}

func TestSearchAllMergesAndSorts(t *testing.T) {
	withBackends(t, []backend.Backend{
		fakeBackend{name: "b2", available: true, pkgs: []backend.Package{
			{Name: "zeta", Source: "b2"}, {Name: "alpha", Source: "b2"},
		}},
		fakeBackend{name: "b1", available: true, pkgs: []backend.Package{
			{Name: "mid", Source: "b1"},
		}},
		fakeBackend{name: "off", available: false, pkgs: []backend.Package{
			{Name: "hidden", Source: "off"},
		}},
	})

	results := searchAll("x")
	if len(results) != 2 {
		t.Fatalf("got %d backend results, want 2 (unavailable skipped)", len(results))
	}
	// Registration order preserved.
	if results[0].name != "b2" || results[1].name != "b1" {
		t.Fatalf("backend order = %s,%s want b2,b1", results[0].name, results[1].name)
	}
	// Packages sorted by name within a backend.
	got := results[0].pkgs
	if got[0].Name != "alpha" || got[1].Name != "zeta" {
		t.Fatalf("pkg order = %s,%s want alpha,zeta", got[0].Name, got[1].Name)
	}
}

func TestSearchAllErrorIsolated(t *testing.T) {
	withBackends(t, []backend.Backend{
		fakeBackend{name: "bad", available: true, searchErr: errors.New("boom")},
		fakeBackend{name: "good", available: true, pkgs: []backend.Package{{Name: "ok"}}},
	})
	results := searchAll("x")
	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}
	if results[0].err == nil {
		t.Error("bad backend error not propagated")
	}
	if results[1].err != nil || len(results[1].pkgs) != 1 {
		t.Errorf("good backend affected by sibling failure: %+v", results[1])
	}
}

func TestSearchAllEmpty(t *testing.T) {
	withBackends(t, nil)
	if results := searchAll("x"); len(results) != 0 {
		t.Fatalf("got %d results, want 0", len(results))
	}
}
