# bk-apigateway framework

## 说明

本框架基于 [蓝鲸开发者中心-Golang Gin](https://github.com/TencentBlueKing/blueapps-go) 开发框架用于快速开发 API 接口，部署在蓝鲸 PaaS 开发者中心，并接入到蓝鲸 API 网关。

能极大地简化开发者对接 API 网关的工作。

## 特性

1. 封装了蓝鲸 PaaS 开发者中心的相关配置，开发者只需要关心 API 的实现以及声明，无需关心部署运行时
2. 只需要按照 [swaggo](https://github.com/swaggo/swag)的方式声明接口(api扩展的蓝鲸api配置需要配合：RegisterBkAPIGatewayRouteWithGroup方法进行注册)自动生成接入蓝鲸 API 网关所需的 definition.yaml 和 resources.yaml
3. 封装了蓝鲸 API 网关的注册流程，开发者只需要在蓝鲸 PaaS 开发者中心进行发布，即可自动注册到蓝鲸 API 网关
4. 封装了蓝鲸 API 网关的调用流程，通过 `Use(middleware.GatewayJWTAuthMiddleware())`实现解析网关请求中的 jwt，并且根据接口配置校验应用或用户。


## 应用场景

1. 网络协议转换，可以将网关的 HTTP 请求转换为其他协议的请求，例如 grpc/thrift
2. 接口组合与编排，可以将多个接口组合成一个接口，也可以在代码中做一些接口编排，例如调用接口 A，根据返回结果调用接口 B 或接口 C
3. 接口协议转换，例如将遗留系统的旧协议封装成新的协议，提供给调用方，或者可以将某些字段做映射/校验/类型转换/配置默认值等


## 开发步骤

1. 在蓝鲸 PaaS 开发者中心新建 API 网关插件，会基于网关插件模板初始化一个项目，项目中包含了部署到蓝鲸 PaaS 开发者中心的相关配置/接入到蓝鲸 API 网关的相关配置。并且包含一个 Demo 示例;
2. 开发者根据需求，开发相关的 API，并且使用 [swaggo](https://github.com/swaggo/swag)的注解进行接口的声明，可以在本地开发环境 swagger ui 查看接口配置渲染的效果，并且确认是否正确，也可以使用相关命令生成 definition.yaml 和 resources.yaml 进行验证;
3. 将代码提交到仓库，并且到蓝鲸 PaaS 开发者中心 - 插件开发进行发布;
4. 发布时，会自动触发构建，将服务部署到蓝鲸 PaaS 开发者中心，并且接入到蓝鲸 API 网关，开发者可以在蓝鲸 API 网关中查看接口的配置，以及调试接口;

## 本地开发

设置环境变量 (可以在项目跟路径新建一个`.envrc`文件，将下面内容放入文件中，启动时会自动加载；也可以在启动命令行终端中手动执行下面的内容)

```bash
export DEBUG=True
export IS_LOCAL=True
export BK_APIGW_NAME="demo"
export BK_API_URL_TMPL=http://bkapi.example.com/api/{api_name}/
export BKPAAS_APP_ID="demo"
export BKPAAS_APP_SECRET=358622d8-d3e7-4522-8f16-b5530776bbb8
export BKPAAS_DEFAULT_PREALLOCATED_URLS='{"dev": "http://0.0.0.0:8080/"}'
export BKPAAS_ENVIRONMENT=dev
export BKPAAS_PROCESS_TYPE=web
```

之后启动命令

```bash
{{cookiecutter.project_name}} webserver
```
可以访问 swagger ui 地址：http://0.0.0.0:8080/swagger-ui/index.html
> 注意这个地址可以查看所有接口的文档


## 配置示例

{% raw %}

### 网关基础信息配置
详情见：[config/gateway.go](pkg/config/gateway.go)

```go
func GetApiConfig(cfg *SvcConfig) *model.APIConfig {
	apiConfig := &model.APIConfig{
		APIGateway: GetGatewayConfig(cfg),
		Stage:      GetStageConfig(cfg),
		Release:    GetReleaseConfig(cfg),
		//主动授权，网关主动给应用，添加访问网关所有资源
		GrantPermissions: model.GrantPermissionConfig{
			// 网关维度授权
			GatewayApps: []string{"app1"},
			// 资源维度
			ResourceApps: map[string][]string{
				"app1": {"get_category_by_id"},
			},
		},
		RelatedApps: []string{cfg.Platform.AppID},
		ResourceDocs: model.ResourceDocConfig{
			// 在项目 docs目录下，通过 markdown文档自动化导入中英文文档;
			// 注意markdown文件名必须等于接口的 operation_id; 见 demo 示例
			BaseDir: envx.Get("BK_APIGW_RELEASE_DOC_LANGUAGE", ""),
			// 通过swagger生成资源文档语言:zh/en, 如果配置了BK_APIGW_RESOURCE_DOCS_BASE_DIR（使用自定义文档）
			// 那么必须将这个变量置空
			Language: envx.Get("BK_APIGW_RELEASE_DOC_LANGUAGE", ""),
		},
	}
	return apiConfig
}

// GetGatewayConfig ... 获取网关配置
func GetGatewayConfig(cfg *SvcConfig) model.GatewayConfig {
	// BK_APIGW_IS_OFFICIAL is True, the BK_APIGW_NAME should be start with `bk-`
	apiType := "10"
	if cast.ToBool(envx.Get("BK_APIGW_IS_OFFICIAL", "False")) {
		apiType = "1"
	}
	return model.GatewayConfig{
		// 网关名称
		Description: fmt.Sprintf("这是应用%s 的 API 网关。由网关开发框架自动注册.", cfg.Platform.AppID),
		DescriptionEn: fmt.Sprintf(
			"This is the API Gateway for app %s."+
			" Registered automatically by the api gateway development framework.",
			cfg.Platform.AppID),
		// 网关是否公开
		IsPublic: cast.ToBool(envx.Get("BK_APIGW_IS_PUBLIC", "true")),
		APIType:  apiType,
		// 网关管理员
		Maintainers: utils.GetEnvList("BK_APIGW_MAINTAINERS", []string{"admin"}),
	}
}

// GetStageConfig ... 获取stage配置
func GetStageConfig(cfg *SvcConfig) *model.StageConfig {
	preallocatedUrls, err := utils.GetEnvJSON("BKPAAS_DEFAULT_PREALLOCATED_URLS", map[string]string{
		"prod": "https://prod.example.com",
		"stag": "https://stag.example.com",
	})
	if err != nil {
		panic(err)
	}
	appAddress := preallocatedUrls[cfg.Platform.RunEnv]
	parsedUrl, err := url.Parse(appAddress)
	if err != nil {
		panic(err)
	}
	description := "预发布环境"
	descriptionEn := "Staging Env"
	if cfg.Platform.RunEnv == "prod" {
		description = "正式环境"
		descriptionEn = "Production Env"
	}

	backendHost := fmt.Sprintf("%s://%s", parsedUrl.Scheme, parsedUrl.Host)
	appSubpath := strings.TrimRight(parsedUrl.Path, "/")

	// 设置环境插件：
	stageHeaderRewritePlugin := model.BuildStagePluginConfigWithType(
		model.PluginTypeHeaderRewrite,
		model.HeaderRewriteConfig{
			Set: []model.HeaderRewriteValue{
				{Key: "X-Real-IP", Value: "123"},
			},
			Remove: []model.HeaderRewriteValue{
				{Key: "X-Forwarded-For"},
			},
		})

	stageConfig := &model.StageConfig{
		Name:           cfg.Platform.RunEnv,
		Description:    description,
		DescriptionEn:  descriptionEn,
		BackendSubPath: appSubpath,
		BackendHost:    backendHost,
		BackendTimeout: 60,
		PluginConfigs: []*model.PluginConfig{
			stageHeaderRewritePlugin,
		},
		EnableMcpServers: cast.ToBool(envx.Get("BK_APIGW_STAGE_ENABLE_MCP_SERVERS", "true")),
	}

	// 设置 mcp server 配置
	if stageConfig.EnableMcpServers {
		stageMcpServers := []*model.McpServer{
			{
				Name:        "mcp-server",
				Description: "mcp-server",
				// 是否公开
				IsPublic: true,
				// 是否启用：0-未启用，1-启用
				Status: 1,
				// mcp server 绑定的资源列表
				Tools: []string{"create_user"},
				// 主动授权
				TargetAppCodes: []string{"app1"},
			},
		}
		stageConfig.McpServerConfigs = stageMcpServers
	}
	return stageConfig
}

// GetReleaseConfig ... 获取发布配置
func GetReleaseConfig(cfg *SvcConfig) model.ReleaseConfig {
	return model.ReleaseConfig{
		// 版本号: v1.0.0+prod
		Version: fmt.Sprintf("%s+%s",
			envx.Get("BK_APIGW_RELEASE_VERSION", "1.0.0"),
			envx.Get("BKPAAS_ENVIRONMENT", "prod")),
		// 版本日志
		Comment: envx.Get("BK_APIGW_RELEASE_COMMENT", ""),
	}
}

```

###  API蓝鲸网关扩展配置
详情见:[demo](pkg/apis/crud/router.go)
```go
    // 使用网关鉴权中间件
	categoryRouter.Use(middleware.GatewayJWTAuthMiddleware())
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
		categoryRouter, "GET", "",
		model.NewAPIGatewayResourceConfig(
			basicConfig,
			basicConfig.WithAuthConfig(model.AuthConfig{
				UserVerifiedRequired: true, // 用户认证
				AppVerifiedRequired:  true, // 应用认证
			}),
			// 设置资源名称，不设置则自动生成
            basicConfig.WithOperationID("resource_name"),
            // 开启mcp
            basicConfig.WithMcpEnable(true),
			// 设置 plugin
			basicConfig.WithPluginConfig(headerWriterPlugin)),
		handler.ListCategories,
	)
```
{% endraw %}

配置完之后，可以本地生成 definition.yaml 和 resources.yaml 进行测试

```bash
{{cookiecutter.project_name}} generate_definition_yaml && cat definition.yaml
{{cookiecutter.project_name}} generate_resources_yaml && cat resources.yaml
```
