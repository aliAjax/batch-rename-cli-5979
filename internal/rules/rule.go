package rules

import (
	"fmt"
	"strings"

	"batch-rename-cli/pkg/pathutil"
)

// Rule describes a single file-name transformation.
type Rule interface {
	Name() string
	Apply(name string, index int) string
}

// PrefixRule prepends a fixed string to the file name.
type PrefixRule struct {
	Prefix string
}

func (r PrefixRule) Name() string {
	return "prefix"
}

func (r PrefixRule) Apply(name string, _ int) string {
	return r.Prefix + name
}

// Replacement is one old/new keyword pair.
type Replacement struct {
	Old string
	New string
}

// ReplaceRule replaces all occurrences of each configured keyword.
type ReplaceRule struct {
	Replacements []Replacement
}

func (r ReplaceRule) Name() string {
	return "replace"
}

func (r ReplaceRule) Apply(name string, _ int) string {
	for _, replacement := range r.Replacements {
		name = strings.ReplaceAll(name, replacement.Old, replacement.New)
	}
	return name
}

// SequenceRule rewrites the file stem to a padded sequence number while
// preserving the original extension.
type SequenceRule struct {
	Start  int
	Step   int
	Digits int
}

func (r SequenceRule) Name() string {
	return "sequence"
}

func (r SequenceRule) Apply(name string, index int) string {
	_, ext := pathutil.SplitExt(name)
	value := r.Start + index*r.Step
	return fmt.Sprintf("%0*d%s", r.Digits, value, ext)
}
