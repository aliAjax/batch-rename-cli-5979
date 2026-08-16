package fileops

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanFilterDoesNotCorruptExistingEntries(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"a.txt", "IMG_a.txt"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("x"), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	result, err := Scan(dir, "a.txt", false)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Candidates) != 1 || result.Candidates[0].Name != "a.txt" {
		t.Fatalf("unexpected candidates: %#v", result.Candidates)
	}

	found := map[string]bool{}
	for _, entry := range result.Existing {
		found[entry.Name] = true
	}
	if !found["a.txt"] || !found["IMG_a.txt"] {
		t.Fatalf("existing entries were corrupted: %#v", result.Existing)
	}
}
