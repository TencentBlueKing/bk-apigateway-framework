// Package version 提供版本信息，构建时通过 ldflags 注入（见 Makefile）
package version

import (
	bkversion "github.com/TencentBlueKing/blueapps-go/pkg/version"
)

var (
	// AppVersion 版本号
	AppVersion = "--"
	// GitCommit CommitID
	GitCommit = "--"
	// BuildTime 二进制构建时间
	BuildTime = "--"
)

// 同步到 blueapps-go，使 /version 接口返回本项目的版本信息
func init() {
	bkversion.AppVersion = AppVersion
	bkversion.GitCommit = GitCommit
	bkversion.BuildTime = BuildTime
}
