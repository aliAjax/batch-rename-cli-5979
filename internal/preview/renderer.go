package preview

import (
	"fmt"
	"io"

	"batch-rename-cli/internal/planner"
)

// Renderer writes a human-readable rename preview.
type Renderer struct {
	Out io.Writer
}

// Render prints the rule, every scanned file, and a summary.
func (r Renderer) Render(plan planner.Plan) error {
	if r.Out == nil {
		return fmt.Errorf("preview output is nil")
	}

	fmt.Fprintf(r.Out, "目录: %s\n", plan.Directory)
	fmt.Fprintf(r.Out, "规则: %s\n", plan.RuleDescription)
	fmt.Fprintf(r.Out, "文件数: %d\n\n", len(plan.Items))

	for _, item := range plan.Items {
		switch item.Status {
		case planner.StatusChanged:
			fmt.Fprintf(r.Out, "  [CHANGE]   %s -> %s\n", item.Source, item.Target)
		case planner.StatusNoChange:
			fmt.Fprintf(r.Out, "  [NO_CHANGE] %s\n", item.Source)
		case planner.StatusConflictExisting:
			fmt.Fprintf(r.Out, "  [SKIP]     %s -> %s (%s)\n", item.Source, item.Target, item.Reason)
		case planner.StatusConflictDuplicate:
			fmt.Fprintf(r.Out, "  [SKIP]     %s -> %s (%s)\n", item.Source, item.Target, item.Reason)
		default:
			fmt.Fprintf(r.Out, "  [UNKNOWN]  %s -> %s\n", item.Source, item.Target)
		}
	}

	fmt.Fprintf(r.Out, "\n汇总: %d 个可重命名, %d 个冲突跳过, %d 个无需变更\n",
		len(plan.Changes()), len(plan.Conflicts()), len(plan.NoChanges()))
	return nil
}
