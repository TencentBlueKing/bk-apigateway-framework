package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/spf13/viper"

	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/utils/envx"
)

var (
	pwd, _     = os.Getwd()
	exePath, _ = os.Executable()
	exeDir     = filepath.Dir(exePath)
	// BaseDir 项目根目录
	BaseDir = lo.Ternary(strings.Contains(exeDir, pwd), exeDir, pwd)
)

func getBkDomainFromEnv() string {
	return envx.Get("BKPAAS_BK_DOMAIN", "example.com")
}

func loadConfigFromFile(cfgFile string) (*Config, error) {
	// 检查配置文件是否存在
	if _, err := os.Stat(cfgFile); err != nil {
		return nil, errors.Errorf("config file %s not found", cfgFile)
	}

	// 使用 viper 从 cfgFile 加载配置
	vp := viper.New()
	vp.SetConfigFile(cfgFile)
	if err := vp.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := vp.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// 从环境变量加载配置
func loadConfigFromEnv() (*Config, error) {
	// 平台配置
	platformCfg, err := loadPlatformConfigFromEnv()
	if err != nil {
		return nil, err
	}

	// 服务配置
	serviceCfg, err := loadServiceConfigFromEnv()
	if err != nil {
		return nil, err
	}

	// 业务配置
	bizCfg, err := loadBizConfigFromEnv()
	if err != nil {
		return nil, err
	}

	return &Config{Platform: platformCfg, Service: serviceCfg, Biz: bizCfg}, nil
}

// 从环境变量加载平台配置
func loadPlatformConfigFromEnv() (PlatformConfig, error) {
	cfg := PlatformConfig{
		AppID:      envx.Get("BKPAAS_APP_ID", "{{cookiecutter.project_name}}"),
		AppSecret:  envx.MustGet("BKPAAS_APP_SECRET"),
		ModuleName: envx.Get("BKPAAS_APP_MODULE_NAME", "default"),
		RunEnv:     envx.Get("BKPAAS_ENVIRONMENT", "dev"),
		Region:     envx.Get("BKPAAS_ENGINE_REGION", "default"),
		CryptoType: envx.Get("BKPAAS_BK_CRYPTO_TYPE", "CLASSIC"),
		BkDomain:   getBkDomainFromEnv(),
		ApiUrlTmpl: envx.Get("BKPAAS_API_URL_TMPL", "http://{api_name}.apigw.example.com"),
		BkPlatUrl:  loadBkPlatUrlFromEnv(),
	}

	var err error
	cfg.Addons, err = loadAddonsConfigFromEnv()
	return cfg, err
}

// 从环境变量读取蓝鲸平台服务地址
func loadBkPlatUrlFromEnv() BkPlatUrlConfig {
	return BkPlatUrlConfig{
		BkPaaS:    strings.TrimRight(envx.Get("BKPAAS_URL", "http://bkpaas.example.com"), "/"),
		BkLogin:   strings.TrimRight(envx.Get("BKPAAS_LOGIN_URL", "http://bklogin.example.com"), "/"),
		BkCompApi: strings.TrimRight(envx.Get("BK_COMPONENT_API_URL", "http://bkapi.example.com"), "/"),
	}
}

// 从环境变量读取增强服务配置
func loadAddonsConfigFromEnv() (addons AddonsConfig, err error) {
	// Mysql
	if addons.Mysql, err = loadMysqlConfigFromEnv(); err != nil {
		return addons, err
	}

	// RabbitMQ
	if addons.RabbitMQ, err = loadRabbitMQConfigFromEnv(); err != nil {
		return addons, err
	}

	// Redis
	if addons.Redis, err = loadRedisConfigFromEnv(); err != nil {
		return addons, err
	}

	// BkRepo
	if addons.BkRepo, err = loadBkRepoConfigFromEnv(); err != nil {
		return addons, err
	}

	// BkOtel
	if addons.BkOtel, err = loadBkOtelConfigFromEnv(); err != nil {
		return addons, err
	}

	return addons, nil
}

// 判断字符串非空
func notEmpty(str string) bool {
	return str != ""
}

// 从环境变量读取 Mysql 增强服务配置
func loadMysqlConfigFromEnv() (*MysqlConfig, error) {
	host := envx.Get("GCS_MYSQL_HOST", "")
	port := envx.Get("GCS_MYSQL_PORT", "")
	name := envx.Get("GCS_MYSQL_NAME", "")
	user := envx.Get("GCS_MYSQL_USER", "")
	passwd := envx.Get("GCS_MYSQL_PASSWORD", "")
	charset := envx.Get("GCS_MYSQL_CHARSET", "utf8")

	if ok := lo.EveryBy([]string{host, port, name, user, passwd}, notEmpty); !ok {
		return nil, nil
	}
	mysqlPort, err := cast.ToIntE(port)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid GCS_MYSQL_PORT: %s", port)
	}

	return &MysqlConfig{Host: host, Port: mysqlPort, Name: name, User: user, Password: passwd, Charset: charset}, nil
}

// 从环境变量读取 RabbitMQ 增强服务配置
func loadRabbitMQConfigFromEnv() (*RabbitMQConfig, error) {
	host := envx.Get("RABBITMQ_HOST", "")
	port := envx.Get("RABBITMQ_PORT", "")
	user := envx.Get("RABBITMQ_USER", "")
	vhost := envx.Get("RABBITMQ_VHOST", "")
	passwd := envx.Get("RABBITMQ_PASSWORD", "")

	if ok := lo.EveryBy([]string{host, port, user, vhost, passwd}, notEmpty); !ok {
		return nil, nil
	}
	mqPort, err := cast.ToIntE(port)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid RABBITMQ_PORT: %s", port)
	}

	return &RabbitMQConfig{Host: host, Port: mqPort, User: user, Vhost: vhost, Password: passwd}, nil
}

// 从环境变量读取 Redis 增强服务配置
func loadRedisConfigFromEnv() (*RedisConfig, error) {
	username := envx.Get("REDIS_USERNAME", "")
	host := envx.Get("REDIS_HOST", "")
	port := envx.Get("REDIS_PORT", "")
	passwd := envx.Get("REDIS_PASSWORD", "")
	db := envx.Get("REDIS_DB", "0")

	if ok := lo.EveryBy([]string{host, port, passwd}, notEmpty); !ok {
		return nil, nil
	}
	rdsPort, err := cast.ToIntE(port)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid REDIS_PORT: %s", port)
	}
	rdsDB, err := cast.ToIntE(db)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid REDIS_DB: %s", db)
	}

	return &RedisConfig{Username: username, Host: host, Port: rdsPort, Password: passwd, DB: rdsDB}, nil
}

// 从环境变量读取 BkRepo 增强服务配置
func loadBkRepoConfigFromEnv() (*BkRepoConfig, error) {
	endpointUrl := envx.Get("BKREPO_ENDPOINT_URL", "")
	project := envx.Get("BKREPO_PROJECT", "")
	username := envx.Get("BKREPO_USERNAME", "")
	passwd := envx.Get("BKREPO_PASSWORD", "")
	bucket := envx.Get("BKREPO_BUCKET", "")
	publicBucket := envx.Get("BKREPO_PUBLIC_BUCKET", "")
	privateBucket := envx.Get("BKREPO_PRIVATE_BUCKET", "")

	if ok := lo.EveryBy(
		[]string{endpointUrl, project, username, passwd, bucket, publicBucket, privateBucket}, notEmpty,
	); !ok {
		return nil, nil
	}
	return &BkRepoConfig{
		EndpointUrl:   endpointUrl,
		Project:       project,
		Username:      username,
		Password:      passwd,
		Bucket:        bucket,
		PublicBucket:  publicBucket,
		PrivateBucket: privateBucket,
	}, nil
}

// 从环境变量读取蓝鲸 Otel 增强服务配置
func loadBkOtelConfigFromEnv() (*BkOtelConfig, error) {
	otelTrace := envx.Get("OTEL_TRACE", "")
	sampler := envx.Get("OTEL_SAMPLER", "")
	bkDataToken := envx.Get("OTEL_BK_DATA_TOKEN", "")
	grpcUrl := envx.Get("OTEL_GRPC_URL", "")

	if ok := lo.EveryBy([]string{otelTrace, sampler, bkDataToken, grpcUrl}, notEmpty); !ok {
		return nil, nil
	}
	return &BkOtelConfig{
		Trace:       strings.ToLower(otelTrace) == "true",
		Sampler:     sampler,
		BkDataToken: bkDataToken,
		GrpcUrl:     grpcUrl,
	}, nil
}

// 从环境变量读取服务配置
func loadServiceConfigFromEnv() (ServiceConfig, error) {
	// 是否为本地开发环境
	isLocalDev := envx.Get("BKPAAS_ENVIRONMENT", "dev") == "dev"

	allowedUsers := []string{}
	if val := envx.Get("ALLOWED_USERS", ""); val != "" {
		// 允许访问的用户在环境变量中格式如 "admin,userAlpha,userBeta"
		allowedUsers = strings.Split(val, ",")
	}
	// 默认允许任意源访问
	allowedOrigins := []string{"*"}
	if val := envx.Get("ALLOWED_ORIGINS", ""); val != "" {
		// 允许访问的源在环境变量中格式如 "http://localhost:8080,http://localhost:8081"
		allowedOrigins = strings.Split(val, ",")
	}
	return ServiceConfig{
		Server: ServerConfig{
			Port:         cast.ToInt(envx.Get("PORT", "8080")),
			GraceTimeout: cast.ToInt(envx.Get("GRACE_TIMEOUT", "30")),
			GinRunMode: envx.Get(
				"GIN_RUN_MODE",
				lo.Ternary[string](isLocalDev, gin.DebugMode, gin.ReleaseMode),
			),
		},
		Log: LogConfig{
			Level: envx.Get(
				"LOG_LEVEL",
				lo.Ternary(isLocalDev, "debug", "info"),
			),
			Dir: envx.Get(
				"LOG_BASE_DIR",
				lo.Ternary(isLocalDev, BaseDir+"/logs/", "/app/v3logs/"),
			),
		},
		AllowedOrigins: allowedOrigins,
		AllowedUsers:   allowedUsers,
		// CSRF cookie domain 默认使用蓝鲸根域名
		CSRFCookieDomain: envx.Get("CSRF_COOKIE_DOMAIN", getBkDomainFromEnv()),
		HealthzToken:     envx.Get("HEALTHZ_TOKEN", "healthz_token"),
		// Metric API Token
		// 需与 app_desc.yaml 中的 `spec.observability.monitoring.metrics.params.token` 保持一致
		// 否则将导致集群中的 serviceMonitor 无法获取服务指标信息，进而导致蓝鲸监控看板无数据
		MetricToken:   envx.Get("METRIC_TOKEN", "metric_token"),
		EnableSwagger: cast.ToBool(envx.Get("ENABLE_SWAGGER", lo.Ternary(isLocalDev, "true", "false"))),
		DocFileBaseDir: envx.Get(
			"DOC_FILE_BASE_DIR",
			lo.Ternary(isLocalDev, BaseDir+"/docs/", "/app/docs/"),
		),
		StaticFileBaseDir: envx.Get(
			"STATIC_FILE_BASE_DIR",
			lo.Ternary(isLocalDev, BaseDir+"/static/", "/app/static/"),
		),
		TmplFileBaseDir: envx.Get(
			"TMPL_FILE_BASE_DIR",
			lo.Ternary(isLocalDev, BaseDir+"/templates/", "/app/templates/"),
		),
		MemoryCacheSize: cast.ToInt(envx.Get("MEMORY_CACHE_SIZE", "100")),
		// DB 加密密钥，若未使用加密功能可不配置
		// 生成方式参见 Readme 文档 - 数据库字段加密
		EncryptSecret: envx.Get("ENCRYPT_SECRET", ""),
	}, nil
}

// 加载业务相关配置
func loadBizConfigFromEnv() (BizConfig, error) {
	return BizConfig{}, nil
}
