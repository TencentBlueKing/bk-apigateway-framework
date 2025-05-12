package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// QueryTokenAuth 中间件用过 query 参数中的 token 验证访问者身份，一般用于 healthz / metrics API
func QueryTokenAuth(token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		queryToken := c.DefaultQuery("token", "")
		if queryToken != token {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}
