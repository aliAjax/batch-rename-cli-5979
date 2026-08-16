package main

import (
	"fmt"
	"os"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) == 0 {
		printRootUsage(os.Stderr)
		return 2
	}

	switch args[0] {
	case "preview":
		return runPreview(args[1:])
	case "apply":
		return runApply(args[1:])
	case "undo":
		return runUndo(args[1:])
	case "help", "--help", "-h":
		if len(args) > 1 {
			printCommandUsage(os.Stdout, args[1])
		} else {
			printRootUsage(os.Stdout)
		}
		return 0
	default:
		fmt.Fprintf(os.Stderr, "未知命令 %q\n\n", args[0])
		printRootUsage(os.Stderr)
		return 2
	}
}
