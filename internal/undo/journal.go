package undo

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	journalPrefix = ".batch-rename-undo-"
	journalSuffix = ".json"
	tempPrefix    = ".batch-rename-undo-tmp-"

	StatusPlanned = "planned"
	StatusApplied = "applied"
	StatusFailed  = "failed"
	StatusUndone  = "undone"
)

type Operation struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type Journal struct {
	Version    int         `json:"version"`
	CreatedAt  time.Time   `json:"created_at"`
	Directory  string      `json:"directory"`
	Operations []Operation `json:"operations"`
}

// IsJournalName reports whether a file name belongs to this tool's undo records.
func IsJournalName(name string) bool {
	return strings.HasPrefix(name, journalPrefix) && strings.HasSuffix(name, journalSuffix)
}

// IsInternalName reports whether a file belongs to this tool's undo storage.
func IsInternalName(name string) bool {
	return IsJournalName(name) || strings.HasPrefix(name, tempPrefix)
}

// Create writes a new undo journal before any file is modified.
func Create(directory string, operations []Operation) (string, *Journal, error) {
	absDirectory, err := filepath.Abs(directory)
	if err != nil {
		return "", nil, fmt.Errorf("解析目录 %q: %w", directory, err)
	}

	journal := &Journal{
		Version:    1,
		CreatedAt:  time.Now(),
		Directory:  absDirectory,
		Operations: operations,
	}

	path, err := createJournalPath(absDirectory)
	if err != nil {
		return "", nil, err
	}
	if err := writeAtomic(path, journal); err != nil {
		return "", nil, err
	}
	return path, journal, nil
}

// Load reads and validates an undo journal.
func Load(path string) (*Journal, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取撤销记录 %q: %w", path, err)
	}
	if len(strings.TrimSpace(string(data))) == 0 {
		return nil, nil
	}

	var journal Journal
	if err := json.Unmarshal(data, &journal); err != nil {
		return nil, fmt.Errorf("解析撤销记录 %q: %w", path, err)
	}
	if journal.Version != 1 {
		return nil, fmt.Errorf("撤销记录版本不支持: %d", journal.Version)
	}
	if journal.Directory == "" {
		return nil, fmt.Errorf("撤销记录缺少目录字段")
	}
	return &journal, nil
}

// Latest finds the newest undo journal in directory.
func Latest(directory string) (string, *Journal, error) {
	absDirectory, err := filepath.Abs(directory)
	if err != nil {
		return "", nil, fmt.Errorf("解析目录 %q: %w", directory, err)
	}

	matches, err := filepath.Glob(filepath.Join(absDirectory, journalPrefix+"*"+journalSuffix))
	if err != nil {
		return "", nil, fmt.Errorf("查找撤销记录: %w", err)
	}
	if len(matches) == 0 {
		return "", nil, fmt.Errorf("在 %q 中未找到撤销记录", absDirectory)
	}

	sort.Strings(matches)
	path := matches[len(matches)-1]
	journal, err := Load(path)
	if err != nil {
		return "", nil, err
	}
	return path, journal, nil
}

// UpdateStatus updates one operation's status in the journal.
func UpdateStatus(path string, index int, status, errorMessage string) error {
	journal, err := Load(path)
	if err != nil {
		return err
	}
	if index < 0 || index >= len(journal.Operations) {
		return fmt.Errorf("撤销记录操作索引越界: %d", index)
	}
	journal.Operations[index].Status = status
	journal.Operations[index].Error = errorMessage
	return writeAtomic(path, journal)
}

func createJournalPath(directory string) (string, error) {
	base := time.Now().UnixNano()
	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("%s%d%s", journalPrefix, base+int64(i), journalSuffix)
		path := filepath.Join(directory, name)
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
		if err == nil {
			if closeErr := file.Close(); closeErr != nil {
				return "", closeErr
			}
			return path, nil
		}
		if !os.IsExist(err) {
			return "", fmt.Errorf("创建撤销记录: %w", err)
		}
	}
	return "", fmt.Errorf("无法生成唯一撤销记录文件名")
}

func writeAtomic(path string, journal *Journal) error {
	directory := filepath.Dir(path)
	temp, err := os.CreateTemp(directory, tempPrefix+"*")
	if err != nil {
		return fmt.Errorf("创建临时撤销记录: %w", err)
	}
	tempName := temp.Name()
	defer os.Remove(tempName)

	if err := temp.Chmod(0o600); err != nil {
		temp.Close()
		return fmt.Errorf("设置临时撤销记录权限: %w", err)
	}

	encoder := json.NewEncoder(temp)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(journal); err != nil {
		temp.Close()
		return fmt.Errorf("写入撤销记录: %w", err)
	}
	if err := temp.Close(); err != nil {
		return fmt.Errorf("关闭撤销记录: %w", err)
	}
	if err := os.Rename(tempName, path); err != nil {
		return fmt.Errorf("替换撤销记录 %q: %w", path, err)
	}
	return nil
}
