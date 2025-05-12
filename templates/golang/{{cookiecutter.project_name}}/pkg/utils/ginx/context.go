// Package ginx 提供一些 Gin 框架相关的工具
package ginx

import (
	"github.com/gin-gonic/gin"

	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/common"
)

// GetRequestID ...
func GetRequestID(c *gin.Context) string {
	return c.GetString(common.RequestIDCtxKey)
}

// SetRequestID ...
func SetRequestID(c *gin.Context, requestID string) {
	c.Set(common.RequestIDCtxKey, requestID)
}

// GetError ...
func GetError(c *gin.Context) (err any, ok bool) {
	return c.Get(common.ErrorCtxKey)
}

// SetError ...
func SetError(c *gin.Context, err error) {
	c.Set(common.ErrorCtxKey, err)
}

// GetUserID ...
func GetUserID(c *gin.Context) string {
	return c.GetString(common.UserIDKey)
}

// SetUserID ...
func SetUserID(c *gin.Context, userID string) {
	c.Set(common.UserIDKey, userID)
}
