# Bug 是什么

`undo.Load` 读取空撤销记录时返回 `(nil, nil)`，调用方未做 nil 检查就直接读取 `journal.Directory`，最终触发 nil pointer dereference。

# 如何触发

创建一个空文件后调用：

```bash
go test ./internal/undo -run TestLoadEmptyJournalReturnsError -count=20
```

# 错误信息

```text
expected an error for an empty journal
```
