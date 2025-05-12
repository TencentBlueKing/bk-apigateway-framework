// Package filex provides utility functions for file operations.
package filex

import "path/filepath"

// GetParentDir returns the parent directory of the given path.
func GetParentDir(path string) string {
	cleaned := filepath.Clean(path)
	parent := filepath.Dir(cleaned)
	return parent
}
