package fileops

import (
	"fmt"
	"os"
	"sync"
)

// Operation is a single planned rename.
type Operation struct {
	Source string
	Target string
}

// Rename renames source to target without overwriting an existing target.
func Rename(source, target string) error {
	if _, err := os.Lstat(target); err == nil {
		return fmt.Errorf("目标已存在，已跳过: %s", target)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("检查目标 %q: %w", target, err)
	}
	if err := os.Rename(source, target); err != nil {
		return fmt.Errorf("重命名 %q -> %q: %w", source, target, err)
	}
	return nil
}

// Apply runs the renames concurrently (one goroutine per operation) and
// invokes observer after each attempt. Each operation is dispatched exactly
// once via its slice index, so no work is skipped or duplicated.
//
// Observer calls are serialized with a mutex: the observer does not need to
// be safe for concurrent use. A failed rename is reported via observer but
// does not abort the remaining operations.
func Apply(operations []Operation, observer func(index int, operation Operation, err error)) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i := range operations {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			operation := operations[index]
			err := Rename(operation.Source, operation.Target)
			if observer != nil {
				mu.Lock()
				observer(index, operation, err)
				mu.Unlock()
			}
		}(i)
	}
	wg.Wait()
	return nil
}
