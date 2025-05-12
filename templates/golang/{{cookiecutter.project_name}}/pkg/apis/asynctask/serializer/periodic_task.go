// Package serializer ...
package serializer

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"

	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/async"
)

// PeriodicTaskListResponse List PeriodicTask API 返回结构
type PeriodicTaskListResponse struct {
	ID      int64  `json:"id"`
	Cron    string `json:"cron"`
	Name    string `json:"name"`
	Args    string `json:"args"`
	Enabled bool   `json:"enabled"`
	Creator string `json:"creator"`
}

// PeriodicTaskCreateRequest Create PeriodicTask API 请求结构
type PeriodicTaskCreateRequest struct {
	Name string `json:"name"`
	Cron string `json:"cron"`
	Args []any  `json:"args"`
}

// Validate ...
func (r *PeriodicTaskCreateRequest) Validate(_ *gin.Context) error {
	// 检查 name 是否合法
	if r.Name == "" {
		return errors.New("Task name required")
	}
	if _, ok := async.RegisteredTasks[r.Name]; !ok {
		return errors.Errorf("Task name %s invalid", r.Name)
	}
	// 检查 cron 表达式是否合法
	if r.Cron == "" {
		return errors.New("cron required")
	}
	if _, err := cron.ParseStandard(r.Cron); err != nil {
		return errors.Wrap(err, "cron invalid")
	}
	return nil
}

// TogglePeriodicTaskEnabledResponse TogglePeriodicTaskEnabled API 返回结构
type TogglePeriodicTaskEnabledResponse struct {
	Enabled bool `json:"enabled"`
}
