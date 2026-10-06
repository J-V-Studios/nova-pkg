package backend

import "testing"

func TestParseSearch(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantN    int
		wantName []string // checked when non-nil
	}{
		{name: "empty", input: "", wantN: 0},
		{
			name:     "golden file",
			input:    readTestdata(t, "pacman_search.txt"),
			wantN:    2,
			wantName: []string{"linux", "firefox"},
		},
		{name: "description only", input: "    orphan description\n", wantN: 0},
		{name: "blank lines", input: "\n\n", wantN: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkgs := parseSearch(tt.input)
			if len(pkgs) != tt.wantN {
				t.Fatalf("got %d packages, want %d", len(pkgs), tt.wantN)
			}
			for i, name := range tt.wantName {
				if pkgs[i].Name != name {
					t.Errorf("pkg %d name = %q, want %q", i, pkgs[i].Name, name)
				}
				if pkgs[i].Source != "pacman" {
					t.Errorf("pkg %d source = %q, want pacman", i, pkgs[i].Source)
				}
			}
		})
	}
}

func TestParseSearchFields(t *testing.T) {
	pkgs := parseSearch(readTestdata(t, "pacman_search.txt"))
	if len(pkgs) != 2 {
		t.Fatalf("got %d packages, want 2", len(pkgs))
	}
	want := []Package{
		{Name: "linux", Version: "6.10.2.arch1-1", Description: "The Linux kernel and modules", Source: "pacman"},
		{Name: "firefox", Version: "130.0-1", Description: "Fast, Private & Safe Web Browser", Source: "pacman"},
	}
	for i, w := range want {
		if pkgs[i] != w {
			t.Errorf("pkg %d = %+v, want %+v", i, pkgs[i], w)
		}
	}
}

func TestParsePacmanQ(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		wantN  int
		source string
	}{
		{name: "empty", input: "", wantN: 0, source: "pacman"},
		{name: "golden file", input: readTestdata(t, "pacman_q.txt"), wantN: 2, source: "pacman"},
		{name: "aur source tag", input: "cava 0.10.4-1\n", wantN: 1, source: "aur"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkgs := parsePacmanQ(tt.input, tt.source)
			if len(pkgs) != tt.wantN {
				t.Fatalf("got %d packages, want %d", len(pkgs), tt.wantN)
			}
			for _, p := range pkgs {
				if p.Name == "" || p.Version == "" {
					t.Errorf("empty field in %+v", p)
				}
				if p.Source != tt.source {
					t.Errorf("source = %q, want %q", p.Source, tt.source)
				}
			}
		})
	}
}
