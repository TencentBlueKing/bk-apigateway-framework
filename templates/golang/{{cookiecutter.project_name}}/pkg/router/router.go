// Package router 是项目 API 服务的主路由入口
package router

import (
	"log/slog"

	"github.com/TencentBlueKing/blueapps-go/pkg/apis/basic"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	slogGin "github.com/samber/slog-gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/TencentBlueKing/blueapps-go/pkg/common"
	"github.com/TencentBlueKing/blueapps-go/pkg/config"
	"github.com/TencentBlueKing/blueapps-go/pkg/infras/otel"
	"github.com/TencentBlueKing/blueapps-go/pkg/middleware"

	"github.com/TencentBlueKing/{{cookiecutter.project_name}}/pkg/apis/user"
)

// New create server router
func New(slogger *slog.Logger) *gin.Engine {
	gin.SetMode(config.G.Service.Server.GinRunMode)
	gin.DisableConsoleColor()

	router := gin.New()

	basic.Register(router)

	// Gin 中间件
	setMiddlewares(router, slogger)
	// 后端 API
	{
		apiRG := router.Group("/api")
		// CRUD
		user.Register(apiRG)
	}

	return router
}

// 为 gin.Engine 设置中间件
// otelgin：OpenTelemetry - Gin Tracing 上报
// RequestID：在 Context，HTTP Header 中设置 Request ID
// slogGin：记录 Gin 框架结构化日志
// CORS / CSRF / CSRFToken：跨域设置 / CSRF 防护
// Recovery：请求 Panic 恢复
func setMiddlewares(router *gin.Engine, slogger *slog.Logger) {
	router.Use(otelgin.Middleware(
		otel.GenServiceName("web"),
		otelgin.WithGinFilter(
			func(c *gin.Context) bool {
				// 忽略部分路径避免过于骚扰
				excludedPaths := []string{"/metrics", "/ping"}
				return !lo.Contains(excludedPaths, c.Request.URL.Path)
			},
		),
	))
	router.Use(middleware.RequestID())
	// 替换 slogGin 配置以保持一致
	slogGin.RequestIDKey = common.RequestIDLogKey
	slogGin.SpanIDKey = common.SpanIDLogKey
	slogGin.TraceIDKey = common.TraceIDLogKey
	cfg := slogGin.Config{WithTraceID: true, WithSpanID: true, WithRequestID: true}
	router.Use(slogGin.NewWithConfig(slogger, cfg))
	router.Use(gin.Recovery())
	// 更多中间件见：https://github.com/TencentBlueKing/blueapps-go/tree/main/pkg/middleware
}
