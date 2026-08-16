package pathutil

import "path/filepath"

// SplitExt splits a file name into its stem and extension.
// A leading-dot file such as ".gitignore" is treated as having no extension.
func SplitExt(name string) (string, string) {
	ext := filepath.Ext(name)
	if ext == name {
		return name, ""
	}
	return name[:len(name)-len(ext)], ext
}
