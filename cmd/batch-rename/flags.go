package main

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"batch-rename-cli/internal/rules"
)

type stringList []string

func (s *stringList) String() string {
	return strings.Join(*s, ",")
}

func (s *stringList) Set(value string) error {
	*s = append(*s, value)
	return nil
}

type renameOptions struct {
	Directory     string
	Config        rules.Config
	Filter        string
	Recursive     bool
	SkipConflicts bool
}

type undoOptions struct {
	Directory string
	Journal   string
}

func parseRenameFlags(command string, args []string, includeSkip bool) (renameOptions, bool, error) {
	options := renameOptions{}
	fs := flag.NewFlagSet(command, flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	prefix := fs.String("prefix", "", "添加前缀")
	var replacements stringList
	fs.Var(&replacements, "replace", "替换关键词，格式 old=new，可重复")
	sequence := fs.Bool("sequence", false, "按序号重命名")
	start := fs.Int("start", 1, "序号起始值")
	step := fs.Int("step", 1, "序号步长")
	digits := fs.Int("digits", 3, "序号位数")
	filter := fs.String("filter", "*", "文件 glob 过滤")
	recursive := fs.Bool("recursive", false, "递归处理子目录")
	help := fs.Bool("help", false, "显示帮助")

	var skipConflicts *bool
	if includeSkip {
		skipConflicts = fs.Bool("skip-conflicts", true, "跳过同名冲突")
	}

	flagArgs, positionals, err := collectFlagArgs(fs, args)
	if err != nil {
		return options, false, err
	}
	if err := fs.Parse(flagArgs); err != nil {
		return options, false, err
	}
	if *help {
		return options, true, nil
	}
	if len(positionals) != 1 {
		return options, false, fmt.Errorf("需要一个目录参数")
	}

	options.Directory = positionals[0]
	options.Filter = *filter
	options.Recursive = *recursive
	if includeSkip {
		options.SkipConflicts = *skipConflicts
	}

	parsedReplacements, err := parseReplacements(replacements)
	if err != nil {
		return options, false, err
	}
	options.Config.Prefix = *prefix
	options.Config.Replacements = parsedReplacements
	if *sequence {
		options.Config.Sequence = &rules.SequenceConfig{
			Start:  *start,
			Step:   *step,
			Digits: *digits,
		}
	}

	if err := options.Config.Validate(); err != nil {
		return options, false, err
	}
	return options, false, nil
}

func parseUndoFlags(args []string) (undoOptions, bool, error) {
	options := undoOptions{}
	fs := flag.NewFlagSet("undo", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	journal := fs.String("journal", "", "指定撤销记录文件")
	help := fs.Bool("help", false, "显示帮助")

	flagArgs, positionals, err := collectFlagArgs(fs, args)
	if err != nil {
		return options, false, err
	}
	if err := fs.Parse(flagArgs); err != nil {
		return options, false, err
	}
	if *help {
		return options, true, nil
	}
	if len(positionals) != 1 {
		return options, false, fmt.Errorf("需要一个目录参数")
	}

	options.Directory = positionals[0]
	options.Journal = *journal
	return options, false, nil
}

func parseReplacements(values []string) ([]rules.Replacement, error) {
	replacements := make([]rules.Replacement, 0, len(values))
	for _, value := range values {
		separator := strings.IndexAny(value, ":=")
		if separator <= 0 {
			return nil, fmt.Errorf("--replace 格式应为 old=new 或 old:new，收到 %q", value)
		}
		replacement := rules.Replacement{
			Old: value[:separator],
			New: value[separator+1:],
		}
		if err := rules.ValidateReplacement(replacement); err != nil {
			return nil, err
		}
		replacements = append(replacements, replacement)
	}
	return replacements, nil
}

// collectFlagArgs moves non-flag positionals to the end so directory can be
// written before or after flags without losing flag parsing.
func collectFlagArgs(fs *flag.FlagSet, args []string) ([]string, []string, error) {
	flagArgs := make([]string, 0, len(args))
	positionals := make([]string, 0, 1)

	for index := 0; index < len(args); index++ {
		arg := args[index]
		if arg == "--" {
			positionals = append(positionals, args[index+1:]...)
			break
		}
		if strings.HasPrefix(arg, "--") && len(arg) > 2 {
			if strings.Contains(arg, "=") {
				flagArgs = append(flagArgs, arg)
				continue
			}
			name := strings.TrimPrefix(arg, "--")
			flagValue := fs.Lookup(name)
			flagArgs = append(flagArgs, arg)
			if flagValue != nil && !isBoolFlag(flagValue.Value) {
				if index+1 >= len(args) {
					return nil, nil, fmt.Errorf("选项 %s 缺少参数", arg)
				}
				index++
				flagArgs = append(flagArgs, args[index])
			}
			continue
		}
		if strings.HasPrefix(arg, "-") && len(arg) > 1 {
			flagArgs = append(flagArgs, arg)
			continue
		}
		positionals = append(positionals, arg)
	}

	return flagArgs, positionals, nil
}

func isBoolFlag(value flag.Value) bool {
	booler, ok := value.(interface{ IsBoolFlag() bool })
	return ok && booler.IsBoolFlag()
}
