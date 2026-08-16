# Bug 是什么

`--replace` 的新关键词未校验路径分隔符。生成的目标名可以包含 `/` 或 `\`，预览会显示为有效变更，实际重命名时目标路径解析到不存在的子目录，导致 apply 失败，且错误信息不能直接指出是哪个替换规则。

# 如何触发

```bash
go test ./internal/rules -run TestConfigRejectsPathSeparatorInReplacement -count=20
```

# 错误信息

```text
expected replacement new value containing a path separator to be rejected
```
