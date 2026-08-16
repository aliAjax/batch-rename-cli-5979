package main

import (
	"batch-rename-cli/internal/fileops"
	"batch-rename-cli/internal/undo"
)

func excludeJournalEntries(entries []fileops.Entry) []fileops.Entry {
	filtered := make([]fileops.Entry, 0, len(entries))
	for _, entry := range entries {
		if undo.IsInternalName(entry.Name) {
			continue
		}
		filtered = append(filtered, entry)
	}
	return filtered
}
