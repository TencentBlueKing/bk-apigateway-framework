package config

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_contrib/model"
	"github.com/spf13/cast"

	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/utils/envx"
)

func GetApiConfig(cfg *Config) *model.APIConfig {
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
func GetGatewayConfig(cfg *Config) model.GatewayConfig {
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
		Maintainers: envx.GetEnvList("BK_APIGW_MAINTAINERS", []string{"admin"}),
	}
}

// GetStageConfig ... 获取stage配置
func GetStageConfig(cfg *Config) model.StageConfig {
	preallocatedUrls, err := envx.MustGetEnvJSON("BKPAAS_DEFAULT_PREALLOCATED_URLS")
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

	// 声明网关不同环境的环境变量
	stagEnvVars := map[string]string{
		// "foo": "bar",
	}
	prodEnvVars := map[string]string{
		// "foo": "bar",
	}
	envVars := stagEnvVars
	if cfg.Platform.RunEnv == "prod" {
		envVars = prodEnvVars
	}

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

	return model.StageConfig{
		Name:           cfg.Platform.RunEnv,
		Description:    description,
		DescriptionEn:  descriptionEn,
		BackendSubPath: appSubpath,
		BackendHost:    backendHost,
		BackendTimeout: 60,
		PluginConfigs: []*model.PluginConfig{
			stageHeaderRewritePlugin,
		},
		EnvVars: envVars,
	}
}

// GetReleaseConfig ... 获取发布配置
func GetReleaseConfig(cfg *Config) model.ReleaseConfig {
	return model.ReleaseConfig{
		// 版本号: v1.0.0+prod
		Version: fmt.Sprintf("%s+%s",
			envx.Get("BK_APIGW_RELEASE_VERSION", "1.0.1"),
			envx.MustGet("BKPAAS_ENVIRONMENT")),
		// 版本日志
		Comment: envx.Get("BK_APIGW_RELEASE_COMMENT", ""),
	}
}
