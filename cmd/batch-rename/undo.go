package main

import (
	"fmt"
	"os"
	"path/filepath"

	"batch-rename-cli/internal/fileops"
	"batch-rename-cli/internal/undo"
)

func runUndo(args []string) int {
	options, showHelp, err := parseUndoFlags(args)
	if showHelp {
		printUndoUsage(os.Stdout)
		return 0
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "参数错误: %v\n\n", err)
		printUndoUsage(os.Stderr)
		return 2
	}

	journalPath := options.Journal
	journal := (*undo.Journal)(nil)
	if journalPath == "" {
		journalPath, journal, err = undo.Latest(options.Directory)
	} else {
		journal, err = undo.Load(journalPath)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "读取撤销记录失败: %v\n", err)
		return 1
	}
	if journal == nil {
		fmt.Fprintf(os.Stderr, "读取撤销记录失败: 记录内容为空\n")
		return 1
	}

	fmt.Printf("撤销记录: %s\n", journalPath)
	fmt.Printf("记录目录: %s\n", journal.Directory)

	undone := 0
	skipped := 0
	for index := len(journal.Operations) - 1; index >= 0; index-- {
		operation := journal.Operations[index]
		if operation.Status != undo.StatusApplied {
			fmt.Printf("  跳过未应用操作: %s -> %s\n", operation.Source, operation.Target)
			skipped++
			continue
		}

		current := filepath.Join(options.Directory, filepath.FromSlash(operation.Target))
		original := filepath.Join(options.Directory, filepath.FromSlash(operation.Source))
		if err := fileops.Rename(current, original); err != nil {
			fmt.Printf("  撤销跳过: %s -> %s (%v)\n", operation.Target, operation.Source, err)
			skipped++
			continue
		}
		if updateErr := undo.UpdateStatus(journalPath, index, undo.StatusUndone, ""); updateErr != nil {
			fmt.Fprintf(os.Stderr, "更新撤销记录失败: %v\n", updateErr)
			return 1
		}
		undone++
		fmt.Printf("  已撤销: %s -> %s\n", operation.Target, operation.Source)
	}

	fmt.Printf("\n撤销完成: 成功 %d 个，跳过 %d 个。\n", undone, skipped)
	return 0
}
