package undo

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadEmptyJournalReturnsError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "empty.json")
	if err := os.WriteFile(path, []byte(""), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Fatal("expected an error for an empty journal")
	}
}
