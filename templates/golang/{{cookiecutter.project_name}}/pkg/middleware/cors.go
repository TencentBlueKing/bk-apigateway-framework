package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS 用于管理跨域请求
func CORS(allowedOrigins []string) gin.HandlerFunc {
	// NOTE：如果需要进行前后端分离开发，需要在 AllowMethods 中加入 `X-CSRF-Token`，并且配置合理的 AllowOrigins
	return cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
