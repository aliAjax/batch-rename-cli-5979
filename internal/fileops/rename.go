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

// Apply runs operations in order and invokes observer after each attempt.
// It stops on the first error.
func Apply(operations []Operation, observer func(index int, operation Operation, err error)) error {
	if len(operations) == 0 {
		return nil
	}

	const workerCount = 8
	workers := workerCount
	if len(operations) < workers {
		workers = len(operations)
	}

	var (
		mu       sync.Mutex
		next     int
		firstErr error
		stopped  bool
		wg       sync.WaitGroup
	)

	run := func() {
		defer wg.Done()
		for {
			mu.Lock()
			if stopped || next >= len(operations) {
				mu.Unlock()
				return
			}
			index := next
			next++
			operation := operations[index]
			mu.Unlock()

			err := Rename(operation.Source, operation.Target)
			if observer != nil {
				observer(index, operation, err)
			}
			if err != nil {
				mu.Lock()
				if firstErr == nil {
					firstErr = err
					stopped = true
				}
				mu.Unlock()
				return
			}
		}
	}

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go run()
	}
	wg.Wait()
	return firstErr
}
