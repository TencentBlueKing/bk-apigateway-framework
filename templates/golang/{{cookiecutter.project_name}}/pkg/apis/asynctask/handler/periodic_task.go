// Package handler ...
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/apis/asynctask/serializer"
	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/infras/database"
	model2 "bk.tencent.com/{{cookiecutter.project_name}}/pkg/model"
	ginx2 "bk.tencent.com/{{cookiecutter.project_name}}/pkg/utils/ginx"
)

// ListPeriodicTasks ...
//
//	@Summary	获取定时任务列表
//	@Tags		async-task
//	@Success	200	{object}	ginx.Response{data=[]serializer.PeriodicTaskListResponse}
//	@Router		/api/periodic-tasks [get]
func ListPeriodicTasks(c *gin.Context) {
	var periodicTasks []model2.PeriodicTask
	if err := database.Client(c.Request.Context()).Order("id DESC").Find(&periodicTasks).Error; err != nil {
		ginx2.SetErrResp(c, http.StatusInternalServerError, err.Error())
		return
	}

	respData := []serializer.PeriodicTaskListResponse{}
	for _, task := range periodicTasks {
		respData = append(respData, serializer.PeriodicTaskListResponse{
			ID:      task.ID,
			Cron:    task.Cron,
			Name:    task.Name,
			Args:    string(task.Args),
			Enabled: task.Enabled,
			Creator: task.Creator,
		})
	}
	ginx2.SetResp(c, http.StatusOK, respData)
}

// CreatePeriodicTask ...
//
//	@Summary	创建定时任务
//	@Tags		async-task
//	@Param		body	body		serializer.PeriodicTaskCreateRequest	true	"定时任务配置"
//	@Success	201		{object}	ginx.Response{data=nil}
//	@Router		/api/periodic-tasks [post]
func CreatePeriodicTask(c *gin.Context) {
	var req serializer.PeriodicTaskCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ginx2.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := req.Validate(c); err != nil {
		ginx2.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}

	args, _ := json.Marshal(req.Args)
	periodicTask := model2.PeriodicTask{
		Cron: req.Cron,
		Name: req.Name,
		Args: args,
		BaseModel: model2.BaseModel{
			Creator: ginx2.GetUserID(c),
			Updater: ginx2.GetUserID(c),
		},
	}
	if err := database.Client(c.Request.Context()).Create(&periodicTask).Error; err != nil {
		ginx2.SetErrResp(c, http.StatusInternalServerError, err.Error())
		return
	}

	ginx2.SetResp(c, http.StatusCreated, nil)
}

// DeletePeriodicTask ...
//
//	@Summary	删除定时任务
//	@Tags		async-task
//	@Param		id	path	int	true	"定时任务 ID"
//	@Success	204	"No Content"
//	@Router		/api/periodic-tasks/{id} [delete]
func DeletePeriodicTask(c *gin.Context) {
	tx := database.Client(c.Request.Context()).Where("id = ?", c.Param("id")).Delete(&model2.PeriodicTask{})
	if tx.Error != nil {
		ginx2.SetErrResp(c, http.StatusInternalServerError, tx.Error.Error())
		return
	}
	ginx2.SetResp(c, http.StatusNoContent, nil)
}

// TogglePeriodicTaskEnabled ...
//
//	@Summary	切换定时任务启用状态
//	@Tags		async-task
//	@Param		id	path	int	true	"定时任务 ID"
//	@Success	204	"No Content"
//	@Router		/api/periodic-tasks/{id}/enabled [put]
func TogglePeriodicTaskEnabled(c *gin.Context) {
	var periodicTask model2.PeriodicTask
	ctx := c.Request.Context()
	tx := database.Client(ctx).Where("id = ?", c.Param("id")).First(&periodicTask)
	if tx.Error != nil {
		ginx2.SetErrResp(c, http.StatusNotFound, tx.Error.Error())
		return
	}

	periodicTask.Enabled = !periodicTask.Enabled
	periodicTask.Updater = ginx2.GetUserID(c)

	tx = database.Client(ctx).Save(&periodicTask)
	if tx.Error != nil {
		ginx2.SetErrResp(c, http.StatusInternalServerError, tx.Error.Error())
		return
	}
	ginx2.SetResp(c, http.StatusOK, serializer.TogglePeriodicTaskEnabledResponse{Enabled: periodicTask.Enabled})
}
