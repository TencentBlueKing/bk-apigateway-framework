package serializer

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/async"
)

// TaskListResponse List Task API 返回结构
type TaskListResponse struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Args      string  `json:"args"`
	Result    string  `json:"result"`
	Creator   string  `json:"creator"`
	StartedAt string  `json:"startedAt"`
	Duration  float64 `json:"duration"`
}

// TaskCreateRequest Create Task API 请求结构
type TaskCreateRequest struct {
	Name string `json:"name"`
	Args []any  `json:"args"`
}

// Validate ...`
func (r *TaskCreateRequest) Validate(_ *gin.Context) error {
	if r.Name == "" {
		return errors.New("Task name required")
	}
	if _, ok := async.RegisteredTasks[r.Name]; !ok {
		return errors.Errorf("Task name %s invalid", r.Name)
	}
	return nil
}
