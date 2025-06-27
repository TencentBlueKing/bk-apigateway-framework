// Package migration stores all database migrations
package migration

import (
	"context"

	log "bk.tencent.com/{{cookiecutter.project_name}}/pkg/logging"
)

// Migration 共享 Context
var ctx = context.Background()

// 打印执行迁移的日志
func logApplying(migrationID string) {
	log.Infof(ctx, "Applying migration %s", migrationID)
}

// 打印回滚迁移的日志
func logRollingBack(migrationID string) {
	log.Infof(ctx, "Rolling back migration %s", migrationID)
}
