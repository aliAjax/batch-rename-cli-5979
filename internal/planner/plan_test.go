package planner

import (
	"testing"

	"batch-rename-cli/internal/fileops"
	"batch-rename-cli/internal/rules"
)

func TestBuildDetectsConflicts(t *testing.T) {
	engine, err := rules.NewEngine(rules.Config{Prefix: "p-"})
	if err != nil {
		t.Fatal(err)
	}
	entries := []fileops.Entry{
		{RelativePath: "a.txt", Name: "a.txt"},
		{RelativePath: "b.txt", Name: "b.txt"},
	}

	plan := Build("/tmp", entries, entries, engine)
	if got := len(plan.Changes()); got != 2 {
		t.Fatalf("expected 2 changes, got %d", got)
	}
	if got := len(plan.Conflicts()); got != 0 {
		t.Fatalf("expected 0 conflicts, got %d", got)
	}
}

func TestBuildDetectsDuplicateTarget(t *testing.T) {
	engine, err := rules.NewEngine(rules.Config{
		Replacements: []rules.Replacement{
			{Old: "a", New: "x"},
			{Old: "b", New: "x"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	entries := []fileops.Entry{
		{RelativePath: "a.txt", Name: "a.txt"},
		{RelativePath: "b.txt", Name: "b.txt"},
	}

	plan := Build("/tmp", entries, entries, engine)
	if got := len(plan.Changes()); got != 1 {
		t.Fatalf("expected 1 change, got %d", got)
	}
	if got := len(plan.Conflicts()); got != 1 {
		t.Fatalf("expected 1 conflict, got %d", got)
	}
	if got := plan.Items[1].Status; got != StatusConflictDuplicate {
		t.Fatalf("expected duplicate conflict for b.txt, got %s", got)
	}
}

func TestBuildDetectsConflictFromNonCandidateFile(t *testing.T) {
	engine, err := rules.NewEngine(rules.Config{Prefix: "IMG_"})
	if err != nil {
		t.Fatal(err)
	}
	candidates := []fileops.Entry{
		{RelativePath: "photo.jpg", Name: "photo.jpg"},
	}
	existing := []fileops.Entry{
		{RelativePath: "photo.jpg", Name: "photo.jpg"},
		{RelativePath: "IMG_photo.jpg", Name: "IMG_photo.jpg"},
	}

	plan := Build("/tmp", candidates, existing, engine)
	if got := len(plan.Conflicts()); got != 1 {
		t.Fatalf("expected 1 conflict from non-candidate file, got %d", got)
	}
	if got := plan.Items[0].Status; got != StatusConflictExisting {
		t.Fatalf("expected existing conflict, got %s", got)
	}
}
