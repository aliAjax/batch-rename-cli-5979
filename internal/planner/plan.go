package planner

import (
	"path"

	"batch-rename-cli/internal/fileops"
	"batch-rename-cli/internal/rules"
)

type Status string

const (
	StatusChanged           Status = "CHANGE"
	StatusNoChange          Status = "NO_CHANGE"
	StatusConflictExisting  Status = "CONFLICT_EXISTING"
	StatusConflictDuplicate Status = "CONFLICT_DUPLICATE"
)

type Item struct {
	Index  int    `json:"index"`
	Source string `json:"source"`
	Target string `json:"target"`
	Status Status `json:"status"`
	Reason string `json:"reason,omitempty"`
}

type Plan struct {
	Directory       string
	RuleDescription string
	Items           []Item
}

// Build computes a deterministic rename plan. entries are the candidate files
// selected by the user filter, while existingEntries are all files in the
// directory and are used to detect target-name conflicts.
func Build(directory string, entries, existingEntries []fileops.Entry, engine *rules.Engine) Plan {
	plan := Plan{
		Directory:       directory,
		RuleDescription: engine.Description(),
		Items:           make([]Item, 0, len(entries)),
	}

	existing := make(map[string]struct{}, len(existingEntries))
	for _, entry := range existingEntries {
		existing[entry.RelativePath] = struct{}{}
	}

	claimed := make(map[string]string, len(entries))
	for index, entry := range entries {
		newName := engine.Rename(entry.Name, index)
		target := joinRelative(path.Dir(entry.RelativePath), newName)
		item := Item{
			Index:  index,
			Source: entry.RelativePath,
			Target: target,
		}

		switch {
		case newName == entry.Name || target == entry.RelativePath:
			item.Status = StatusNoChange
			item.Reason = "名称未发生变化"
		case target == "" || path.Base(target) != newName:
			item.Status = StatusConflictExisting
			item.Reason = "生成的目标名称无效"
		default:
			if _, ok := existing[target]; ok && target != entry.RelativePath {
				item.Status = StatusConflictExisting
				item.Reason = "目标名称已存在于目录中"
			} else if previous, ok := claimed[target]; ok {
				item.Status = StatusConflictDuplicate
				item.Reason = "与 " + previous + " 生成的目标名称重复"
			} else {
				claimed[target] = entry.RelativePath
				item.Status = StatusChanged
			}
		}

		plan.Items = append(plan.Items, item)
	}

	return plan
}

func (p Plan) Changes() []Item {
	return filterItems(p.Items, func(item Item) bool { return item.Status == StatusChanged })
}

func (p Plan) Conflicts() []Item {
	return filterItems(p.Items, func(item Item) bool {
		return item.Status == StatusConflictExisting || item.Status == StatusConflictDuplicate
	})
}

func (p Plan) NoChanges() []Item {
	return filterItems(p.Items, func(item Item) bool { return item.Status == StatusNoChange })
}

func filterItems(items []Item, keep func(Item) bool) []Item {
	result := make([]Item, 0, len(items))
	for _, item := range items {
		if keep(item) {
			result = append(result, item)
		}
	}
	return result
}

func joinRelative(dir, base string) string {
	if dir == "." || dir == "" {
		return base
	}
	return dir + "/" + base
}
