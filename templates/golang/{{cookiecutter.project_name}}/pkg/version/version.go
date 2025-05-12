// Package version 提供版本信息
package version

import (
	"fmt"
	"runtime"
)

var (
	// AppVersion 版本号
	AppVersion = "--"
	// GitCommit CommitID
	GitCommit = "--"
	// BuildTime 二进制构建时间
	BuildTime = "--"
	// TmplVersion 开发框架模板版本（不建议 SaaS 开发者修改该值）
	TmplVersion = "1.4.3"
	// GoVersion Go 版本号
	GoVersion = runtime.Version()
)

// Version 获取版本信息
func Version() string {
	return fmt.Sprintf(
		"\nVersion  : %s\nGitCommit: %s\nBuildTime: %s\nTmplVersion: %s\nGoVersion: %s\n",
		AppVersion, GitCommit, BuildTime, TmplVersion, GoVersion,
	)
}
