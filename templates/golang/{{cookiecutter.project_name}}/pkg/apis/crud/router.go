// Package crud ...
package crud

import (
	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_contrib/middleware"
	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_contrib/util"
	"github.com/gin-gonic/gin"

	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_contrib/model"

	handler2 "bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/apis/crud/handler"
)

// Register ...
func Register(rg *gin.RouterGroup) {
	// category
	categoryRouter := rg.Group("/categories")

	// 使用网关鉴权中间件
	categoryRouter.Use(middleware.GatewayJWTAuthMiddleware())

	basicConfig := model.ResourceBasicConfig{
		IsPublic:             true,
		AllowApplyPermission: true,
		MatchSubpath:         true,
		EnableWebsocket:      false,
	}

	// 共用的插件配置
	{% raw %}
	headerWriterPlugin := model.BuildResourcePluginConfigWithType(
		model.PluginTypeHeaderRewrite, model.HeaderRewriteConfig{
			Set:    []model.HeaderRewriteValue{{Key: "X-Test", Value: "test"}},
			Remove: []model.HeaderRewriteValue{{Key: "X-Test2"}},
		})
	{% endraw %}

	util.RegisterBkAPIGatewayRouteWithGroup(
		categoryRouter, "GET", "",
		model.NewAPIGatewayResourceConfig(
			basicConfig,
			basicConfig.WithAuthConfig(model.AuthConfig{
				UserVerifiedRequired: true, // 用户认证
				AppVerifiedRequired:  true, // 应用认证
			}),
			basicConfig.WithPluginConfig(headerWriterPlugin)),
		handler2.ListCategories,
	)

	util.RegisterBkAPIGatewayRouteWithGroup(
		categoryRouter, "POST", "",
		model.NewAPIGatewayResourceConfig(basicConfig, basicConfig.WithAuthConfig(
			model.AuthConfig{
				UserVerifiedRequired: false, // 用户认证
				AppVerifiedRequired:  true,  // 应用认证
			}),
			basicConfig.WithPluginConfig(headerWriterPlugin),
			basicConfig.WithBackend(model.BackendConfig{
				Timeout: 10,
				Path:    "/api/v1/categories",
			})),
		handler2.CreateCategory,
	)

	categoryRouter.GET("/:id", handler2.RetrieveCategory)
	categoryRouter.PUT("/:id", handler2.UpdateCategory)
	categoryRouter.DELETE("/:id", handler2.DestroyCategory)

	// entry
	entryRouter := rg.Group("/entries")
	entryRouter.GET("", handler2.ListEntries)
	entryRouter.POST("", handler2.CreateEntry)
	entryRouter.GET("/:id", handler2.RetrieveEntry)
	entryRouter.PUT("/:id", handler2.UpdateEntry)
	entryRouter.DELETE("/:id", handler2.DestroyEntry)
}
