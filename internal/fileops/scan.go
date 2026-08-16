package fileops

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
)

// Entry is one file discovered inside the target directory.
type Entry struct {
	RelativePath string
	Name         string
}

// ScanResult separates candidate files from the full file set. The full set
// is needed for conflict detection even when a filter excludes some files.
type ScanResult struct {
	Candidates []Entry
	Existing   []Entry
}

// Scan returns regular files under root, sorted by their slash-separated
// relative path. When recursive is false, only the top level is scanned.
func Scan(root, pattern string, recursive bool) (ScanResult, error) {
	result := ScanResult{}
	if _, err := path.Match(pattern, "probe"); err != nil {
		return result, fmt.Errorf("过滤规则 %q 无效: %w", pattern, err)
	}

	absRoot, err := filepath.Abs(root)
	if err != nil {
		return result, fmt.Errorf("解析目录路径 %q: %w", root, err)
	}

	info, err := os.Stat(absRoot)
	if err != nil {
		return result, fmt.Errorf("访问目录 %q: %w", root, err)
	}
	if !info.IsDir() {
		return result, fmt.Errorf("%q 不是目录", root)
	}

	allEntries := make([]Entry, 0)
	add := func(absPath string, entry fs.DirEntry) error {
		if entry.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(absRoot, absPath)
		if err != nil {
			return fmt.Errorf("计算相对路径 %q: %w", absPath, err)
		}
		allEntries = append(allEntries, Entry{
			RelativePath: filepath.ToSlash(rel),
			Name:         entry.Name(),
		})
		return nil
	}

	if recursive {
		err = filepath.WalkDir(absRoot, func(absPath string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if absPath == absRoot {
				return nil
			}
			return add(absPath, entry)
		})
		if err != nil {
			return result, fmt.Errorf("扫描目录 %q: %w", root, err)
		}
	} else {
		dirEntries, readErr := os.ReadDir(absRoot)
		if readErr != nil {
			return result, fmt.Errorf("读取目录 %q: %w", root, readErr)
		}
		for _, entry := range dirEntries {
			absPath := filepath.Join(absRoot, entry.Name())
			if err := add(absPath, entry); err != nil {
				return result, err
			}
		}
	}

	sort.Slice(allEntries, func(i, j int) bool {
		return allEntries[i].RelativePath < allEntries[j].RelativePath
	})

	candidates := allEntries[:0]
	for _, entry := range allEntries {
		matched, matchErr := path.Match(pattern, entry.Name)
		if matchErr != nil {
			return result, fmt.Errorf("过滤规则 %q 无效: %w", pattern, matchErr)
		}
		if matched {
			candidates = append(candidates, entry)
		}
	}

	result.Candidates = candidates
	result.Existing = allEntries
	return result, nil
}
