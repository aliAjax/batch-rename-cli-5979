package main

import (
	"fmt"
	"os"

	"batch-rename-cli/internal/fileops"
	"batch-rename-cli/internal/planner"
	"batch-rename-cli/internal/preview"
	"batch-rename-cli/internal/rules"
)

func runPreview(args []string) int {
	options, showHelp, err := parseRenameFlags("preview", args, false)
	if showHelp {
		printPreviewUsage(os.Stdout)
		return 0
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "参数错误: %v\n\n", err)
		printPreviewUsage(os.Stderr)
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
	return 0
}
