package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"

	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/utils/ginx"
)

// AccessControl 用户访问控制（重要：应该在 UserAuth 中间件后使用）
func AccessControl(allowedUsers []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 没有配置 -> 允许所有用户访问
		if len(allowedUsers) == 0 {
			c.Next()
			return
		}

		// 检查用户是否可访问
		userID := ginx.GetUserID(c)
		if lo.Contains(allowedUsers, userID) {
			c.Next()
			return
		}
		c.AbortWithStatus(http.StatusForbidden)
	}
}
