package fileops

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestApplyRenamesAllOperations(t *testing.T) {
	dir := t.TempDir()
	const count = 120
	operations := make([]Operation, 0, count)
	for i := 0; i < count; i++ {
		source := filepath.Join(dir, fmt.Sprintf("f%03d.txt", i))
		target := filepath.Join(dir, fmt.Sprintf("renamed-%03d.txt", i))
		if err := os.WriteFile(source, []byte(fmt.Sprintf("%d", i)), 0o600); err != nil {
			t.Fatal(err)
		}
		operations = append(operations, Operation{Source: source, Target: target})
	}

	if err := Apply(operations, nil); err != nil {
		t.Fatalf("Apply returned error: %v", err)
	}

	for _, operation := range operations {
		if _, err := os.Stat(operation.Target); err != nil {
			t.Fatalf("target %s missing after apply: %v", operation.Target, err)
		}
	}
}
