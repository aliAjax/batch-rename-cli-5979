# Bug 是什么

`apply` 在创建撤销记录后无条件 `defer os.Remove(journalPath)`，成功执行后仍会把撤销记录删掉，导致随后的 `undo` 找不到记录。

# 如何触发

```bash
go test ./cmd/batch-rename -run TestApplyKeepsJournalAfterSuccess -count=20
```

# 错误信息

```text
expected one undo journal after successful apply, got 0
```
