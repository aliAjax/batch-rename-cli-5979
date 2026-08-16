package main

import (
	"fmt"
	"os"
	"path/filepath"

	"batch-rename-cli/internal/fileops"
	"batch-rename-cli/internal/planner"
	"batch-rename-cli/internal/preview"
	"batch-rename-cli/internal/rules"
	"batch-rename-cli/internal/undo"
)

func runApply(args []string) int {
	options, showHelp, err := parseRenameFlags("apply", args, true)
	if showHelp {
		printApplyUsage(os.Stdout)
		return 0
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "参数错误: %v\n\n", err)
		printApplyUsage(os.Stderr)
		return 2
	}

	scanResult, err := fileops.Scan(options.Directory, options.Filter, options.Recursive)
	if err != nil {
		fmt.Fprintf(os.Stderr, "扫描失败: %v\n", err)
		return 1
	}
	entries := excludeJournalEntries(scanResult.Candidates)
	existingEntries := excludeJournalEntries(scanResult.Existing)

	engine, err := rules.NewEngine(options.Config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "规则错误: %v\n", err)
		return 2
	}
	plan := planner.Build(options.Directory, entries, existingEntries, engine)
	if err := (preview.Renderer{Out: os.Stdout}).Render(plan); err != nil {
		fmt.Fprintf(os.Stderr, "预览失败: %v\n", err)
		return 1
	}

	changes := plan.Changes()
	if len(changes) == 0 {
		fmt.Println("\n没有需要执行的变更。")
		return 0
	}

	fileOperations := make([]fileops.Operation, 0, len(changes))
	undoOperations := make([]undo.Operation, 0, len(changes))
	for _, item := range changes {
		source := filepath.Join(options.Directory, filepath.FromSlash(item.Source))
		target := filepath.Join(options.Directory, filepath.FromSlash(item.Target))
		fileOperations = append(fileOperations, fileops.Operation{Source: source, Target: target})
		undoOperations = append(undoOperations, undo.Operation{
			Source: item.Source,
			Target: item.Target,
			Status: undo.StatusPlanned,
		})
	}

	journalPath, _, err := undo.Create(options.Directory, undoOperations)
	if err != nil {
		fmt.Fprintf(os.Stderr, "创建撤销记录失败: %v\n", err)
		return 1
	}

	fmt.Printf("\n撤销记录: %s\n", journalPath)
	defer os.Remove(journalPath)
	applied := 0
	applyErr := fileops.Apply(fileOperations, func(index int, operation fileops.Operation, opErr error) {
		if opErr != nil {
			if updateErr := undo.UpdateStatus(journalPath, index, undo.StatusFailed, opErr.Error()); updateErr != nil {
				fmt.Fprintf(os.Stderr, "更新撤销记录失败: %v\n", updateErr)
			}
			fmt.Fprintf(os.Stderr, "应用失败: %v\n", opErr)
			return
		}
		if updateErr := undo.UpdateStatus(journalPath, index, undo.StatusApplied, ""); updateErr != nil {
			fmt.Fprintf(os.Stderr, "更新撤销记录失败: %v\n", updateErr)
			return
		}
		applied++
		fmt.Printf("  已重命名: %s -> %s\n", operation.Source, operation.Target)
	})
	if applyErr != nil {
		fmt.Fprintf(os.Stderr, "\n应用中断: %v\n", applyErr)
		return 1
	}

	fmt.Printf("\n应用完成: 成功 %d 个，跳过冲突 %d 个。\n", applied, len(plan.Conflicts()))
	return 0
}
