package rules

import "fmt"

// SequenceConfig controls sequential renaming.
type SequenceConfig struct {
	Start  int
	Step   int
	Digits int
}

// Config is the parsed rule set for one rename command.
type Config struct {
	Prefix       string
	Replacements []Replacement
	Sequence     *SequenceConfig
}

// Validate rejects empty rule sets and ambiguous or invalid combinations.
func (c Config) Validate() error {
	hasPrefix := c.Prefix != ""
	hasReplacements := len(c.Replacements) > 0
	hasSequence := c.Sequence != nil

	if !hasPrefix && !hasReplacements && !hasSequence {
		return fmt.Errorf("至少需要指定一种规则：--prefix、--replace 或 --sequence")
	}

	if hasSequence && (hasPrefix || hasReplacements) {
		return fmt.Errorf("--sequence 不能与 --prefix/--replace 同时使用")
	}

	for _, replacement := range c.Replacements {
		if replacement.Old == "" {
			return fmt.Errorf("--replace 的旧关键词不能为空")
		}
	}

	if hasSequence {
		if c.Sequence.Digits < 1 || c.Sequence.Digits > 20 {
			return fmt.Errorf("--digits 必须在 1 到 20 之间")
		}
		if c.Sequence.Step == 0 {
			return fmt.Errorf("--step 不能为 0")
		}
	}

	return nil
}
