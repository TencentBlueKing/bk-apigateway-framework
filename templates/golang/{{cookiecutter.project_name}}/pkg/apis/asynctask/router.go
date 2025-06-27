// Package asynctask ...
package asynctask

import (
	"github.com/gin-gonic/gin"

	handler2 "bk.tencent.com/{{cookiecutter.project_name}}/pkg/apis/asynctask/handler"
)

// Register ...
func Register(rg *gin.RouterGroup) {
	// task
	taskRouter := rg.Group("/tasks")
	taskRouter.GET("", handler2.ListTasks)
	taskRouter.POST("", handler2.CreateTask)

	// periodic task
	periodicTaskRouter := rg.Group("/periodic-tasks")
	periodicTaskRouter.GET("", handler2.ListPeriodicTasks)
	periodicTaskRouter.POST("", handler2.CreatePeriodicTask)
	periodicTaskRouter.DELETE("/:id", handler2.DeletePeriodicTask)
	periodicTaskRouter.PUT("/:id/enabled", handler2.TogglePeriodicTaskEnabled)
}
