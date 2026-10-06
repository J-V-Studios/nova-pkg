package backend

import (
	"os"
	"path/filepath"
	"testing"
)

// readTestdata loads a golden sample file from testdata/.
func readTestdata(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("read testdata/%s: %v", name, err)
	}
	return string(data)
}
