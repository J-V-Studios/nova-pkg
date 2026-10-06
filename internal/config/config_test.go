package config

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []string
		wantErr bool
	}{
		{name: "empty", input: "", want: nil},
		{name: "priority", input: `priority = ["flatpak", "pacman"]`, want: []string{"flatpak", "pacman"}},
		{name: "bad toml", input: "priority = [", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := Parse(strings.NewReader(tt.input))
			if tt.wantErr {
				if err == nil {
					t.Fatal("want error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			if len(c.Priority) != len(tt.want) {
				t.Fatalf("priority = %v, want %v", c.Priority, tt.want)
			}
			for i := range tt.want {
				if c.Priority[i] != tt.want[i] {
					t.Errorf("priority[%d] = %q, want %q", i, c.Priority[i], tt.want[i])
				}
			}
		})
	}
}
