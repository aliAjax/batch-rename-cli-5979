package rules

import (
	"fmt"
	"strings"
)

// Engine applies a validated Config to file names.
type Engine struct {
	rules []Rule
}

// NewEngine validates the configuration and builds its rule pipeline.
func NewEngine(config Config) (*Engine, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	pipeline := make([]Rule, 0, 3)
	if config.Prefix != "" {
		pipeline = append(pipeline, PrefixRule{Prefix: config.Prefix})
	}
	if len(config.Replacements) > 0 {
		pipeline = append(pipeline, ReplaceRule{Replacements: config.Replacements})
	}
	if config.Sequence != nil {
		pipeline = append(pipeline, SequenceRule{
			Start:  config.Sequence.Start,
			Step:   config.Sequence.Step,
			Digits: config.Sequence.Digits,
		})
	}

	return &Engine{rules: pipeline}, nil
}

// Rename applies the configured rules to name. index is used only by the
// sequence rule and is zero-based within a batch.
func (e *Engine) Rename(name string, index int) string {
	for _, rule := range e.rules {
		name = rule.Apply(name, index)
	}
	return name
}

// Description returns a compact human-readable summary of the active rules.
func (e *Engine) Description() string {
	if len(e.rules) == 0 {
		return "none"
	}

	parts := make([]string, 0, len(e.rules))
	for _, rule := range e.rules {
		switch r := rule.(type) {
		case PrefixRule:
			parts = append(parts, fmt.Sprintf("prefix=%q", r.Prefix))
		case ReplaceRule:
			pairs := make([]string, 0, len(r.Replacements))
			for _, replacement := range r.Replacements {
				pairs = append(pairs, fmt.Sprintf("%q->%q", replacement.Old, replacement.New))
			}
			parts = append(parts, "replace=["+strings.Join(pairs, ", ")+"]")
		case SequenceRule:
			parts = append(parts, fmt.Sprintf("sequence start=%d step=%d digits=%d", r.Start, r.Step, r.Digits))
		}
	}
	return strings.Join(parts, " ")
}
