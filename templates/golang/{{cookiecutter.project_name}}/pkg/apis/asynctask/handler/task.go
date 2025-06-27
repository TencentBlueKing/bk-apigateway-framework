package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"

	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/apis/asynctask/serializer"
	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/async"
	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/infras/database"
	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/model"
	ginx2 "bk.tencent.com/{{cookiecutter.project_name}}/pkg/utils/ginx"
)

// ListTasks ...
//
//	@Summary	获取任务列表
//	@Tags		async-task
//	@Success	200	{object}	ginx.Response{data=ginx.PaginatedResp{results=[]serializer.TaskListResponse}}
//	@Router		/api/tasks [get]
func ListTasks(c *gin.Context) {
	tx := database.Client(c.Request.Context()).Order("created_at desc").Model(&model.Task{})

	// 总条目数量
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		ginx2.SetErrResp(c, http.StatusInternalServerError, err.Error())
	}

	var executedTasks []model.Task
	if err := tx.Offset(ginx2.GetOffset(c)).Limit(ginx2.GetLimit(c)).Find(&executedTasks).Error; err != nil {
		ginx2.SetErrResp(c, http.StatusInternalServerError, err.Error())
		return
	}

	respData := []serializer.TaskListResponse{}
	for _, task := range executedTasks {
		respData = append(respData, serializer.TaskListResponse{
			ID:        task.ID,
			Name:      task.Name,
			Args:      string(task.Args),
			Result:    string(task.Result),
			Creator:   task.Creator,
			StartedAt: lo.Ternary(task.StartedAt.IsZero(), "", task.StartedAt.Format(time.RFC3339)),
			Duration:  task.Duration.Seconds(),
		})
	}
	ginx2.SetResp(c, http.StatusOK, ginx2.NewPaginatedRespData(total, respData))
}

// CreateTask ...
//
//	@Summary	创建异步任务
//	@Tags		async-task
//	@Param		body	body		serializer.TaskCreateRequest	true	"异步任务配置"
//	@Success	201		{object}	ginx.Response{data=nil}
//	@Router		/api/tasks [post]
func CreateTask(c *gin.Context) {
	var req serializer.TaskCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ginx2.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := req.Validate(c); err != nil {
		ginx2.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}

	// 异步任务执行，不使用 c.Request.Context() 以避免提前 cancel
	async.ApplyTask(context.Background(), req.Name, req.Args)
	ginx2.SetResp(c, http.StatusCreated, nil)
}
