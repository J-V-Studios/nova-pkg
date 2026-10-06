package backend

import "testing"

func TestParseFlatpakSearch(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantN    int
		wantName []string
	}{
		{name: "empty", input: "", wantN: 0},
		{
			// flatpak omits the header when stdout is not a terminal.
			name:     "no header",
			input:    readTestdata(t, "flatpak_search.txt"),
			wantN:    2,
			wantName: []string{"org.cava.Cava", "de.haeckerfelix.Shortwave"},
		},
		{
			// On a terminal flatpak prints a header row; it must be skipped.
			name:     "with header",
			input:    readTestdata(t, "flatpak_search_header.txt"),
			wantN:    1,
			wantName: []string{"org.cava.Cava"},
		},
		{name: "malformed line", input: "justoneword\n", wantN: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkgs := parseFlatpakSearch(tt.input)
			if len(pkgs) != tt.wantN {
				t.Fatalf("got %d packages, want %d", len(pkgs), tt.wantN)
			}
			for i, name := range tt.wantName {
				if pkgs[i].Name != name {
					t.Errorf("pkg %d name = %q, want %q", i, pkgs[i].Name, name)
				}
				if pkgs[i].Source != "flatpak" {
					t.Errorf("pkg %d source = %q, want flatpak", i, pkgs[i].Source)
				}
			}
		})
	}
}

func TestParseFlatpakList(t *testing.T) {
	tests := []struct {
		name  string
		input string
		wantN int
	}{
		{name: "empty", input: "", wantN: 0},
		{name: "no header", input: readTestdata(t, "flatpak_list.txt"), wantN: 1},
		{name: "with header", input: readTestdata(t, "flatpak_list_header.txt"), wantN: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkgs := parseFlatpakList(tt.input)
			if len(pkgs) != tt.wantN {
				t.Fatalf("got %d packages, want %d", len(pkgs), tt.wantN)
			}
			for _, p := range pkgs {
				if p.Name == "" || p.Source != "flatpak" {
					t.Errorf("bad package: %+v", p)
				}
			}
		})
	}
}
