package main

import (
	"fmt"
	"io"
)

func printRootUsage(w io.Writer) {
	fmt.Fprintln(w, "批量文件重命名工具")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "用法:")
	fmt.Fprintln(w, "  batch-rename <command> [选项] <目录>")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "命令:")
	fmt.Fprintln(w, "  preview  预览重命名结果，不修改文件")
	fmt.Fprintln(w, "  apply    应用重命名，修改前自动生成撤销记录")
	fmt.Fprintln(w, "  undo     根据撤销记录回滚最近一次应用")
	fmt.Fprintln(w, "  help     查看命令帮助")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "规则:")
	fmt.Fprintln(w, "  --prefix string      添加前缀")
	fmt.Fprintln(w, "  --replace old=new    替换关键词，可重复使用")
	fmt.Fprintln(w, "  --sequence           按序号重命名并保留扩展名")
}

func printPreviewUsage(w io.Writer) {
	fmt.Fprintln(w, "用法: batch-rename preview [选项] <目录>")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "选项:")
	fmt.Fprintln(w, "  --prefix string      添加前缀")
	fmt.Fprintln(w, "  --replace old=new    替换关键词，可重复使用")
	fmt.Fprintln(w, "  --sequence           按序号重命名并保留扩展名")
	fmt.Fprintln(w, "  --start int          序号起始值（默认 1）")
	fmt.Fprintln(w, "  --step int           序号步长（默认 1）")
	fmt.Fprintln(w, "  --digits int         序号位数（默认 3）")
	fmt.Fprintln(w, "  --filter pattern     文件 glob 过滤（默认 *）")
	fmt.Fprintln(w, "  --recursive          递归处理子目录")
	fmt.Fprintln(w, "  --help               显示帮助")
}

func printApplyUsage(w io.Writer) {
	fmt.Fprintln(w, "用法: batch-rename apply [选项] <目录>")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "选项:")
	fmt.Fprintln(w, "  --prefix string      添加前缀")
	fmt.Fprintln(w, "  --replace old=new    替换关键词，可重复使用")
	fmt.Fprintln(w, "  --sequence           按序号重命名并保留扩展名")
	fmt.Fprintln(w, "  --start int          序号起始值（默认 1）")
	fmt.Fprintln(w, "  --step int           序号步长（默认 1）")
	fmt.Fprintln(w, "  --digits int         序号位数（默认 3）")
	fmt.Fprintln(w, "  --filter pattern     文件 glob 过滤（默认 *）")
	fmt.Fprintln(w, "  --recursive          递归处理子目录")
	fmt.Fprintln(w, "  --skip-conflicts     跳过同名/重复目标冲突（默认 true）")
	fmt.Fprintln(w, "  --help               显示帮助")
}

func printUndoUsage(w io.Writer) {
	fmt.Fprintln(w, "用法: batch-rename undo [选项] <目录>")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "选项:")
	fmt.Fprintln(w, "  --journal path       指定撤销记录文件；缺省使用目录中最新的记录")
	fmt.Fprintln(w, "  --help               显示帮助")
}

func printCommandUsage(w io.Writer, command string) {
	switch command {
	case "preview":
		printPreviewUsage(w)
	case "apply":
		printApplyUsage(w)
	case "undo":
		printUndoUsage(w)
	default:
		fmt.Fprintf(w, "未知命令 %q\n\n", command)
		printRootUsage(w)
	}
}
