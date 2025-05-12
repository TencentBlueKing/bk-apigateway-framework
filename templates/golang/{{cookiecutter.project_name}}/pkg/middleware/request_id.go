package middleware

import (
	"context"

	"github.com/gin-gonic/gin"

	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/common"
	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/utils/ginx"
	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/utils/uuidx"
)

// RequestID 中间件用于注入 RequestID
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(common.RequestIDHeaderKey)
		// 若 RequestID 不存在或不是 32 位随机字符串，则需要重新生成并写入到 Request Header
		if len(requestID) != 32 {
			requestID = uuidx.New()
			// Request 的 Header 中需要注入 RequestID，方便 slog 中获取
			c.Request.Header.Set(common.RequestIDHeaderKey, requestID)
		}

		// 在 context 中设置 RequestID
		ctx := context.WithValue(c.Request.Context(), common.RequestIDCtxKey, requestID)
		c.Request = c.Request.WithContext(ctx)
		// 在 gin.Context 中设置 RequestID
		ginx.SetRequestID(c, requestID)
		// Writer 的 Header 中需要注入 RequestID，用于提供给请求方
		c.Writer.Header().Set(common.RequestIDHeaderKey, requestID)

		c.Next()
	}
}
