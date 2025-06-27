package web

import (
	"github.com/gin-gonic/gin"

	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/web/handler"
)

// Register ...
func Register(rg *gin.RouterGroup) {
	rg.GET("", handler.GetIndexPage)
	rg.GET("home", handler.GetHomePage)
	rg.GET("crud", handler.GetCRUDPage)
	rg.GET("cache", handler.GetCachePage)
	rg.GET("cloud-api", handler.GetCloudAPIPage)
	rg.GET("async-task", handler.GetAsyncTaskPage)
	rg.GET("obj-storage", handler.GetObjStoragePage)
}
