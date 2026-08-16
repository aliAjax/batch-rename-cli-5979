# Bug 是什么

批量重命名应用并发执行时，`Apply` 使用无同步的共享游标分配操作，多个 goroutine 同时读取和写入 `next`，导致 data race、重复处理、漏处理，最终一部分目标文件没有被重命名。

# 如何触发

```bash
go test -race ./internal/fileops -run TestApplyRenamesAllOperations -count=20
```

# 错误信息

`WARNING: DATA RACE`，随后测试失败并提示若干目标文件不存在，例如：

```text
target .../renamed-004.txt missing after apply: stat ...: no such file or directory
```
