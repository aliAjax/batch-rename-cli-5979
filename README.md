# batch-rename-cli

一个使用 Go 标准库实现的批量文件重命名命令行工具。工具会先扫描目录并预览重命名结果，确认后再执行修改；执行前会自动生成撤销记录，发生误操作时可用 `undo` 回滚。

## 项目简介

支持以下重命名规则：

- 添加前缀：`--prefix "IMG_"`
- 替换关键词：`--replace "draft=final"`，可重复指定多个替换
- 按序号重命名：`--sequence --start 1 --step 1 --digits 3`

核心流程：

1. `preview` 只读取目录并显示每个文件的变更、冲突和无需变更状态，不修改文件。
2. `apply` 先写入 `.batch-rename-undo-*.json` 撤销记录，再执行重命名；同名或重复目标会自动跳过。
3. `undo` 根据最新撤销记录逆序回滚已经执行的重命名。

## 目录结构

```text
.
├── cmd
│   └── batch-rename
│       ├── main.go       # 根命令入口
│       ├── preview.go    # preview 子命令
│       ├── apply.go      # apply 子命令
│       ├── undo.go       # undo 子命令
│       ├── flags.go      # 参数解析
│       ├── usage.go      # 帮助信息
│       └── common.go     # 命令行公共过滤逻辑
├── internal
│   ├── fileops           # 目录扫描和文件重命名
│   ├── planner           # 冲突检测与重命名计划
│   ├── preview           # 预览输出
│   ├── rules             # 规则解析、校验和执行
│   └── undo              # 撤销记录写入、读取和状态更新
├── pkg
│   └── pathutil          # 可复用的路径工具
├── Dockerfile
└── go.mod
```

## 构建

```bash
docker build -t batch-rename-cli .
```

## 运行

将需要处理的目录挂载到容器内的 `/data`：

```bash
docker run --rm -v /absolute/path/to/dir:/data batch-rename-cli preview /data --prefix "IMG_"
docker run --rm -v /absolute/path/to/dir:/data batch-rename-cli apply /data --prefix "IMG_"
docker run --rm -v /absolute/path/to/dir:/data batch-rename-cli undo /data
```

也可以先查看帮助：

```bash
docker run --rm batch-rename-cli help
docker run --rm batch-rename-cli preview --help
docker run --rm batch-rename-cli apply --help
docker run --rm batch-rename-cli undo --help
```

## 命令帮助

```text
批量文件重命名工具

用法:
  batch-rename <command> [选项] <目录>

命令:
  preview  预览重命名结果，不修改文件
  apply    应用重命名，修改前自动生成撤销记录
  undo     根据撤销记录回滚最近一次应用
  help     查看命令帮助

规则:
  --prefix string      添加前缀
  --replace old=new    替换关键词，可重复使用
  --sequence           按序号重命名并保留扩展名
```

`preview` 和 `apply` 共用以下规则选项：

```text
--prefix string
--replace old=new
--sequence
--start int
--step int
--digits int
--filter pattern
--recursive
--skip-conflicts
```

## 示例

假设目录 `/tmp/files` 内有 `a.txt`、`b.txt`：

```bash
# 预览添加前缀
docker run --rm -v /tmp/files:/data batch-rename-cli preview /data --prefix "IMG_"

# 执行添加前缀
docker run --rm -v /tmp/files:/data batch-rename-cli apply /data --prefix "IMG_"

# 撤销最近一次应用
docker run --rm -v /tmp/files:/data batch-rename-cli undo /data
```

替换关键词示例：

```bash
docker run --rm -v /tmp/files:/data batch-rename-cli apply /data \
  --replace "draft=final" --replace "raw=edited"
```

序号重命名示例：

```bash
docker run --rm -v /tmp/files:/data batch-rename-cli apply /data \
  --sequence --start 10 --step 5 --digits 3
```

## 验证

在开发机没有安装 Go 的情况下，使用 Docker 完成以下验证：

```bash
docker run --rm -v "$PWD":/src -w /src golang:1.23-alpine go test ./...
docker run --rm -v "$PWD":/src -w /src golang:1.23-alpine go vet ./...
docker build -t batch-rename-cli .
```

同时在临时目录真实创建文件，验证了前缀、替换、序号、冲突跳过和撤销流程。
