package main

import (
	"testing"

	"batch-rename-cli/internal/fileops"
)

func TestExcludeJournalEntriesDoesNotReuseBackingArray(t *testing.T) {
	entries := []fileops.Entry{
		{RelativePath: ".batch-rename-undo-1.json", Name: ".batch-rename-undo-1.json"},
		{RelativePath: "a.txt", Name: "a.txt"},
	}
	_ = excludeJournalEntries(entries)
	if entries[0].Name != ".batch-rename-undo-1.json" {
		t.Fatalf("original entries were mutated: %#v", entries)
	}
}
