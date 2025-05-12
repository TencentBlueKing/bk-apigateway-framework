// Package uuidx 提供一些 uuid 相关的工具
package uuidx

import (
	"strings"

	"github.com/google/uuid"
)

// New returns a new uuid (length 32 without '-')
func New() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
