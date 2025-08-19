// Package crud ...
package user

import (
	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_contrib/middleware"
	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_contrib/util"
	"github.com/gin-gonic/gin"

	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_contrib/model"

	"github.com/TencentBlueKing/{{cookiecutter.project_name}}/pkg/apis/user/handler"
)

// Register ...
func Register(rg *gin.RouterGroup) {
	// userRouter
	userRouter := rg.Group("/users")

	// 使用网关鉴权中间件
	userRouter.Use(middleware.GatewayJWTAuthMiddleware())

	basicConfig := model.ResourceBasicConfig{
		IsPublic:             true,
		AllowApplyPermission: true,
		MatchSubpath:         true,
		EnableWebsocket:      false,
	}

	// 共用的插件配置
	headerWriterPlugin := model.BuildResourcePluginConfigWithType(
		model.PluginTypeHeaderRewrite, model.HeaderRewriteConfig{
			Set:    []model.HeaderRewriteValue{{Key: "X-Test", Value: "test"}},
			Remove: []model.HeaderRewriteValue{{Key: "X-Test2"}},
		})

	util.RegisterBkAPIGatewayRouteWithGroup(
		userRouter, "GET", "/list",
		model.NewAPIGatewayResourceConfig(
			basicConfig,
			basicConfig.WithAuthConfig(model.AuthConfig{
				UserVerifiedRequired: true, // 用户认证
				AppVerifiedRequired:  true, // 应用认证
			}),
			basicConfig.WithPluginConfig(headerWriterPlugin)),
		handler.ListUsers,
	)

	util.RegisterBkAPIGatewayRouteWithGroup(
		userRouter, "POST", "",
		model.NewAPIGatewayResourceConfig(basicConfig, basicConfig.WithAuthConfig(
			model.AuthConfig{
				UserVerifiedRequired: false, // 用户认证
				AppVerifiedRequired:  true,  // 应用认证
			}),
			basicConfig.WithPluginConfig(headerWriterPlugin),
			basicConfig.WithBackend(model.BackendConfig{
				Timeout: 10,
				Path:    "/api/v1/users",
			}),
			// 开启mcp
			basicConfig.WithMcpEnable(true),
		),
		handler.CreateUser,
	)

}
