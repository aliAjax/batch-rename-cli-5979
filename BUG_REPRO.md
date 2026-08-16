# Bug 是什么

扫描时把候选切片复用为 `allEntries[:0]`，过滤时又把 `excludeJournalEntries` 的返回切片复用为入参切片，导致非候选文件集合和原始 entries 被覆盖。冲突检测使用的 `Existing` 因此漏掉文件，可能把本应冲突的改名误判为可执行。

# 如何触发

```bash
go test ./internal/fileops -run TestScanFilterDoesNotCorruptExistingEntries
go test ./cmd/batch-rename -run TestExcludeJournalEntriesDoesNotReuseBackingArray
```

# 错误信息

```text
existing entries were corrupted: [{a.txt a.txt} {a.txt a.txt}]
original entries were mutated: [{a.txt a.txt} {a.txt a.txt}]
```
