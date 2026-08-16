package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestApplyKeepsJournalAfterSuccess(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}

	if code := runApply([]string{dir, "--prefix", "IMG_"}); code != 0 {
		t.Fatalf("runApply exited with %d", code)
	}

	matches, err := filepath.Glob(filepath.Join(dir, ".batch-rename-undo-*.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 1 {
		t.Fatalf("expected one undo journal after successful apply, got %d", len(matches))
	}
}
