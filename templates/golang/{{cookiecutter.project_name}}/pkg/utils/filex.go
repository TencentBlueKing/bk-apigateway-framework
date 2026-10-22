// Package filex provides utility functions for file operations.
package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetParentDir returns the parent directory of the given path.
func GetParentDir(path string) string {
	cleaned := filepath.Clean(path)
	parent := filepath.Dir(cleaned)
	return parent
}

// EnsureWithinDir 校验 target（解析符号链接后）位于 baseDir 之内，防止路径遍历导致越界读写
func EnsureWithinDir(baseDir, target string) error {
	realBase, err := realPath(baseDir)
	if err != nil {
		return err
	}
	realTarget, err := realPath(target)
	if err != nil {
		return err
	}
	rel, err := filepath.Rel(realBase, realTarget)
	if err != nil {
		return err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return fmt.Errorf("path %q is outside of %q", target, baseDir)
	}
	return nil
}

// EnsureWithinWorkDir 校验 target（解析符号链接后）位于当前工作目录之内
func EnsureWithinWorkDir(target string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	return EnsureWithinDir(cwd, target)
}

// 返回解析符号链接后的绝对路径
func realPath(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(absPath)
}
