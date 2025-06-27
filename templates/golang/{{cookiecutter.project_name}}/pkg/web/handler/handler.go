// Package handler ...
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"

	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/config"
	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/infras/objstorage"
	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/utils/ginx"
)

// GetIndexPage 首页
func GetIndexPage(c *gin.Context) {
	renderHTML(c, "index.html", nil)
}

// GetHomePage 主页
func GetHomePage(c *gin.Context) {
	renderHTML(c, "home.html", nil)
}

// GetCRUDPage CRUD 示例页面
func GetCRUDPage(c *gin.Context) {
	renderHTML(c, "crud.html", nil)
}

// GetCachePage 缓存示例页面
func GetCachePage(c *gin.Context) {
	renderHTML(c, "cache.html", nil)
}

// GetCloudAPIPage 云 API 调用示例页面
func GetCloudAPIPage(c *gin.Context) {
	renderHTML(c, "cloud_api.html", nil)
}

// GetAsyncTaskPage 异步任务示例页面
func GetAsyncTaskPage(c *gin.Context) {
	renderHTML(c, "async_task.html", nil)
}

// GetObjStoragePage 对象存储示例页面
func GetObjStoragePage(c *gin.Context) {
	renderHTML(c, "obj_storage.html", gin.H{"objectStorageEnabled": objstorage.IsBkRepoAvailable()})
}

func renderHTML(c *gin.Context, name string, data gin.H) {
	data = lo.Assign(data, gin.H{
		"appID": config.G.Platform.AppID,
		"user":  ginx.GetUserID(c),
	})
	c.HTML(http.StatusOK, name, data)
}
